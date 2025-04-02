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