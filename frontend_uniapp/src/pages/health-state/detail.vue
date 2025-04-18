<template>
  <view class="container">
    <!-- 自定义导航栏 -->
    <CustomNavBar 
      title="健康状态详情" 
      :showBack="true"
      :showHome="true"
    />
    
    <view class="header">
      <text class="title">健康状态详情</text>
      <text class="date">{{ formatDate(record.CreatedAt) }}</text>
    </view>
    
    <view class="content">
      <!-- 身体测量指标 -->
      <view class="section">
        <text class="section-title">身体测量指标</text>
        
        <view class="data-grid">
          <view class="data-item">
            <text class="data-label">身高</text>
            <text class="data-value">{{ record.height }} cm</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">体重</text>
            <text class="data-value">{{ record.weight }} kg</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">BMI</text>
            <text class="data-value">{{ record.bmi }}</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">体脂率</text>
            <text class="data-value">{{ record.body_fat_percentage || '未记录' }} %</text>
          </view>
        </view>
      </view>
      
      <!-- 基本生命体征 -->
      <view class="section">
        <text class="section-title">基本生命体征</text>
        
        <view class="data-grid">
          <view class="data-item">
            <text class="data-label">体温</text>
            <text class="data-value">{{ record.temperature || '未记录' }} ℃</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">心率</text>
            <text class="data-value">{{ record.heart_rate || '未记录' }} 次/分钟</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">呼吸频率</text>
            <text class="data-value">{{ record.respiratory_rate || '未记录' }} 次/分钟</text>
          </view>
        </view>
      </view>
      
      <!-- 血糖和血脂 -->
      <view class="section">
        <text class="section-title">血糖和血脂</text>
        
        <view class="data-grid">
          <view class="data-item">
            <text class="data-label">空腹血糖</text>
            <text class="data-value">{{ record.fasting_glucose || '未记录' }} mmol/L</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">餐后血糖</text>
            <text class="data-value">{{ record.postprandial_glucose || '未记录' }} mmol/L</text>
          </view>
          
          <view class="data-item">
            <text class="data-label">总胆固醇</text>
            <text class="data-value">{{ record.total_cholesterol || '未记录' }} mmol/L</text>
          </view>
        </view>
      </view>
      
      <!-- 备注 -->
      <view class="section" v-if="record.notes">
        <text class="section-title">备注</text>
        
        <view class="notes-container">
          <text class="notes-text">{{ record.notes }}</text>
        </view>
      </view>
    </view>
    
    <!-- 操作按钮 -->
    <view class="action-buttons">
      <button 
        class="action-button edit-button" 
        @click="editRecord"
      >
        编辑记录
      </button>
      
      <button 
        class="action-button delete-button" 
        @click="confirmDelete"
      >
        删除记录
      </button>
    </view>
  </view>
</template>

<script>
import { getHealthState, deleteHealthState } from '@/api/userstate';

export default {
  data() {
    return {
      record: {},
      loading: false
    };
  },
  
  onLoad(options) {
    if (options.id) {
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
    // 格式化日期为 YYYY年MM月DD日
    formatDate(dateString) {
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}年${month}月${day}日`;
    },
    
    // 加载记录详情
    async loadRecord(id) {
      this.loading = true;
      
      try {
        const res = await getHealthState(id);
        if (res && res.record) {
          this.record = res.record;
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
        console.error('获取健康状态记录详情失败', error);
        uni.showToast({
          title: '获取记录失败',
          icon: 'none'
        });
      } finally {
        this.loading = false;
      }
    },
    
    // 编辑记录
    editRecord() {
      uni.navigateTo({
        url: `/pages/health-state/edit?id=${this.record.ID}`
      });
    },
    
    // 确认删除
    confirmDelete() {
      uni.showModal({
        title: '确认删除',
        content: '确定要删除这条健康状态记录吗？此操作不可恢复。',
        success: (res) => {
          if (res.confirm) {
            this.deleteRecord();
          }
        }
      });
    },
    
    // 删除记录
    async deleteRecord() {
      try {
        await deleteHealthState(this.record.ID);
        
        uni.showToast({
          title: '记录已删除',
          icon: 'success'
        });
        
        setTimeout(() => {
          uni.navigateBack();
        }, 1500);
      } catch (error) {
        console.error('删除健康状态记录失败', error);
        uni.showToast({
          title: '删除失败，请重试',
          icon: 'none'
        });
      }
    }
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

.date {
  font-size: 28rpx;
  color: #666;
}

.content {
  background-color: #fff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
  margin-bottom: 40rpx;
}

.section {
  margin-bottom: 40rpx;
}

.section:last-child {
  margin-bottom: 0;
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

.data-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.data-item {
  background-color: #f5f5f5;
  border-radius: 10rpx;
  padding: 20rpx;
}

.data-label {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 10rpx;
  display: block;
}

.data-value {
  font-size: 32rpx;
  color: #333;
  font-weight: bold;
}

.notes-container {
  background-color: #f5f5f5;
  border-radius: 10rpx;
  padding: 20rpx;
}

.notes-text {
  font-size: 28rpx;
  color: #333;
  line-height: 1.5;
}

.action-buttons {
  display: flex;
  justify-content: space-between;
}

.action-button {
  flex: 1;
  height: 90rpx;
  line-height: 90rpx;
  font-size: 32rpx;
  border-radius: 45rpx;
  font-weight: bold;
  margin: 0 10rpx;
}

.edit-button {
  background-color: #007AFF;
  color: #fff;
}

.delete-button {
  background-color: #ff3b30;
  color: #fff;
}
</style> 