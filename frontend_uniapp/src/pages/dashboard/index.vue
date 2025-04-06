<template>
  <view class="dashboard-container">
    <view class="function-grid">
      <view class="function-item" @click="navigateTo('/pages/food-analysis/index')">
        <text class="function-title">食物分析</text>
        <view class="cyber-line"></view>
      </view>
      
      <view class="function-item" @click="navigateTo('/pages/food-history/index')">
        <text class="function-title">饮食历史</text>
        <view class="cyber-line"></view>
      </view>
      
      <view class="function-item" @click="navigateTo('/pages/food-history/visualization')">
        <text class="function-title">饮食可视化分析</text>
        <view class="cyber-line"></view>
      </view>
      
      <view class="function-item" @click="navigateTo('/pages/health/analysis')">
        <text class="function-title">健康分析</text>
        <view class="cyber-line"></view>
      </view>
      
      <view class="function-item" @click="navigateTo('/pages/health-state/index')">
        <text class="function-title">健康状态记录</text>
        <view class="cyber-line"></view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '@/api';

const user = ref<any>(null);
const loading = ref(false);
const error = ref('');

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    loading.value = true;
    error.value = '';
    
    // 从storage获取用户信息
    const userStr = uni.getStorageSync('user');
    if (!userStr) {
      throw new Error('未登录');
    }
    
    const storedUser = JSON.parse(userStr);
    console.log('当前用户信息:', storedUser);
    
    // 验证用户权限
    if (storedUser.Role === 'admin') {
      uni.reLaunch({
        url: '/pages/admin/dashboard/index'
      });
      return;
    }
    
    // 获取最新的用户信息
    const response = await api.auth.getCurrentUser();
    user.value = response.user;
    
  } catch (err: any) {
    console.error('获取用户信息失败:', err);
    error.value = '获取用户信息失败，请重新登录';
    
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

// 页面跳转
const navigateTo = (url: string) => {
  uni.navigateTo({
    url,
    fail: (err) => {
      console.error('页面跳转失败:', err);
      uni.showToast({
        title: '页面跳转失败',
        icon: 'none'
      });
    }
  });
};

// 页面加载时获取用户信息
onMounted(() => {
  fetchUserInfo();
});
</script>

<style>
.dashboard-container {
  min-height: 100vh;
  background-image: url('/static/images/user-dashboard-bg.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  padding: 20px;
  position: relative;
}

.dashboard-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1;
}

.function-grid {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 30px;
  padding: 20px;
  margin-top: 60px;
}

.function-item {
  background: rgba(0, 223, 255, 0.1);
  border: 2px solid rgba(0, 223, 255, 0.5);
  border-radius: 8px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.2);
  transition: all 0.3s ease;
  aspect-ratio: 1 / 1; /* 使按钮成为正方形 */
}

.function-item::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(45deg, 
    rgba(0, 223, 255, 0.5), 
    rgba(255, 0, 255, 0.5), 
    rgba(0, 223, 255, 0.5)
  );
  z-index: -1;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.function-item:active {
  transform: scale(0.98);
  box-shadow: 0 0 30px rgba(0, 223, 255, 0.4);
}

.function-item:active::before {
  opacity: 1;
}

.cat-icon {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 30px;
  height: 30px;
  transform: rotate(15deg);
  animation: cat-bounce 2s ease-in-out infinite;
}

.cat-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

@keyframes cat-bounce {
  0%, 100% {
    transform: rotate(15deg) translateY(0);
  }
  50% {
    transform: rotate(15deg) translateY(-5px);
  }
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

.function-title {
  font-size: 18px;
  color: #00DFFF;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5),
               0 0 20px rgba(0, 223, 255, 0.3),
               0 0 30px rgba(0, 223, 255, 0.1);
  margin-bottom: 10px;
  padding: 0 20px;
}
</style> 