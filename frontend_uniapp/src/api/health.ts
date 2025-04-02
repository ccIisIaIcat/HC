import request from '../utils/axios';

interface HealthAnalysisParams {
  start_date: string;
  end_date: string;
  analysis_type: 'comprehensive' | 'nutrition' | 'calories' | 'macros';
  description?: string;
}

interface HealthAnalysisResult {
  analysis: string;
  date: string;
}

interface HealthCheckResponse {
  status: string;
  message: string;
}

export default {
  /**
   * 健康检查接口
   */
  check() {
    return request({
      method: 'GET',
      url: '/api/health'
    });
  },

  /**
   * 进行健康分析
   * @param params 健康分析参数
   */
  analyzeHealth(params: HealthAnalysisParams) {
    return request({
      method: 'POST',
      url: '/api/health-analysis',
      data: params
    });
  },
  
  /**
   * 获取用户的健康分析历史
   */
  getHealthAnalysisHistory() {
    return request({
      method: 'GET',
      url: '/api/health-analysis/history'
    });
  }
}; 