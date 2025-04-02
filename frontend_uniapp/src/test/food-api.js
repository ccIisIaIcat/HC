const axios = require('axios');

// 设置 axios 默认配置
axios.defaults.baseURL = 'http://localhost:8080';

// 模拟数据
const mockFoodRecord = {
  food_name: '测试食物',
  weight: 200,
  record_time: '2024-03-20T10:00:00Z',
  meal_type: '早餐',
  notes: '测试记录',
  calories: 300,
  protein: 20,
  total_fat: 10,
  saturated_fat: 3,
  trans_fat: 0.5,
  unsaturated_fat: 6.5,
  carbohydrates: 40,
  sugar: 5,
  fiber: 3,
  vitamin_a: 100,
  vitamin_c: 30,
  vitamin_d: 2,
  vitamin_b1: 0.5,
  vitamin_b2: 0.6,
  calcium: 200,
  iron: 3,
  sodium: 400,
  potassium: 500
};

// 登录并获取 token
const login = async () => {
  try {
    console.log('开始登录...');
    const loginData = {
      email: '965377515@qq.com',
      password: '123456'
    };
    const response = await axios.post('/api/login', loginData);
    const token = response.data.token;
    console.log('登录成功，获取到 token');
    
    // 设置 axios 默认 headers
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    return token;
  } catch (error) {
    console.error('登录失败:', error);
    if (error.response) {
      console.error('错误响应数据:', error.response.data);
      console.error('错误状态码:', error.response.status);
    }
    throw error;
  }
};

// 测试食物记录 API
const testFoodAPI = async () => {
  try {
    // 先登录获取 token
    await login();
    
    console.log('开始测试食物记录 API...\n');

    // 1. 测试获取食物记录列表
    console.log('测试获取食物记录列表...');
    const params = {
      start_date: '2024-03-01',
      end_date: '2024-03-20',
      meal_type: '早餐'
    };
    const recordsResponse = await axios.get('/api/food-records', { params });
    console.log('获取食物记录列表成功:', recordsResponse.data);
    console.log('----------------------------------------\n');

    // 2. 测试创建食物记录
    console.log('测试创建食物记录...');
    const createResponse = await axios.post('/api/food-records', mockFoodRecord);
    console.log('创建食物记录成功:', createResponse.data);
    console.log('----------------------------------------\n');

    // 从创建响应中获取记录 ID
    const recordId = createResponse.data.record.ID || createResponse.data.record.id;
    if (!recordId) {
      throw new Error('无法获取创建的记录 ID');
    }

    // 3. 测试获取单个食物记录
    console.log('测试获取单个食物记录...');
    const detailResponse = await axios.get(`/api/food-records/${recordId}`);
    console.log('获取单个食物记录原始响应:', JSON.stringify(detailResponse.data));
    console.log('响应数据类型:', typeof detailResponse.data);
    console.log('响应是否包含record字段:', detailResponse.data.hasOwnProperty('record'));
    if (typeof detailResponse.data === 'object') {
      console.log('响应对象的所有字段:', Object.keys(detailResponse.data));
    }
    if (!detailResponse.data.record) {
      throw new Error('返回数据中没有 record 字段');
    }
    console.log('获取单个食物记录成功:', detailResponse.data.record);
    console.log('----------------------------------------\n');

    // 4. 测试更新食物记录
    console.log('测试更新食物记录...');
    // 使用获取到的记录作为基础，只更新需要修改的字段
    const existingRecord = detailResponse.data.record;
    const updateData = {
      ...existingRecord,
      notes: '更新后的测试记录',
      record_time: new Date().toISOString()
    };
    // 删除不需要的字段
    delete updateData.created_at;
    delete updateData.updated_at;
    delete updateData.deleted_at;
    
    console.log('更新数据:', updateData);
    const updateResponse = await axios.put(`/api/food-records/${recordId}`, updateData);
    if (!updateResponse.data.record) {
      throw new Error('返回数据中没有 record 字段');
    }
    console.log('更新食物记录成功:', updateResponse.data.record);
    console.log('----------------------------------------\n');

    // 5. 测试删除食物记录
    console.log('测试删除食物记录...');
    const deleteResponse = await axios.delete(`/api/food-records/${recordId}`);
    console.log('删除食物记录成功:', deleteResponse.data);
    console.log('----------------------------------------\n');

    // 6. 测试食物分析
    console.log('测试食物分析...');
    const formData = new FormData();
    const mockImage = new Blob(['mock image data'], { type: 'image/jpeg' });
    formData.append('image', mockImage, 'test.jpg');
    formData.append('image_description', '测试食物图片');
    
    const analysisResponse = await axios.post('/api/analyze-food', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    console.log('食物分析成功:', analysisResponse.data);
    console.log('----------------------------------------\n');

    console.log('所有测试完成！');
  } catch (error) {
    console.error('测试过程中发生错误:', error);
    if (error.response) {
      console.error('错误响应数据:', error.response.data);
      console.error('错误状态码:', error.response.status);
      if (error.response.status === 401) {
        console.error('未授权访问，请检查 token');
      } else if (error.response.status === 403) {
        console.error('无权访问此记录');
      } else if (error.response.status === 404) {
        console.error('记录不存在');
      }
    }
  }
};

// 执行测试
testFoodAPI(); 