import axiosInstance from './axios';

/**
 * 用户注册
 * @param {string} username 用户名
 * @param {string} email 用户邮箱
 * @param {string} password 用户密码
 * @returns {Promise}
 */
export const register = async (username, email, password) => {
  try {
    const response = await axiosInstance.post('/register', { 
      name: username, // 将username作为name字段发送到后端
      email, 
      password 
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

/**
 * 用户登录
 * @param {string} email 用户邮箱
 * @param {string} password 用户密码
 * @returns {Promise}
 */
export const login = async (email, password) => {
  try {
    const response = await axiosInstance.post('/login', { email, password });
    return response.data;
  } catch (error) {
    throw error;
  }
};

/**
 * 验证邮箱
 * @param {string} email 用户邮箱
 * @param {string} code 验证码
 * @returns {Promise}
 */
export const verifyEmail = async (email, code) => {
  try {
    const response = await axiosInstance.post('/verify-email', { email, code });
    return response.data;
  } catch (error) {
    throw error;
  }
};

/**
 * 重新发送验证邮件
 * @param {string} email 用户邮箱
 * @returns {Promise}
 */
export const resendVerificationCode = async (email) => {
  try {
    const response = await axiosInstance.post('/resend-verification', { email });
    return response.data;
  } catch (error) {
    throw error;
  }
};

/**
 * 获取当前登录用户信息
 * @returns {Promise}
 */
export const getCurrentUser = async () => {
  try {
    const response = await axiosInstance.get('/me');
    return response.data;
  } catch (error) {
    throw error;
  }
};

/**
 * 退出登录
 * @returns {Promise}
 */
export const logout = async () => {
  try {
    // 清除本地存储的token和用户信息
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    return { success: true };
  } catch (error) {
    throw error;
  }
};

export default {
  register,
  login,
  verifyEmail,
  resendVerificationCode,
  getCurrentUser,
  logout
}; 