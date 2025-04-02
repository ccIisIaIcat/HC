import { apiConfig } from '@/config';

/**
 * 上传文件到服务器
 * @param options 上传选项
 * @returns Promise<T> 上传响应
 */
export const uploadFile = <T = any>(options: {
  url: string;              // API 路径
  filePath: string;         // 本地文件路径
  name: string;             // 文件字段名称
  formData?: Record<string, any>; // 附加表单数据
  showLoading?: boolean;    // 是否显示加载提示
}): Promise<T> => {
  return new Promise((resolve, reject) => {
    // 从配置中获取基础URL
    const apiBaseURL = apiConfig.baseURL;
    const token = uni.getStorageSync('token');
    
    // 构建完整URL - 确保能处理相对路径
    let fullUrl = options.url;
    if (!fullUrl.startsWith('http') && apiBaseURL) {
      fullUrl = `${apiBaseURL}${fullUrl.startsWith('/') ? '' : '/'}${fullUrl}`;
    }
    
    console.log('准备上传文件到:', fullUrl);
    
    // 显示加载提示
    if (options.showLoading !== false) {
      uni.showLoading({
        title: '上传中...',
        mask: true
      });
    }
    
    uni.uploadFile({
      url: fullUrl,
      filePath: options.filePath,
      name: options.name,
      formData: options.formData || {},
      header: token ? { 'Authorization': `Bearer ${token}` } : {},
      success: (res) => {
        console.log('上传响应状态码:', res.statusCode);
        
        // 隐藏加载提示
        if (options.showLoading !== false) {
          uni.hideLoading();
        }
        
        if (res.statusCode >= 200 && res.statusCode < 300) {
          try {
            const parsedData = JSON.parse(res.data);
            console.log('解析后的响应数据:', parsedData);
            resolve(parsedData as T);
          } catch (e) {
            console.error('解析响应数据失败:', e, res.data);
            reject(new Error('解析响应数据失败'));
          }
        } else {
          console.error('上传失败，状态码:', res.statusCode);
          console.error('响应内容:', res.data);
          
          let errorMessage = `上传失败，状态码: ${res.statusCode}`;
          try {
            const errorData = JSON.parse(res.data);
            if (errorData.message) {
              errorMessage = errorData.message;
            }
          } catch (e) {
            // 解析失败时使用默认错误消息
          }
          
          reject(new Error(errorMessage));
        }
      },
      fail: (err) => {
        console.error('上传请求失败:', err);
        
        // 隐藏加载提示
        if (options.showLoading !== false) {
          uni.hideLoading();
        }
        
        reject(err);
      }
    });
  });
};

export default {
  uploadFile
}; 