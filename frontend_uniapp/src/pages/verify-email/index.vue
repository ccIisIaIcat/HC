<template>
  <view class="verify-email-container">
    <view class="verify-email-box">
      <!-- 返回按钮 -->
      <view class="back-link" @tap="goToRegister">
        <text>← 返回注册</text>
      </view>

      <!-- 头部 -->
      <view class="header">
        <image class="logo" src="/static/logo.png" mode="aspectFit" />
        <text class="title">验证您的邮箱</text>
      </view>

      <!-- 提示信息 -->
      <view class="info-message">
        <text>验证码已发送至邮箱</text>
        <text class="email">{{ email }}</text>
      </view>

      <!-- 错误提示 -->
      <view v-if="error" class="error-message">
        <text>{{ error }}</text>
      </view>

      <!-- 验证码输入框 -->
      <view class="form">
        <view class="form-item">
          <view class="input-wrapper">
            <text class="input-icon">🔢</text>
            <input
              class="input verification-code-input"
              type="text"
              v-model="verificationCode"
              placeholder="请输入6位验证码"
              maxlength="6"
              @input="handleCodeInput"
              :disabled="loading"
            />
          </view>
        </view>

        <!-- 验证按钮 -->
        <button 
          class="submit-btn" 
          :disabled="loading || verificationCode.length !== 6" 
          @tap="handleSubmit"
        >
          {{ loading ? '验证中...' : '验证邮箱' }}
        </button>

        <!-- 重新发送验证码 -->
        <view class="resend-section">
          <text>没有收到验证码？</text>
          <button 
            class="resend-btn"
            :disabled="!canResend || loading"
            @tap="handleResend"
          >
            {{ canResend ? '重新发送' : `${countdown}秒后可重新发送` }}
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import api from '@/api';

// 获取路由参数中的邮箱
const email = ref('');
const options = uni.getLaunchOptionsSync();
// 从本地存储获取邮箱
const pendingUser = uni.getStorageSync('pendingUser');
if (pendingUser) {
  try {
    const { email: storedEmail } = JSON.parse(pendingUser);
    email.value = storedEmail;
  } catch (err) {
    console.error('解析本地存储数据失败:', err);
  }
}

// 状态变量
const verificationCode = ref('');
const loading = ref(false);
const error = ref('');
const countdown = ref(60);
const canResend = ref(false);
let timer: number | null = null;

// 倒计时逻辑
const startCountdown = () => {
  countdown.value = 60;
  canResend.value = false;
  timer = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--;
    } else {
      canResend.value = true;
      if (timer) {
        clearInterval(timer);
        timer = null;
      }
    }
  }, 1000);
};

// 验证码输入处理
const handleCodeInput = (e: any) => {
  // 只允许输入数字
  const input = e.detail.value.replace(/\D/g, '');
  verificationCode.value = input;
};

// 重新发送验证码
const handleResend = async () => {
  if (countdown.value > 0 || loading.value) return;

  try {
    loading.value = true;
    error.value = '';

    console.log('重新发送验证邮件:', email.value);
    await api.auth.resendVerification(email.value);
    console.log('重新发送成功');

    // 显示成功提示
    uni.showToast({
      title: '验证邮件已重新发送',
      icon: 'success',
      duration: 2000
    });

    // 开始倒计时
    countdown.value = 60;
    const timer = setInterval(() => {
      countdown.value--;
      if (countdown.value <= 0) {
        clearInterval(timer);
      }
    }, 1000);

  } catch (err: any) {
    console.error('重新发送失败:', err);
    if (err.response) {
      console.error('错误响应:', err.response.data);
      error.value = err.response.data?.message || '重新发送失败，请重试';
    } else if (err.request) {
      console.error('请求错误:', err.request);
      error.value = '网络错误，请检查网络连接';
    } else {
      console.error('其他错误:', err);
      error.value = '重新发送失败，请重试';
    }
  } finally {
    loading.value = false;
  }
};

// 提交验证
const handleSubmit = async () => {
  if (verificationCode.value.length !== 6 || loading.value) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    await api.auth.verifyEmail(email.value, verificationCode.value);
    
    // 验证成功，跳转到登录页
    uni.redirectTo({
      url: '/pages/login/index?message=' + encodeURIComponent('邮箱验证成功，请登录您的账号')
    });
  } catch (err: any) {
    console.error('验证失败:', err);
    if (err.response) {
      switch (err.response.status) {
        case 400:
          error.value = '验证码无效或已过期';
          break;
        case 404:
          error.value = '用户不存在';
          break;
        case 410:
          error.value = '验证码已过期，请重新发送';
          break;
        default:
          error.value = err.response.data?.message || '验证失败，请重试';
      }
    } else {
      error.value = '网络错误，请检查网络连接';
    }
  } finally {
    loading.value = false;
  }
};

// 返回注册页
const goToRegister = () => {
  uni.redirectTo({
    url: '/pages/register/index'
  });
};

// 生命周期钩子
onMounted(() => {
  if (!email.value) {
    uni.redirectTo({
      url: '/pages/register/index'
    });
    return;
  }
  startCountdown();
});

onUnmounted(() => {
  if (timer) {
    clearInterval(timer);
  }
});
</script>

<style>
.verify-email-container {
  min-height: 100vh;
  background-color: #f5f7fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.verify-email-box {
  width: 70%;
  max-width: 260px;
  background: #fff;
  border-radius: 10px;
  padding: 20px 15px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.back-link {
  margin-bottom: 20px;
  color: #1890ff;
  font-size: 14px;
}

.header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
}

.logo {
  width: 50px;
  height: 50px;
  margin-bottom: 8px;
}

.title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.info-message {
  text-align: center;
  margin-bottom: 20px;
  color: #666;
  font-size: 14px;
}

.email {
  display: block;
  margin-top: 4px;
  color: #333;
  font-weight: 500;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-item {
  display: flex;
  flex-direction: column;
}

.input-wrapper {
  display: flex;
  align-items: center;
  background: #f8f9fa;
  border-radius: 6px;
  padding: 10px;
  border: 1px solid #eee;
}

.input-icon {
  font-size: 16px;
  margin-right: 8px;
}

.verification-code-input {
  flex: 1;
  height: 36px;
  font-size: 13px;
  color: #333;
  letter-spacing: 4px;
}

.error-message {
  background: #fff2f0;
  border: 1px solid #ffccc7;
  border-radius: 6px;
  padding: 10px;
  margin-bottom: 12px;
}

.error-message text {
  color: #ff4d4f;
  font-size: 13px;
}

.submit-btn {
  height: 36px;
  background: #1890ff;
  color: #fff;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  margin-top: 12px;
}

.submit-btn[disabled] {
  opacity: 0.7;
  background: #1890ff;
}

.resend-section {
  margin-top: 16px;
  text-align: center;
  font-size: 12px;
  color: #666;
}

.resend-btn {
  display: inline;
  background: none;
  border: none;
  color: #1890ff;
  font-size: 12px;
  padding: 0;
  margin-left: 4px;
}

.resend-btn[disabled] {
  color: #999;
}
</style> 