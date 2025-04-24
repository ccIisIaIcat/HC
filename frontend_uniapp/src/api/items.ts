import request from '@/utils/request';
import { apiConfig } from '@/config'

// 物品类型定义
export interface Item {
  id: number;
  ID?: number; // 后端可能返回大写的ID
  name: string;
  Name?: string; // 后端可能返回大写的Name
  description: string;
  Description?: string;
  source: string;
  Source?: string;
  icon_url: string;
  IconURL?: string;
  image_url: string;
  ImageURL?: string;
  total_quantity: number;
  obtained_from: string;
  obtained_at?: string;
}

// API响应类型
interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 获取带认证的请求头
const getHeaders = () => {
  const token = uni.getStorageSync('token');
  if (!token) {
    throw new Error('用户未登录');
  }
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  };
};

// 处理图片URL
function normalizeImageUrl(url: string | null | undefined): string {
  if (!url) return '/static/collections/cat-locked.png'
  if (url.startsWith('http')) return url
  return `${apiConfig.baseURL}${url}`
}

// 标准化物品数据
function normalizeItem(item: any): Item {
  return {
    id: item.ID || item.id,
    name: item.Name || item.name,
    description: item.Description || item.description,
    source: item.Source || item.source,
    icon_url: normalizeImageUrl(item.IconURL || item.icon_url),
    image_url: normalizeImageUrl(item.ImageURL || item.image_url),
    total_quantity: item.TotalQuantity || item.total_quantity || 0,
    obtained_from: item.ObtainedFrom || item.obtained_from,
    obtained_at: item.ObtainedAt || item.obtained_at
  }
}

// 获取用户的所有物品
export const getUserItems = async () => {
  try {
    // 获取token
    const token = uni.getStorageSync('token');
    if (!token) {
      throw new Error('用户未登录');
    }
    
    // 获取用户ID
    let userId;
    const userInfo = uni.getStorageSync('user');
    
    if (userInfo) {
      if (typeof userInfo === 'object') {
        userId = userInfo.ID || userInfo.id;
      } else if (typeof userInfo === 'string') {
        try {
          const parsedUser = JSON.parse(userInfo);
          userId = parsedUser.ID || parsedUser.id;
        } catch (e) {
          console.error('无法解析用户信息:', e);
        }
      }
    }
    
    // 如果上面的方式未能获取到userId，尝试从userId存储中获取
    if (!userId) {
      userId = uni.getStorageSync('userId');
    }
    
    // 如果仍然没有userId，报错
    if (!userId) {
      throw new Error('用户未登录');
    }
    
    // 执行请求
    const response = await request.get<ApiResponse<any[]>>(`/api/user-items/${userId}`, {
      header: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    // 检查返回数据格式
    if (response.data === null) {
      response.data = [];
    }
    
    // 确保数据存在并规范化每个物品的数据
    if (response.data && Array.isArray(response.data) && response.data.length > 0) {
      response.data = response.data.map(normalizeItem);
    }
    
    // 添加code和message字段，如果API没有返回
    if (response.code === undefined) {
      response.code = 200; // 默认成功
    }
    
    if (response.message === undefined) {
      response.message = "获取成功";
    }
    
    return response;
  } catch (error: any) {
    console.error('获取物品失败:', error);
    throw error; // 让调用方处理错误
  }
};

// 获取所有物品信息
export const getAllItems = () => {
  return request.get<ApiResponse<Item[]>>('/api/items', {
    header: getHeaders()
  });
};

// 获取单个物品信息
export const getItemInfo = (itemId: number) => {
  return request.get<ApiResponse<Item>>(`/api/items/${itemId}`, {
    header: getHeaders()
  });
};

// 检查用户是否拥有某个物品
export const checkUserItem = (itemId: number) => {
  try {
    // 获取token
    const token = uni.getStorageSync('token');
    
    // 直接从user对象中获取ID
    let userId;
    const userInfo = uni.getStorageSync('user');
    
    if (userInfo) {
      if (typeof userInfo === 'string') {
        try {
          const parsedUser = JSON.parse(userInfo);
          userId = parsedUser.ID || parsedUser.id;
        } catch (e) {
          console.error('无法解析用户信息:', e);
        }
      } else {
        // 已经是对象，直接使用
        userId = userInfo.ID || userInfo.id;
      }
    }
    
    // 如果上面的方式未能获取到userId，尝试从userId存储中获取
    if (!userId) {
      userId = uni.getStorageSync('userId');
    }
    
    // 如果仍然没有userId，报错
    if (!userId) {
      console.log('未找到用户ID，需要登录');
      throw new Error('用户未登录');
    }
    
    // 执行请求
    console.log(`准备检查用户物品，用户ID: ${userId}, 物品ID: ${itemId}`);
    return request.get<ApiResponse<{has_item: boolean}>>(`/api/user-items/${userId}/check/${itemId}`, {
      header: getHeaders()
    });
  } catch (error: any) {
    console.error('检查物品错误:', error);
    throw error;
  }
}; 