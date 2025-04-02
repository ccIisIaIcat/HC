<template>
  <view class="register-container">
    <view class="register-box">
      <!-- 头部 -->
      <view class="register-header">
        <image class="logo-image" src="/static/logo.png" mode="aspectFit" />
        <text class="system-title">健康检查系统</text>
        <text class="subtitle">创建新账户</text>
      </view>

      <!-- 错误提示 -->
      <view v-if="error" class="error-message">
        <text>{{ error }}</text>
      </view>

      <!-- 注册表单 -->
      <view class="register-form">
        <!-- 用户名 -->
        <view class="form-item">
          <view class="input-wrapper">
            <input
              class="form-input"
              type="text"
              v-model="formData.name"
              placeholder="请输入用户名"
              placeholder-class="input-placeholder"
              :disabled="loading"
            />
          </view>
        </view>

        <!-- 邮箱 -->
        <view class="form-item">
          <view class="input-wrapper">
            <input
              class="form-input"
              type="text"
              v-model="formData.email"
              placeholder="请输入邮箱"
              placeholder-class="input-placeholder"
              :disabled="loading"
            />
          </view>
        </view>

        <!-- 密码 -->
        <view class="form-item">
          <view class="input-wrapper">
            <input
              class="form-input"
              :type="showPassword ? 'text' : 'password'"
              v-model="formData.password"
              placeholder="请输入密码"
              placeholder-class="input-placeholder"
              :disabled="loading"
            />
            <text 
              class="password-toggle" 
              @tap="showPassword = !showPassword"
            >
              {{ showPassword ? '👁️' : '👁️‍🗨️' }}
            </text>
          </view>
        </view>

        <!-- 确认密码 -->
        <view class="form-item">
          <view class="input-wrapper">
            <input
              class="form-input"
              :type="showConfirmPassword ? 'text' : 'password'"
              v-model="formData.confirmPassword"
              placeholder="请确认密码"
              placeholder-class="input-placeholder"
              :disabled="loading"
            />
            <text 
              class="password-toggle" 
              @tap="showConfirmPassword = !showConfirmPassword"
            >
              {{ showConfirmPassword ? '👁️' : '👁️‍🗨️' }}
            </text>
          </view>
        </view>

        <!-- 注册按钮 -->
        <button 
          class="register-button" 
          :disabled="loading || !isFormValid" 
          :class="{ 'button-disabled': loading || !isFormValid }"
          @tap="handleSubmit"
        >
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </view>

      <!-- 登录链接 -->
      <view class="register-footer">
        <text>已有账号？</text>
        <text class="footer-link" @tap="goToLogin">立即登录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import api from '@/api';

// 表单数据
const formData = ref({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
});

// 状态变量
const loading = ref(false);
const error = ref('');
const showPassword = ref(false);
const showConfirmPassword = ref(false);

// 表单验证
const isFormValid = computed(() => {
  const { name, email, password, confirmPassword } = formData.value;
  return (
    name.length >= 2 &&
    email.includes('@') &&
    password.length >= 6 &&
    password === confirmPassword
  );
});

// 处理注册
const handleSubmit = async () => {
  if (!isFormValid.value || loading.value) return;

  try {
    loading.value = true;
    error.value = '';

    const { name, email, password } = formData.value;
    console.log('开始注册，参数:', { name, email, password });
    
    // 发送注册请求
    await api.auth.register({ name, email, password });
    console.log('注册请求成功，等待后端发送验证邮件');

    // 显示成功提示
    uni.showToast({
      title: '验证邮件已发送，请查收',
      icon: 'success',
      duration: 2000
    });

    // 延迟跳转，让用户看到提示
    setTimeout(() => {
      console.log('准备跳转到验证邮箱页面');
      // 存储邮箱到本地存储
      uni.setStorageSync('pendingUser', JSON.stringify({ email: email }));
      // 跳转到验证页面
      uni.navigateTo({
        url: '/pages/verify-email/index',
        success: () => {
          console.log('跳转成功');
        },
        fail: (err) => {
          console.error('跳转失败:', err);
          error.value = '页面跳转失败，请重试';
        }
      });
    }, 2000);

  } catch (err: any) {
    console.error('注册失败:', err);
    if (err.response) {
      console.error('错误响应:', err.response.data);
      error.value = err.response.data?.message || '注册失败，请重试';
    } else if (err.request) {
      console.error('请求错误:', err.request);
      error.value = '网络错误，请检查网络连接';
    } else {
      console.error('其他错误:', err);
      error.value = '注册失败，请重试';
    }
  } finally {
    loading.value = false;
  }
};

// 跳转到登录页
const goToLogin = () => {
  uni.redirectTo({
    url: '/pages/login/index'
  });
};
</script>

<style>
.register-container {
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
}

.register-container::before {
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

.register-box {
  width: 75%;
  max-width: 280px;
  padding: 18px 15px;
  background: rgba(13, 14, 20, 0.85);
  border-radius: 8px;
  box-shadow: 0 0 20px rgba(0, 223, 255, 0.2),
              0 0 40px rgba(255, 0, 98, 0.2);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  z-index: 1;
}

.register-box::before {
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

.register-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 15px;
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

.subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4px;
  text-shadow: 0 0 10px rgba(255, 0, 98, 0.5);
}

.register-form {
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

.password-toggle {
  font-size: 16px;
  padding: 0 8px;
  color: #666;
}

.register-button {
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

.register-button::before {
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

.register-button:not(:disabled):active {
  transform: translateY(1px);
}

.register-button:not(:disabled):hover::before {
  left: 100%;
}

.button-disabled {
  background: #d9d9d9 !important;
  cursor: not-allowed;
  opacity: 0.7;
}

.register-footer {
  margin-top: 16px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.footer-link {
  color: #00dfff;
  margin-left: 4px;
  display: inline;
  text-shadow: 0 0 10px rgba(0, 223, 255, 0.5);
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
</style> 