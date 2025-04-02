import axios from 'axios';

// 创建axios实例
const instance = axios.create({
  baseURL: '/api', // 使用相对路径，依赖package.json中的proxy配置
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器 - 添加token
instance.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理错误
instance.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response) {
      // 处理服务器返回的错误
      const status = error.response.status;
      
      // 401未授权，可能token过期
      if (status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        // 可以在这里添加跳转到登录页面的逻辑
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export default instance; 