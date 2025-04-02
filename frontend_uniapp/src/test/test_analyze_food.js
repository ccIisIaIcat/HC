const fs = require('fs');
const path = require('path');
const axios = require('axios');
const FormData = require('form-data');

// API 配置
const BASE_URL = 'http://localhost:8080/api';
const TEST_IMAGE_PATH = path.join(__dirname, '../static/images/tangyuan.png');

// 登录信息
const LOGIN_CREDENTIALS = {
  email: '965377515@qq.com',
  password: '123456'
};

// 登录函数
async function login() {
  try {
    console.log('正在登录...');
    const response = await axios.post(`${BASE_URL}/login`, LOGIN_CREDENTIALS);
    
    if (response.status === 200) {
      const token = response.data.token;
      console.log('登录成功！');
      return token;
    } else {
      throw new Error('登录失败');
    }
  } catch (error) {
    console.error('登录失败：', error.message);
    if (error.response) {
      console.error('错误状态码：', error.response.status);
      console.error('错误信息：', error.response.data);
    }
    throw error;
  }
}

// 测试函数
async function testAnalyzeFood() {
  try {
    console.log('开始测试食物分析 API...');
    
    // 先登录获取 token
    const token = await login();
    
    // 检查图片文件是否存在
    if (!fs.existsSync(TEST_IMAGE_PATH)) {
      console.error('错误：测试图片不存在！');
      console.log('请确保图片文件位于：', TEST_IMAGE_PATH);
      return;
    }
    
    // 创建 FormData
    const formData = new FormData();
    formData.append('image', fs.createReadStream(TEST_IMAGE_PATH));
    formData.append('image_description', '这是一碗汤圆，里面有芝麻馅');
    
    // 发送请求
    console.log('发送分析请求...');
    const response = await axios.post(`${BASE_URL}/analyze-food`, formData, {
      headers: {
        ...formData.getHeaders(),
        'Authorization': `Bearer ${token}`
      }
    });
    
    // 检查响应
    if (response.status === 200) {
      console.log('分析成功！');
      console.log('响应数据：', JSON.stringify(response.data, null, 2));
      
      // 打印关键信息
      const result = response.data;
      if (result.hasFood) {
        console.log('\n食物信息：');
        console.log(`- 食物类型：${result.foodType}`);
        console.log(`- 估计重量：${result.weight}克`);
        
        console.log('\n营养成分：');
        const nutrition = result.nutrition;
        console.log(`- 热量：${nutrition.calories}卡路里`);
        console.log(`- 蛋白质：${nutrition.protein}克`);
        console.log(`- 碳水化合物：${nutrition.carbohydrates}克`);
        console.log(`- 脂肪：${nutrition.totalFat}克`);
      } else {
        console.log('- 未检测到食物');
      }
    } else {
      console.error('请求失败：', response.status);
      console.error('错误信息：', response.data);
    }
    
  } catch (error) {
    console.error('测试失败：', error.message);
    if (error.response) {
      console.error('错误状态码：', error.response.status);
      console.error('错误信息：', error.response.data);
    }
  }
}

// 运行测试
console.log('========== 食物分析 API 测试开始 ==========\n');
testAnalyzeFood().then(() => {
  console.log('\n========== 测试完成 ==========');
}); 