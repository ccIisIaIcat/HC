<template>
  <view class="achievements-container">
    <!-- 自定义导航栏 -->
    <CustomNavBar title="我的收藏" />
    
    <!-- 页面内容区域 -->
    <view class="content-area">
      <!-- 九宫格展示区 -->
      <view class="grid-container" v-if="userItems.length > 0">
        <view 
          v-for="(item, index) in currentPageItems" 
          :key="index"
          class="grid-item"
          @click="handleItemClick(item)"
        >
          <image 
            :src="getFullImageUrl(item.icon_url)" 
            mode="aspectFit" 
            class="item-icon"
          ></image>
          <view class="item-name">{{ item.name }}</view>
          <view class="item-quantity">x{{ item.total_quantity }}</view>
        </view>
      </view>

      <!-- 空状态提示 -->
      <view class="empty-state" v-if="!loading && userItems.length === 0">
        <image src="/static/comunity/box.png" mode="aspectFit" class="empty-icon"></image>
        <view class="empty-text">暂无收藏物品</view>
        <view class="empty-subtext">完成任务和打卡可以获得更多物品哦~</view>
      </view>

      <!-- 分页控制 -->
      <view class="pagination" v-if="userItems.length > 0">
        <view 
          class="page-button" 
          :class="{ disabled: currentPage === 1 }"
          @click="handlePageChange(currentPage - 1)"
        >
          上一页
        </view>
        <view class="page-info">{{ currentPage }} / {{ totalPages }}</view>
        <view 
          class="page-button" 
          :class="{ disabled: currentPage === totalPages }"
          @click="handlePageChange(currentPage + 1)"
        >
          下一页
        </view>
      </view>
    </view>

    <!-- 加载提示 -->
    <view class="loading" v-if="loading">
      <text>加载中...</text>
    </view>
    
    <!-- 物品详情弹窗 -->
    <view class="item-detail-popup" v-if="showItemDetail" @click.stop="closeItemDetail">
      <view class="popup-content" @click.stop>
        <view class="popup-close" @click="closeItemDetail">×</view>
        <view class="popup-header">
          <text class="popup-title">{{ selectedItem?.name }}</text>
          <text class="popup-quantity">拥有: {{ selectedItem?.total_quantity }}</text>
        </view>
        <image 
          :src="getFullImageUrl(selectedItem?.image_url)" 
          mode="aspectFit" 
          class="popup-image"
        ></image>
        <view class="popup-info">
          <view class="info-row">
            <text class="info-label">来源:</text>
            <text class="info-value">{{ selectedItem?.source }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">获取方式:</text>
            <text class="info-value">{{ selectedItem?.obtained_from }}</text>
          </view>
          <view class="info-desc">
            <text class="desc-title">描述</text>
            <text class="desc-content">{{ selectedItem?.description }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import CustomNavBar from '@/components/CustomNavBar.vue';
import { getUserItems, type Item } from '../../api/items';
import { apiConfig } from '@/config';

// 状态变量
const loading = ref(false);
const userItems = ref<Item[]>([]);
const currentPage = ref(1);
const itemsPerPage = 9;

// 物品详情弹窗
const showItemDetail = ref(false);
const selectedItem = ref<Item | null>(null);

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(userItems.value.length / itemsPerPage) || 1;
});

// 计算当前页的物品
const currentPageItems = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return userItems.value.slice(start, end);
});

// 获取完整的图片URL
const getFullImageUrl = (url: string | undefined) => {
  if (!url) return '/static/collections/cat-locked.png';
  if (url.startsWith('http')) return url;
  return `${apiConfig.baseURL}${url}`;
};

// 处理页面切换
const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
};

// 点击物品
const handleItemClick = (item: Item) => {
  selectedItem.value = item;
  showItemDetail.value = true;
};

// 关闭物品详情
const closeItemDetail = () => {
  showItemDetail.value = false;
  setTimeout(() => {
    selectedItem.value = null;
  }, 300);
};

// 获取用户物品数据
const fetchUserItems = async () => {
  try {
    // 检查登录状态
    const token = uni.getStorageSync('token');
    
    const userInfo = uni.getStorageSync('user');
    
    // 获取用户ID
    let userId;
    if (userInfo) {
      if (typeof userInfo === 'string') {
        try {
          const parsedUser = JSON.parse(userInfo);
          userId = parsedUser.ID || parsedUser.id;
        } catch (e) {
          console.error('解析userInfo字符串失败:', e);
        }
      } else {
        userId = userInfo.ID || userInfo.id;
      }
    }
    
    if (!userId) {
      userId = uni.getStorageSync('userId');
    }
    
    // 检查登录状态
    if (!token || !userInfo) {
      uni.showToast({
        title: '请先登录',
        icon: 'none',
        duration: 3000
      });
      return;
    }
    
    // 确保有userId
    if (!userId) {
      throw new Error('无法获取用户ID');
    }
    
    // 保存userId到storage
    uni.setStorageSync('userId', userId);
    
    loading.value = true;
    
    // 调用获取物品的API
    const response = await getUserItems();
    
    // 处理返回的数据
    userItems.value = response.data || [];
  } catch (error: any) {
    console.error('获取物品失败:', error);
    
    userItems.value = [];
    
    // 处理登录相关错误
    if (error.message === '请先登录' || error.message.includes('登录') || error.statusCode === 401) {
      uni.showToast({
        title: '请先登录',
        icon: 'none',
        duration: 3000
      });
    } else {
      uni.showToast({
        title: error.message || '获取物品失败，请稍后重试',
        icon: 'none',
        duration: 3000
      });
    }
  } finally {
    loading.value = false;
  }
};

// 页面加载时获取数据
onMounted(() => {
  fetchUserItems();
});
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

/* 分页样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
  padding: 10px;
}

.page-button {
  padding: 8px 15px;
  background: rgba(0, 223, 255, 0.2);
  border: 1px solid #00DFFF;
  border-radius: 5px;
  color: #FFFFFF;
  font-size: 14px;
  margin: 0 10px;
  cursor: pointer;
}

.page-button.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  color: #FFFFFF;
  font-size: 14px;
  margin: 0 15px;
}

/* 数量标签 */
.item-quantity {
  position: absolute;
  bottom: 5px;
  right: 5px;
  background: rgba(0, 0, 0, 0.7);
  color: #00DFFF;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 12px;
}

/* 加载提示 */
.loading {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.7);
  padding: 15px 30px;
  border-radius: 10px;
  z-index: 999;
}

.loading text {
  color: #FFFFFF;
  font-size: 14px;
}

/* 空状态样式 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 50px 20px;
  text-align: center;
}

.empty-icon {
  width: 100px;
  height: 100px;
  margin-bottom: 20px;
  opacity: 0.7;
}

.empty-text {
  color: #FFFFFF;
  font-size: 18px;
  margin-bottom: 10px;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.5);
}

.empty-subtext {
  color: #CCCCCC;
  font-size: 14px;
  max-width: 250px;
  line-height: 1.5;
}

/* 物品详情弹窗 */
.item-detail-popup {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.8);
  z-index: 999;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.popup-content {
  width: 85%;
  max-width: 500px;
  background-color: rgba(0, 10, 30, 0.95);
  border: 2px solid #00DFFF;
  border-radius: 15px;
  padding: 20px;
  position: relative;
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.3);
  max-height: 80vh;
  overflow-y: auto;
}

.popup-close {
  position: absolute;
  top: 10px;
  right: 10px;
  color: #FFFFFF;
  font-size: 24px;
  cursor: pointer;
  width: 30px;
  height: 30px;
  line-height: 30px;
  text-align: center;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 50%;
  z-index: 10;
}

.popup-header {
  margin-bottom: 15px;
  padding-bottom: 10px;
  padding-right: 40px; /* 为关闭按钮留出空间 */
  border-bottom: 1px solid rgba(0, 223, 255, 0.3);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.popup-title {
  color: #00DFFF;
  font-size: 20px;
  font-weight: bold;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.5);
}

.popup-quantity {
  color: #FFFFFF;
  font-size: 14px;
  background: rgba(0, 223, 255, 0.2);
  padding: 3px 8px;
  border-radius: 10px;
}

.popup-image {
  width: 100%;
  height: 200px;
  object-fit: contain;
  margin-bottom: 15px;
  filter: drop-shadow(0 0 8px rgba(0, 223, 255, 0.3));
  background: rgba(0, 0, 0, 0.3);
  border-radius: 10px;
  padding: 10px;
}

.popup-info {
  color: #FFFFFF;
}

.info-row {
  display: flex;
  margin-bottom: 10px;
}

.info-label {
  width: 80px;
  color: #AAAAAA;
  font-size: 14px;
}

.info-value {
  flex: 1;
  font-size: 14px;
}

.info-desc {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid rgba(0, 223, 255, 0.3);
}

.desc-title {
  display: block;
  color: #AAAAAA;
  font-size: 14px;
  margin-bottom: 8px;
}

.desc-content {
  font-size: 15px;
  line-height: 1.5;
  display: block;
  color: #FFFFFF;
}
</style> 