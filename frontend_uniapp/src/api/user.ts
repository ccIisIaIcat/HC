import request from '../utils/axios';

interface UserProfile {
  id?: number;
  name: string;
  email: string;
  avatar?: string;
  role?: string;
  preferences?: {
    theme?: string;
    notifications?: boolean;
    language?: string;
  };
}

interface UserResponse {
  user: UserProfile;
}

interface PreferencesResponse {
  preferences: UserProfile['preferences'];
}

interface AvatarResponse {
  avatarUrl: string;
}

/**
 * 获取当前用户信息
 * @returns {Promise<UserProfile>} 用户信息
 */
export const getCurrentUser = async (): Promise<UserResponse> => {
  try {
    const response = await request({
      method: 'GET',
      url: '/api/me'
    }) as UserResponse;
    return response;
  } catch (error) {
    console.error('获取用户信息失败:', error);
    throw error;
  }
};

/**
 * 更新用户信息
 * @param {Partial<UserProfile>} data 要更新的用户数据
 * @returns {Promise<UserProfile>} 更新后的用户信息
 */
export const updateProfile = async (data: Partial<UserProfile>): Promise<UserResponse> => {
  try {
    const response = await request({
      method: 'PUT',
      url: '/me',
      data
    }) as UserResponse;
    return response;
  } catch (error) {
    console.error('更新用户信息失败:', error);
    throw error;
  }
};

/**
 * 更新用户密码
 * @param {string} oldPassword 旧密码
 * @param {string} newPassword 新密码
 * @returns {Promise<void>}
 */
export const updatePassword = async (oldPassword: string, newPassword: string): Promise<{ message: string }> => {
  try {
    const response = await request({
      method: 'PUT',
      url: '/me/password',
      data: {
        oldPassword,
        newPassword
      }
    }) as { message: string };
    return response;
  } catch (error) {
    console.error('更新密码失败:', error);
    throw error;
  }
};

/**
 * 更新用户头像
 * @param {FormData} formData 包含头像文件的FormData对象
 * @returns {Promise<{avatarUrl: string}>} 新的头像URL
 */
export const updateAvatar = async (formData: FormData): Promise<AvatarResponse> => {
  try {
    const response = await request({
      method: 'PUT',
      url: '/me/avatar',
      data: formData,
      header: {
        'Content-Type': 'multipart/form-data'
      }
    }) as AvatarResponse;
    return response;
  } catch (error) {
    console.error('更新头像失败:', error);
    throw error;
  }
};

/**
 * 获取用户偏好设置
 * @returns {Promise<Object>} 用户偏好设置
 */
export const getPreferences = async (): Promise<PreferencesResponse> => {
  try {
    const response = await request({
      method: 'GET',
      url: '/me/preferences'
    }) as PreferencesResponse;
    return response;
  } catch (error) {
    console.error('获取用户偏好设置失败:', error);
    throw error;
  }
};

/**
 * 更新用户偏好设置
 * @param {Object} preferences 新的偏好设置
 * @returns {Promise<Object>} 更新后的偏好设置
 */
export const updatePreferences = async (preferences: UserProfile['preferences']): Promise<PreferencesResponse> => {
  try {
    const response = await request({
      method: 'PUT',
      url: '/me/preferences',
      data: preferences
    }) as PreferencesResponse;
    return response;
  } catch (error) {
    console.error('更新用户偏好设置失败:', error);
    throw error;
  }
};

export default {
  getCurrentUser,
  updateProfile,
  updatePassword,
  updateAvatar,
  getPreferences,
  updatePreferences
}; 