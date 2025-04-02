<template>
  <view class="admin-dashboard">
    <view class="dashboard-header">
      <text class="page-title">管理员控制台</text>
    </view>

    <!-- 系统统计卡片 -->
    <view class="stats-section">
      <view class="stats-grid">
        <view class="stat-card">
          <text class="stat-title">总用户数</text>
          <text class="stat-value">{{ stats?.totalUsers || 0 }}</text>
          <text class="stat-desc">系统注册用户总数</text>
        </view>
        <view class="stat-card">
          <text class="stat-title">今日分析</text>
          <text class="stat-value">{{ stats?.todayAnalysis || 0 }}</text>
          <text class="stat-desc">今日食物分析次数</text>
        </view>
        <view class="stat-card">
          <text class="stat-title">总分析</text>
          <text class="stat-value">{{ stats?.totalAnalysis || 0 }}</text>
          <text class="stat-desc">累计食物分析次数</text>
        </view>
      </view>
    </view>

    <!-- 用户列表 -->
    <view class="users-section">
      <view class="section-header">
        <text class="section-title">用户管理</text>
        <text class="section-subtitle">系统用户列表</text>
      </view>

      <!-- 用户列表 -->
      <view class="users-list">
        <view v-if="loading" class="loading-state">
          <text>加载中...</text>
        </view>
        <view v-else-if="error" class="error-state">
          <text>{{ error }}</text>
        </view>
        <view v-else>
          <view v-for="user in users" :key="user.id" class="user-item">
            <view class="user-info">
              <view class="user-avatar">
                <text>{{ user.name.charAt(0) }}</text>
              </view>
              <view class="user-details">
                <text class="user-name">{{ user.name }}</text>
                <text class="user-email">{{ user.email }}</text>
              </view>
            </view>
            <view class="user-meta">
              <text class="user-role">{{ user.role === 'admin' ? '管理员' : '用户' }}</text>
              <text class="user-date">{{ formatDate(user.createdAt) }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '@/api';

const users = ref<any[]>([]);
const stats = ref<any>(null);
const loading = ref(true);
const error = ref('');

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '未知';
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

// 获取用户列表和统计数据
const fetchDashboardData = async () => {
  try {
    loading.value = true;
    error.value = '';
    
    // 从storage获取用户信息
    const userStr = uni.getStorageSync('user');
    if (!userStr) {
      throw new Error('未登录');
    }
    
    const storedUser = JSON.parse(userStr);
    
    // 验证管理员权限
    if (storedUser.role !== 'admin') {
      uni.reLaunch({
        url: '/pages/dashboard/index'
      });
      return;
    }
    
    // 并行获取用户列表和统计数据
    const [usersList, statsData] = await Promise.all([
      api.admin.getUsers(),
      api.admin.getStats()
    ]);
    
    users.value = usersList;
    stats.value = statsData;
    
  } catch (err: any) {
    console.error('获取数据失败:', err);
    error.value = '获取数据失败，请重试';
    
    // 如果是未登录或token失效，跳转到登录页
    if (err.message === '未登录' || (err.response && err.response.status === 401)) {
      uni.reLaunch({
        url: '/pages/login/index'
      });
    }
  } finally {
    loading.value = false;
  }
};

// 页面加载时获取数据
onMounted(() => {
  fetchDashboardData();
});
</script>

<style>
.admin-dashboard {
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
}

.dashboard-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stats-section {
  margin-bottom: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #1890ff;
  margin-bottom: 4px;
}

.stat-desc {
  font-size: 12px;
  color: #999;
}

.users-section {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.section-header {
  margin-bottom: 20px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
  display: block;
}

.section-subtitle {
  font-size: 14px;
  color: #666;
  display: block;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 20px;
  color: #666;
}

.user-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  flex: 1;
}

.user-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.user-avatar text {
  color: #fff;
  font-size: 16px;
  font-weight: bold;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  display: block;
}

.user-email {
  font-size: 14px;
  color: #666;
  display: block;
}

.user-meta {
  text-align: right;
}

.user-role {
  font-size: 14px;
  color: #1890ff;
  margin-bottom: 4px;
  display: block;
}

.user-date {
  font-size: 12px;
  color: #999;
  display: block;
}
</style> 