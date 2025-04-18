<template>
  <view class="login-container" :style="{ paddingBottom: keyboardHeight + 'px' }">
    <view class="login-box" :style="{ transform: `translateY(${scrollOffset}px)` }">
      <view class="login-header">
        <image src="@/static/logo.png" class="logo-image" mode="aspectFill" />
        <text class="system-title">饭团猫</text>
      </view>

      <view v-if="message" class="success-message">
        <text>✅ {{ message }}</text>
      </view>
      
      <view v-if="error" class="error-message">
        <text>{{ error }}</text>
      </view>

      <view class="login-form">
        <view class="form-item">
          <view class="input-wrapper">
            <input 
              id="email-input"
              type="text"
              class="form-input"
              placeholder="请输入邮箱"
              placeholder-class="input-placeholder"
              v-model="email"
              @focus="handleInputFocus('email-input')"
              @blur="handleInputBlur"
            />
          </view>
        </view>
        
        <view class="form-item">
          <view class="input-wrapper">
            <input 
              id="password-input"
              type="password"
              class="form-input"
              placeholder="请输入密码"
              placeholder-class="input-placeholder"
              v-model="password"
              @focus="handleInputFocus('password-input')"
              @blur="handleInputBlur"
            />
          </view>
        </view>

        <button 
          class="login-button"
          :disabled="loading"
          :class="{ 'button-disabled': loading }"
          @click="handleSubmit"
        >
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </view>

      <view class="login-footer">
        <text>还没有账号？</text>
        <navigator url="/pages/register/index" hover-class="link-hover">立即注册</navigator>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import api from '@/api';

// 响应式状态
const email = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');
const message = ref('');
const scrollOffset = ref(0); // 初始没有滚动偏移
const keyboardHeight = ref(0);

// 处理输入框获得焦点
const handleInputFocus = (inputId: string) => {
  console.log('输入框获得焦点:', inputId);
  
  // 由于整体窗口已经上移，减小初始滚动偏移量
  if (inputId === 'password-input') {
    scrollOffset.value = -80; // 密码输入框滚动偏移
  } else {
    scrollOffset.value = -30; // 邮箱输入框滚动偏移
  }
  
  // 使用 setTimeout 确保在下一个渲染周期执行滚动
  setTimeout(() => {
    uni.createSelectorQuery()
      .select(`#${inputId}`)
      .boundingClientRect(data => {
        if (data && 'bottom' in data && typeof data.bottom === 'number') {
          console.log('输入框位置信息:', data);
          
          // 获取可视区域高度
          const systemInfo = uni.getSystemInfoSync();
          const windowHeight = systemInfo.windowHeight;
          
          // 计算输入框底部到屏幕底部的距离
          const distanceToBottom = windowHeight - data.bottom;
          
          // 如果距离小于键盘高度的一半，增加滚动偏移
          if (distanceToBottom < 150) {
            scrollOffset.value = -(150 - distanceToBottom + 20); // 减小额外偏移
          }
        }
      })
      .exec();
  }, 300);
};

// 处理输入框失去焦点
const handleInputBlur = () => {
  console.log('输入框失去焦点');
  // 重置滚动偏移
  setTimeout(() => {
    scrollOffset.value = 0;
  }, 200);
};

// 监听键盘高度变化
onMounted(() => {
  // #ifdef APP-PLUS
  uni.onKeyboardHeightChange((res) => {
    console.log('键盘高度变化:', res.height);
    keyboardHeight.value = res.height;
    
    // 如果键盘收起，重置滚动
    if (res.height === 0) {
      scrollOffset.value = 0;
    }
  });
  // #endif
});

// 组件卸载时移除监听
onUnmounted(() => {
  // #ifdef APP-PLUS
  uni.offKeyboardHeightChange();
  // #endif
});

// 在页面加载时初始化
onLoad((option: any) => {
  console.log('=== 登录页面开始加载 ===');
  console.log('当前页面路径:', getCurrentPages());
  console.log('页面参数:', option);
  
  if (option.message) {
    console.log('收到消息:', option.message);
    message.value = decodeURIComponent(option.message);
  }
  
  // 检查是否已经登录
  checkAutoLogin();
});

// 检查自动登录
const checkAutoLogin = async () => {
  try {
    // 检查本地存储中是否有token
    const token = uni.getStorageSync('token');
    let userInfo = uni.getStorageSync('user');
    
    // 处理可能的字符串格式的用户信息
    if (userInfo && typeof userInfo === 'string') {
      try {
        userInfo = JSON.parse(userInfo);
      } catch (e) {
        console.error('解析用户信息失败:', e);
        // 如果解析失败，清除可能损坏的数据
        uni.removeStorageSync('user');
        return;
      }
    }
    
    if (token && userInfo) {
      console.log('发现已保存的登录信息，尝试自动登录');
      
      // 验证token是否有效
      try {
        // 显示加载提示
        uni.showLoading({
          title: '自动登录中...',
          mask: true
        });
        
        // 调用API验证token
        const response = await api.auth.getCurrentUser();
        console.log('验证登录状态成功:', response);
        
        // 登录成功，根据用户角色跳转到相应页面
        const userData = response.user || userInfo;
        
        setTimeout(() => {
          uni.hideLoading();
          if (userData.Role === 'admin') {
            console.log('管理员用户，跳转到管理后台');
            uni.reLaunch({
              url: '/pages/admin/dashboard/index'
            });
          } else {
            console.log('普通用户，跳转到个人中心');
            uni.reLaunch({
              url: '/pages/dashboard/index'
            });
          }
        }, 500);
      } catch (err) {
        // token已过期或无效
        console.error('自动登录失败，token可能已过期:', err);
        uni.hideLoading();
        // 清除无效的登录信息
        uni.removeStorageSync('token');
        uni.removeStorageSync('user');
        // 显示需要重新登录的提示
        uni.showToast({
          title: '登录已过期，请重新登录',
          icon: 'none',
          duration: 2000
        });
      }
    } else {
      console.log('未找到保存的登录信息，需要手动登录');
    }
  } catch (err) {
    console.error('检查自动登录时出错:', err);
  }
};

// 处理表单提交
const handleSubmit = async () => {
  console.log('=== 点击登录按钮 ===');
  console.log('邮箱:', email.value);
  console.log('密码:', password.value);
  
  if (!email.value || !password.value) {
    console.log('表单验证失败');
    error.value = '请填写完整信息';
    return;
  }

  try {
    console.log('开始登录流程');
    loading.value = true;
    error.value = '';
    
    // 调用登录API
    console.log('准备发送登录请求');
    await api.auth.login(email.value, password.value);
    console.log('登录请求成功');

    // 获取用户信息
    try {
      console.log('开始获取用户信息');
      const response = await api.auth.getCurrentUser();
      console.log('获取用户信息成功:', response);
      
      // 从响应中获取用户信息
      const userInfo = response.user;
      console.log('用户信息:', userInfo);
      
      // 保存用户信息到本地存储，确保下次自动登录
      uni.setStorageSync('user', userInfo);
      console.log('用户信息已保存到本地存储');
      
      // 根据用户角色直接跳转到对应页面
      if (userInfo.Role === 'admin') {
        console.log('管理员用户，跳转到管理后台');
        uni.reLaunch({
          url: '/pages/admin/dashboard/index'
        });
      } else {
        console.log('普通用户，跳转到个人中心');
        uni.reLaunch({
          url: '/pages/dashboard/index'
        });
      }
    } catch (userError) {
      console.error('获取用户信息失败:', userError);
      error.value = '获取用户信息失败，请重试';
      // 清除token，因为获取用户信息失败
      uni.removeStorageSync('token');
      uni.removeStorageSync('user');
    }

  } catch (err: any) {
    console.error('登录失败:', err);
    if (err.response) {
      console.error('错误响应:', err.response.data);
      error.value = err.response.data?.message || '登录失败，请重试';
    } else if (err.request) {
      console.error('请求错误:', err.request);
      error.value = '网络错误，请检查网络连接';
    } else {
      console.error('其他错误:', err);
      error.value = '登录失败，请重试: ' + err.message;
    }
  } finally {
    loading.value = false;
  }
};

// 处理注册链接点击
const handleRegisterClick = () => {
  uni.navigateTo({
    url: '/pages/register/index'
  });
};
</script>

<style>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #000;
  background-image: linear-gradient(rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.7)),
    url('@/static/images/login-bg.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  padding: 20px;
  position: relative;
  overflow: hidden;
  transition: padding-bottom 0.3s ease;
  align-items: flex-start;
  padding-top: 15vh;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, 
    rgba(255, 0, 98, 0.1), 
    rgba(0, 223, 255, 0.1));
  animation: gradientMove 8s ease infinite;
}

@keyframes gradientMove {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.login-box {
  width: 65%;
  max-width: 240px;
  padding: 18px 12px;
  background: rgba(13, 14, 20, 0.9); /* 稍微提高背景不透明度 */
  border-radius: 8px;
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.3),
              0 0 40px rgba(255, 0, 98, 0.3);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  position: relative;
  z-index: 1;
  transition: transform 0.3s ease-out;
}

.login-box::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(45deg, #ff0062, #00dfff);
  border-radius: 9px;
  z-index: -1;
  animation: borderGlow 3s ease infinite;
  opacity: 0.5;
}

@keyframes borderGlow {
  0% { opacity: 0.5; }
  50% { opacity: 0.8; }
  100% { opacity: 0.5; }
}

.login-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 18px;
}

.logo-image {
  width: 45px;
  height: 45px;
  margin-bottom: 6px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #00dfff;
  box-shadow: 0 0 15px rgba(0, 223, 255, 0.5);
  animation: logoGlow 2s ease infinite;
}

@keyframes logoGlow {
  0% { box-shadow: 0 0 15px rgba(0, 223, 255, 0.5); }
  50% { box-shadow: 0 0 25px rgba(0, 223, 255, 0.8); }
  100% { box-shadow: 0 0 15px rgba(0, 223, 255, 0.5); }
}

.system-title {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  text-align: center;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.8);
  letter-spacing: 1px;
}

.login-form {
  margin: 10px 0;
}

.form-item {
  margin-bottom: 10px;
}

.input-wrapper {
  display: flex;
  align-items: center;
  background-color: #ffffff;
  border: none;
  border-radius: 4px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  margin-bottom: 5px;
}

.input-wrapper:focus-within {
  box-shadow: 0 0 0 2px rgba(0, 223, 255, 0.5);
}

.form-input {
  flex: 1;
  height: 34px;
  font-size: 12px;
  background: #ffffff;
  border: none;
  padding: 0 8px;
  color: #333;
  outline: none;
  width: 100%;
}

.input-placeholder {
  color: #999;
  font-size: 12px;
}

.login-button {
  width: 100%;
  height: 34px;
  background: linear-gradient(45deg, #ff0062, #00dfff);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  margin-top: 10px;
  transition: all 0.3s ease;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
  position: relative;
  overflow: hidden;
}

.login-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.2),
    transparent
  );
  transition: 0.5s;
}

.login-button:not(:disabled):active {
  transform: translateY(1px);
}

.login-button:not(:disabled):hover::before {
  left: 100%;
}

.button-disabled {
  background: #d9d9d9 !important;
  cursor: not-allowed;
  opacity: 0.7;
}

.login-footer {
  margin-top: 16px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.login-footer navigator {
  color: #00dfff;
  margin-left: 4px;
  display: inline;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
}

.link-hover {
  opacity: 0.8;
}

.error-message {
  background-color: rgba(255, 0, 98, 0.1);
  border: 1px solid rgba(255, 0, 98, 0.3);
  color: #ff4d4f;
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 15px;
  text-align: center;
  font-size: 13px;
  text-shadow: 0 0 10px rgba(255, 0, 98, 0.5);
}

.success-message {
  background-color: rgba(0, 223, 255, 0.1);
  border: 1px solid rgba(0, 223, 255, 0.3);
  color: #00dfff;
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 15px;
  text-align: center;
  font-size: 13px;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
}
</style> 