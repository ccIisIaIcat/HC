import axios from './axios';

/**
 * 获取所有用户列表（管理员接口）
 * @returns {Promise}
 */
export const getAllUsers = () => {
  return axios.get('/admin/users');
};

/**
 * 获取系统统计信息（管理员接口）
 * @returns {Promise}
 */
export const getSystemStats = () => {
  return axios.get('/admin/stats');
};

export default {
  getAllUsers,
  getSystemStats
}; 