import axiosInstance from './axios';
import * as auth from './auth';
import * as admin from './admin';
import * as food from './food';
import * as health from './health';
import * as user from './user';

// 导出模块
export default {
  axios: axiosInstance,
  auth,
  admin,
  food,
  health,
  user
}; 