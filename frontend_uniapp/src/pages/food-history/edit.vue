<template>
  <view class="food-edit-container">
    <!-- 加载状态 -->
    <view class="loading-container" v-if="loading">
      <text class="loading-text">正在加载数据...</text>
    </view>

    <!-- 错误信息 -->
    <view class="error-container" v-if="!loading && error">
      <text class="error-text">{{ error }}</text>
      <button class="back-button" @tap="handleBack">返回</button>
    </view>

    <!-- 编辑表单 -->
    <view class="edit-form" v-if="!loading && !error && editableRecord">
      <!-- 基本信息部分 -->
      <view class="form-card">
        <view class="food-name-weight">
          <view class="food-name">
            <input 
              type="text" 
              v-model="editableRecord.food_name" 
              class="name-input"
              placeholder="输入食物名称"
            />
          </view>
          <view class="food-weight">
            <input 
              type="number" 
              v-model="editableRecord.weight" 
              @input="handleWeightChange"
              class="weight-input"
            />
            <text class="weight-unit">克</text>
          </view>
        </view>
      </view>
      
      <!-- 营养信息部分 -->
      <view class="form-section">
        <view class="section-title">基本营养素</view>
        
        <view class="nutrition-grid">
          <view class="nutrition-item">
            <text class="item-label">热量</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="editableRecord.calories" 
                class="nutrition-input"
              />
              <text class="unit-text">卡路里</text>
            </view>
          </view>
          
          <view class="nutrition-item">
            <text class="item-label">蛋白质</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="editableRecord.protein" 
                class="nutrition-input"
              />
              <text class="unit-text">克</text>
            </view>
          </view>
          
          <view class="nutrition-item">
            <text class="item-label">脂肪</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="editableRecord.total_fat" 
                class="nutrition-input"
              />
              <text class="unit-text">克</text>
            </view>
          </view>
          
          <view class="nutrition-item">
            <text class="item-label">碳水化合物</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="editableRecord.carbohydrates" 
                class="nutrition-input"
              />
              <text class="unit-text">克</text>
            </view>
          </view>
          
          <view class="nutrition-item">
            <text class="item-label">糖分</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="editableRecord.sugar" 
                class="nutrition-input"
              />
              <text class="unit-text">克</text>
            </view>
          </view>
          
          <view class="nutrition-item">
            <text class="item-label">纤维素</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="editableRecord.fiber" 
                class="nutrition-input"
              />
              <text class="unit-text">克</text>
            </view>
          </view>
        </view>
        
        <!-- 更多营养物质下拉框 -->
        <view class="more-nutrients">
          <view class="more-nutrients-header" @tap="toggleMoreNutrients">
            <text class="more-nutrients-title">更多营养物质</text>
            <text class="more-nutrients-icon">{{ showMoreNutrients ? '▲' : '▼' }}</text>
          </view>
          
          <view class="more-nutrients-content" v-if="showMoreNutrients">
            <view class="nutrition-grid">
              <!-- 脂肪详情 -->
              <view class="nutrition-item">
                <text class="item-label">饱和脂肪</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.saturated_fat" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">反式脂肪</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.trans_fat" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">不饱和脂肪</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.unsaturated_fat" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">克</text>
                </view>
              </view>
              
              <!-- 维生素 -->
              <view class="nutrition-item">
                <text class="item-label">维生素A</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.vitamin_a" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">微克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">维生素C</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.vitamin_c" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">毫克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">维生素D</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.vitamin_d" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">微克</text>
                </view>
              </view>
              
              <!-- 矿物质 -->
              <view class="nutrition-item">
                <text class="item-label">钙</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.calcium" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">毫克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">铁</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.iron" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">毫克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">钠</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.sodium" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">毫克</text>
                </view>
              </view>
              
              <view class="nutrition-item">
                <text class="item-label">钾</text>
                <view class="item-input-wrapper">
                  <input 
                    type="number" 
                    v-model="editableRecord.potassium" 
                    class="nutrition-input"
                  />
                  <text class="unit-text">毫克</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 其他信息部分 -->
      <view class="form-section">
        <view class="section-title">其他信息</view>
        
        <view class="record-detail">
          <text class="detail-label">餐食类型</text>
          <picker 
            class="meal-picker" 
            @change="handleMealTypeChange" 
            :value="mealTypeIndex" 
            :range="mealTypes"
          >
            <view class="picker-value">{{ mealTypes[mealTypeIndex] }}</view>
          </picker>
        </view>
        
        <view class="record-detail">
          <text class="detail-label">备注</text>
          <textarea 
            class="notes-textarea" 
            v-model="notes" 
            placeholder="添加关于这顿饭的备注..."
          />
        </view>
      </view>
      
      <!-- 保存按钮 -->
      <button class="save-button" @tap="handleSave">保存更改</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import api from '@/api';

// 导入uni-app的生命周期钩子
import { onLoad } from '@dcloudio/uni-app';

// 页面配置
defineOptions({
  navigationStyle: 'custom' // 自定义导航栏，不使用默认导航栏
});

interface PageOptions {
  id?: string;
}

interface FoodRecord {
  id?: number;
  [key: string]: any; // 添加索引签名允许任意属性包括ID
  food_name: string;
  weight: number;
  record_time: string;
  meal_type: string;
  notes: string;
  calories: number;
  protein: number;
  total_fat: number;
  saturated_fat: number;
  trans_fat: number;
  unsaturated_fat: number;
  carbohydrates: number;
  sugar: number;
  fiber: number;
  vitamin_a: number;
  vitamin_c: number;
  vitamin_d: number;
  vitamin_b1: number;
  vitamin_b2: number;
  calcium: number;
  iron: number;
  sodium: number;
  potassium: number;
}

// 状态变量
const loading = ref(true);
const error = ref('');
const recordId = ref<number | null>(null);
const record = ref<FoodRecord | null>(null);
const editableRecord = ref<FoodRecord | null>(null);
const originalWeight = ref<number>(0);
const notes = ref('');
const mealTypeIndex = ref(1); // 默认选择"午餐"
const mealTypes = ['早餐', '午餐', '晚餐', '加餐'];
const showMoreNutrients = ref(false); // 控制更多营养物质的显示

// 切换更多营养物质的显示状态
const toggleMoreNutrients = () => {
  showMoreNutrients.value = !showMoreNutrients.value;
};

// 直接使用onLoad钩子，无需defineExpose
onLoad((options: any) => {
  console.log('编辑页面加载，参数:', options);
  if (options && options.id) {
    recordId.value = Number(options.id);
    console.log('设置要编辑的记录ID:', recordId.value);
    fetchFoodRecord();
  } else {
    error.value = '无效的记录ID';
    loading.value = false;
  }
});

// 加载食物记录数据
const fetchFoodRecord = async () => {
  console.log('开始执行 fetchFoodRecord 函数');
  // 确保使用 recordId.value
  const id = recordId.value;
  
  if (!id || isNaN(id)) {
    console.error('无效的记录ID:', id);
    error.value = '无效的记录ID';
    loading.value = false;
    return;
  }
  
  try {
    loading.value = true;
    error.value = '';
    
    console.log('正在获取食物记录，ID:', id);
    
    // 添加加载超时保护
    const timeout = setTimeout(() => {
      if (loading.value) {
        console.error('获取食物记录超时');
        error.value = '获取数据超时，请重试';
        loading.value = false;
      }
    }, 15000); // 15秒超时
    
    // 使用API封装调用获取记录
    console.log('调用API: getFoodRecordDetail，ID:', id);
    const response = await api.food.getFoodRecordDetail(id);
    
    // 清除超时计时器
    clearTimeout(timeout);
    
    console.log('API响应数据:', JSON.stringify(response).substring(0, 500));
    
    // 检查响应是否为空
    if (!response) {
      console.error('API返回空响应');
      error.value = 'API返回空响应';
      loading.value = false;
      return;
    }
    
    // 无论响应是否包含record字段，尝试获取有效数据
    const recordData = response.record || response;
    
    if (recordData) {
      console.log('提取的记录数据:', JSON.stringify(recordData).substring(0, 500));
      
      // 确保有ID字段
      if (!recordData.id && (recordData as any).ID) {
        recordData.id = (recordData as any).ID;
        console.log('将大写ID赋值给小写id:', recordData.id);
      }
      
      // 映射数据字段
      const mappedData: FoodRecord = {
        id: recordData.id || (recordData as any).ID,
        food_name: recordData.food_name || '',
        weight: recordData.weight || 100,
        record_time: recordData.record_time || new Date().toISOString(),
        meal_type: recordData.meal_type || '午餐',
        notes: recordData.notes || '',
        calories: recordData.calories || 0,
        protein: recordData.protein || 0,
        total_fat: recordData.total_fat || 0,
        saturated_fat: recordData.saturated_fat || 0,
        trans_fat: recordData.trans_fat || 0,
        unsaturated_fat: recordData.unsaturated_fat || 0,
        carbohydrates: recordData.carbohydrates || 0,
        sugar: recordData.sugar || 0,
        fiber: recordData.fiber || 0,
        vitamin_a: recordData.vitamin_a || 0,
        vitamin_c: recordData.vitamin_c || 0,
        vitamin_d: recordData.vitamin_d || 0,
        vitamin_b1: recordData.vitamin_b1 || 0,
        vitamin_b2: recordData.vitamin_b2 || 0,
        calcium: recordData.calcium || 0,
        iron: recordData.iron || 0,
        sodium: recordData.sodium || 0,
        potassium: recordData.potassium || 0
      };
      
      console.log('映射后的数据:', JSON.stringify(mappedData));
      record.value = mappedData;
      editableRecord.value = JSON.parse(JSON.stringify(mappedData));
      originalWeight.value = mappedData.weight;
      notes.value = mappedData.notes || '';
      const mealTypeIdx = mealTypes.indexOf(mappedData.meal_type);
      mealTypeIndex.value = mealTypeIdx >= 0 ? mealTypeIdx : 1;
    } else {
      console.error('API返回数据中没有有效记录:', response);
      error.value = '未找到记录数据';
    }
  } catch (err: any) {
    console.error('获取食物记录失败:', err);
    error.value = '获取食物记录失败: ' + (err.message || '未知错误');
  } finally {
    loading.value = false;
  }
};

// 处理返回按钮
const handleBack = () => {
  uni.navigateBack();
};

// 处理重量变化
const handleWeightChange = (e: any) => {
  if (!editableRecord.value || !originalWeight.value) return;
  
  const newWeight = parseFloat(e.detail.value) || 0;
  if (newWeight <= 0) return;
  
  // 计算比例因子
  const ratio = newWeight / editableRecord.value.weight;
  
  // 按比例调整基本营养素
  const nutrientFields = [
    'calories', 'protein', 'total_fat', 'saturated_fat', 
    'trans_fat', 'unsaturated_fat', 'carbohydrates', 
    'sugar', 'fiber', 'vitamin_a', 'vitamin_c', 
    'vitamin_d', 'vitamin_b1', 'vitamin_b2', 'calcium', 
    'iron', 'sodium', 'potassium'
  ];
  
  nutrientFields.forEach(field => {
    if (editableRecord.value && typeof editableRecord.value[field as keyof FoodRecord] === 'number') {
      const oldValue = editableRecord.value[field as keyof FoodRecord] as number;
      const newValue = +(oldValue * ratio).toFixed(2);
      (editableRecord.value as any)[field] = newValue;
    }
  });
};

// 处理餐食类型变化
const handleMealTypeChange = (e: any) => {
  mealTypeIndex.value = e.detail.value;
};

// 保存更新的记录
const handleSave = async () => {
  if (!editableRecord.value) {
    error.value = '没有记录数据可保存';
    return;
  }
  
  // 确保使用正确的 ID
  const id = recordId.value;
  if (!id || isNaN(id)) {
    console.error('保存时 ID 无效:', id);
    error.value = '记录 ID 无效，无法保存';
    return;
  }
  
  try {
    loading.value = true;
    
    // 准备更新的数据
    const updateData = {
      ...editableRecord.value,
      id: id, // 使用本地变量
      meal_type: mealTypes[mealTypeIndex.value],
      notes: notes.value,
      // 确保包含record_time字段
      record_time: editableRecord.value.record_time || new Date().toISOString()
    };
    
    console.log('准备更新数据:', JSON.stringify(updateData).substring(0, 500), '使用ID:', id);
    
    // 使用API封装调用更新记录
    console.log('调用API: updateFoodRecord，ID:', id);
    await api.food.updateFoodRecord(id, updateData);
    
    uni.showToast({
      title: '食物记录更新成功',
      icon: 'success'
    });
    
    // 返回到记录列表
    setTimeout(() => {
      uni.navigateBack();
    }, 1500);
    
  } catch (err: any) {
    console.error('更新食物记录失败:', err);
    error.value = '更新食物记录失败: ' + (err.message || '未知错误');
    uni.showToast({
      title: '更新失败，请重试',
      icon: 'none'
    });
  } finally {
    loading.value = false;
  }
};
</script>

<style>
.food-edit-container {
  min-height: 100vh;
  background-color: #121212;
  padding: 0;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 60vh;
}

.loading-text, .error-text {
  color: #00DFFF;
  font-size: 16px;
  margin-bottom: 20px;
}

.error-text {
  color: #ff4d4f;
}

.back-button {
  background: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  color: #00DFFF;
  border-radius: 4px;
  padding: 8px 16px;
  font-size: 14px;
}

.edit-form {
  margin-top: 0;
  padding: 10px 20px 20px;
}

.form-card {
  background-color: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 15px;
}

.food-name-weight {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.food-name {
  flex: 1;
  margin-right: 15px;
}

.name-input {
  width: 100%;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #FFFFFF;
  font-size: 16px;
  padding: 0 10px;
}

.food-weight {
  display: flex;
  align-items: center;
}

.weight-input {
  width: 80px;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #FFFFFF;
  text-align: center;
  font-size: 16px;
}

.weight-unit {
  color: #FFFFFF;
  margin-left: 5px;
  font-size: 16px;
}

.form-section {
  background-color: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 15px;
}

.section-title {
  color: #00DFFF;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 12px;
  border-bottom: 1px solid rgba(0, 223, 255, 0.3);
  padding-bottom: 6px;
}

.nutrition-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.nutrition-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  padding: 12px;
}

.item-label {
  color: #FFFFFF;
  font-size: 16px;
}

.item-input-wrapper {
  display: flex;
  align-items: center;
}

.nutrition-input {
  width: 70px;
  height: 36px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #FFFFFF;
  text-align: center;
  font-size: 16px;
}

.unit-text {
  color: rgba(255, 255, 255, 0.7);
  margin-left: 5px;
  font-size: 14px;
  width: 40px;
}

.record-detail {
  margin-bottom: 16px;
}

.detail-label {
  color: #00DFFF;
  font-size: 16px;
  margin-bottom: 8px;
  display: block;
}

.meal-picker {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 10px;
  width: 100%;
}

.picker-value {
  color: #FFFFFF;
  font-size: 16px;
}

.notes-textarea {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 10px;
  width: 100%;
  height: 80px;
  color: #FFFFFF;
  font-size: 16px;
}

.save-button {
  background: linear-gradient(45deg, #00DFFF, #00BFFF);
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  padding: 15px 0;
  font-size: 18px;
  font-weight: bold;
  margin-top: 20px;
  margin-bottom: 40px;
  width: 100%;
}

.more-nutrients {
  margin-top: 16px;
}

.more-nutrients-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(0, 223, 255, 0.15);
  border: 1px solid rgba(0, 223, 255, 0.4);
  border-radius: 6px;
  padding: 12px 16px;
  cursor: pointer;
}

.more-nutrients-title {
  color: #00DFFF;
  font-size: 16px;
  font-weight: bold;
}

.more-nutrients-icon {
  color: #00DFFF;
  font-size: 18px;
  transition: transform 0.3s ease;
}

.more-nutrients-content {
  margin-top: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style> 
