const axios = require('axios');

// 设置 axios 默认配置
axios.defaults.baseURL = 'http://localhost:8080';

// 测试注册功能
const testRegister = async () => {
  try {
    console.log('开始测试注册...');
    const testData = {
      name: '测试用户',
      email: 'test@example.com',
      password: 'password123'
    };
    
    console.log('发送注册请求，数据:', testData);
    const response = await axios.post('/api/register', testData);
    console.log('注册成功，返回结果:', response.data);
  } catch (error) {
    console.error('注册测试失败:', error);
    if (error.response) {
      console.error('错误响应数据:', error.response.data);
      console.error('错误状态码:', error.response.status);
    }
  }
};

// 执行测试
testRegister(); 