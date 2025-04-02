import { apiConfig } from '@/config/index';

const BASE_URL = apiConfig.baseURL || '';  // 改为空字符串默认值，允许使用相对路径
const TIMEOUT = apiConfig.timeout || 30000; // 默认30秒超时
const RETRY_COUNT = apiConfig.retryCount || 0; // 重试次数
const RETRY_DELAY = apiConfig.retryDelay || 1000; // 重试延迟

// 获取token
function getToken() {
  return uni.getStorageSync('token');
}

// 请求拦截器
function requestInterceptor(config) {
  console.log('发起请求:', config.url, '方法:', config.method, '数据:', config.data);
  const token = getToken();
  if (token) {
    config.header = {
      ...config.header,
      'Authorization': `Bearer ${token}`
    };
  }
  return config;
}

// 响应拦截器
function responseInterceptor(response) {
  const { statusCode, data } = response;
  
  console.log(`响应状态码: ${statusCode}`, '响应数据:', JSON.stringify(data).substring(0, 500));
  
  // 请求成功
  if (statusCode >= 200 && statusCode < 300) {
    // 直接返回响应数据
    return data;
  }
  
  // 处理错误
  const error = new Error(data.message || '请求失败');
  error.response = response;
  throw error;
}

// 延迟函数
function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// 重试函数
async function retryRequest(config, retryCount, retryDelay) {
  let lastError = null;
  let currentRetry = 0;
  
  while (currentRetry <= retryCount) {
    try {
      if (currentRetry > 0) {
        console.log(`重试请求 (${currentRetry}/${retryCount}): ${config.url}`);
        await delay(retryDelay);
      }
      
      return await executeRequest(config);
    } catch (error) {
      lastError = error;
      currentRetry++;
      
      // 如果是网络错误或超时错误，则重试
      const shouldRetry = error.errMsg && (
        error.errMsg.includes('timeout') ||
        error.errMsg.includes('Failed to fetch') ||
        error.errMsg.includes('request:fail')
      );
      
      if (!shouldRetry || currentRetry > retryCount) {
        break;
      }
    }
  }
  
  throw lastError;
}

// 执行请求
function executeRequest(config) {
  console.log('准备执行请求:', JSON.stringify({
    url: config.url,
    method: config.method,
    header: config.header,
    data: config.data
  }, null, 2));
  
  return new Promise((resolve, reject) => {
    try {
      const requestTask = uni.request({
        ...config,
        success: (res) => {
          try {
            console.log('请求成功，状态码:', res.statusCode);
            const result = responseInterceptor(res);
            resolve(result);
          } catch (error) {
            console.error('请求处理失败:', error);
            reject(error);
          }
        },
        fail: (error) => {
          console.error('请求发送失败:', error, '请求配置:', config.url);
          // 增强错误信息
          if (error.errMsg && error.errMsg.includes('timeout')) {
            error.message = '请求超时，请检查网络连接';
          } else if (error.errMsg && error.errMsg.includes('abort')) {
            error.message = '请求已取消';
          } else if (error.errMsg && error.errMsg.includes('Failed to fetch')) {
            error.message = '网络连接失败，请检查API服务是否可用';
          } else {
            error.message = error.errMsg || '网络请求失败';
          }
          reject(error);
        },
        complete: () => {
          console.log('请求完成:', config.url);
        }
      });
      
      console.log('请求任务已创建');
    } catch (err) {
      console.error('创建请求任务失败:', err);
      reject(new Error('创建请求任务失败: ' + (err.message || String(err))));
    }
  });
}

// 请求函数
function request(options) {
  // 确保 options 是一个对象
  options = options || {};
  
  // 确保 options.url 有值
  if (!options.url) {
    console.error('请求URL不能为空');
    return Promise.reject(new Error('请求URL不能为空'));
  }
  
  // 处理URL
  let url = options.url;
  let originalUrl = options.url;
  let baseUrlUsed = '未使用';
  
  // 如果 URL 已经是完整的 HTTP URL，则直接使用
  if (url.startsWith('http://') || url.startsWith('https://')) {
    // URL 已经完整，不需要处理
    console.log('使用完整URL:', url);
  } 
  // 如果不是完整URL但有 BASE_URL
  else if (BASE_URL) {
    baseUrlUsed = BASE_URL;
    // 确保 baseURL 以 http:// 或 https:// 开头
    let baseURL = BASE_URL;
    if (!baseURL.startsWith('http://') && !baseURL.startsWith('https://')) {
      baseURL = 'https://' + baseURL.replace(/^[^:]+:\/\//, '');
      baseUrlUsed = baseURL;
    }
    
    // 确保 URL 路径正确连接
    if (url.startsWith('/') && baseURL.endsWith('/')) {
      url = baseURL + url.substring(1);
    } else if (!url.startsWith('/') && !baseURL.endsWith('/')) {
      url = baseURL + '/' + url;
    } else {
      url = baseURL + url;
    }
    
    console.log('构建完整URL:', url, '基础URL:', baseURL, '路径:', options.url);
  }
  // 如果没有 BASE_URL，使用相对路径（这可能在某些环境下无法工作）
  else {
    // 尝试获取 VITE_API_URL 或 API_URL
    const apiURL = uni.getStorageSync('API_URL') || 'https://jesuvukndxpo.sealoshzh.site';
    if (apiURL) {
      baseUrlUsed = apiURL;
      if (url.startsWith('/') && apiURL.endsWith('/')) {
        url = apiURL + url.substring(1);
      } else if (!url.startsWith('/') && !apiURL.endsWith('/')) {
        url = apiURL + '/' + url;
      } else {
        url = apiURL + url;
      }
      console.log('使用存储的API URL:', apiURL, '构建完整URL:', url);
    } else {
      console.warn('没有配置baseURL，使用相对路径可能导致请求失败:', url);
    }
  }
  
  // 避免重复添加已经添加的查询参数
  // GET 请求的参数已经由 get() 方法处理过
  let queryParamsInfo = '';
  if (options.method === 'GET' && options.data && !url.includes('?')) {
    const queryParams = options.data;
    let queryString = '';
    
    if (Object.keys(queryParams).length > 0) {
      try {
        queryString = Object.entries(queryParams)
          .filter(([_, value]) => value !== undefined && value !== null)
          .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
          .join('&');
        
        if (queryString) {
          url += (url.includes('?') ? '&' : '?') + queryString;
          queryParamsInfo = JSON.stringify(queryParams);
        }
      } catch (err) {
        console.error('处理查询参数失败:', err, queryParams);
      }
    }
  }
  
  // 记录详细的请求信息（仅日志，无弹窗）
  console.log('请求详细信息:', {
    完整URL: url,
    原始路径: originalUrl,
    使用的baseURL: baseUrlUsed,
    方法: options.method || 'GET',
    查询参数: queryParamsInfo || '无',
    请求体: options.data && options.method !== 'GET' ? JSON.stringify(options.data).substring(0, 200) : '无'
  });
  
  // 合并配置
  const config = {
    url, // 使用处理后的完整 URL
    method: options.method || 'GET',
    header: {
      'Content-Type': 'application/json',
      ...options.header
    },
    data: options.method === 'GET' ? undefined : options.data, // GET 请求不在 body 中发送数据
    timeout: options.timeout || TIMEOUT,
    dataType: 'json'  // 明确指定返回数据类型
  };
  
  // 请求拦截
  const interceptedConfig = requestInterceptor(config);
  
  // 执行请求（带重试）
  const retryCount = options.retryCount !== undefined ? options.retryCount : RETRY_COUNT;
  const retryDelay = options.retryDelay !== undefined ? options.retryDelay : RETRY_DELAY;
  
  return retryRequest(interceptedConfig, retryCount, retryDelay);
}

// 导出请求方法
export default {
  get(url, options = {}) {
    console.log('GET 请求开始处理:', url, options);
    
    // 确保 options 是一个对象
    const safeOptions = options || {};
    
    // 从 options 中获取 data 属性作为查询参数
    const queryParams = safeOptions.data || {};
    
    // 如果没有查询参数，直接发送请求
    if (!queryParams || Object.keys(queryParams).length === 0) {
      console.log('GET 请求没有查询参数:', url);
      return request({ ...safeOptions, url, method: 'GET' });
    }
    
    // 构建查询字符串
    let queryString = '';
    try {
      queryString = Object.entries(queryParams)
        .filter(([_, value]) => value !== undefined && value !== null)
        .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
        .join('&');
    } catch (err) {
      console.error('构建查询字符串失败:', err, queryParams);
    }
    
    // 如果有查询参数，添加到 URL
    const fullUrl = queryString ? 
      (url.includes('?') ? `${url}&${queryString}` : `${url}?${queryString}`) : 
      url;
    
    console.log('GET 请求的完整 URL:', fullUrl);
    
    // 删除 data 属性，避免重复添加查询参数
    const { data, ...restOptions } = safeOptions;
    
    // 发送最终请求
    return request({ ...restOptions, url: fullUrl, method: 'GET' });
  },
  
  post(url, data, options = {}) {
    return request({ ...options, url, method: 'POST', data });
  },
  
  put(url, data, options = {}) {
    return request({ ...options, url, method: 'PUT', data });
  },
  
  delete(url, options = {}) {
    return request({ ...options, url, method: 'DELETE' });
  }
}; 