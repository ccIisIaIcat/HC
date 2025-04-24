<template>
  <view class="food-analysis-container">
    <!-- 添加自定义提示框 -->
    <view class="custom-modal" v-if="showModal" @click="closeModal">
      <view class="modal-content" @click.stop>
        <view class="modal-title"> ฅ^•ﻌ•^ฅ</view>
        <view class="modal-text">确定要放弃当前的食物记录吗？</view>
        <view class="modal-buttons">
          <view class="modal-button cancel" @click="closeModal">再想想</view>
          <view class="modal-button confirm" @click="confirmDiscard">确定</view>
        </view>
      </view>
    </view>
    <!-- 自定义导航栏 -->
    <CustomNavBar title="饭团猫" :showBack="false" />
    
    <!-- 页面内容区域，添加顶部内边距给导航栏留出空间 -->
    <view class="content-area">
      <view class="upload-section">
        <view class="upload-area" @tap="chooseImage">
          <block v-if="preview">
            <image :src="preview" mode="aspectFit" class="preview-image"/>
            <view class="preview-overlay">
              <text>点击更换图片</text>
            </view>
          </block>
          <block v-else>
            <view class="upload-placeholder">
              <text class="upload-icon">📷</text>
              <text class="upload-text">点击上传食物图片</text>
              <text class="upload-tip">支持JPG和PNG格式</text>
            </view>
          </block>
        </view>

        <view class="description-input" v-if="preview">
          <textarea
            v-model="imageDescription"
            placeholder="请描述图片中的食物，例如：这是一盘炒饭，里面有鸡蛋、胡萝卜和青豆..."
            :disabled="loading"
            class="description-textarea"
          />
          <text class="description-tip">提供准确的食物描述可以帮助AI更精确地分析食物营养成分</text>
        </view>

        <view class="action-buttons">
          <button 
            class="analyze-button" 
            :disabled="!preview || loading"
            @tap="handleAnalyze"
          >
            <text>{{ loading ? '分析中...' : '开始分析' }}</text>
          </button>
          <button 
            class="reset-button" 
            v-if="preview" 
            @tap="handleReset"
          >
            <text>清除</text>
          </button>
        </view>
      </view>

      <!-- 分析结果部分 -->
      <view class="result-section" v-if="result">
        <view class="result-card">
          <view class="result-header">
            <view class="food-name-input">
              <input 
                type="text" 
                v-model="result.foodType"
                class="name-input"
                placeholder="输入食物名称"
              />
            </view>
            <view class="food-weight">
              <input 
                type="number" 
                v-model="result.weight"
                class="weight-input"
                @input="(e: any) => handleWeightChange(Number(e.detail.value))"
              />
              <text class="weight-unit">克</text>
            </view>
          </view>

          <view class="nutrition-info">
            <!-- 基本营养素 -->
            <view class="nutrition-section">
              <text class="section-title">基本营养素</text>
              <view class="nutrition-grid">
                <view class="nutrition-item" v-for="(value, key) in basicNutrients" :key="key">
                  <text class="item-label">{{ nutritionLabels[key] }}</text>
                  <view class="item-value">
                    <input 
                      type="number" 
                      v-model="result.nutrition[key]"
                      class="value-input"
                    />
                    <text class="unit">{{ nutritionUnits[key] }}</text>
                  </view>
                </view>
              </view>
            </view>

            <!-- 展开/收起按钮 -->
            <button class="expand-button-custom" @click="toggleExpand">
              <text>{{ isExpanded ? '收起更多' : '展开更多' }}</text>
              <text class="expand-icon" :class="{ 'expanded': isExpanded }">▼</text>
            </button>
            
            <!-- 更多营养素（展开时显示） -->
            <view class="expanded-nutrition" v-if="isExpanded">
              <!-- 维生素 -->
              <view class="nutrition-section">
                <text class="section-title">维生素</text>
                <view class="nutrition-grid">
                  <view class="nutrition-item" v-for="(label, key) in vitaminNutrients" :key="key">
                    <text class="item-label">{{ label }}</text>
                    <view class="item-value">
                      <input 
                        type="number" 
                        v-model="result.nutrition.vitamins[key]"
                        class="value-input"
                      />
                      <text class="unit">{{ nutritionUnits[key] }}</text>
                    </view>
                  </view>
                </view>
              </view>

              <!-- 矿物质 -->
              <view class="nutrition-section">
                <text class="section-title">矿物质</text>
                <view class="nutrition-grid">
                  <view class="nutrition-item" v-for="(label, key) in mineralNutrients" :key="key">
                    <text class="item-label">{{ label }}</text>
                    <view class="item-value">
                      <input 
                        type="number" 
                        v-model="result.nutrition.minerals[key]"
                        class="value-input"
                      />
                      <text class="unit">{{ nutritionUnits[key] }}</text>
                    </view>
                  </view>
                </view>
              </view>
            </view>
          </view>

          <!-- 记录表单 -->
          <view class="record-form">
            <text class="form-title">保存记录</text>
            
            <view class="form-item meal-type-container">
              <text class="form-label">餐食类型</text>
              <picker 
                :value="mealTypeIndex" 
                :range="mealTypes"
                @change="handleMealTypeChange"
                class="meal-type-picker"
              >
                <view class="picker-value">{{ mealTypes[mealTypeIndex] }}</view>
              </picker>
            </view>

            <view class="form-item">
              <text class="form-label">备注</text>
              <textarea
                v-model="notes"
                placeholder="添加关于这顿饭的备注..."
                class="notes-textarea"
              />
            </view>

            <view class="form-buttons">
              <button 
                class="save-button" 
                @tap="handleSaveRecord"
                :disabled="loading"
              >
                <text>{{ loading ? '保存中...' : '保存记录' }}</text>
              </button>
              <button 
                class="discard-button" 
                @tap="handleDiscardRecord"
              >
                <text>放弃</text>
              </button>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import api from '@/api';
import CustomNavBar from '@/components/CustomNavBar.vue';

interface NutritionData {
  calories: number;
  protein: number;
  totalFat: number;
  saturatedFat: number;
  transFat: number;
  unsaturatedFat: number;
  carbohydrates: number;
  sugar: number;
  fiber: number;
  vitamins: {
    vitaminA: number;
    vitaminC: number;
    vitaminD: number;
    vitaminB1: number;
    vitaminB2: number;
  };
  minerals: {
    calcium: number;
    iron: number;
    sodium: number;
    potassium: number;
  };
}

interface FoodAnalysisResult {
  foodType: string;
  weight: number;
  nutrition: NutritionData;
}

interface FoodRecord {
  // 基本信息
  id: number; // 必需字段
  food_name: string;
  weight: number;
  record_time: string;
  meal_type: string;
  notes: string;

  // 基本营养素
  calories: number;
  protein: number;
  total_fat: number;
  saturated_fat: number;
  trans_fat: number;
  unsaturated_fat: number;
  carbohydrates: number;
  sugar: number;
  fiber: number;

  // 维生素
  vitamin_a: number;
  vitamin_c: number;
  vitamin_d: number;
  vitamin_b1: number;
  vitamin_b2: number;

  // 矿物质
  calcium: number;
  iron: number;
  sodium: number;
  potassium: number;
}

// 创建食物记录时的数据结构（不需要ID）
interface CreateFoodRecordData {
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

const preview = ref('');
const loading = ref(false);
const imageDescription = ref('');
const result = ref<FoodAnalysisResult | null>(null);
const originalResult = ref<FoodAnalysisResult | null>(null); // 保存原始分析结果
const notes = ref('');
const mealTypeIndex = ref(1); // 默认午餐
const mealTypes = ['早餐', '午餐', '晚餐', '加餐'];
const isExpanded = ref(false);

// 添加提示框相关的响应式变量
const showModal = ref(false);

const basicNutrients = {
  calories: '热量',
  protein: '蛋白质',
  totalFat: '总脂肪',
  saturatedFat: '饱和脂肪',
  transFat: '反式脂肪',
  unsaturatedFat: '不饱和脂肪',
  carbohydrates: '碳水化合物',
  sugar: '糖',
  fiber: '膳食纤维'
} as const;

const vitaminNutrients = {
  vitaminA: '维生素A',
  vitaminC: '维生素C',
  vitaminD: '维生素D',
  vitaminB1: '维生素B1',
  vitaminB2: '维生素B2'
} as const;

const mineralNutrients = {
  calcium: '钙',
  iron: '铁',
  sodium: '钠',
  potassium: '钾'
} as const;

const nutritionLabels = {
  calories: '热量',
  protein: '蛋白质',
  totalFat: '总脂肪',
  saturatedFat: '饱和脂肪',
  transFat: '反式脂肪',
  unsaturatedFat: '不饱和脂肪',
  carbohydrates: '碳水化合物',
  sugar: '糖',
  fiber: '膳食纤维',
  vitaminA: '维生素A',
  vitaminC: '维生素C',
  vitaminD: '维生素D',
  vitaminB1: '维生素B1',
  vitaminB2: '维生素B2',
  calcium: '钙',
  iron: '铁',
  sodium: '钠',
  potassium: '钾'
};

const nutritionUnits = {
  calories: '卡路里',
  protein: '克',
  totalFat: '克',
  saturatedFat: '克',
  transFat: '克',
  unsaturatedFat: '克',
  carbohydrates: '克',
  sugar: '克',
  fiber: '克',
  vitaminA: '微克',
  vitaminC: '毫克',
  vitaminD: '微克',
  vitaminB1: '毫克',
  vitaminB2: '毫克',
  calcium: '毫克',
  iron: '毫克',
  sodium: '毫克',
  potassium: '毫克'
};

// 选择图片
const chooseImage = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['original', 'compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      preview.value = res.tempFilePaths[0];
    },
    fail: (err) => {
      console.error('选择图片失败:', err);
      uni.showToast({
        title: '选择图片失败',
        icon: 'none'
      });
    }
  });
};

// 切换展开/收起状态
const toggleExpand = () => {
  try {
    console.log('切换展开状态，当前状态:', isExpanded.value);
    
    // 显示点击反馈
    uni.showToast({
      title: '正在切换展开状态...',
      icon: 'none',
      duration: 1000
    });
    
    // 检查结果对象
    console.log('当前结果对象:', JSON.stringify(result.value));
    
    // 检查结果对象中是否包含必要的嵌套属性
    if (result.value) {
      // 检查维生素属性是否存在
      if (!result.value.nutrition.vitamins) {
        console.warn('维生素数据不存在，正在创建...');
        // 使用Vue的响应式API确保更新被检测到
        const newNutrition = { ...result.value.nutrition };
        newNutrition.vitamins = {
          vitaminA: 0, vitaminC: 0, vitaminD: 0, vitaminB1: 0, vitaminB2: 0
        };
        result.value.nutrition = newNutrition;
      }
      
      // 检查矿物质属性是否存在
      if (!result.value.nutrition.minerals) {
        console.warn('矿物质数据不存在，正在创建...');
        // 使用Vue的响应式API确保更新被检测到
        const newNutrition = { ...result.value.nutrition };
        newNutrition.minerals = {
          calcium: 0, iron: 0, sodium: 0, potassium: 0
        };
        result.value.nutrition = newNutrition;
      }
    }
    
    // 切换状态
    isExpanded.value = !isExpanded.value;
    
    console.log('切换后状态:', isExpanded.value);
    
  } catch (error) {
    console.error('切换展开状态时出错:', error);
    uni.showToast({
      title: '操作失败，请查看控制台',
      icon: 'none'
    });
  }
};

// 添加重量变化处理函数
const handleWeightChange = (newWeight: number) => {
  if (!result.value || !originalResult.value || newWeight <= 0) {
    console.log('无法更新重量：缺少必要数据');
    return;
  }
  
  console.log('更新前的重量:', result.value.weight);
  console.log('新重量:', newWeight);
  console.log('原始重量:', originalResult.value.weight);
  
  const ratio = newWeight / originalResult.value.weight;
  console.log('计算的比率:', ratio);
  
  // 更新当前结果的重量
  result.value.weight = newWeight;
  
  // 更新基本营养素
  result.value.nutrition.calories = Number((originalResult.value.nutrition.calories * ratio).toFixed(1));
  result.value.nutrition.protein = Number((originalResult.value.nutrition.protein * ratio).toFixed(1));
  result.value.nutrition.totalFat = Number((originalResult.value.nutrition.totalFat * ratio).toFixed(1));
  result.value.nutrition.saturatedFat = Number((originalResult.value.nutrition.saturatedFat * ratio).toFixed(1));
  result.value.nutrition.transFat = Number((originalResult.value.nutrition.transFat * ratio).toFixed(1));
  result.value.nutrition.unsaturatedFat = Number((originalResult.value.nutrition.unsaturatedFat * ratio).toFixed(1));
  result.value.nutrition.carbohydrates = Number((originalResult.value.nutrition.carbohydrates * ratio).toFixed(1));
  result.value.nutrition.sugar = Number((originalResult.value.nutrition.sugar * ratio).toFixed(1));
  result.value.nutrition.fiber = Number((originalResult.value.nutrition.fiber * ratio).toFixed(1));
  
  // 更新维生素
  result.value.nutrition.vitamins.vitaminA = Number((originalResult.value.nutrition.vitamins.vitaminA * ratio).toFixed(1));
  result.value.nutrition.vitamins.vitaminC = Number((originalResult.value.nutrition.vitamins.vitaminC * ratio).toFixed(1));
  result.value.nutrition.vitamins.vitaminD = Number((originalResult.value.nutrition.vitamins.vitaminD * ratio).toFixed(1));
  result.value.nutrition.vitamins.vitaminB1 = Number((originalResult.value.nutrition.vitamins.vitaminB1 * ratio).toFixed(1));
  result.value.nutrition.vitamins.vitaminB2 = Number((originalResult.value.nutrition.vitamins.vitaminB2 * ratio).toFixed(1));
  
  // 更新矿物质
  result.value.nutrition.minerals.calcium = Number((originalResult.value.nutrition.minerals.calcium * ratio).toFixed(1));
  result.value.nutrition.minerals.iron = Number((originalResult.value.nutrition.minerals.iron * ratio).toFixed(1));
  result.value.nutrition.minerals.sodium = Number((originalResult.value.nutrition.minerals.sodium * ratio).toFixed(1));
  result.value.nutrition.minerals.potassium = Number((originalResult.value.nutrition.minerals.potassium * ratio).toFixed(1));
  
  console.log('更新后的结果:', result.value);
};

// 修改分析图片函数
const handleAnalyze = async () => {
  if (!preview.value) return;
  
  loading.value = true;
  isExpanded.value = false; // 重置展开状态
  
  try {
    console.log('开始分析图片:', preview.value);
    
    // 使用API模块上传分析图片
    const apiResult = await api.food.analyzeFoodImage({
      filePath: preview.value,
      description: imageDescription.value
    });
    
    console.log('API分析结果:', apiResult);
    
    // 创建完整的默认结果对象
    const defaultResult: FoodAnalysisResult = {
      foodType: apiResult?.foodType || '未知食物',
      weight: apiResult?.weight || 100,
      nutrition: {
        calories: apiResult?.nutrition?.calories || 0,
        protein: apiResult?.nutrition?.protein || 0,
        totalFat: apiResult?.nutrition?.totalFat || 0,
        saturatedFat: 0,  // 这些字段API中没有，设置默认值
        transFat: 0,      // 这些字段API中没有，设置默认值
        unsaturatedFat: 0, // 这些字段API中没有，设置默认值
        carbohydrates: apiResult?.nutrition?.carbohydrates || 0,
        sugar: apiResult?.nutrition?.sugar || 0,
        fiber: apiResult?.nutrition?.fiber || 0,
        vitamins: {
          vitaminA: apiResult?.nutrition?.vitamins?.vitaminA || 0,
          vitaminC: apiResult?.nutrition?.vitamins?.vitaminC || 0,
          vitaminD: apiResult?.nutrition?.vitamins?.vitaminD || 0,
          vitaminB1: 0,  // API中没有这个字段，设置默认值
          vitaminB2: 0   // API中没有这个字段，设置默认值
        },
        minerals: {
          calcium: apiResult?.nutrition?.minerals?.calcium || 0,
          iron: apiResult?.nutrition?.minerals?.iron || 0,
          sodium: apiResult?.nutrition?.minerals?.sodium || 0,
          potassium: apiResult?.nutrition?.minerals?.potassium || 0
        }
      }
    };
    
    // 保存原始结果
    originalResult.value = JSON.parse(JSON.stringify(defaultResult));
    // 设置当前结果
    result.value = defaultResult;
    
    console.log('分析结果已保存:', result.value);
    console.log('原始结果已保存:', originalResult.value);
    
  } catch (err) {
    console.error('分析失败:', err);
    uni.showToast({
      title: '分析失败，请重试',
      icon: 'none'
    });
  } finally {
    loading.value = false;
  }
};

// 重置
const handleReset = () => {
  preview.value = '';
  imageDescription.value = '';
  result.value = null;
  notes.value = '';
};

// 处理餐食类型变化
const handleMealTypeChange = (e: { detail: { value: number } }) => {
  mealTypeIndex.value = e.detail.value;
};

// 保存记录
const handleSaveRecord = async () => {
  try {
    // 检查token是否存在
    const token = uni.getStorageSync('token');
    if (!token) {
      uni.showToast({
        title: '请先登录',
        icon: 'none'
      });
      setTimeout(() => {
        uni.navigateTo({
          url: '/pages/login/index'
        });
      }, 1500);
      return;
    }

    // 检查分析结果是否存在
    if (!result.value) {
      uni.showToast({
        title: '请先分析食物',
        icon: 'none'
      });
      return;
    }

    // 构建食物记录数据，确保所有字段都是正确的类型
    const foodRecord = {
      food_name: String(result.value.foodType),
      weight: Number(result.value.weight),
      calories: Number(result.value.nutrition.calories),
      protein: Number(result.value.nutrition.protein),
      total_fat: Number(result.value.nutrition.totalFat),
      saturated_fat: Number(result.value.nutrition.saturatedFat),
      trans_fat: Number(result.value.nutrition.transFat),
      unsaturated_fat: Number(result.value.nutrition.unsaturatedFat),
      carbohydrates: Number(result.value.nutrition.carbohydrates),
      sugar: Number(result.value.nutrition.sugar),
      fiber: Number(result.value.nutrition.fiber),
      vitamin_a: Number(result.value.nutrition.vitamins.vitaminA),
      vitamin_c: Number(result.value.nutrition.vitamins.vitaminC),
      vitamin_d: Number(result.value.nutrition.vitamins.vitaminD),
      vitamin_b1: Number(result.value.nutrition.vitamins.vitaminB1),
      vitamin_b2: Number(result.value.nutrition.vitamins.vitaminB2),
      calcium: Number(result.value.nutrition.minerals.calcium),
      iron: Number(result.value.nutrition.minerals.iron),
      sodium: Number(result.value.nutrition.minerals.sodium),
      potassium: Number(result.value.nutrition.minerals.potassium),
      meal_type: mealTypes[mealTypeIndex.value],
      record_time: new Date().toISOString(),
      notes: String(notes.value || '')
    };

    console.log('准备发送的数据:', foodRecord);

    // 调用API保存记录
    const response = await api.food.createFoodRecord(foodRecord);
    console.log('API响应:', response);
    
    if (response && response.record) {
      uni.showToast({
        title: '保存成功',
        icon: 'success'
      });
      
      // 重置表单
      handleReset();
      
      // 返回上一页
      setTimeout(() => {
        uni.navigateBack();
      }, 1500);
    } else {
      throw new Error('保存失败：响应数据格式错误');
    }
  } catch (error: any) {
    console.error('保存食物记录失败:', error);
    uni.showToast({
      title: error?.message || '保存失败，请重试',
      icon: 'none'
    });
  }
};

// 关闭提示框
const closeModal = () => {
  showModal.value = false;
};

// 确认放弃记录
const confirmDiscard = () => {
  result.value = null;
  notes.value = '';
  closeModal();
};

// 修改放弃记录函数
const handleDiscardRecord = () => {
  showModal.value = true;
};
</script>

<style>
.food-analysis-container {
  min-height: 100vh;
  background-color: #000000;
  background-image: linear-gradient(rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.7)),
    url('@/static/images/food-record-bj.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.food-analysis-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, 
    rgba(255, 0, 98, 0.1), 
    rgba(0, 223, 255, 0.1));
  animation: gradientMove 8s ease infinite;
}

@keyframes gradientMove {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.content-area {
  padding-top: 74px; /* 状态栏高度 + 导航栏高度 (约20px + 44px + 10px边距) */
  padding-bottom: 20px;
}

.upload-section {
  background: rgba(0, 223, 255, 0.1);
  border: 2px solid rgba(0, 223, 255, 0.5);
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.upload-area {
  width: 100%;
  height: 200px;
  border: 2px dashed rgba(0, 223, 255, 0.5);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.preview-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-overlay text {
  color: #fff;
  font-size: 16px;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.upload-icon {
  font-size: 40px;
}

.upload-text {
  color: #00DFFF;
  font-size: 16px;
}

.upload-tip {
  color: rgba(0, 223, 255, 0.6);
  font-size: 12px;
}

.description-input {
  margin-top: 20px;
}

.description-textarea {
  width: 100%;
  height: 80px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 10px;
  color: #fff;
}

.description-tip {
  color: rgba(0, 223, 255, 0.6);
  font-size: 12px;
  margin-top: 5px;
}

.action-buttons {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.analyze-button, .reset-button {
  flex: 1;
  height: 40px;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-size: 16px;
}

.analyze-button {
  background: linear-gradient(45deg, #00DFFF, #00BFFF);
}

.reset-button {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.result-section {
  margin-top: 20px;
}

.result-card {
  background: rgba(0, 223, 255, 0.1);
  border: 2px solid rgba(0, 223, 255, 0.5);
  border-radius: 8px;
  padding: 20px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.food-name-input {
  flex: 1;
  margin-right: 20px;
}

.name-input {
  width: 100%;
  height: 30px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #00DFFF;
  font-size: 20px;
  font-weight: bold;
  padding: 0 10px;
}

.food-weight {
  display: flex;
  align-items: center;
  gap: 5px;
}

.weight-input {
  width: 60px;
  height: 30px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #fff;
  text-align: center;
}

.weight-unit {
  color: #fff;
  font-size: 14px;
}

.nutrition-section {
  margin-bottom: 20px;
}

.section-title {
  color: #00DFFF;
  font-size: 18px;
  margin-bottom: 10px;
}

.nutrition-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.nutrition-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.item-label {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.item-value {
  display: flex;
  align-items: center;
  gap: 5px;
}

.value-input {
  width: 60px;
  height: 30px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #fff;
  text-align: center;
}

.unit {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}

.record-form {
  margin-top: 20px;
  padding: 15px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  border: 1px solid rgba(0, 223, 255, 0.3);
}

.form-title {
  color: var(--primary-color);
  font-size: 16px;
  margin-bottom: 15px;
  font-weight: bold;
}

.form-item {
  margin-bottom: 15px;
}

.meal-type-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.meal-type-picker {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 8px 12px;
  margin-top: 4px;
}

.picker-value {
  color: #fff;
  font-size: 14px;
  line-height: 1.4;
}

.notes-textarea {
  width: 100%;
  height: 80px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 8px 12px;
  color: #fff;
  font-size: 14px;
  margin-top: 4px;
}

.form-buttons {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.save-button, .discard-button {
  flex: 1;
  height: 40px;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-size: 16px;
}

.save-button {
  background: linear-gradient(45deg, #00DFFF, #00BFFF);
}

.discard-button {
  background: rgba(255, 77, 79, 0.8);
}

button[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
}

.expand-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 10px;
  background: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  margin: 10px 0;
  cursor: pointer;
}

.expand-button text {
  color: #00DFFF;
  font-size: 14px;
}

.expand-button-custom {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 10px;
  width: 100%;
  background: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  margin: 10px 0;
  cursor: pointer;
}

.expand-button-custom text {
  color: #00DFFF;
  font-size: 14px;
}

.expand-icon {
  font-size: 12px;
  transition: transform 0.3s ease;
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

.expanded-nutrition {
  animation: slideDown 0.3s ease;
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

/* 加载动画 */
.analyze-button {
  position: relative;
  overflow: hidden;
}

.analyze-button::after {
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
  animation: loading 1.5s infinite;
}

@keyframes loading {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

/* 添加自定义提示框样式 */
.custom-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal-content {
  background: rgba(0, 223, 255, 0.1);
  border: 2px solid rgba(0, 223, 255, 0.5);
  border-radius: 15px;
  padding: 20px;
  width: 80%;
  max-width: 300px;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  animation: modalFadeIn 0.3s ease;
}

.modal-title {
  color: #00DFFF;
  font-size: 20px;
  text-align: center;
  margin-bottom: 15px;
}

.modal-text {
  color: #fff;
  font-size: 16px;
  text-align: center;
  margin-bottom: 20px;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.modal-button {
  flex: 1;
  padding: 10px;
  text-align: center;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.modal-button.cancel {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.modal-button.confirm {
  background: linear-gradient(45deg, #00DFFF, #00BFFF);
  color: #fff;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style> 