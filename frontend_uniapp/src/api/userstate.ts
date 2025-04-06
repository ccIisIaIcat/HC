import request from '@/utils/request';

/**
 * 用户健康状态接口
 */

// 创建健康状态记录
export function createHealthState(data: any) {
  return request.post('/api/health-states', data);
}

// 获取健康状态记录列表
export function getHealthStates(params?: { start_date?: string; end_date?: string }) {
  return request.get('/api/health-states', { data: params });
}

// 获取单个健康状态记录
export function getHealthState(id: number) {
  return request.get(`/api/health-states/${id}`);
}

// 获取最新的健康状态记录
export function getLatestHealthState() {
  return request.get('/api/health-states/latest');
}

// 更新健康状态记录
export function updateHealthState(id: number, data: any) {
  return request.put(`/api/health-states/${id}`, data);
}

// 删除健康状态记录
export function deleteHealthState(id: number) {
  return request.delete(`/api/health-states/${id}`);
}

// 健康状态记录类型定义
export interface HealthState {
  id?: number;
  user_id?: number;
  record_time?: string;
  height?: number;
  weight?: number;
  bmi?: number;
  body_fat_percentage?: number;
  temperature?: number;
  heart_rate?: number;
  respiratory_rate?: number;
  fasting_glucose?: number;
  postprandial_glucose?: number;
  total_cholesterol?: number;
  notes?: string;
  created_at?: string;
  updated_at?: string;
}
