import request from '@/utils/request';
import { uploadFile } from '@/utils/upload';

interface FoodRecord {
  // 基本信息
  id?: number;  // 创建时不需要提供id
  food_name: string;
  weight: number;
  record_time: string;
  meal_type: string;
  notes: string;

  // 基本营养素
  calories: number;
  protein: number;
  total_fat: number;
  saturated_fat: number;
  trans_fat: number;
  unsaturated_fat: number;
  carbohydrates: number;
  sugar: number;
  fiber: number;

  // 维生素
  vitamin_a: number;
  vitamin_c: number;
  vitamin_d: number;
  vitamin_b1: number;
  vitamin_b2: number;

  // 矿物质
  calcium: number;
  iron: number;
  sodium: number;
  potassium: number;
}

interface FoodAnalysisResult {
  foodType: string;
  weight: number;
  nutrition: {
    calories: number;
    protein: number;
    totalFat: number;
    carbohydrates: number;
    sugar: number;
    fiber: number;
    vitamins: {
      vitaminA: number;
      vitaminC: number;
      vitaminD: number;
    };
    minerals: {
      calcium: number;
      iron: number;
      sodium: number;
      potassium: number;
    };
  };
}

interface FoodRecordResponse {
  record: FoodRecord;
}

interface FoodRecordsResponse {
  records: FoodRecord[];
  total: number;
}

export default {
  /**
   * 分析食物图片
   * @param options 图片上传选项
   */
  analyzeFoodImage(options: {
    filePath: string;
    description?: string;
  }): Promise<FoodAnalysisResult> {
    const formData: Record<string, any> = {};
    if (options.description) {
      formData.image_description = options.description;
    }
    
    return uploadFile<FoodAnalysisResult>({
      url: '/api/analyze-food',
      filePath: options.filePath,
      name: 'image',
      formData: formData,
      showLoading: true
    });
  },

  /**
   * 创建食物记录
   * @param data 食物记录数据
   */
  createFoodRecord(data: FoodRecord) {
    return request.post('/api/food-records', data);
  },

  /**
   * 获取用户的食物记录列表
   * @param params 查询参数
   */
  getFoodRecords(params?: { startDate?: string; endDate?: string; mealType?: string }): Promise<FoodRecordsResponse> {
    console.log('调用 food.ts 中的 getFoodRecords', params);
    
    // 将参数转换为查询字符串 - 不使用 URLSearchParams
    let url = '/api/food-records';
    if (params) {
      const queryParts: string[] = [];
      
      // 手动构建查询参数
      if (params.startDate !== undefined) {
        queryParts.push(`startDate=${encodeURIComponent(params.startDate)}`);
      }
      
      if (params.endDate !== undefined) {
        queryParts.push(`endDate=${encodeURIComponent(params.endDate)}`);
      }
      
      if (params.mealType !== undefined) {
        queryParts.push(`mealType=${encodeURIComponent(params.mealType)}`);
      }
      
      // 组合查询字符串
      if (queryParts.length > 0) {
        url += `?${queryParts.join('&')}`;
      }
    }
    
    console.log('请求URL:', url);
    return request.get(url) as Promise<FoodRecordsResponse>;
  },

  /**
   * 获取单个食物记录详情
   * @param recordId 记录ID
   */
  getFoodRecordDetail(recordId: number): Promise<FoodRecordResponse> {
    return request.get(`/api/food-records/${recordId}`) as Promise<FoodRecordResponse>;
  },

  /**
   * 更新食物记录
   * @param recordId 记录ID
   * @param data 更新的数据
   */
  updateFoodRecord(recordId: number, data: Partial<FoodRecord>): Promise<FoodRecordResponse> {
    return request.put(`/api/food-records/${recordId}`, data) as Promise<FoodRecordResponse>;
  },

  /**
   * 删除食物记录
   * @param recordId 记录ID
   */
  deleteFoodRecord(recordId: number): Promise<{ message: string }> {
    return request.delete(`/api/food-records/${recordId}`) as Promise<{ message: string }>;
  }
}; 