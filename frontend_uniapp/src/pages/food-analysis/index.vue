<template>
  <view class="food-analysis-container">
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
          
          <view class="form-item">
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
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import api from '@/api';

interface NutritionData {
  calories: number;
  protein: number;
  totalFat: number;
  saturatedFat: number;
  transFat: number;
  unsaturatedFat: number;
  cholesterol: number;
  carbohydrates: number;
  sugar: number;
  fiber: number;
  water: number;
  vitamins: {
    vitaminA: number;
    vitaminC: number;
    vitaminD: number;
    vitaminE: number;
    vitaminK: number;
    vitaminB1: number;
    vitaminB2: number;
    vitaminB3: number;
    vitaminB5: number;
    vitaminB6: number;
    vitaminB7: number;
    vitaminB9: number;
    vitaminB12: number;
  };
  minerals: {
    calcium: number;
    iron: number;
    magnesium: number;
    phosphorus: number;
    potassium: number;
    sodium: number;
    zinc: number;
    copper: number;
    manganese: number;
    selenium: number;
    chromium: number;
    iodine: number;
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
const notes = ref('');
const mealTypeIndex = ref(1); // 默认午餐
const mealTypes = ['早餐', '午餐', '晚餐', '加餐'];
const isExpanded = ref(false);

const basicNutrients = {
  calories: '热量',
  protein: '蛋白质',
  totalFat: '总脂肪',
  saturatedFat: '饱和脂肪',
  transFat: '反式脂肪',
  cholesterol: '胆固醇',
  carbohydrates: '碳水化合物',
  sugar: '糖',
  fiber: '膳食纤维',
  water: '水分'
} as const;

const vitaminNutrients = {
  vitaminA: '维生素A',
  vitaminC: '维生素C',
  vitaminD: '维生素D',
  vitaminE: '维生素E',
  vitaminK: '维生素K',
  vitaminB1: '维生素B1',
  vitaminB2: '维生素B2',
  vitaminB3: '维生素B3',
  vitaminB5: '维生素B5',
  vitaminB6: '维生素B6',
  vitaminB7: '维生素B7',
  vitaminB9: '维生素B9',
  vitaminB12: '维生素B12'
} as const;

const mineralNutrients = {
  calcium: '钙',
  iron: '铁',
  magnesium: '镁',
  phosphorus: '磷',
  potassium: '钾',
  sodium: '钠',
  zinc: '锌',
  copper: '铜',
  manganese: '锰',
  selenium: '硒',
  chromium: '铬',
  iodine: '碘'
} as const;

const nutritionLabels = {
  calories: '热量',
  protein: '蛋白质',
  totalFat: '总脂肪',
  saturatedFat: '饱和脂肪',
  transFat: '反式脂肪',
  cholesterol: '胆固醇',
  carbohydrates: '碳水化合物',
  sugar: '糖',
  fiber: '膳食纤维',
  water: '水分',
  vitaminA: '维生素A',
  vitaminC: '维生素C',
  vitaminD: '维生素D',
  vitaminE: '维生素E',
  vitaminK: '维生素K',
  vitaminB1: '维生素B1',
  vitaminB2: '维生素B2',
  vitaminB3: '维生素B3',
  vitaminB5: '维生素B5',
  vitaminB6: '维生素B6',
  vitaminB7: '维生素B7',
  vitaminB9: '维生素B9',
  vitaminB12: '维生素B12',
  calcium: '钙',
  iron: '铁',
  magnesium: '镁',
  phosphorus: '磷',
  potassium: '钾',
  sodium: '钠',
  zinc: '锌',
  copper: '铜',
  manganese: '锰',
  selenium: '硒',
  chromium: '铬',
  iodine: '碘'
};

const nutritionUnits = {
  calories: '卡路里',
  protein: '克',
  totalFat: '克',
  saturatedFat: '克',
  transFat: '克',
  cholesterol: '毫克',
  carbohydrates: '克',
  sugar: '克',
  fiber: '克',
  water: '克',
  vitaminA: '微克',
  vitaminC: '毫克',
  vitaminD: '微克',
  vitaminE: '毫克',
  vitaminK: '微克',
  vitaminB1: '毫克',
  vitaminB2: '毫克',
  vitaminB3: '毫克',
  vitaminB5: '毫克',
  vitaminB6: '毫克',
  vitaminB7: '微克',
  vitaminB9: '微克',
  vitaminB12: '微克',
  calcium: '毫克',
  iron: '毫克',
  magnesium: '毫克',
  phosphorus: '毫克',
  potassium: '毫克',
  sodium: '毫克',
  zinc: '毫克',
  copper: '毫克',
  manganese: '毫克',
  selenium: '微克',
  chromium: '微克',
  iodine: '微克'
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
          vitaminA: 0, vitaminC: 0, vitaminD: 0, vitaminE: 0, vitaminK: 0,
          vitaminB1: 0, vitaminB2: 0, vitaminB3: 0, vitaminB5: 0, vitaminB6: 0,
          vitaminB7: 0, vitaminB9: 0, vitaminB12: 0
        };
        result.value.nutrition = newNutrition;
      }
      
      // 检查矿物质属性是否存在
      if (!result.value.nutrition.minerals) {
        console.warn('矿物质数据不存在，正在创建...');
        // 使用Vue的响应式API确保更新被检测到
        const newNutrition = { ...result.value.nutrition };
        newNutrition.minerals = {
          calcium: 0, iron: 0, magnesium: 0, phosphorus: 0, potassium: 0,
          sodium: 0, zinc: 0, copper: 0, manganese: 0, selenium: 0,
          chromium: 0, iodine: 0
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

// 分析图片
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
      foodType: '未知食物',
      weight: 100,
      nutrition: {
        calories: 0,
        protein: 0,
        totalFat: 0,
        saturatedFat: 0,
        transFat: 0,
        unsaturatedFat: 0,
        cholesterol: 0,
        carbohydrates: 0,
        sugar: 0,
        fiber: 0,
        water: 0,
        
        vitamins: {
          vitaminA: 0,
          vitaminC: 0,
          vitaminD: 0,
          vitaminE: 0,
          vitaminK: 0,
          vitaminB1: 0,
          vitaminB2: 0,
          vitaminB3: 0,
          vitaminB5: 0,
          vitaminB6: 0,
          vitaminB7: 0,
          vitaminB9: 0,
          vitaminB12: 0
        },
        
        minerals: {
          calcium: 0,
          iron: 0,
          magnesium: 0,
          phosphorus: 0,
          potassium: 0,
          sodium: 0,
          zinc: 0,
          copper: 0,
          manganese: 0,
          selenium: 0,
          chromium: 0,
          iodine: 0
        }
      }
    };
    
    // 将API返回的结果合并到默认结果对象中
    if (apiResult) {
      // 基本属性
      if (apiResult.foodType) defaultResult.foodType = apiResult.foodType;
      if (apiResult.weight) defaultResult.weight = apiResult.weight;
      
      // 营养素
      if (apiResult.nutrition) {
        const nutrition = apiResult.nutrition;
        
        // 基本营养素
        if (nutrition.calories !== undefined) defaultResult.nutrition.calories = nutrition.calories;
        if (nutrition.protein !== undefined) defaultResult.nutrition.protein = nutrition.protein;
        if (nutrition.totalFat !== undefined) defaultResult.nutrition.totalFat = nutrition.totalFat;
        if (nutrition.carbohydrates !== undefined) defaultResult.nutrition.carbohydrates = nutrition.carbohydrates;
        if (nutrition.sugar !== undefined) defaultResult.nutrition.sugar = nutrition.sugar;
        if (nutrition.fiber !== undefined) defaultResult.nutrition.fiber = nutrition.fiber;
        
        // 维生素
        if (nutrition.vitamins) {
          const vitamins = nutrition.vitamins;
          if (vitamins.vitaminA !== undefined) defaultResult.nutrition.vitamins.vitaminA = vitamins.vitaminA;
          if (vitamins.vitaminC !== undefined) defaultResult.nutrition.vitamins.vitaminC = vitamins.vitaminC;
          if (vitamins.vitaminD !== undefined) defaultResult.nutrition.vitamins.vitaminD = vitamins.vitaminD;
        }
        
        // 矿物质
        if (nutrition.minerals) {
          const minerals = nutrition.minerals;
          if (minerals.calcium !== undefined) defaultResult.nutrition.minerals.calcium = minerals.calcium;
          if (minerals.iron !== undefined) defaultResult.nutrition.minerals.iron = minerals.iron;
          if (minerals.sodium !== undefined) defaultResult.nutrition.minerals.sodium = minerals.sodium;
          if (minerals.potassium !== undefined) defaultResult.nutrition.minerals.potassium = minerals.potassium;
        }
      }
    }
    
    console.log('创建的结果对象:', defaultResult);
    
    // 确保结果对象中的所有嵌套对象都已创建
    result.value = defaultResult;
    
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
  if (!result.value) return;
  
  loading.value = true;
  try {
    const recordData: CreateFoodRecordData = {
      // 基础信息
      food_name: result.value.foodType,
      weight: result.value.weight,
      record_time: new Date().toISOString(),
      meal_type: mealTypes[mealTypeIndex.value],
      notes: notes.value,

      // 基本营养素
      calories: result.value.nutrition.calories,
      protein: result.value.nutrition.protein,
      total_fat: result.value.nutrition.totalFat,
      saturated_fat: result.value.nutrition.saturatedFat,
      trans_fat: result.value.nutrition.transFat,
      unsaturated_fat: result.value.nutrition.unsaturatedFat,
      carbohydrates: result.value.nutrition.carbohydrates,
      sugar: result.value.nutrition.sugar,
      fiber: result.value.nutrition.fiber,

      // 维生素
      vitamin_a: result.value.nutrition.vitamins.vitaminA,
      vitamin_c: result.value.nutrition.vitamins.vitaminC,
      vitamin_d: result.value.nutrition.vitamins.vitaminD,
      vitamin_b1: result.value.nutrition.vitamins.vitaminB1,
      vitamin_b2: result.value.nutrition.vitamins.vitaminB2,

      // 矿物质
      calcium: result.value.nutrition.minerals.calcium,
      iron: result.value.nutrition.minerals.iron,
      sodium: result.value.nutrition.minerals.sodium,
      potassium: result.value.nutrition.minerals.potassium
    };
    
    // 使用类型断言将数据转为API需要的格式
    await api.food.createFoodRecord(recordData as any);
    
    uni.showToast({
      title: '记录保存成功',
      icon: 'success'
    });
    
    handleReset();
    
  } catch (err) {
    console.error('保存失败:', err);
    uni.showToast({
      title: '保存失败，请重试',
      icon: 'none'
    });
  } finally {
    loading.value = false;
  }
};

// 放弃记录
const handleDiscardRecord = () => {
  uni.showModal({
    title: '提示',
    content: '确定要放弃当前的食物记录吗？',
    success: (res) => {
      if (res.confirm) {
        result.value = null;
        notes.value = '';
      }
    }
  });
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
  padding-top: 20px;
  border-top: 1px solid rgba(0, 223, 255, 0.3);
}

.form-title {
  color: #00DFFF;
  font-size: 18px;
  margin-bottom: 15px;
}

.form-item {
  margin-bottom: 15px;
}

.form-label {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  margin-bottom: 5px;
}

.meal-type-picker {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 10px;
}

.picker-value {
  color: #fff;
}

.notes-textarea {
  width: 100%;
  height: 80px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 10px;
  color: #fff;
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
</style> 