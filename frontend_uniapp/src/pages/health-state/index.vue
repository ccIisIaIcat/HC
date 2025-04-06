<template>
  <view class="health-state-container">
    <!-- 背景图片 -->
    <image class="background-image" src="/static/images/health-ana-bj.jpg" mode="aspectFill"></image>
    
    <!-- 加载状态 -->
    <view class="loading-container" v-if="loading">
      <text class="loading-text">正在加载数据...</text>
    </view>

    <!-- 编辑表单 -->
    <view class="edit-form" v-if="!loading">
      <!-- 身体测量指标部分 -->
      <view class="form-card">
        <view class="section-title">身体测量指标</view>
        <view class="measurement-grid">
          <view class="measurement-item">
            <text class="item-label">身高</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.height" 
                class="measurement-input"
                placeholder="输入身高"
              />
              <text class="unit">cm</text>
            </view>
          </view>
          
          <view class="measurement-item">
            <text class="item-label">体重</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.weight" 
                class="measurement-input"
                placeholder="输入体重"
              />
              <text class="unit">kg</text>
            </view>
          </view>
          
          <view class="measurement-item">
            <text class="item-label">BMI</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.bmi" 
                class="measurement-input"
                disabled
                placeholder="自动计算"
              />
            </view>
          </view>
          
          <view class="measurement-item">
            <text class="item-label">体脂率</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.body_fat_percentage" 
                class="measurement-input"
                placeholder="输入体脂率"
              />
              <text class="unit">%</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 基本生命体征部分 -->
      <view class="form-card">
        <view class="section-title">基本生命体征</view>
        <view class="vital-signs-grid">
          <view class="vital-sign-item">
            <text class="item-label">体温</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.temperature" 
                class="vital-sign-input"
                placeholder="输入体温"
              />
              <text class="unit">℃</text>
            </view>
          </view>
          
          <view class="vital-sign-item">
            <text class="item-label">心率</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="formData.heart_rate" 
                class="vital-sign-input"
                placeholder="输入心率"
              />
              <text class="unit">次/分</text>
            </view>
          </view>
          
          <view class="vital-sign-item">
            <text class="item-label">呼吸频率</text>
            <view class="item-input-wrapper">
              <input 
                type="number" 
                v-model="formData.respiratory_rate" 
                class="vital-sign-input"
                placeholder="输入呼吸频率"
              />
              <text class="unit">次/分</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 血糖和血脂部分 -->
      <view class="form-card">
        <view class="section-title">血糖和血脂</view>
        <view class="blood-grid">
          <view class="blood-item">
            <text class="item-label">空腹血糖</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.fasting_glucose" 
                class="blood-input"
                placeholder="输入空腹血糖"
              />
              <text class="unit">mmol/L</text>
            </view>
          </view>
          
          <view class="blood-item">
            <text class="item-label">餐后血糖</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.postprandial_glucose" 
                class="blood-input"
                placeholder="输入餐后血糖"
              />
              <text class="unit">mmol/L</text>
            </view>
          </view>
          
          <view class="blood-item">
            <text class="item-label">总胆固醇</text>
            <view class="item-input-wrapper">
              <input 
                type="digit" 
                v-model="formData.total_cholesterol" 
                class="blood-input"
                placeholder="输入总胆固醇"
              />
              <text class="unit">mmol/L</text>
            </view>
          </view>
        </view>
      </view>
      
      <!-- 备注部分 -->
      <view class="form-card">
        <view class="section-title">备注</view>
        <view class="notes-section">
          <textarea 
            v-model="formData.notes" 
            class="notes-input"
            placeholder="添加备注信息（可选）"
          />
        </view>
      </view>
      
      <!-- 提交按钮 -->
      <view class="button-container">
        <button 
          class="submit-button" 
          :loading="loading"
          @tap="submitHealthState"
        >
          提交健康状态
        </button>
      </view>
    </view>
    
    <!-- 历史记录入口 -->
    <view class="history-entry" @tap="navigateToHistory">
      <text class="history-text">查看历史记录</text>
      <text class="iconfont icon-arrow-right"></text>
    </view>
  </view>
</template>

<script>
import { createHealthState, getLatestHealthState } from '@/api/userstate';

export default {
  data() {
    return {
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
  
  onLoad() {
    this.loadLatestHealthState();
  },
  
  methods: {
    // 加载最新的健康状态记录
    async loadLatestHealthState() {
      this.loading = true;
      try {
        const res = await getLatestHealthState();
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
        }
      } catch (error) {
        console.error('获取最新健康状态记录失败', error);
        uni.showToast({
          title: '获取最新记录失败',
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
        const bmi = (this.formData.weight / (heightInMeters * heightInMeters)).toFixed(2);
        this.formData.bmi = bmi;
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
        
        const res = await createHealthState(data);
        
        uni.showToast({
          title: '健康状态记录已提交',
          icon: 'success'
        });
        
        // 延迟返回上一页
        setTimeout(() => {
          uni.navigateBack();
        }, 1500);
        
      } catch (error) {
        console.error('提交健康状态记录失败', error);
        uni.showToast({
          title: '提交失败，请重试',
          icon: 'none'
        });
      } finally {
        this.loading = false;
      }
    },
    
    // 跳转到历史记录页面
    navigateToHistory() {
      uni.navigateTo({
        url: '/pages/health-state/history'
      });
    }
  },
  
  watch: {
    'formData.height': 'watchHeightWeight',
    'formData.weight': 'watchHeightWeight'
  }
};
</script>

<style>
.health-state-container {
  min-height: 100vh;
  background-color: #1a1a1a;
  padding: 20px;
  position: relative;
  z-index: 1;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
  opacity: 0.4;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.loading-text {
  color: #00dfff;
  font-size: 16px;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-card {
  background-color: rgba(0, 0, 0, 0.7);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  padding: 16px;
  backdrop-filter: blur(5px);
}

.section-title {
  color: #00dfff;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
}

.measurement-grid,
.vital-signs-grid,
.blood-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.measurement-item,
.vital-sign-item,
.blood-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.item-label {
  color: #ffffff;
  font-size: 14px;
}

.item-input-wrapper {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  padding: 0 10px;
}

.measurement-input,
.vital-sign-input,
.blood-input {
  flex: 1;
  height: 40px;
  background: transparent;
  color: #ffffff;
  font-size: 16px;
}

.unit {
  color: #00dfff;
  font-size: 14px;
  margin-left: 5px;
}

.notes-section {
  margin-top: 10px;
}

.notes-input {
  width: 100%;
  height: 100px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 4px;
  color: #ffffff;
  font-size: 14px;
  padding: 10px;
}

.button-container {
  margin-top: 20px;
}

.submit-button {
  width: 100%;
  height: 50px;
  background: linear-gradient(135deg, #00dfff, #0066ff);
  color: #ffffff;
  font-size: 18px;
  font-weight: bold;
  border-radius: 25px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.history-entry {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 15px;
  background: rgba(0, 0, 0, 0.7);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  backdrop-filter: blur(5px);
}

.history-text {
  color: #00dfff;
  font-size: 16px;
  margin-right: 10px;
}

.icon-arrow-right {
  color: #00dfff;
  font-size: 16px;
}
</style>