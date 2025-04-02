/**
 * 请求工具封装
 * 
 * 注意：为了保证兼容性，避免使用以下API：
 * 1. URLSearchParams - 某些低版本设备不支持
 * 2. Promise.finally - 有些老旧设备不支持
 * 3. 可选链操作符(?.) - 低版本设备可能不支持
 * 4. 空值合并操作符(??) - 低版本设备可能不支持
 */

import { ref } from 'vue';
import { apiConfig } from '@/config';

// 基础配置
const baseURL = apiConfig.baseURL;
const timeout = apiConfig.timeout;

// 请求计数器
const pendingRequests = ref(0);

// 显示加载状态
const showLoading = () => {
  if (pendingRequests.value === 0) {
    uni.showLoading({
      title: '加载中...',
      mask: true
    });
  }
  pendingRequests.value++;
};

// 隐藏加载状态
const hideLoading = () => {
  pendingRequests.value--;
  if (pendingRequests.value === 0) {
    uni.hideLoading();
  }
};

// 处理错误
const handleError = (error: any) => {
  // 隐藏加载状态
  hideLoading();

  // 处理错误信息
  const message = error.errMsg || '请求失败';
  uni.showToast({
    title: message,
    icon: 'none',
    duration: 2000
  });

  return Promise.reject(error);
};

// 创建请求方法
const request = <T = any>(options: UniApp.RequestOptions): Promise<T> => {
  // 显示加载状态
  showLoading();

  // 合并选项
  const requestOptions: UniApp.RequestOptions = {
    ...options,
    url: `${baseURL}${options.url}`,
    timeout,
    header: {
      'Content-Type': 'application/json',
      ...options.header
    }
  };

  // 添加 token
  const token = uni.getStorageSync('token');
  if (token) {
    requestOptions.header = {
      ...requestOptions.header,
      Authorization: `Bearer ${token}`
    };
  }

  return new Promise((resolve, reject) => {
    uni.request({
      ...requestOptions,
      success: (response) => {
        hideLoading();
        
        // 处理响应
        const { statusCode, data } = response;
        
        if (statusCode >= 200 && statusCode < 300) {
          resolve(data as T);
        } else {
          // 处理业务错误
          const errorMessage = typeof data === 'object' && data !== null ? (data as any).message || '请求失败' : '请求失败';
          const error = new Error(errorMessage);
          reject(error);
        }
      },
      fail: (error) => {
        handleError(error);
        reject(error);
      }
    });
  });
};

// 创建 HTTP 方法
const http = {
  get<T = any>(url: string, options?: Omit<UniApp.RequestOptions, 'url' | 'method'>) {
    return request<T>({
      url,
      method: 'GET',
      ...options
    });
  },

  post<T = any>(url: string, data?: any, options?: Omit<UniApp.RequestOptions, 'url' | 'method' | 'data'>) {
    return request<T>({
      url,
      method: 'POST',
      data,
      ...options
    });
  },

  put<T = any>(url: string, data?: any, options?: Omit<UniApp.RequestOptions, 'url' | 'method' | 'data'>) {
    return request<T>({
      url,
      method: 'PUT',
      data,
      ...options
    });
  },

  delete<T = any>(url: string, options?: Omit<UniApp.RequestOptions, 'url' | 'method'>) {
    return request<T>({
      url,
      method: 'DELETE',
      ...options
    });
  }
};

export default http; 