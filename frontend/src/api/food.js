import request from '@/utils/request';

/**
 * 上传并分析食物图片
 * @param {FormData} formData 包含图片文件的FormData对象
 * @returns {Promise} 分析结果
 */
export function analyzeFoodImage(formData) {
  return request.post('/api/analyze-food', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}

/**
 * 创建食物记录
 * @param {Object} recordData 食物记录数据
 * @returns {Promise} 创建结果
 */
export function createFoodRecord(data) {
  return request.post('/api/food-records', data);
}

/**
 * 获取食物记录列表
 * @param {String} startDate 开始日期 (可选，格式：YYYY-MM-DD)
 * @param {String} endDate 结束日期 (可选，格式：YYYY-MM-DD)
 * @returns {Promise} 食物记录列表
 */
export function getFoodRecords(params) {
  return request.get('/api/food-records', { params });
}

/**
 * 获取单个食物记录详情
 * @param {Number} id 食物记录ID
 * @returns {Promise} 食物记录详情
 */
export function getFoodRecordDetail(recordId) {
  return request.get(`/api/food-records/${recordId}`);
}

/**
 * 更新食物记录
 * @param {Number} id 食物记录ID
 * @param {Object} recordData 更新的食物记录数据
 * @returns {Promise} 更新结果
 */
export function updateFoodRecord(recordId, data) {
  return request.put(`/api/food-records/${recordId}`, data);
}

/**
 * 删除食物记录
 * @param {Number} id 食物记录ID
 * @returns {Promise} 删除结果
 */
export function deleteFoodRecord(recordId) {
  return request.delete(`/api/food-records/${recordId}`);
}

/**
 * 分析食物图片并直接保存为食物记录
 * @param {FormData} formData 包含图片和记录信息的FormData对象
 * @returns {Promise} 分析和保存结果
 */
export function analyzeAndSaveFoodRecord(formData) {
  return request.post('/api/analyze-and-save', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}

// 处理 API 响应数据
export function handleResponse(response) {
  if (!response || !response.data) {
    throw new Error('无效的响应数据');
  }

  // 如果响应中包含 record 字段，直接返回 record 数据
  if (response.data.record) {
    return response.data.record;
  }

  // 否则返回整个响应数据
  return response.data;
}

// 处理 API 错误
export function handleError(error) {
  if (error.response) {
    const status = error.response.status;
    const data = error.response.data;

    switch (status) {
      case 401:
        throw new Error('未授权访问，请重新登录');
      case 403:
        throw new Error('无权访问此记录');
      case 404:
        throw new Error('记录不存在');
      case 400:
        throw new Error(data.error || '请求参数错误');
      case 500:
        throw new Error('服务器内部错误');
      default:
        throw new Error(data.error || '请求失败');
    }
  }
  throw error;
}

export default {
  analyzeFoodImage,
  createFoodRecord,
  getFoodRecords,
  getFoodRecordDetail,
  updateFoodRecord,
  deleteFoodRecord,
  analyzeAndSaveFoodRecord,
  handleResponse,
  handleError
};