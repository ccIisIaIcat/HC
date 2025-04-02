import request from '../utils/axios';

interface User {
  id: number;
  name: string;
  email: string;
  Role: string;
  createdAt: string;
  updatedAt: string;
}

interface LoginResponse {
  token: string;
  user: User;
}

interface RegisterResponse {
  message: string;
}

interface VerifyEmailResponse {
  message: string;
}

interface GetCurrentUserResponse {
  user: User;
}

/**
 * 用户登录
 * @param email 邮箱
 * @param password 密码
 * @returns {Promise<LoginResponse>} 登录响应
 */
export const login = async (email: string, password: string): Promise<LoginResponse> => {
  try {
    const response = await request({
      method: 'POST',
      url: '/api/login',
      data: {
        email,
        password
      }
    }) as LoginResponse;
    
    // 保存token和用户信息
    if (response.token) {
      uni.setStorageSync('token', response.token);
      uni.setStorageSync('user', response.user);
    }
    
    return response;
  } catch (error) {
    console.error('登录失败:', error);
    throw error;
  }
};

/**
 * 用户注册
 * @param {string} name 用户名
 * @param {string} email 用户邮箱
 * @param {string} password 用户密码
 * @returns {Promise<RegisterResponse>} 注册结果
 */
export const register = async (data: { name: string; email: string; password: string }): Promise<RegisterResponse> => {
  try {
    console.log('发送注册请求:', data);
    const response = await request({
      method: 'POST',
      url: '/api/register',
      data: {
        name: data.name,
        email: data.email,
        password: data.password
      }
    }) as RegisterResponse;
    console.log('注册响应:', response);
    return response;
  } catch (error) {
    console.error('注册失败:', error);
    throw error;
  }
};

/**
 * 验证邮箱
 * @param {string} email 用户邮箱
 * @param {string} code 验证码
 * @returns {Promise<VerifyEmailResponse>} 验证结果
 */
export const verifyEmail = async (email: string, code: string): Promise<VerifyEmailResponse> => {
  try {
    const response = await request({
      method: 'POST',
      url: '/api/verify-email',
      data: { email, code }
    }) as VerifyEmailResponse;
    return response;
  } catch (error) {
    console.error('验证邮箱失败:', error);
    throw error;
  }
};

/**
 * 退出登录
 */
export const logout = () => {
  uni.removeStorageSync('token');
  uni.removeStorageSync('user');
  // 重定向到登录页
  uni.reLaunch({
    url: '/pages/login/index'
  });
};

/**
 * 获取当前登录用户信息
 * @returns {Promise<GetCurrentUserResponse>} 用户信息
 */
export const getCurrentUser = async (): Promise<GetCurrentUserResponse> => {
  try {
    const response = await request({
      method: 'GET',
      url: '/api/me'
    }) as GetCurrentUserResponse;
    return response;
  } catch (error) {
    console.error('获取用户信息失败:', error);
    throw error;
  }
};

/**
 * 重新发送验证邮件
 * @param {string} email 用户邮箱
 * @returns {Promise<VerifyEmailResponse>} 发送结果
 */
export const resendVerification = async (email: string): Promise<VerifyEmailResponse> => {
  try {
    const response = await request({
      method: 'POST',
      url: '/api/resend-verification',
      data: { email }
    }) as VerifyEmailResponse;
    return response;
  } catch (error) {
    console.error('重新发送验证邮件失败:', error);
    throw error;
  }
};

export default {
  login,
  register,
  logout,
  getCurrentUser,
  verifyEmail,
  resendVerification
}; 