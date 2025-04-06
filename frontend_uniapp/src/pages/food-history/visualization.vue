<template>
  <view class="visualization-container">
    <!-- 筛选区域 -->
    <view class="filter-section">
      <view class="filter-grid">
        <view class="filter-item">
          <text class="label">起始时间</text>
          <picker mode="date" :value="startDate" @change="onStartDateChange" class="picker">
            <view class="picker-text">{{startDate}}</view>
          </picker>
        </view>
        <view class="filter-item">
          <text class="label">结束时间</text>
          <picker mode="date" :value="endDate" @change="onEndDateChange" class="picker">
            <view class="picker-text">{{endDate}}</view>
          </picker>
        </view>
        <view class="filter-item">
          <text class="label">营养成分</text>
          <picker :range="nutritionOptions" :value="nutritionIndex" @change="onNutritionChange" class="picker">
            <view class="picker-text">{{nutritionOptions[nutritionIndex]}}</view>
          </picker>
        </view>
        <view class="filter-item btn-item">
          <button class="confirm-btn" @tap="onConfirm">确定</button>
        </view>
      </view>
    </view>

    <!-- 图表区域 -->
    <view class="chart-container">
      <canvas canvas-id="foodChart" id="foodChart" class="chart-canvas" @tap="tap" />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import uCharts from '@qiun/ucharts/u-charts'
import foodApi from '@/api/food';

// 日期辅助函数
// 获取多少天前的日期
function getDateBefore(days: number): Date {
  const date = new Date();
  date.setDate(date.getDate() - days);
  return date;
}

// 格式化日期为YYYY-MM-DD
function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

// 将日期字符串解析成日期对象（只保留年月日）
function parseDate(dateStr: string): Date {
  const parts = dateStr.split('-');
  if (parts.length !== 3) {
    console.error('日期格式错误:', dateStr);
    return new Date(); // 返回当前日期
  }
  // 创建只含年月日的日期对象（不含时分秒）
  return new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]), 0, 0, 0);
}

// 计算两个日期之间的天数差
function getDaysBetween(start: Date, end: Date): number {
  // 将两个日期转换为毫秒数并计算差值，然后转换为天数
  // 使用Math.floor确保得到整数天数
  const oneDay = 24 * 60 * 60 * 1000; // 一天的毫秒数
  return Math.floor((end.getTime() - start.getTime()) / oneDay);
}

// 初始化日期范围 - 默认最近7天
function initializeDefaultDateRange() {
  const today = new Date();
  const endDate = formatDate(today);
  
  // 计算7天前的日期
  const startDateObj = new Date(today);
  startDateObj.setDate(today.getDate() - 6); // 7天范围（包括今天）
  const startDate = formatDate(startDateObj);
  
  return { startDate, endDate };
}

const uChartsInstance = ref<Record<string, uCharts>>({});
const cWidth = ref(0);
const cHeight = ref(0);

// 日期选择
const startDate = ref(formatDate(getDateBefore(6))); // 默认显示近7天数据（包括今天）
const endDate = ref(formatDate(new Date()));

// 营养成分选择
const nutritionOptions = ['热量', '蛋白质', '脂肪', '碳水化合物'];
const nutritionIndex = ref(0);
const selectedNutrition = computed(() => nutritionOptions[nutritionIndex.value]);

// 加载状态
const loading = ref(false);
const error = ref('');

// 原始数据
const rawData = ref<any[]>([]);

// 事件处理函数
const onStartDateChange = (e: any) => {
  startDate.value = e.detail.value;
};

const onEndDateChange = (e: any) => {
  endDate.value = e.detail.value;
};

const onNutritionChange = (e: any) => {
  nutritionIndex.value = e.detail.value;
};

const onConfirm = () => {
  fetchFoodData();
};

// 从后端获取数据
const fetchFoodData = async () => {
  try {
    loading.value = true;
    error.value = '';
    
    // 确保开始日期不晚于结束日期
    if (new Date(startDate.value) > new Date(endDate.value)) {
      error.value = '开始日期不能晚于结束日期';
      generateMockData();
      return;
    }
    
    // 调用API获取食物数据
    const response = await foodApi.getFoodRecords({
      startDate: startDate.value,
      endDate: endDate.value
    });
    
    // 输出API响应以便调试
    console.log('API返回数据结构:', JSON.stringify(response).substring(0, 500) + '...');
    
    if (response && response.records) {
      rawData.value = response.records || [];
      
      // 检查记录中的字段
      if (rawData.value.length > 0) {
        const sampleRecord = rawData.value[0];
        console.log('第一条记录示例:', sampleRecord);
        console.log('记录日期格式:', {
          'record_time字段': sampleRecord.record_time || '不存在',
          'created_at字段': sampleRecord.created_at || '不存在',
          '日期格式示例': sampleRecord.record_time || sampleRecord.created_at || '未找到日期字段'
        });
        console.log('记录字段检查:', {
          '热量字段': {
            'calories存在?': 'calories' in sampleRecord,
            '值': sampleRecord.calories
          },
          '蛋白质字段': {
            'protein存在?': 'protein' in sampleRecord,
            '值': sampleRecord.protein
          },
          '脂肪字段': {
            'total_fat存在?': 'total_fat' in sampleRecord,
            'TotalFat存在?': 'TotalFat' in sampleRecord,
            'totalFat存在?': 'totalFat' in sampleRecord,
            'fat存在?': 'fat' in sampleRecord,
            '值': sampleRecord.total_fat || sampleRecord.TotalFat || sampleRecord.totalFat || sampleRecord.fat
          },
          '碳水化合物字段': {
            'carbohydrates存在?': 'carbohydrates' in sampleRecord,
            'Carbohydrates存在?': 'Carbohydrates' in sampleRecord,
            'carbs存在?': 'carbs' in sampleRecord,
            '值': sampleRecord.carbohydrates || sampleRecord.Carbohydrates || sampleRecord.carbs
          },
          '所有可能的脂肪相关字段': Object.keys(sampleRecord).filter(key => 
            key.toLowerCase().includes('fat')
          ),
          '所有可能的碳水相关字段': Object.keys(sampleRecord).filter(key => 
            key.toLowerCase().includes('carb')
          )
        });
      }
      
      // 检查是否有数据
      if (rawData.value.length === 0) {
        error.value = '所选时间段内没有饮食记录';
        generateMockData();
      } else {
        processAndDisplayData();
      }
    } else {
      error.value = '无法获取数据，将使用模拟数据';
      // 使用模拟数据
      generateMockData();
    }
  } catch (err: any) {
    console.error('获取食物数据失败:', err);
    error.value = err.message || '网络错误，无法获取数据';
    // 使用模拟数据
    generateMockData();
  } finally {
    loading.value = false;
  }
};

// 生成模拟数据（当API请求失败时使用）
const generateMockData = () => {
  const start = parseDate(startDate.value);
  const end = parseDate(endDate.value);
  const days = getDaysBetween(start, end);
  
  if (days < 0) {
    console.error('日期范围错误：结束日期早于开始日期');
    return;
  }
  
  // 模拟生成日期区间内的数据
  const mockData = [];
  const currentDate = new Date(start);
  
  // 一天生成多条记录
  for (let i = 0; i <= days; i++) {
    // 一天3餐
    const mealTypes = ['早餐', '午餐', '晚餐'];
    
    for (let j = 0; j < mealTypes.length; j++) {
      const recordTime = new Date(currentDate);
      
      // 设置不同餐点的时间
      switch (j) {
        case 0: recordTime.setHours(8, 0, 0); break;  // 早餐 8:00
        case 1: recordTime.setHours(12, 0, 0); break; // 午餐 12:00
        case 2: recordTime.setHours(18, 0, 0); break; // 晚餐 18:00
      }
      
      // 格式化日期为后端格式 "YYYY-MM-DD HH:MM:SS.SSS"
      const formattedDate = recordTime.toISOString().replace('T', ' ').split('.')[0] + '.000';
      
      // 随机生成食物记录，使用与后端一致的字段名
      mockData.push({
        id: mockData.length + 1,
        food_name: `模拟食物 ${i+1}-${j+1}`,
        weight: Math.floor(Math.random() * 300) + 100, // 100-400g
        calories: Math.floor(Math.random() * 500) + 300, // 300-800 kcal
        protein: Math.floor(Math.random() * 20) + 10, // 10-30g
        carbohydrates: Math.floor(Math.random() * 40) + 30, // 30-70g
        total_fat: Math.floor(Math.random() * 15) + 5, // 5-20g
        meal_type: mealTypes[j],
        created_at: formattedDate,  // 使用后端格式
        record_time: formattedDate,  // 保留原格式以兼容
        notes: ''
      });
    }
    
    currentDate.setDate(currentDate.getDate() + 1);
  }
  
  rawData.value = mockData;
  console.log('生成的模拟数据:', mockData[0]);
  processAndDisplayData();
};

// 处理数据并展示
const processAndDisplayData = () => {
  const start = parseDate(startDate.value);
  const end = parseDate(endDate.value);
  const days = getDaysBetween(start, end);
  
  // 根据日期区间长度决定处理方式
  let processedData;
  
  // 无论时间跨度多长，都限制最多显示7个数据点
  if (days <= 7) {
    // 7天内，直接显示每天数据
    processedData = processDataByDay();
  } else if (days <= 30) {
    // 8-30天，按照周分组
    processedData = processDataByWeek();
  } else {
    // 30天以上，按月分组
    processedData = processDataByMonth();
  }
  
  // 如果数据点仍然超过7个，进行采样
  if (processedData.categories.length > 7) {
    processedData = sampleDataPoints(processedData, 7);
  }
  
  drawCharts('foodChart', processedData);
};

// 数据采样函数，将数据点数量限制为指定数量
const sampleDataPoints = (data: any, maxPoints: number) => {
  const { categories, series } = data;
  
  // 如果数据点数量不超过最大值，直接返回
  if (categories.length <= maxPoints) {
    return data;
  }
  
  // 计算采样间隔
  const interval = Math.ceil(categories.length / maxPoints);
  
  // 采样后的类别和数据
  const sampledCategories: string[] = [];
  const sampledSeries = series.map((serie: any) => ({
    ...serie,
    data: []
  }));
  
  // 进行采样
  for (let i = 0; i < categories.length; i += interval) {
    sampledCategories.push(categories[i]);
    
    // 对每个系列的数据进行采样
    series.forEach((serie: any, serieIndex: number) => {
      sampledSeries[serieIndex].data.push(serie.data[i]);
    });
  }
  
  // 确保我们有足够的点（可能由于四舍五入导致采样点不足）
  // 添加最后一个点（如果尚未包含）
  if (sampledCategories.length < maxPoints && categories.length > interval) {
    const lastIndex = categories.length - 1;
    if (sampledCategories[sampledCategories.length - 1] !== categories[lastIndex]) {
      sampledCategories.push(categories[lastIndex]);
      series.forEach((serie: any, serieIndex: number) => {
        sampledSeries[serieIndex].data.push(serie.data[lastIndex]);
      });
    }
  }
  
  return {
    categories: sampledCategories,
    series: sampledSeries
  };
};

// 根据选择的营养成分获取对应的数据键
const getNutritionKey = (): string => {
  switch (nutritionOptions[nutritionIndex.value]) {
    case '蛋白质': return 'protein';
    case '脂肪': return 'total_fat';  // 修改为后端对应的字段名
    case '碳水化合物': return 'carbohydrates';  // 修改为后端对应的字段名
    case '热量':
    default: return 'calories';
  }
};

// 获取记录中的营养值，考虑可能的字段名变体
const getNutritionValue = (item: any, nutritionKey: string): number => {
  // 直接匹配
  if (item[nutritionKey] !== undefined && item[nutritionKey] !== null) {
    return Number(item[nutritionKey]);
  }
  
  // 尝试不同的字段名变体，针对后端模型进行调整
  const alternativeKeys: Record<string, string[]> = {
    'total_fat': ['fat', 'TotalFat', 'totalFat', '脂肪', 'fatContent', 'fat_content'],
    'protein': ['Protein', '蛋白质', 'proteinContent', 'protein_content'],
    'carbohydrates': ['carbs', 'Carbohydrates', '碳水化合物', 'carbohydrate', 'carb', 'carbContent', 'carb_content'],
    'calories': ['Calories', '热量', 'calorie', 'calorieContent', 'calorie_content', 'energy']
  };
  
  // 检查替代字段名
  const alternatives = alternativeKeys[nutritionKey] || [];
  for (const altKey of alternatives) {
    if (item[altKey] !== undefined && item[altKey] !== null) {
      return Number(item[altKey]);
    }
  }
  
  // 尝试不区分大小写的匹配
  const lowerCaseKey = nutritionKey.toLowerCase();
  for (const key of Object.keys(item)) {
    if (key.toLowerCase() === lowerCaseKey || 
        key.toLowerCase().includes(lowerCaseKey)) {
      return Number(item[key]);
    }
  }
  
  // 如果都找不到，返回0
  console.warn(`未能找到字段 ${nutritionKey} 的任何匹配项`, item);
  return 0;
};

// 按天处理数据
const processDataByDay = () => {
  const data = rawData.value;
  const nutritionKey = getNutritionKey();
  
  // 打印原始数据，帮助调试
  console.log('选择的营养成分:', nutritionKey);
  console.log('原始数据示例:', data.slice(0, 3));
  
  // 获取选定日期范围
  const start = parseDate(startDate.value);
  const end = parseDate(endDate.value);
  
  // 过滤出日期范围内的数据
  const filteredData = data.filter(item => {
    // 检查日期字段可能存在的不同格式
    const recordTimeStr = item.record_time || item.created_at || '';
    
    // 处理后端格式的日期 ("2025-03-29 16:53:07.162")
    let itemDate;
    if (recordTimeStr.includes(' ') && !recordTimeStr.includes('T')) {
      // 后端格式日期处理
      const dateTimeParts = recordTimeStr.split(' ');
      itemDate = new Date(dateTimeParts[0] + 'T' + dateTimeParts[1]);
    } else {
      // ISO格式或其他格式
      itemDate = new Date(recordTimeStr);
    }
    
    const itemDateOnly = new Date(itemDate.getFullYear(), itemDate.getMonth(), itemDate.getDate());
    return itemDateOnly >= start && itemDateOnly <= end;
  });
  
  console.log(`过滤后数据数量: ${filteredData.length}`);
  
  // 如果过滤后没有数据，返回空数据结构
  if (filteredData.length === 0) {
    return { categories: [], series: [{ name: `每日${selectedNutrition.value}总量`, data: [], type: 'line' }] };
  }
  
  // 按日期排序
  filteredData.sort((a, b) => new Date(a.record_time).getTime() - new Date(b.record_time).getTime());
  
  // 按日期分组并累计每日总摄入量
  const dailyTotals: Record<string, number> = {};
  
  filteredData.forEach(item => {
    // 获取日期字符串
    const recordTimeStr = item.record_time || item.created_at || '';
    let dateStr;
    
    // 处理后端格式的日期 ("2025-03-29 16:53:07.162")
    if (recordTimeStr.includes(' ') && !recordTimeStr.includes('T')) {
      dateStr = recordTimeStr.split(' ')[0];
    } else {
      dateStr = recordTimeStr.split('T')[0];
    }
    
    if (!dailyTotals[dateStr]) {
      dailyTotals[dateStr] = 0;
    }
    
    // 使用兼容函数获取营养值并累加
    const value = getNutritionValue(item, nutritionKey);
    dailyTotals[dateStr] += value;
    
    console.log(`日期: ${dateStr}, 单次值: ${value}, 累计: ${dailyTotals[dateStr]}`);
  });
  
  // 打印处理后的数据，帮助调试
  console.log('按日期累计的数据:', dailyTotals);
  
  // 转换为数组格式
  const sortedDates = Object.keys(dailyTotals).sort();
  const categories = sortedDates.map(date => date.substring(5)); // 只显示MM-DD格式
  
  const nutritionData = sortedDates.map(date => {
    return Math.round(dailyTotals[date]);
  });
  
  console.log('图表数据:', {
    categories: categories,
    values: nutritionData
  });
  
  return {
    categories,
    series: [{
      name: `每日${selectedNutrition.value}总量`,
      data: nutritionData,
      type: 'line'
    }]
  };
};

// 按周分组处理数据
const processDataByWeek = () => {
  const data = rawData.value;
  const nutritionKey = getNutritionKey();
  
  // 获取选定日期范围
  const start = parseDate(startDate.value);
  const end = parseDate(endDate.value);
  
  // 过滤出日期范围内的数据
  const filteredData = data.filter(item => {
    // 检查日期字段可能存在的不同格式
    const recordTimeStr = item.record_time || item.created_at || '';
    
    // 处理后端格式的日期 ("2025-03-29 16:53:07.162")
    let itemDate;
    if (recordTimeStr.includes(' ') && !recordTimeStr.includes('T')) {
      // 后端格式日期处理
      const dateTimeParts = recordTimeStr.split(' ');
      itemDate = new Date(dateTimeParts[0] + 'T' + dateTimeParts[1]);
    } else {
      // ISO格式或其他格式
      itemDate = new Date(recordTimeStr);
    }
    
    const itemDateOnly = new Date(itemDate.getFullYear(), itemDate.getMonth(), itemDate.getDate());
    return itemDateOnly >= start && itemDateOnly <= end;
  });
  
  // 如果过滤后没有数据，返回空数据结构
  if (filteredData.length === 0) {
    return { categories: [], series: [{ name: `周${selectedNutrition.value}总量`, data: [], type: 'line' }] };
  }
  
  // 首先按日期组织数据
  const dailyTotals: Record<string, number> = {};
  
  filteredData.forEach(item => {
    // 获取日期字符串
    const recordTimeStr = item.record_time || item.created_at || '';
    let dateStr;
    
    // 处理后端格式的日期 ("2025-03-29 16:53:07.162")
    if (recordTimeStr.includes(' ') && !recordTimeStr.includes('T')) {
      dateStr = recordTimeStr.split(' ')[0];
    } else {
      dateStr = recordTimeStr.split('T')[0];
    }
    
    if (!(dateStr in dailyTotals)) {
      dailyTotals[dateStr] = 0;
    }
    dailyTotals[dateStr] += getNutritionValue(item, nutritionKey);
  });
  
  console.log('每日总量:', dailyTotals);
  
  // 然后按周分组
  const weekTotals: Record<string, { total: number, dates: string[] }> = {};
  
  Object.entries(dailyTotals).forEach(([dateStr, dailyTotal]) => {
    const date = new Date(dateStr);
    
    // 获取周数 (从当年的第一天开始计算)
    const firstDayOfYear = new Date(date.getFullYear(), 0, 1);
    const pastDaysOfYear = (date.getTime() - firstDayOfYear.getTime()) / 86400000;
    const weekNum = Math.ceil((pastDaysOfYear + firstDayOfYear.getDay() + 1) / 7);
    
    // 生成周标识，如"2023-W12"
    const weekKey = `${date.getFullYear()}-W${weekNum}`;
    
    if (!(weekKey in weekTotals)) {
      weekTotals[weekKey] = { total: 0, dates: [] };
    }
    
    weekTotals[weekKey].total += dailyTotal;
    weekTotals[weekKey].dates.push(dateStr);
  });
  
  console.log('周累计:', weekTotals);
  
  // 排序周标识
  const sortedWeeks = Object.keys(weekTotals).sort();
  
  // 生成可读的周标签（例如："第1周"）
  const categories = sortedWeeks.map((key, index) => {
    const firstDate = weekTotals[key].dates.sort()[0];
    const lastDate = weekTotals[key].dates.sort().pop();
    return `第${index + 1}周 (${firstDate.substring(5)}-${lastDate?.substring(5)})`;
  });
  
  const nutritionData = sortedWeeks.map(key => {
    return Math.round(weekTotals[key].total);
  });
  
  console.log('图表数据:', {
    categories: categories,
    values: nutritionData
  });
  
  return {
    categories,
    series: [{
      name: `周${selectedNutrition.value}总量`,
      data: nutritionData,
      type: 'line'
    }]
  };
};

// 按月分组处理数据
const processDataByMonth = () => {
  const data = rawData.value;
  const nutritionKey = getNutritionKey();
  
  // 获取选定日期范围
  const start = parseDate(startDate.value);
  const end = parseDate(endDate.value);
  
  // 过滤出日期范围内的数据
  const filteredData = data.filter(item => {
    // 检查日期字段可能存在的不同格式
    const recordTimeStr = item.record_time || item.created_at || '';
    
    // 处理后端格式的日期 ("2025-03-29 16:53:07.162")
    let itemDate;
    if (recordTimeStr.includes(' ') && !recordTimeStr.includes('T')) {
      // 后端格式日期处理
      const dateTimeParts = recordTimeStr.split(' ');
      itemDate = new Date(dateTimeParts[0] + 'T' + dateTimeParts[1]);
    } else {
      // ISO格式或其他格式
      itemDate = new Date(recordTimeStr);
    }
    
    const itemDateOnly = new Date(itemDate.getFullYear(), itemDate.getMonth(), itemDate.getDate());
    return itemDateOnly >= start && itemDateOnly <= end;
  });
  
  // 如果过滤后没有数据，返回空数据结构
  if (filteredData.length === 0) {
    return { categories: [], series: [{ name: `月${selectedNutrition.value}总量`, data: [], type: 'line' }] };
  }
  
  // 首先按日期组织数据
  const dailyTotals: Record<string, number> = {};
  
  filteredData.forEach(item => {
    // 获取日期字符串
    const recordTimeStr = item.record_time || item.created_at || '';
    let dateStr;
    
    // 处理后端格式的日期 ("2025-03-29 16:53:07.162")
    if (recordTimeStr.includes(' ') && !recordTimeStr.includes('T')) {
      dateStr = recordTimeStr.split(' ')[0];
    } else {
      dateStr = recordTimeStr.split('T')[0];
    }
    
    if (!dailyTotals[dateStr]) {
      dailyTotals[dateStr] = 0;
    }
    
    // 使用兼容函数获取营养值并累加
    const value = getNutritionValue(item, nutritionKey);
    dailyTotals[dateStr] += value;
    
    console.log(`日期: ${dateStr}, 单次值: ${value}, 累计: ${dailyTotals[dateStr]}`);
  });
  
  // 然后按月分组
  const monthTotals: Record<string, number> = {};
  
  Object.entries(dailyTotals).forEach(([dateStr, dailyTotal]) => {
    // 获取年月 (YYYY-MM)
    const monthKey = dateStr.substring(0, 7);
    
    if (!(monthKey in monthTotals)) {
      monthTotals[monthKey] = 0;
    }
    
    monthTotals[monthKey] += dailyTotal;
  });
  
  console.log('月累计:', monthTotals);
  
  // 排序月标识
  const sortedMonths = Object.keys(monthTotals).sort();
  const categories = sortedMonths.map(key => {
    // 将YYYY-MM转换为更友好的格式，如"2023年5月"
    const [year, month] = key.split('-');
    return `${year}年${parseInt(month)}月`;
  });
  
  const nutritionData = sortedMonths.map(key => {
    return Math.round(monthTotals[key]);
  });
  
  console.log('图表数据:', {
    categories: categories,
    values: nutritionData
  });
  
  return {
    categories,
    series: [{
      name: `月${selectedNutrition.value}总量`,
      data: nutritionData,
      type: 'line'
    }]
  };
};

// 初始化
onMounted(() => {
  const systemInfo = uni.getSystemInfoSync();
  // 获取屏幕宽度
  cWidth.value = systemInfo.windowWidth;
  // 设置图表高度为屏幕宽度的0.8倍
  cHeight.value = systemInfo.windowWidth * 0.8;
  
  // 初始加载数据
  fetchFoodData();
});

// 图表渲染函数
const drawCharts = (id: string, data: any) => {
  const ctx = uni.createCanvasContext(id);
  uChartsInstance.value[id] = new uCharts({
    type: 'line', // 默认使用折线图
    context: ctx,
    width: cWidth.value * 0.92,
    height: cHeight.value * 0.9,
    categories: data.categories,
    series: data.series,
    padding: [45, 30, 30, 30],
    background: '#000000',
    pixelRatio: 1,
    rotate: false,
    rotateLock: true,
    fontColor: "#FFFFFF",
    xAxis: {
      disableGrid: false,
      gridColor: '#666666',
      gridType: 'dash',
      dashLength: 4,
      labelCount: Math.min(data.categories.length, 7), // 最多显示7个标签
      boundaryGap: 'center',
      axisLine: true,
      axisLineColor: '#666666',
      fontSize: 11,
      rotateLabel: data.categories.length > 7 // 数据多时旋转标签
    },
    yAxis: {
      disableGrid: false,
      gridType: 'dash',
      gridColor: '#666666',
      dashLength: 4,
      splitNumber: 5,
      format: (val: number) => {
        // 根据不同营养成分添加单位
        const unit = getNutritionUnit();
        return val + unit;
      },
      min: 0,
      axisLine: true,
      axisLineColor: '#666666'
    },
    extra: {
      line: {
        type: 'curve', // 曲线
        width: 3,
        activeType: 'hollow',
        linearType: 'custom',
        color: ['#00DFFF', '#8A2BE2']
      }
    },
    legend: {
      show: true,
      position: 'top',
      float: 'center',
      padding: 15,
      margin: 10,
      backgroundColor: 'rgba(0,0,0,0.3)',
      borderColor: 'rgba(255,255,255,0.1)',
      borderWidth: 1,
      itemGap: 15,
      fontSize: 14,
      lineHeight: 25,
      fontColor: "#FFFFFF"
    }
  });
};

// 获取营养成分对应的单位
const getNutritionUnit = (): string => {
  switch (nutritionOptions[nutritionIndex.value]) {
    case '蛋白质': return 'g';
    case '脂肪': return 'g';
    case '碳水化合物': return 'g';
    case '热量':
    default: return 'kcal';
  }
};

const tap = (e: any) => {
  uChartsInstance.value[e.target.id].touchLegend(e);
  uChartsInstance.value[e.target.id].showToolTip(e);
};
</script>

<style>
.visualization-container {
  min-height: 100vh;
  background-image: url('/static/images/food-record-bj.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  padding: 30rpx;
  position: relative;
  box-sizing: border-box;
}

.visualization-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1;
}

.filter-section {
  position: relative;
  z-index: 2;
  background: rgba(40, 44, 52, 0.7);
  border-radius: 12rpx;
  padding: 30rpx;
  margin-bottom: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.5);
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300rpx, 1fr));
  gap: 30rpx;
}

.filter-item {
  display: flex;
  flex-direction: column;
}

.label {
  color: #fff;
  font-size: 28rpx;
  margin-bottom: 16rpx;
  font-weight: 500;
}

.picker {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 8rpx;
  padding: 16rpx 20rpx;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.picker-text {
  color: #fff;
  font-size: 28rpx;
}

.btn-item {
  justify-content: flex-end;
}

.confirm-btn {
  min-width: 160rpx;
  height: 72rpx;
  line-height: 72rpx;
  text-align: center;
  background: #00DFFF;
  color: #fff;
  border-radius: 36rpx;
  font-size: 28rpx;
  font-weight: bold;
  padding: 0 40rpx;
  border: none;
  box-shadow: 0 4rpx 12rpx rgba(0, 223, 255, 0.4);
  margin-top: 44rpx;
}

.chart-container {
  position: relative;
  z-index: 2;
  width: 100%;
  height: calc(100vh - 280rpx);
  background: rgba(40, 44, 52, 0.7);
  border-radius: 12rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.5);
  padding: 20rpx 20rpx 40rpx 20rpx;
  display: flex;
  justify-content: center;
  align-items: center;
}

.chart-canvas {
  width: 100%;
  height: 100%;
}

/* 添加uni-canvas-canvas样式 */
.uni-canvas-canvas {
  width: 100% !important;
  height: 100% !important;
  display: block !important;
  margin: 0 auto;
}

/* 媒体查询适配不同设备 */
@media screen and (min-width: 768px) {
  .filter-grid {
    grid-template-columns: repeat(4, 1fr);
  }
  
  .filter-item:last-child {
    justify-content: center;
    align-items: flex-end;
  }
}

@media screen and (max-width: 767px) {
  .filter-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .btn-item {
    grid-column: 2;
    align-items: flex-end;
  }
}

@media screen and (max-width: 375px) {
  .filter-grid {
    grid-template-columns: 1fr;
  }
  
  .btn-item {
    grid-column: 1;
    align-items: center;
  }
  
  .confirm-btn {
    width: 100%;
  }
}
</style>
