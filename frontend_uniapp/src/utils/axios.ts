import { apiConfig } from '@/config';

// 创建请求函数
const request = (config: any) => {
  return new Promise((resolve, reject) => {
    // 从storage获取token
    const token = uni.getStorageSync('token');
    
    // 构建请求头
    const headers = {
      'Content-Type': 'application/json',
      ...config.headers
    };
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    // 构建完整URL
    let url = config.url || '';
    // 检查 URL 是否已经是完整的 HTTP/HTTPS URL
    if (!url.startsWith('http') && apiConfig.baseURL) {
      url = `${apiConfig.baseURL}${url}`;
    }
    
    // 打印请求详情
    console.log('发送请求:', {
      baseURL: apiConfig.baseURL,
      originalUrl: config.url,
      finalUrl: url,
      method: config.method,
      headers,
      data: config.data,
      params: config.params
    });

    // 如果有查询参数，将其添加到URL中
    if (config.params) {
      const queryString = Object.entries(config.params)
        .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
        .join('&');
      if (queryString) {
        url += (url.includes('?') ? '&' : '?') + queryString;
      }
    }

    // 检查网络状态
    uni.getNetworkType({
      success: (res) => {
        if (res.networkType === 'none') {
          reject(new Error('网络连接已断开，请检查网络设置'));
          return;
        }
        
        // 发送请求
        uni.request({
          url,
          method: config.method,
          data: config.data,
          header: headers,
          timeout: 60000, // 设置更长的超时时间：60秒
          success: (res: any) => {
            // 处理401未授权错误
            if (res.statusCode === 401) {
              uni.removeStorageSync('token');
              uni.removeStorageSync('user');
              uni.reLaunch({
                url: '/pages/login/index'
              });
              reject(new Error('未授权'));
              return;
            }
            
            // 处理403禁止访问错误
            if (res.statusCode === 403) {
              uni.removeStorageSync('token');
              uni.removeStorageSync('user');
              uni.showToast({
                title: '登录信息已失效，请重新登录',
                icon: 'none',
                duration: 2000
              });
              setTimeout(() => {
                uni.reLaunch({
                  url: '/pages/login/index'
                });
              }, 1500);
              reject(new Error('登录信息已失效，请重新登录'));
              return;
            }
            
            // 处理其他状态码
            if (res.statusCode >= 200 && res.statusCode < 300) {
              resolve(res.data);
            } else {
              reject(new Error(res.data?.message || `请求失败 (${res.statusCode})`));
            }
          },
          fail: (err) => {
            console.error('请求失败:', err);
            
            // 检查网络状态
            uni.getNetworkType({
              success: (networkRes) => {
                if (networkRes.networkType === 'none') {
                  reject(new Error('网络连接已断开，请检查网络设置'));
                  return;
                }
                
                if (err.errMsg.includes('abort')) {
                  // 请求被取消时，检查是否是网络问题
                  if (networkRes.networkType === 'wifi' || networkRes.networkType === '4g') {
                    reject(new Error('请求被中断，可能是页面切换或应用状态变化导致'));
                  } else {
                    reject(new Error('网络连接不稳定，请检查网络设置'));
                  }
                } else if (err.errMsg.includes('timeout')) {
                  reject(new Error('请求超时，请检查网络连接或稍后重试'));
                } else {
                  reject(new Error(err.errMsg || '网络请求失败'));
                }
              },
              fail: () => {
                reject(new Error('无法获取网络状态，请检查网络设置'));
              }
            });
          }
        });
      },
      fail: () => {
        reject(new Error('无法获取网络状态，请检查网络设置'));
      }
    });
  });
};

export default request; 