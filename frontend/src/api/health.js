import axiosInstance from './axios'; 
 
/**
 * 发送健康分析请求
 * @param {Object} data - 分析请求数据
 * @param {string} data.start_date - 分析起始日期 (YYYY-MM-DD)
 * @param {string} data.end_date - 分析结束日期 (YYYY-MM-DD)
 * @param {string} data.analysis_type - 分析类型 (comprehensive|nutrition|calories|macros)
 * @param {string} data.description - 用户描述（可选）
 * @returns {Promise<Object>} - 包含analysis字段的健康分析结果
 */
export const analyzeHealth = async (data) => {
  try { 
    const response = await axiosInstance.post('/health-analysis', data); 
    return response.data;
  } catch (error) {
    console.error('健康分析请求失败:', error);
    throw error;
  }
};
