<template>
  <view class="health-analysis-container">
    <view class="health-analysis-content">
      
      <view class="health-analysis-section">
        <text class="section-description">
          根据您的饮食记录，我们可以为您提供个性化的健康分析报告。请选择您希望分析的时间范围和分析类型。
        </text>
        
        <view class="analysis-form">
          <view class="form-group date-picker-group">
            <text class="form-label">起始日期</text>
            <view class="cyber-input-wrapper">
              <picker 
                mode="date" 
                :value="startDate"
                @change="handleStartDateChange"
                class="cyber-date-picker"
              >
                <view class="picker-content">
                  <text class="picker-value">{{ startDate || '选择日期' }}</text>
                  <text class="picker-icon">📅</text>
                </view>
              </picker>
              <view class="cyber-input-border"></view>
            </view>
          </view>
          
          <view class="form-group date-picker-group">
            <text class="form-label">结束日期</text>
            <view class="cyber-input-wrapper">
              <picker 
                mode="date" 
                :value="endDate"
                @change="handleEndDateChange"
                class="cyber-date-picker"
              >
                <view class="picker-content">
                  <text class="picker-value">{{ endDate || '选择日期' }}</text>
                  <text class="picker-icon">📅</text>
                </view>
              </picker>
              <view class="cyber-input-border"></view>
            </view>
          </view>
          
          <view class="form-group">
            <text class="form-label">分析类型</text>
            <view class="cyber-input-wrapper">
              <picker 
                mode="selector" 
                :range="analysisTypeOptions" 
                :value="analysisTypeIndex"
                @change="handleAnalysisTypeChange"
                class="cyber-select-picker"
              >
                <view class="picker-content">
                  <text class="picker-value">{{ analysisTypeOptions[analysisTypeIndex] }}</text>
                  <text class="picker-icon">▼</text>
                </view>
              </picker>
              <view class="cyber-input-border"></view>
            </view>
          </view>
          
          <view class="form-group">
            <text class="form-label">附加描述 (可选)</text>
            <view class="cyber-input-wrapper textarea-wrapper">
              <textarea 
                v-model="description"
                placeholder="请输入您的健康关注点或饮食习惯，例如：我正在减肥，每天运动30分钟..."
                class="cyber-textarea"
              />
              <view class="cyber-input-border"></view>
            </view>
          </view>
          
          <view class="error-message" v-if="error">{{ error }}</view>
          
          <button 
            class="analyze-button" 
            :disabled="loading"
            @tap="handleAnalyze"
          >
            <view class="analyze-button-content">
              <text class="button-text">{{ loading ? '分析中...' : '开始分析' }}</text>
              <view class="button-glitch"></view>
            </view>
          </button>
        </view>

        <!-- 加载状态 -->
        <view class="loading-container" v-if="loading">
          <view class="loading-spinner"></view>
          <text class="loading-text">正在分析您的健康数据，请稍候...</text>
        </view>

        <!-- 分析结果显示 -->
        <view class="analysis-result" v-if="analysisResult">
          <view class="result-header">
            <text class="result-title">您的健康分析报告</text>
            <text class="analysis-date" v-if="analysisDate">分析时间: {{ analysisDate }}</text>
          </view>
          
          <view class="result-content">
            <rich-text :nodes="formattedResult"></rich-text>
          </view>
        </view>
      </view>
      
      <view class="cyber-line"></view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import api from '@/api';

// 定义接口
interface HealthAnalysisResponse {
  analysis: string;
  date?: string;
}

// 状态变量
const startDate = ref<string>('');
const endDate = ref<string>('');
const description = ref<string>('');
const analysisTypeIndex = ref<number>(0);
const analysisTypeOptions = ['全面分析', '营养均衡分析', '热量摄入分析', '宏量营养素分析'];
const analysisTypes = ['comprehensive', 'nutrition', 'calories', 'macros'] as const;
const analysisResult = ref<string>('');
const loading = ref<boolean>(false);
const error = ref<string>('');
const analysisDate = ref<string>('');

// 日期选择器事件处理函数
const handleStartDateChange = (e: any) => {
  startDate.value = e.detail.value;
};

const handleEndDateChange = (e: any) => {
  endDate.value = e.detail.value;
};

// 格式化分析结果
const formattedResult = computed(() => {
  if (!analysisResult.value) return '';
  
  try {
    // 处理Markdown语法 (简单实现)
    let processedText = analysisResult.value
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // 粗体
      .replace(/\*(.*?)\*/g, '<em>$1</em>') // 斜体
      .replace(/###\s+(.*?)(?:\n|$)/g, '<h4 style="color:#00DFFF;margin:10px 0;">$1</h4>') // 三级标题
      .replace(/##\s+(.*?)(?:\n|$)/g, '<h3 style="color:#00DFFF;margin:12px 0;">$1</h3>') // 二级标题
      .replace(/#\s+(.*?)(?:\n|$)/g, '<h2 style="color:#00DFFF;margin:15px 0;">$1</h2>'); // 一级标题
    
    // 将换行符转换为<br>
    processedText = processedText.replace(/\n/g, '<br>');
    
    return processedText;
  } catch (err) {
    console.error('格式化分析结果出错:', err);
    return analysisResult.value;
  }
});

// 获取分析类型值
const getAnalysisType = () => {
  return analysisTypes[analysisTypeIndex.value] as "comprehensive" | "nutrition" | "calories" | "macros";
};

// 处理分析类型选择变化
const handleAnalysisTypeChange = (e: any) => {
  analysisTypeIndex.value = parseInt(e.detail.value);
};

// 处理分析请求
const handleAnalyze = async () => {
  // 重置状态
  analysisResult.value = '';
  error.value = '';
  analysisDate.value = '';
  
  // 表单验证
  if (!startDate.value || !endDate.value) {
    error.value = '请选择起始和结束日期';
    return;
  }

  if (new Date(startDate.value) > new Date(endDate.value)) {
    error.value = '起始日期不能晚于结束日期';
    return;
  }

  loading.value = true;

  try {
    // 准备请求数据
    const analysisData = {
      start_date: startDate.value,
      end_date: endDate.value,
      analysis_type: getAnalysisType(),
      description: description.value
    };

    const result = await api.health.analyzeHealth(analysisData) as HealthAnalysisResponse;
    
    if (result && typeof result.analysis === 'string') {
      analysisResult.value = result.analysis;
      // 设置分析日期，确保是字符串类型
      analysisDate.value = typeof result.date === 'string' 
        ? result.date 
        : new Date().toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false
          });
    } else {
      analysisResult.value = "分析完成，但未返回详细结果。";
      analysisDate.value = new Date().toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      });
    }
  } catch (err: any) {
    error.value = err?.message || '分析请求失败，请重试';
    console.error('健康分析失败:', err);
  } finally {
    loading.value = false;
  }
};

// 生成测试数据（仅用于开发测试）
const getMockAnalysisResult = (params: any) => {
  const analysisType = params.analysis_type;
  const startDate = new Date(params.start_date).toLocaleDateString('zh-CN');
  const endDate = new Date(params.end_date).toLocaleDateString('zh-CN');
  
  let mockAnalysis = '';
  
  switch (analysisType) {
    case 'comprehensive':
      mockAnalysis = `# 全面健康分析报告
      
## 分析周期: ${startDate} 至 ${endDate}

### 总体评估
您的饮食状况整体表现为**中等水平**。在此期间，您摄入了充足的蛋白质，但是碳水化合物和膳食纤维的摄入略低于建议值。

### 热量摄入分析
- 平均每日热量摄入: **1850千卡**
- 推荐热量摄入: **2000千卡**
- 评估: 您的热量摄入略低于推荐值，如果您正在减重，这是一个合理的范围。

### 营养素平衡
- 蛋白质: **优** (每日平均85克，占总热量的18%)
- 脂肪: **良** (每日平均65克，占总热量的32%)
- 碳水化合物: **中** (每日平均210克，占总热量的50%)
- 膳食纤维: **中** (每日平均18克，建议25克)
- 维生素摄入: **良**，但维生素D摄入不足
- 矿物质: **中**，钙和镁摄入不足

### 改进建议
1. 增加全谷物、豆类和新鲜蔬果的摄入，提高膳食纤维含量
2. 增加富含维生素D的食物，如鱼类、蛋黄和强化乳制品
3. 适当增加钙和镁的摄入，可以考虑添加坚果、绿叶蔬菜和乳制品`;
      break;
    
    case 'nutrition':
      mockAnalysis = `# 营养均衡分析报告
      
## 分析周期: ${startDate} 至 ${endDate}

### 微量营养素摄入情况
您的微量营养素摄入总体状况**需要改善**。以下是详细分析：

#### 维生素摄入评估
- 维生素A: **良好** (120% 日推荐量)
- 维生素B群: **良好** (平均达到推荐量的105%)
- 维生素C: **优秀** (150% 日推荐量)
- 维生素D: **不足** (仅达到推荐量的40%)
- 维生素E: **中等** (80% 日推荐量)
- 叶酸: **良好** (110% 日推荐量)

#### 矿物质摄入评估
- 钙: **不足** (仅达到推荐量的65%)
- 铁: **中等** (90% 日推荐量)
- 锌: **良好** (100% 日推荐量)
- 镁: **不足** (75% 日推荐量)
- 钾: **中等** (80% 日推荐量)
- 硒: **良好** (110% 日推荐量)

### 食物多样性评分
您的食物多样性评分为 **7.5/10**，这是一个不错的分数，但仍有提升空间。

### 改进建议
1. 增加乳制品和绿叶蔬菜的摄入，提高钙的摄入量
2. 考虑适当补充维生素D，或增加阳光照射时间
3. 增加全谷物和坚果的摄入，提高镁的摄入量
4. 尝试每周至少摄入5种不同颜色的蔬菜和水果，增加食物多样性`;
      break;
    
    case 'calories':
      mockAnalysis = `# 热量摄入分析报告
      
## 分析周期: ${startDate} 至 ${endDate}

### 热量摄入概况
- 平均每日热量摄入: **1850千卡**
- 最高单日热量摄入: **2300千卡** (${startDate})
- 最低单日热量摄入: **1500千卡** (${endDate})
- 基础代谢率(BMR)估计: **1650千卡**
- 总能量消耗估计: **2100千卡**

### 热量平衡状况
根据您的热量摄入和估计消耗，您处于**轻微热量赤字**状态，平均每日赤字约为**250千卡**。按照这个趋势，您预计每月可以减重约0.75公斤。

### 热量来源分析
- 碳水化合物: **50%** (920千卡)
- 脂肪: **32%** (590千卡)
- 蛋白质: **18%** (340千卡)

### 改进建议
1. 您目前的热量摄入适合缓慢减重，如果这是您的目标，请保持当前摄入
2. 建议不要将热量摄入降低至1500千卡以下，以确保摄入足够的营养素
3. 可以适当增加蛋白质占比，有助于保持肌肉量和增加饱腹感
4. 注意保持均衡的三餐，避免过度依赖零食获取热量`;
      break;
    
    case 'macros':
      mockAnalysis = `# 宏量营养素分析报告
      
## 分析周期: ${startDate} 至 ${endDate}

### 宏量营养素比例
- 当前摄入比例: 碳水化合物 **50%** : 蛋白质 **18%** : 脂肪 **32%**
- 推荐健康比例: 碳水化合物 45-65% : 蛋白质 10-35% : 脂肪 20-35%
- 评估: 您的宏量营养素比例处于**健康范围内**

### 碳水化合物质量评估
- 总碳水化合物: **210克/天**
- 膳食纤维: **18克/天** (推荐25克/天)
- 添加糖: **42克/天** (建议限制在25克/天以下)
- 碳水化合物质量评分: **中等**

### 蛋白质质量评估
- 总蛋白质: **85克/天**
- 优质蛋白来源: **65%** (动物蛋白、大豆蛋白)
- 所有必需氨基酸充足度: **良好**
- 蛋白质质量评分: **良好**

### 脂肪质量评估
- 总脂肪: **65克/天**
- 饱和脂肪: **20克/天** (略高，建议限制在16克/天以下)
- 不饱和脂肪: **40克/天** (良好)
- 反式脂肪: **<1克/天** (优秀)
- Omega-3脂肪酸: **1.2克/天** (略低，建议1.6克/天)
- 脂肪质量评分: **中等偏上**

### 改进建议
1. 减少添加糖的摄入，选择天然甜味的水果代替甜点
2. 增加膳食纤维摄入，多食用全谷物、豆类和蔬菜
3. 适当减少饱和脂肪的摄入，选择更多植物油和鱼类
4. 增加富含Omega-3的食物，如亚麻籽、核桃和深海鱼`;
      break;
    
    default:
      mockAnalysis = `# 健康分析报告

无法识别您请求的分析类型"${analysisType}"，请选择有效的分析类型。`;
  }
  
  return {
    analysis: mockAnalysis,
    date: new Date().toLocaleString('zh-CN')
  };
};

// 初始化日期
const initializeDates = () => {
  // 获取当前日期
  const today = new Date();
  const endDateValue = today.toISOString().split('T')[0];
  
  // 获取7天前的日期
  const sevenDaysAgo = new Date();
  sevenDaysAgo.setDate(today.getDate() - 7);
  const startDateValue = sevenDaysAgo.toISOString().split('T')[0];
  
  // 设置默认日期
  startDate.value = startDateValue;
  endDate.value = endDateValue;
};

// 页面加载时初始化日期
initializeDates();
</script>

<style>
.health-analysis-container {
  min-height: 100vh;
  background-color: #121212;
  position: relative;
  background-image: url('static/images/food-record-bj.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.health-analysis-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  z-index: 1;
}

.health-analysis-content {
  position: relative;
  z-index: 2;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 30px;
}

.page-title {
  color: #00DFFF;
  font-size: 28px;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 2px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5),
               0 0 20px rgba(0, 223, 255, 0.3),
               0 0 30px rgba(0, 223, 255, 0.1);
  position: relative;
}

.page-title::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent, 
    #00DFFF, 
    transparent
  );
}

.health-analysis-section {
  background: rgba(0, 223, 255, 0.05);
  border: 2px solid rgba(0, 223, 255, 0.2);
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  position: relative;
  box-shadow: 0 0 15px rgba(0, 223, 255, 0.1);
}

.section-header {
  margin-bottom: 15px;
  border-bottom: 1px solid rgba(0, 223, 255, 0.2);
  padding-bottom: 10px;
}

.section-title {
  color: #00DFFF;
  font-size: 18px;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 1px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.3);
}

.section-description {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin-bottom: 20px;
  display: block;
}

.analysis-form {
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  color: #00DFFF;
  font-size: 14px;
  margin-bottom: 8px;
  display: block;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.cyber-input-wrapper {
  position: relative;
  background: rgba(0, 0, 0, 0.4);
  border-radius: 4px;
  overflow: hidden;
}

.cyber-input-border {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent, 
    #00DFFF, 
    #ff00ea, 
    transparent
  );
  animation: border-flow 3s linear infinite;
}

@keyframes border-flow {
  0% {
    background-position: -300px 0;
  }
  100% {
    background-position: 300px 0;
  }
}

.cyber-date-picker, .cyber-select-picker {
  width: 100%;
}

.picker-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
}

.picker-value {
  color: #FFFFFF;
  font-size: 14px;
}

.picker-icon {
  color: #00DFFF;
  font-size: 16px;
  text-shadow: 0 0 5px #00DFFF;
}

.cyber-textarea {
  width: 100%;
  height: 100px;
  background: transparent;
  color: #FFFFFF;
  padding: 12px 15px;
  border: none;
}

.textarea-wrapper {
  border: 1px solid rgba(0, 223, 255, 0.3);
  background: rgba(0, 0, 0, 0.4);
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

.analyze-button {
  background: none;
  border: none;
  padding: 0;
  width: 100%;
}

.analyze-button-content {
  position: relative;
  background: linear-gradient(45deg, #00DFFF, #ff00ea);
  color: #FFFFFF;
  border-radius: 4px;
  padding: 15px;
  font-size: 16px;
  text-transform: uppercase;
  letter-spacing: 1px;
  overflow: hidden;
  box-shadow: 0 0 15px rgba(0, 223, 255, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}

.button-text {
  position: relative;
  z-index: 2;
  font-weight: bold;
}

.button-glitch {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.2);
  opacity: 0;
}

.analyze-button-content:active .button-glitch {
  opacity: 1;
  animation: glitch 0.3s linear;
}

@keyframes glitch {
  0% {
    transform: translate(0);
  }
  20% {
    transform: translate(-2px, 2px);
  }
  40% {
    transform: translate(-2px, -2px);
  }
  60% {
    transform: translate(2px, 2px);
  }
  80% {
    transform: translate(2px, -2px);
  }
  100% {
    transform: translate(0);
  }
}

.loading-container {
  text-align: center;
  padding: 30px;
}

.loading-spinner {
  display: inline-block;
  width: 40px;
  height: 40px;
  border: 3px solid rgba(0, 223, 255, 0.3);
  border-radius: 50%;
  border-top-color: #00DFFF;
  animation: spin 1s linear infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  color: #00DFFF;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
}

.analysis-result {
  background: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  padding: 20px;
  margin-top: 20px;
}

.result-header {
  border-bottom: 1px solid rgba(0, 223, 255, 0.3);
  padding-bottom: 10px;
  margin-bottom: 15px;
}

.result-title {
  color: #00DFFF;
  font-size: 18px;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 1px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.3);
  display: block;
  margin-bottom: 5px;
}

.analysis-date {
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
  display: block;
}

.result-content {
  color: #FFFFFF;
  line-height: 1.6;
}

.cyber-line {
  position: relative;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent, 
    rgba(0, 223, 255, 0.8), 
    transparent
  );
  margin: 20px auto;
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
</style> 