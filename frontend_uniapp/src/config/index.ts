// 环境配置
const ENV = {
  development: 'development',
  production: 'production',
  test: 'test'
};

// 当前环境 - 默认使用开发环境
const currentEnv = ENV.development;

// API配置
const API_CONFIG = {
  [ENV.development]: {
    baseURL: 'https://jesuvukndxpo.sealoshzh.site',  // 开发环境API地址
    timeout: 15000,
    retryCount: 2,
    retryDelay: 1000
  },
  [ENV.test]: {
    baseURL: 'https://jesuvukndxpo.sealoshzh.site',  // 测试环境API地址
    timeout: 15000,
    retryCount: 1,
    retryDelay: 1000
  },
  [ENV.production]: {
    baseURL: 'https://jesuvukndxpo.sealoshzh.site',  // 生产环境API地址
    timeout: 15000,
    retryCount: 3,
    retryDelay: 1500
  }
};

// 导出当前环境的API配置
export const apiConfig = API_CONFIG[currentEnv];

// 导出环境常量
export const env = currentEnv;

// 导出是否为开发环境
export const isDev = currentEnv === ENV.development;

// 导出是否为生产环境
export const isProd = currentEnv === ENV.production;

// 导出是否为测试环境
export const isTest = currentEnv === ENV.test;

// API基础URL
export const API_BASE_URL = 'https://jesuvukndxpo.sealoshzh.site/api';
export const STATIC_BASE_URL = 'https://jesuvukndxpo.sealoshzh.site';

// 其他配置项
export const CONFIG = {
  // 每页显示的物品数量
  ITEMS_PER_PAGE: 9,
  
  // 默认图片
  DEFAULT_ITEM_IMAGE: '/static/collections/cat-locked.png'
}; 