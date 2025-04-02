import axiosInstance from './axios';

/**
 * 获取所有用户列表
 * @returns {Promise} 返回用户列表
 */
export const getUsers = async () => {
  try {
    console.log('发送请求: GET /admin/users');
    const response = await axiosInstance.get('/admin/users');
    console.log('API响应 (getUsers):', response.data);
    
    // 确保响应数据格式正确
    if (!response.data.users) {
      // 如果响应中没有users字段，但有其他数据，则构建预期格式
      if (Array.isArray(response.data)) {
        return { users: response.data };
      } else {
        return { users: [] }; // 返回空数组作为默认值
      }
    }
    
    return response.data;
  } catch (error) {
    console.error('获取用户列表失败:', error);
    throw error;
  }
};

/**
 * 获取系统统计信息
 * @returns {Promise} 返回系统统计信息
 */
export const getStats = async () => {
  try {
    console.log('发送请求: GET /admin/stats');
    const response = await axiosInstance.get('/admin/stats');
    console.log('API响应 (getStats):', response.data);
    
    // 确保响应数据格式正确
    if (!response.data.stats) {
      // 如果响应中没有stats字段，但有其他数据，则构建预期格式
      return { stats: response.data };
    }
    
    return response.data;
  } catch (error) {
    console.error('获取统计信息失败:', error);
    throw error;
  }
};

export default {
  getUsers,
  getStats
}; 