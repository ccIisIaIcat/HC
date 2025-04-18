import request from '@/utils/request';

interface CheckInData {
  id?: number;
  user_id: number;
  check_in_at: string;
  content: string;
  image_url?: string;
}

interface CheckInResponse {
  success: boolean;
  message: string;
  has_food_record: boolean;
  already_checked_in: boolean;
  check_in_data?: CheckInData;
}

interface TodayCheckInResponse {
  success: boolean;
  data: {
    has_checked_in: boolean;
    has_food_record: boolean;
    check_in_days: number;  // 累计打卡天数
  };
}

interface CheckInsListResponse {
  success: boolean;
  message: string;
  data: {
    total: number;
    records: CheckInData[];
    page: string;
    page_size: string;
  };
}

export default {
  /**
   * 获取今日签到状态
   */
  getTodayStatus(): Promise<TodayCheckInResponse> {
    return request.get('/api/check-in/today');
  },

  /**
   * 执行签到
   * @param content 签到内容
   */
  checkIn(content: string): Promise<CheckInResponse> {
    const requestData = {
      content: content || ' '
    };
    return request.post('/api/check-in', requestData);
  },

  /**
   * 获取签到记录列表
   * @param params 分页参数
   */
  getCheckInList(params?: { page?: number; page_size?: number }): Promise<CheckInsListResponse> {
    let url = '/api/check-ins';
    if (params) {
      const queryParts: string[] = [];
      if (params.page !== undefined) {
        queryParts.push(`page=${params.page}`);
      }
      if (params.page_size !== undefined) {
        queryParts.push(`page_size=${params.page_size}`);
      }
      if (queryParts.length > 0) {
        url += `?${queryParts.join('&')}`;
      }
    }
    return request.get(url);
  }
}; 