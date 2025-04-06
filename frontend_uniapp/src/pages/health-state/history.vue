<template>
  <view class="container">
    <view class="header">
      <text class="title">健康状态历史记录</text>
      <text class="subtitle">查看您的健康指标变化趋势</text>
    </view>
    
    <!-- 日期选择器 -->
    <view class="date-picker">
      <picker 
        mode="date" 
        :value="startDate" 
        start="2020-01-01" 
        :end="endDate" 
        @change="onStartDateChange"
      >
        <view class="picker-item">
          <text class="picker-label">开始日期</text>
          <text class="picker-value">{{ startDate }}</text>
        </view>
      </picker>
      
      <text class="separator">至</text>
      
      <picker 
        mode="date" 
        :value="endDate" 
        :start="startDate" 
        :end="today" 
        @change="onEndDateChange"
      >
        <view class="picker-item">
          <text class="picker-label">结束日期</text>
          <text class="picker-value">{{ endDate }}</text>
        </view>
      </picker>
    </view>
    
    <!-- 记录列表 -->
    <view class="records-container">
      <view v-if="loading" class="loading-container">
        <text class="loading-text">加载中...</text>
      </view>
      
      <view v-else-if="records.length === 0" class="empty-container">
        <text class="empty-text">暂无健康状态记录</text>
      </view>
      
      <view v-else class="records-list">
        <view 
          v-for="(record, index) in records" 
          :key="record.ID" 
          class="record-item"
          @click="viewRecordDetail(record)"
        >
          <view class="record-header">
            <text class="record-date">{{ formatDate(record.CreatedAt) }}</text>
            <text class="record-bmi">BMI: {{ record.bmi }}</text>
          </view>
          
          <view class="record-body">
            <view class="record-row">
              <text class="record-label">身高:</text>
              <text class="record-value">{{ record.height }} cm</text>
            </view>
            
            <view class="record-row">
              <text class="record-label">体重:</text>
              <text class="record-value">{{ record.weight }} kg</text>
            </view>
            
            <view class="record-row">
              <text class="record-label">体脂率:</text>
              <text class="record-value">{{ record.body_fat_percentage || '未记录' }} %</text>
            </view>
          </view>
          
          <view class="record-footer">
            <text class="record-notes" v-if="record.notes">{{ record.notes }}</text>
            <text class="iconfont icon-arrow-right"></text>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 加载更多 -->
    <view v-if="hasMore && !loading" class="load-more" @click="loadMore">
      <text class="load-more-text">加载更多</text>
    </view>
  </view>
</template>

<script>
import { getHealthStates } from '@/api/userstate';

export default {
  data() {
    return {
      records: [],
      loading: false,
      hasMore: true,
      page: 1,
      pageSize: 10,
      startDate: '',
      endDate: '',
      today: ''
    };
  },
  
  onLoad() {
    // 设置默认日期范围（最近30天）
    const today = new Date();
    const thirtyDaysAgo = new Date();
    thirtyDaysAgo.setDate(today.getDate() - 30);
    
    this.today = this.formatDateForPicker(today);
    this.endDate = this.formatDateForPicker(today);
    this.startDate = this.formatDateForPicker(thirtyDaysAgo);
    
    this.loadRecords();
  },
  
  methods: {
    // 格式化日期为 YYYY-MM-DD
    formatDateForPicker(date) {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    },
    
    // 格式化日期为 YYYY年MM月DD日
    formatDate(dateString) {
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}年${month}月${day}日`;
    },
    
    // 开始日期变化
    onStartDateChange(e) {
      this.startDate = e.detail.value;
      this.resetAndReload();
    },
    
    // 结束日期变化
    onEndDateChange(e) {
      this.endDate = e.detail.value;
      this.resetAndReload();
    },
    
    // 重置并重新加载
    resetAndReload() {
      this.page = 1;
      this.records = [];
      this.hasMore = true;
      this.loadRecords();
    },
    
    // 加载记录
    async loadRecords() {
      if (this.loading || !this.hasMore) return;
      
      this.loading = true;
      
      try {
        const res = await getHealthStates({
          page: this.page,
          page_size: this.pageSize,
          start_date: this.startDate,
          end_date: this.endDate
        });
        
        if (res && res.records) {
          // 第一页直接替换，其他页追加
          if (this.page === 1) {
            this.records = res.records;
          } else {
            this.records = [...this.records, ...res.records];
          }
          
          // 判断是否还有更多
          this.hasMore = res.records.length === this.pageSize;
          this.page++;
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
    
    // 加载更多
    loadMore() {
      this.loadRecords();
    },
    
    // 查看记录详情
    viewRecordDetail(record) {
      uni.navigateTo({
        url: `/pages/health-state/detail?id=${record.ID}`
      });
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

.subtitle {
  font-size: 28rpx;
  color: #666;
}

.date-picker {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  border-radius: 20rpx;
  padding: 20rpx;
  margin-bottom: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.picker-item {
  flex: 1;
}

.picker-label {
  font-size: 24rpx;
  color: #666;
  margin-bottom: 10rpx;
  display: block;
}

.picker-value {
  font-size: 28rpx;
  color: #333;
  font-weight: bold;
}

.separator {
  margin: 0 20rpx;
  color: #999;
}

.records-container {
  background-color: #fff;
  border-radius: 20rpx;
  padding: 20rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
  min-height: 300rpx;
}

.loading-container, .empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300rpx;
}

.loading-text, .empty-text {
  font-size: 28rpx;
  color: #999;
}

.record-item {
  padding: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.record-item:last-child {
  border-bottom: none;
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.record-date {
  font-size: 28rpx;
  color: #333;
  font-weight: bold;
}

.record-bmi {
  font-size: 28rpx;
  color: #007AFF;
  font-weight: bold;
}

.record-body {
  margin-bottom: 20rpx;
}

.record-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10rpx;
}

.record-label {
  font-size: 26rpx;
  color: #666;
}

.record-value {
  font-size: 26rpx;
  color: #333;
}

.record-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.record-notes {
  font-size: 24rpx;
  color: #999;
  flex: 1;
  margin-right: 20rpx;
}

.icon-arrow-right {
  font-size: 24rpx;
  color: #999;
}

.load-more {
  text-align: center;
  padding: 30rpx 0;
}

.load-more-text {
  font-size: 28rpx;
  color: #007AFF;
}
</style> 