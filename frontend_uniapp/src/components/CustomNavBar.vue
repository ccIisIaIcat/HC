<template>
  <view class="custom-navbar" :style="{ paddingTop: statusBarHeight + 'px' }">
    <view class="navbar-content">
      <view class="left-area">
        <view v-if="showBack" class="back-btn" @click="goBack">
          <text class="back-icon">◀</text>
        </view>
        <view v-if="showHome" class="home-btn" @click="goHome">
          <text class="iconfont icon-home"></text>
        </view>
      </view>
      
      <view class="title-area">
        <text class="title">{{ title }}</text>
      </view>
      
      <view class="right-area">
        <slot name="right"></slot>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  showBack: {
    type: Boolean,
    default: true
  },
  showHome: {
    type: Boolean,
    default: false
  }
});

const statusBarHeight = ref(20);

onMounted(() => {
  // 获取状态栏高度
  const sysInfo = uni.getSystemInfoSync();
  statusBarHeight.value = sysInfo.statusBarHeight || 20;
});

// 返回上一页
const goBack = () => {
  uni.navigateBack({
    delta: 1,
    fail: () => {
      uni.switchTab({
        url: '/pages/dashboard/index'
      });
    }
  });
};

// 返回首页
const goHome = () => {
  uni.switchTab({
    url: '/pages/dashboard/index'
  });
};
</script>

<style scoped>
.custom-navbar {
  width: 100%;
  background-color: #000000;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 999;
}

.navbar-content {
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
}

.left-area, .right-area {
  width: 80px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.left-area {
  justify-content: flex-start;
}

.right-area {
  justify-content: flex-end;
}

.back-btn, .home-btn {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-icon {
  font-size: 16px;
  color: #00DFFF;
}

.icon-home {
  font-size: 18px;
  color: #00DFFF;
}

.title-area {
  flex: 1;
  text-align: center;
}

.custom-navbar .navbar-content .title-area .title {
  font-size: 18px;
  color: #00DFFF !important;
  font-weight: bold;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5),
               0 0 20px rgba(0, 223, 255, 0.3),
               0 0 30px rgba(0, 223, 255, 0.1);
}
</style> 