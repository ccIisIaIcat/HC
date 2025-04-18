<template>
  <view class="community-container">
    <!-- 页面内容区域 -->
    <view class="content-area">
      <!-- 用户信息卡片 -->
      <view class="user-info-card">
        <view class="user-avatar" @click="handleAvatarClick">
          <image src="/static/touxiang.png" mode="aspectFill"></image>
        </view>
        <view class="user-details">
          <text class="username">{{ userInfo.name || '加载中...' }}</text>
          <text class="user-email">{{ userInfo.email }}</text>
        </view>
      </view>

      <!-- 打卡状态展示 -->
      <view class="check-in-status" v-if="checkInStatus">
        <text class="status-text" :class="{ 'checked': checkInStatus.has_checked_in }">
          {{ checkInStatus.has_checked_in ? '今日已打卡' : '今日未打卡' }}
        </text>
        <text class="days-text">已累计打卡 <text class="days-count">{{ checkInStatus.check_in_days || 0 }}</text> 天</text>
      </view>
      
      <!-- 功能按钮区域 -->
      <view class="function-buttons">
        <view class="function-button" @click="handleTreasureBoxClick">
          <view class="button-icon-wrapper">
            <image src="/static/comunity/box.png" mode="aspectFit" class="button-icon"></image>
          </view>
          <text class="button-text">百宝箱</text>
        </view>
        <view class="function-button" @click="handleFriendsClick">
          <view class="button-icon-wrapper">
            <image src="/static/comunity/fs.png" mode="aspectFit" class="button-icon"></image>
          </view>
          <text class="button-text">朋友</text>
        </view>
        <view class="function-button" @click="handleSquareClick">
          <view class="button-icon-wrapper">
            <image src="/static/comunity/sq.png" mode="aspectFit" class="button-icon"></image>
          </view>
          <text class="button-text">广场</text>
        </view>
      </view>
    </view>

    <!-- 打卡按钮 -->
    <view class="check-in-button" @click="handleCheckIn">
      <image src="/static/images/check_in.png" mode="aspectFit" class="check-in-icon"></image>
      <text class="check-in-text">打卡</text>
    </view>

    <!-- 自定义弹窗 -->
    <view class="custom-modal" v-if="showModal" @click="closeModal">
      <view class="modal-content" @click.stop>
        <view class="modal-title">喵呜 ฅ^•ﻌ•^ฅ</view>
        <view class="modal-text">猫好人坏，乖乖当一个饭团</view>
        <view class="modal-text">（还没有更换头像功能）</view>
        <view class="modal-button" @click="closeModal">好的~</view>
      </view>
    </view>

    <!-- 打卡弹窗 -->
    <view class="custom-modal" v-if="showCheckInModal" @click="closeCheckInModal">
      <view class="modal-content" @click.stop>
        <view class="modal-title">今日打卡 ฅ^•ﻌ•^ฅ</view>
        <view class="modal-text" v-if="!checkInStatus?.has_food_record">今天还没有记录饮食哦，先去记录一下吧！</view>
        <view class="modal-text" v-else-if="checkInStatus?.has_checked_in">今天已经打过卡啦，明天再来哦！</view>
        <view v-else>
          <view class="modal-input-area">
            <textarea 
              class="check-in-input" 
              v-model="checkInContent" 
              placeholder="记录一下今天的心情吧..."
              maxlength="500"
            ></textarea>
          </view>
          <view class="modal-buttons">
            <view class="modal-button cancel" @click="closeCheckInModal">取消</view>
            <view class="modal-button confirm" @click="submitCheckIn">打卡</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import CustomNavBar from '@/components/CustomNavBar.vue';
import checkInApi from '@/api/checkin';

interface UserData {
  ID: number;
  Name: string;
  Email: string;
  Role: string;
  Verified: boolean;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: null | string;
}

interface ApiResponse {
  user: UserData;
}

// 用户信息数据
const userInfo = ref({
  name: '未知用户',
  email: '无邮箱'
});
const loading = ref(false);

// 获取用户信息
const fetchUserData = async () => {
  try {
    loading.value = true;
    
    // 检查是否有token
    const token = uni.getStorageSync('token');
    if (!token) return;
    
    // 直接使用uni.request测试API连接
    uni.request({
      url: 'https://jesuvukndxpo.sealoshzh.site/api/me',
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (res: any) => {
        if (res.statusCode === 200 && res.data) {
          // 从嵌套的user字段中获取用户数据
          const apiResponse = res.data as ApiResponse;
          if (apiResponse.user && apiResponse.user.Name) {
            userInfo.value = {
              name: apiResponse.user.Name,
              email: apiResponse.user.Email
            };
          }
        } else if (res.statusCode === 401) {
          // 未授权，可能是token过期
          uni.removeStorageSync('token');
        }
      },
      fail: (err: any) => {
        console.error('请求失败:', err);
      }
    });
    
  } catch (error: any) {
    console.error('获取用户信息失败:', error);
  } finally {
    loading.value = false;
  }
};

const showModal = ref(false);

const handleAvatarClick = () => {
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

// 打卡相关的状态
const showCheckInModal = ref(false);
const checkInContent = ref('');
const checkInStatus = ref<{ 
  has_checked_in: boolean; 
  has_food_record: boolean;
  check_in_days?: number;
} | null>(null);

// 获取打卡天数
const fetchCheckInDays = async () => {
  try {
    const response = await checkInApi.getCheckInList({ page_size: 1000 });
    if (response.success && checkInStatus.value) {
      // 计算不同打卡日期的数量
      const uniqueDates = new Set();
      response.data.records.forEach(record => {
        if (record.check_in_at) {
          // 提取日期部分（YYYY-MM-DD）
          const dateOnly = record.check_in_at.split('T')[0];
          uniqueDates.add(dateOnly);
        }
      });
      
      // 更新打卡天数
      checkInStatus.value.check_in_days = uniqueDates.size;
    }
  } catch (error) {
    console.error('获取打卡天数失败:', error);
  }
};

// 获取打卡状态
const fetchCheckInStatus = async () => {
  try {
    const response = await checkInApi.getTodayStatus();
    if (response.success) {
      checkInStatus.value = {
        ...response.data,
        check_in_days: 0
      };
      // 获取打卡天数
      await fetchCheckInDays();
    }
  } catch (error) {
    console.error('获取打卡状态失败:', error);
  }
};

// 处理打卡按钮点击
const handleCheckIn = () => {
  showCheckInModal.value = true;
};

// 关闭打卡弹窗
const closeCheckInModal = () => {
  showCheckInModal.value = false;
  checkInContent.value = '';
};

// 提交打卡
const submitCheckIn = async () => {
  try {
    const content = checkInContent.value.trim();
    
    const response = await checkInApi.checkIn(content);
    if (response.success) {
      uni.showToast({
        title: '打卡成功',
        icon: 'success'
      });
      closeCheckInModal();
      await fetchCheckInStatus(); // 刷新打卡状态
    } else {
      uni.showToast({
        title: response.message,
        icon: 'none'
      });
    }
  } catch (error: any) {
    uni.showToast({
      title: error.message || '打卡失败',
      icon: 'none'
    });
  }
};

// 功能按钮点击处理
const handleTreasureBoxClick = () => {
  uni.navigateTo({
    url: '/pages/achievements/index'
  });
};

const handleFriendsClick = () => {
  uni.showToast({
    title: '朋友功能即将上线',
    icon: 'none'
  });
};

const handleSquareClick = () => {
  uni.showToast({
    title: '广场功能即将上线',
    icon: 'none'
  });
};

onMounted(() => {
  fetchUserData();
  fetchCheckInStatus();
});
</script>

<style>
.community-container {
  min-height: 100vh;
  background-color: #000000;
  background-image: url('/static/images/user-dashboard-bg.jpg');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
  position: relative;
}

.community-container::before {
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

.user-info-card {
  background: rgba(0, 0, 0, 0.8);
  border: 1px solid #00DFFF;
  border-radius: 15px;
  padding: 20px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.2);
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 30px;
  overflow: hidden;
  margin-right: 15px;
  border: 2px solid #00DFFF;
}

.user-avatar image {
  width: 100%;
  height: 100%;
}

.user-details {
  flex: 1;
}

.username {
  font-size: 18px;
  color: #FFFFFF;
  font-weight: bold;
  margin-bottom: 5px;
  display: block;
}

.user-email {
  font-size: 12px;
  color: #CCCCCC;
  display: block;
  margin-top: 5px;
}

/* 自定义弹窗样式 */
.custom-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal-content {
  background-color: rgba(0, 0, 0, 0.9);
  border: 1px solid #00DFFF;
  border-radius: 15px;
  padding: 20px;
  width: 80%;
  max-width: 300px;
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.3);
}

.modal-title {
  color: #FFFFFF;
  font-size: 18px;
  text-align: center;
  margin-bottom: 15px;
  font-weight: bold;
}

.modal-text {
  color: #FFFFFF;
  font-size: 16px;
  text-align: center;
  margin-bottom: 20px;
  line-height: 1.5;
}

.modal-button {
  color: #00DFFF;
  font-size: 16px;
  text-align: center;
  padding: 10px;
  border: 1px solid #00DFFF;
  border-radius: 8px;
  margin: 0 auto;
  width: 120px;
  cursor: pointer;
}

.modal-button:active {
  opacity: 0.8;
}

/* 打卡按钮样式 */
.check-in-button {
  position: fixed;
  right: 20px;
  bottom: 100px;
  width: 60px;
  height: 60px;
  background: rgba(0, 223, 255, 0.1);
  border: 2px solid rgba(0, 223, 255, 0.8);
  border-radius: 30px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 15px rgba(0, 223, 255, 0.5),
              inset 0 0 15px rgba(0, 223, 255, 0.3);
  z-index: 10;
  transition: all 0.3s ease;
  backdrop-filter: blur(4px);
}

.check-in-button:active {
  transform: scale(0.95);
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.7),
              inset 0 0 20px rgba(0, 223, 255, 0.4);
  border-color: rgba(0, 223, 255, 1);
}

.check-in-icon {
  width: 30px;
  height: 30px;
  margin-bottom: 2px;
  filter: drop-shadow(0 0 3px rgba(0, 223, 255, 0.5));
}

.check-in-text {
  color: #FFFFFF;
  font-size: 12px;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.8);
}

/* 打卡状态样式 */
.check-in-status {
  margin-top: 10px;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.status-text {
  color: #CCCCCC;
  font-size: 14px;
}

.status-text.checked {
  color: #00DFFF;
}

.days-text {
  color: #CCCCCC;
  font-size: 14px;
}

.days-count {
  color: #00DFFF;
  font-weight: bold;
  font-size: 16px;
}

/* 打卡输入框样式 */
.modal-input-area {
  margin: 15px auto;
  width: 90%;
}

.check-in-input {
  width: 100%;
  height: 100px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  border-radius: 8px;
  padding: 12px;
  color: #FFFFFF;
  font-size: 14px;
  box-sizing: border-box;
  display: block;
  margin: 0 auto;
}

.check-in-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
  text-align: center;
  line-height: 76px;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  margin: 20px auto 0;
  width: 90%;
}

.modal-button.cancel {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.3);
  flex: 1;
}

.modal-button.confirm {
  background: linear-gradient(45deg, #00DFFF, #00BFFF);
  color: #fff;
  flex: 1;
}

/* 功能按钮样式 */
.function-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 15px;
  flex-wrap: wrap;
}

.function-button {
  background: rgba(0, 0, 0, 0.8);
  border: 1px solid #00DFFF;
  border-radius: 12px;
  padding: 12px;
  width: 70px;
  height: 70px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 10px rgba(0, 223, 255, 0.2);
  transition: all 0.3s ease;
}

.function-button:active {
  transform: scale(0.95);
  box-shadow: 0 0 15px rgba(0, 223, 255, 0.5);
  background: rgba(0, 0, 0, 0.9);
}

.button-icon-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
}

.button-icon {
  width: 40px;
  height: 40px;
  object-fit: contain;
}

.button-text {
  color: #FFFFFF;
  font-size: 14px;
  text-align: center;
  text-shadow: 0 0 5px rgba(0, 223, 255, 0.5);
}
</style> 