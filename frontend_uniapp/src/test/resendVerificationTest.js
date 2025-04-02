const axios = require('axios');

// 设置 axios 默认配置
axios.defaults.baseURL = 'http://localhost:8080';

// 测试重新发送验证邮件功能
const testResendVerification = async () => {
  try {
    console.log('开始测试重新发送验证邮件...');
    const testEmail = 'zwj19980717@gmail.com';
    
    console.log('发送重新发送验证邮件请求，邮箱:', testEmail);
    const response = await axios.post('/api/resend-verification', {
      email: testEmail
    });
    console.log('重新发送验证邮件成功，返回结果:', response.data);
  } catch (error) {
    console.error('重新发送验证邮件测试失败:', error);
    if (error.response) {
      console.error('错误响应数据:', error.response.data);
      console.error('错误状态码:', error.response.status);
    }
  }
};

// 执行测试
testResendVerification(); 