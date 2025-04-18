<template>
  <view class="container">
    <view class="header">
      <text class="title">编辑健康状态</text>
      <text class="subtitle">修改您的健康指标记录</text>
    </view>

    <view class="form-container">
      <form @submit="submitHealthState">
        <!-- 身体测量指标 -->
        <view class="section">
          <text class="section-title">身体测量指标</text>
          
          <view class="form-item">
            <text class="label">身高 (cm)</text>
            <input 
              type="digit" 
              v-model="formData.height" 
              placeholder="请输入身高" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">体重 (kg)</text>
            <input 
              type="digit" 
              v-model="formData.weight" 
              placeholder="请输入体重" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">BMI</text>
            <input 
              type="digit" 
              v-model="formData.bmi" 
              placeholder="输入或自动计算" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">体脂率 (%)</text>
            <input 
              type="digit" 
              v-model="formData.body_fat_percentage" 
              placeholder="请输入体脂率" 
              class="input"
            />
          </view>
        </view>
        
        <!-- 基本生命体征 -->
        <view class="section">
          <text class="section-title">基本生命体征</text>
          
          <view class="form-item">
            <text class="label">体温 (℃)</text>
            <input 
              type="digit" 
              v-model="formData.temperature" 
              placeholder="请输入体温" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">心率 (次/分钟)</text>
            <input 
              type="number" 
              v-model="formData.heart_rate" 
              placeholder="请输入心率" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">呼吸频率 (次/分钟)</text>
            <input 
              type="number" 
              v-model="formData.respiratory_rate" 
              placeholder="请输入呼吸频率" 
              class="input"
            />
          </view>
        </view>
        
        <!-- 血糖和血脂 -->
        <view class="section">
          <text class="section-title">血糖和血脂</text>
          
          <view class="form-item">
            <text class="label">空腹血糖 (mmol/L)</text>
            <input 
              type="digit" 
              v-model="formData.fasting_glucose" 
              placeholder="请输入空腹血糖" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">餐后血糖 (mmol/L)</text>
            <input 
              type="digit" 
              v-model="formData.postprandial_glucose" 
              placeholder="请输入餐后血糖" 
              class="input"
            />
          </view>
          
          <view class="form-item">
            <text class="label">总胆固醇 (mmol/L)</text>
            <input 
              type="digit" 
              v-model="formData.total_cholesterol" 
              placeholder="请输入总胆固醇" 
              class="input"
            />
          </view>
        </view>
        
        <!-- 备注 -->
        <view class="section">
          <text class="section-title">备注</text>
          
          <view class="form-item">
            <textarea 
              v-model="formData.notes" 
              placeholder="添加备注信息（可选）" 
              class="textarea"
            />
          </view>
        </view>
        
        <!-- 提交按钮 -->
        <view class="button-container">
          <button 
            type="primary" 
            form-type="submit" 
            class="submit-button"
            :loading="loading"
          >
            保存修改
          </button>
        </view>
      </form>
    </view>
  </view>
</template>

<script>
import { getHealthState, updateHealthState } from '@/api/userstate';

export default {
  data() {
    return {
      recordId: null,
      formData: {
        height: '',
        weight: '',
        bmi: '',
        body_fat_percentage: '',
        temperature: '',
        heart_rate: '',
        respiratory_rate: '',
        fasting_glucose: '',
        postprandial_glucose: '',
        total_cholesterol: '',
        notes: ''
      },
      loading: false
    };
  },
  
  onLoad(options) {
    if (options.id) {
      this.recordId = options.id;
      this.loadRecord(options.id);
    } else {
      uni.showToast({
        title: '记录ID不存在',
        icon: 'none'
      });
      setTimeout(() => {
        uni.navigateBack();
      }, 1500);
    }
  },
  
  methods: {
    // 加载记录
    async loadRecord(id) {
      this.loading = true;
      
      try {
        const res = await getHealthState(id);
        if (res && res.record) {
          // 填充表单数据
          const record = res.record;
          this.formData = {
            height: record.height || '',
            weight: record.weight || '',
            bmi: record.bmi || '',
            body_fat_percentage: record.body_fat_percentage || '',
            temperature: record.temperature || '',
            heart_rate: record.heart_rate || '',
            respiratory_rate: record.respiratory_rate || '',
            fasting_glucose: record.fasting_glucose || '',
            postprandial_glucose: record.postprandial_glucose || '',
            total_cholesterol: record.total_cholesterol || '',
            notes: record.notes || ''
          };
          
          // 如果没有BMI，计算BMI
          if (this.formData.height && this.formData.weight && !this.formData.bmi) {
            this.calculateBMI();
          }
        } else {
          uni.showToast({
            title: '记录不存在',
            icon: 'none'
          });
          setTimeout(() => {
            uni.navigateBack();
          }, 1500);
        }
      } catch (error) {
        console.error('获取健康状态记录失败', error);
        uni.showToast({
          title: '获取记录失败',
          icon: 'none'
        });
      } finally {
        this.loading = false;
      }
    },
    
    // 计算BMI
    calculateBMI() {
      if (this.formData.height && this.formData.weight) {
        const heightInMeters = this.formData.height / 100;
        this.formData.bmi = Math.round((this.formData.weight / (heightInMeters * heightInMeters)) * 100) / 100;
      }
    },
    
    // 监听身高和体重变化，自动计算BMI
    watchHeightWeight() {
      if (this.formData.height && this.formData.weight) {
        this.calculateBMI();
      }
    },
    
    // 提交健康状态
    async submitHealthState() {
      // 表单验证
      if (!this.formData.height || !this.formData.weight) {
        uni.showToast({
          title: '请至少填写身高和体重',
          icon: 'none'
        });
        return;
      }
      
      // 计算BMI
      this.calculateBMI();
      
      this.loading = true;
      
      try {
        // 转换数值类型
        const data = {
          height: parseFloat(this.formData.height),
          weight: parseFloat(this.formData.weight),
          bmi: parseFloat(this.formData.bmi),
          body_fat_percentage: this.formData.body_fat_percentage ? parseFloat(this.formData.body_fat_percentage) : null,
          temperature: this.formData.temperature ? parseFloat(this.formData.temperature) : null,
          heart_rate: this.formData.heart_rate ? parseInt(this.formData.heart_rate) : null,
          respiratory_rate: this.formData.respiratory_rate ? parseInt(this.formData.respiratory_rate) : null,
          fasting_glucose: this.formData.fasting_glucose ? parseFloat(this.formData.fasting_glucose) : null,
          postprandial_glucose: this.formData.postprandial_glucose ? parseFloat(this.formData.postprandial_glucose) : null,
          total_cholesterol: this.formData.total_cholesterol ? parseFloat(this.formData.total_cholesterol) : null,
          notes: this.formData.notes
        };
        
        await updateHealthState(this.recordId, data);
        
        uni.showToast({
          title: '健康状态记录已更新',
          icon: 'success'
        });
        
        // 延迟返回上一页
        setTimeout(() => {
          uni.navigateBack();
        }, 1500);
        
      } catch (error) {
        console.error('更新健康状态记录失败', error);
        uni.showToast({
          title: '更新失败，请重试',
          icon: 'none'
        });
      } finally {
        this.loading = false;
      }
    }
  },
  
  watch: {
    'formData.height': 'watchHeightWeight',
    'formData.weight': 'watchHeightWeight'
  }
};
</script>

<style>
.container {
  padding: 30rpx;
  background-color: #f8f8f8;
  min-height: 100vh;
}

.header {
  margin-bottom: 40rpx;
}

.title {
  font-size: 40rpx;
  font-weight: bold;
  color: #333;
  display: block;
  margin-bottom: 10rpx;
}

.subtitle {
  font-size: 28rpx;
  color: #666;
}

.form-container {
  background-color: #fff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.section {
  margin-bottom: 40rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
  position: relative;
  padding-left: 20rpx;
}

.section-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 8rpx;
  height: 32rpx;
  background-color: #007AFF;
  border-radius: 4rpx;
}

.form-item {
  margin-bottom: 30rpx;
}

.label {
  font-size: 28rpx;
  color: #333;
  margin-bottom: 10rpx;
  display: block;
}

.input {
  width: 100%;
  height: 80rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  color: #333;
}

.textarea {
  width: 100%;
  height: 160rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  padding: 20rpx;
  font-size: 28rpx;
  color: #333;
}

.button-container {
  margin-top: 40rpx;
}

.submit-button {
  width: 100%;
  height: 90rpx;
  line-height: 90rpx;
  background-color: #007AFF;
  color: #fff;
  font-size: 32rpx;
  border-radius: 45rpx;
  font-weight: bold;
}
</style> 