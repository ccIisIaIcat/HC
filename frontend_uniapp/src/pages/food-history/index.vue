<template>
  <view class="food-history-container">
    <view class="history-header">
      <text class="title">我的饮食记录</text>
      <view class="header-actions">
        <button class="back-button" @tap="handleBack">返回仪表盘</button>
        <button class="add-button" @tap="handleAddRecord">添加记录</button>
      </view>
    </view>

    <view class="history-content">
      <!-- 加载状态 -->
      <view class="loading-indicator" v-if="loading">
        <text>正在加载数据...</text>
      </view>

      <!-- 错误信息 -->
      <view class="error-message" v-if="error">
        <text>{{ error }}</text>
      </view>

      <!-- 空状态 -->
      <view class="empty-records" v-if="!loading && !error && records.length === 0">
        <text>暂无饮食记录</text>
        <button class="add-button" @tap="handleAddRecord">开始记录饮食</button>
      </view>

      <!-- 记录列表 -->
      <block v-if="!loading && !error && records.length > 0">
        <view class="records-summary">
          <text>共 {{ totalRecords }} 条记录，显示 {{ (currentPage - 1) * recordsPerPage + 1 }} - {{ Math.min(currentPage * recordsPerPage, totalRecords) }} 条</text>
        </view>

        <view class="records-list">
          <view class="record-item" v-for="(record, index) in records" :key="record.id || record.ID || `record-${index}`">
            <view class="record-time">{{ formatDateTime(record.record_time) }}</view>
            <view class="record-info">
              <text class="record-food-name">{{ record.food_name }}</text>
              <text class="record-meal-type">{{ record.meal_type }}</text>
            </view>
            <view class="record-calories">{{ (record.calories || 0).toFixed(0) }} 卡路里</view>
            <view class="record-actions">
              <button class="edit-button" @tap="handleEdit(record)">编辑</button>
              <button class="delete-button" @tap="handleDeleteRecord(record.id || record.ID)">
                <text class="delete-icon">删除</text>
              </button>
            </view>
            <view class="cyber-line"></view>
          </view>
        </view>

        <!-- 分页控件 -->
        <view class="pagination">
          <button 
            class="pagination-button" 
            :disabled="currentPage === 1"
            @tap="handlePageChange(1)"
          >首页</button>
          <button 
            class="pagination-button" 
            :disabled="currentPage === 1"
            @tap="handlePageChange(currentPage - 1)"
          >上一页</button>
          <text class="pagination-info">{{ currentPage }} / {{ totalPages }}</text>
          <button 
            class="pagination-button" 
            :disabled="currentPage === totalPages"
            @tap="handlePageChange(currentPage + 1)"
          >下一页</button>
          <button 
            class="pagination-button" 
            :disabled="currentPage === totalPages"
            @tap="handlePageChange(totalPages)"
          >末页</button>
        </view>
      </block>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
// import api from '@/api';
import foodApi from '@/api/food';

// 调试代码 - 检查导入的 API 模块
console.log('导入的 food API:', foodApi);
console.log('food API 类型:', typeof foodApi);
console.log('food API 方法:', Object.keys(foodApi));

interface FoodRecord {
  id: number;
  ID?: number; // 添加可选的大写ID字段
  food_name: string;
  weight: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  meal_type: string;
  record_time: string;
  notes: string;
}

interface FoodRecordsResponse {
  records: FoodRecord[];
  total: number;
}

const records = ref<FoodRecord[]>([]);
const loading = ref(true);
const error = ref('');
const currentPage = ref(1);
const totalPages = ref(1);
const totalRecords = ref(0);
const recordsPerPage = 10;

// 获取食物记录列表
const fetchFoodRecords = async (retryCount = 3) => {
  try {
    loading.value = true;
    error.value = '';
    
    // 计算六个月前的日期
    const today = new Date();
    const sixMonthsAgo = new Date();
    sixMonthsAgo.setMonth(today.getMonth() - 6);
    
    // 格式化为YYYY-MM-DD
    const startDate = sixMonthsAgo.toISOString().split('T')[0];
    const endDate = today.toISOString().split('T')[0];
    
    // 记录将要发送的请求参数
    console.log(`将发送请求到: /api/food-records，参数: startDate=${startDate}, endDate=${endDate}`);
    
    // 从API获取数据
    try {
      // 使用食物API获取记录
      const response = await foodApi.getFoodRecords({
        startDate,
        endDate
      }) as unknown as FoodRecordsResponse;
      
      // 获取记录数组
      const recordsArray = response.records || [];
      
      // 计算总页数
      const totalItems = response.total || recordsArray.length;
      totalRecords.value = totalItems;
      totalPages.value = Math.ceil(totalItems / recordsPerPage);
      
      // 获取当前页的数据
      const startIndex = (currentPage.value - 1) * recordsPerPage;
      const endIndex = startIndex + recordsPerPage;
      const currentPageData = recordsArray.slice(startIndex, endIndex);
      
      records.value = currentPageData;
    } catch (apiError: any) {
      console.error('API调用错误:', apiError);
      
      // 检查是否是 URLSearchParams 不存在的错误
      if (apiError.message && apiError.message.includes('URLSearchParams is not defined')) {
        // 向用户提示错误，并建议更新应用
        throw new Error('您的设备不支持某些功能，请尝试更新应用或使用较新的设备');
      } else {
        // 其他错误继续抛出
        throw apiError;
      }
    }
  } catch (err: any) {
    console.error('获取食物记录失败:', err);
    // 构建详细的错误信息
    let errorMessage = '获取食物记录失败: ';
    
    if (err.errMsg) {
      errorMessage += err.errMsg;
    } else if (err.message) {
      errorMessage += err.message;
    }
    
    // 添加更多错误详情
    if (err.statusCode) {
      errorMessage += ` (状态码: ${err.statusCode})`;
    }
    if (err.data && err.data.message) {
      errorMessage += ` - ${err.data.message}`;
    }
    
    error.value = errorMessage;
    
    // 如果是请求取消且还有重试次数，则延迟后重试
    if ((err.errMsg?.includes('abort') || err.message?.includes('请求已取消')) && retryCount > 0) {
      console.log(`将在1秒后重试，剩余重试次数: ${retryCount - 1}`, '完整错误信息:', err);
      setTimeout(() => {
        fetchFoodRecords(retryCount - 1);
      }, 1000);
    }
  } finally {
    loading.value = false;
  }
};

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return '未知时间';
  
  try {
    const date = new Date(dateString);
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      return '日期无效';
    }
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (error) {
    console.error('日期格式化错误:', error);
    return '格式错误';
  }
};

// 处理页面导航
const handlePageChange = (newPage: number) => {
  if (newPage >= 1 && newPage <= totalPages.value) {
    currentPage.value = newPage;
    fetchFoodRecords();
  }
};

// 返回仪表盘
const handleBack = () => {
  uni.navigateBack();
};

// 编辑食物记录
const handleEdit = (record: FoodRecord) => {
  const recordId = record.id || (record as any).ID;
  if (!recordId) {
    console.error('编辑错误：记录没有ID', record);
    uni.showToast({
      title: '无法编辑：记录ID不存在',
      icon: 'none'
    });
    return;
  }
  console.log('编辑记录:', record, '使用ID:', recordId);
  uni.navigateTo({
    url: `/pages/food-history/edit?id=${recordId}`
  });
};

// 删除食物记录
const handleDeleteRecord = async (id: number | undefined) => {
  if (!id) {
    uni.showToast({
      title: '无法删除：记录ID不存在',
      icon: 'none'
    });
    return;
  }
  
  console.log('准备删除记录，ID:', id);
  
  uni.showModal({
    title: '提示',
    content: '确定要删除这条食物记录吗？此操作不可撤销。',
    success: async (res) => {
      if (res.confirm) {
        try {
          loading.value = true;
          await foodApi.deleteFoodRecord(id);
          uni.showToast({
            title: '记录已成功删除',
            icon: 'success'
          });
          // 重新获取数据
          fetchFoodRecords();
        } catch (err) {
          console.error('删除记录失败:', err);
          uni.showToast({
            title: '删除记录失败，请稍后再试',
            icon: 'none'
          });
        } finally {
          loading.value = false;
        }
      }
    }
  });
};

// 添加新记录（跳转到食物分析页面）
const handleAddRecord = () => {
  uni.navigateTo({
    url: '/pages/food-analysis/index'
  });
};

// 页面加载时获取数据
onMounted(() => {
  // 延迟一段时间再获取数据，等待页面和网络就绪
  setTimeout(() => {
    fetchFoodRecords();
  }, 500);
});
</script>

<style>
.food-history-container {
  min-height: 100vh;
  background-color: #000000;
  padding: 20px;
  position: relative;
  background-image: url('/static/images/food-record-bj.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.food-history-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  z-index: 1;
}

.history-header, .history-content {
  position: relative;
  z-index: 2;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.title {
  color: #00DFFF;
  font-size: 24px;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 2px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5),
               0 0 20px rgba(0, 223, 255, 0.3),
               0 0 30px rgba(0, 223, 255, 0.1);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.back-button, .add-button {
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  color: #fff;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.back-button {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
  position: relative;
  overflow: hidden;
}

.back-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.2),
    transparent
  );
  transition: 0.5s;
}

.back-button:active::before {
  left: 100%;
}

.add-button {
  background: linear-gradient(45deg, #00DFFF, #00BFFF);
  border: none;
  position: relative;
  overflow: hidden;
  box-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
}

.add-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  transition: 0.5s;
}

.add-button:active::before {
  left: 100%;
}

.loading-indicator {
  text-align: center;
  padding: 20px;
  color: #00DFFF;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
}

.error-message {
  background: rgba(255, 77, 79, 0.1);
  border: 1px solid rgba(255, 77, 79, 0.3);
  color: #ff4d4f;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 20px;
  text-shadow: 0 0 10px rgba(255, 0, 98, 0.5);
}

.empty-records {
  text-align: center;
  padding: 40px 20px;
  color: rgba(255, 255, 255, 0.6);
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.3);
}

.records-summary {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  margin-bottom: 10px;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.3);
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.record-item {
  background: rgba(0, 223, 255, 0.1);
  border: 2px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  padding: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  overflow: hidden;
  box-shadow: 0 0 15px rgba(0, 223, 255, 0.1);
  transition: all 0.3s ease;
}

.record-item:hover {
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.2);
  border-color: rgba(0, 223, 255, 0.5);
}

.cyber-line {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent, 
    rgba(0, 223, 255, 0.8), 
    transparent
  );
  animation: cyber-line 2s linear infinite;
}

@keyframes cyber-line {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.record-time {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.2);
}

.record-info {
  flex: 1;
  margin: 0 15px;
}

.record-food-name {
  color: #00DFFF;
  font-size: 16px;
  font-weight: bold;
  display: block;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
  letter-spacing: 0.5px;
}

.record-meal-type {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  text-transform: uppercase;
}

.record-calories {
  color: #fff;
  font-size: 14px;
  margin-right: 15px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.3);
}

.record-actions {
  display: flex;
  gap: 8px;
}

.edit-button, .delete-button {
  background: none;
  border: none;
  padding: 5px 10px;
  font-size: 14px;
  border-radius: 4px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  transition: all 0.3s ease;
}

.edit-button {
  color: #00DFFF;
  background: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
}

.edit-button:active {
  background: rgba(0, 223, 255, 0.2);
}

.delete-button {
  color: #ff4d4f;
  background: rgba(255, 77, 79, 0.1);
  border: 1px solid rgba(255, 77, 79, 0.3);
}

.delete-button:active {
  background: rgba(255, 77, 79, 0.2);
}

.delete-icon {
  font-size: 14px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 25px;
}

.pagination-button {
  padding: 5px 10px;
  background: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #00DFFF;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.3);
}

.pagination-button:active {
  background: rgba(0, 223, 255, 0.2);
}

.pagination-button[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-info {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.3);
}
</style> 