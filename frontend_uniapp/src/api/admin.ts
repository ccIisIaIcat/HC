import request from '../utils/axios';

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  createdAt: string;
  updatedAt: string;
}

interface Stats {
  totalUsers: number;
  activeUsers: number;
  totalRecords: number;
  dailyActiveUsers: number;
}

interface UsersResponse {
  users: User[];
}

interface StatsResponse {
  stats: Stats;
}

/**
 * 获取所有用户列表
 * @returns {Promise<UsersResponse>} 返回用户列表
 */
export const getUsers = async (): Promise<UsersResponse> => {
  try {
    console.log('发送请求: GET /admin/users');
    const response = await request({
      method: 'GET',
      url: '/admin/users'
    }) as UsersResponse;
    console.log('API响应 (getUsers):', response);
    
    // 确保响应数据格式正确
    if (!response.users) {
      // 如果响应中没有users字段，但有其他数据，则构建预期格式
      if (Array.isArray(response)) {
        return { users: response };
      } else {
        return { users: [] }; // 返回空数组作为默认值
      }
    }
    
    return response;
  } catch (error) {
    console.error('获取用户列表失败:', error);
    throw error;
  }
};

/**
 * 获取系统统计信息
 * @returns {Promise<StatsResponse>} 返回系统统计信息
 */
export const getStats = async (): Promise<StatsResponse> => {
  try {
    console.log('发送请求: GET /admin/stats');
    const response = await request({
      method: 'GET',
      url: '/admin/stats'
    }) as StatsResponse;
    console.log('API响应 (getStats):', response);
    
    // 确保响应数据格式正确
    if (!response.stats) {
      // 如果响应中没有stats字段，但有其他数据，则构建预期格式
      return { stats: response as unknown as Stats };
    }
    
    return response;
  } catch (error) {
    console.error('获取统计信息失败:', error);
    throw error;
  }
};

export default {
  getUsers,
  getStats
}; 