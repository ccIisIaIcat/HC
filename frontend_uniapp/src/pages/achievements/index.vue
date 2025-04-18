<template>
  <view class="achievements-container">
    <!-- 自定义导航栏 -->
    <CustomNavBar title="我的收藏" />
    
    <!-- 页面内容区域 -->
    <view class="content-area">
      <!-- 九宫格展示区 -->
      <view class="grid-container">
        <view 
          v-for="(item, index) in collectionItems" 
          :key="index"
          class="grid-item"
          :class="{ 'locked': !item.unlocked }"
          @click="handleItemClick(item)"
        >
          <image 
            :src="item.unlocked ? item.icon : item.lockedIcon" 
            mode="aspectFit" 
            class="item-icon"
          ></image>
          <view class="item-name">{{ item.name }}</view>
          <view class="lock-icon" v-if="!item.unlocked">
            <text class="iconfont icon-lock"></text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import CustomNavBar from '@/components/CustomNavBar.vue';

// 收藏品数据
const collectionItems = ref([
  {
    id: 1,
    name: '饭团猫',
    icon: '/static/collections/cat-normal.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: true
  },
  {
    id: 2,
    name: '运动猫',
    icon: '/static/collections/cat-sport.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 3,
    name: '厨师猫',
    icon: '/static/collections/cat-chef.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 4,
    name: '医生猫',
    icon: '/static/collections/cat-doctor.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 5,
    name: '学者猫',
    icon: '/static/collections/cat-scholar.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 6,
    name: '艺术猫',
    icon: '/static/collections/cat-artist.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 7,
    name: '音乐猫',
    icon: '/static/collections/cat-music.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 8,
    name: '魔法猫',
    icon: '/static/collections/cat-magic.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  },
  {
    id: 9,
    name: '国王猫',
    icon: '/static/collections/cat-king.png',
    lockedIcon: '/static/collections/cat-locked.png',
    unlocked: false
  }
]);

// 点击收藏品
const handleItemClick = (item: any) => {
  if (!item.unlocked) {
    uni.showToast({
      title: '继续努力解锁该收藏品吧！',
      icon: 'none'
    });
    return;
  }
  
  uni.showToast({
    title: `这是${item.name}`,
    icon: 'none'
  });
};
</script>

<style>
.achievements-container {
  min-height: 100vh;
  background-color: #000000;
  background-image: url('/static/images/user-dashboard-bg.jpg');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
  position: relative;
}

.achievements-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  z-index: 1;
}

.content-area {
  position: relative;
  z-index: 2;
  padding: 20px;
  padding-top: calc(var(--status-bar-height) + 44px + 20px);
}

/* 九宫格样式 */
.grid-container {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 15px;
  padding: 15px;
}

.grid-item {
  aspect-ratio: 1;
  background: rgba(0, 0, 0, 0.8);
  border: 2px solid #00DFFF;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px;
  position: relative;
  transition: all 0.3s ease;
}

.grid-item.locked {
  border-color: rgba(0, 223, 255, 0.3);
  opacity: 0.7;
}

.grid-item:active {
  transform: scale(0.95);
}

.item-icon {
  width: 70%;
  height: 70%;
  object-fit: contain;
  margin-bottom: 8px;
  filter: drop-shadow(0 0 5px rgba(0, 223, 255, 0.3));
}

.item-name {
  color: #FFFFFF;
  font-size: 14px;
  text-align: center;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.5);
}

.lock-icon {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 20px;
  height: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.icon-lock {
  color: rgba(255, 255, 255, 0.5);
  font-size: 16px;
}

.grid-item.locked .item-icon {
  filter: grayscale(100%) brightness(0.7);
}

.grid-item.locked .item-name {
  color: rgba(255, 255, 255, 0.5);
}
</style> 