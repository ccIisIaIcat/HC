import request from '../utils/axios';

interface AppUpdate {
  id: number;
  version: string;
  platform: string;
  update_url: string;
  release_notes: string;
  is_force: boolean;
  created_at: string;
  updated_at: string;
}

interface CheckUpdateResponse {
  has_update: boolean;
  update?: AppUpdate;
}

interface AppUpdatesResponse {
  updates: AppUpdate[];
}

/**
 * 显示更新提示对话框
 * @param update 更新信息
 * @returns {Promise<boolean>} 用户是否同意更新
 */
const showUpdateDialog = (update: AppUpdate): Promise<boolean> => {
  return new Promise((resolve) => {
    // 如果是强制更新，不显示取消按钮
    const showCancel = !update.is_force;
    
    uni.showModal({
      title: '发现新版本',
      content: `最新版本：${update.version}\n${update.release_notes}`,
      confirmText: '立即更新',
      cancelText: '稍后再说',
      showCancel,
      success: (res) => {
        if (res.confirm) {
          resolve(true);
        } else {
          resolve(false);
        }
      },
      fail: () => {
        resolve(false);
      }
    });
  });
};

/**
 * 检查并提示更新
 * @param platform 平台类型
 * @param current_version 当前版本号
 * @returns {Promise<void>}
 */
export const checkAndPromptUpdate = async (
  platform: string,
  current_version: string
): Promise<void> => {
  try {
    const response = await checkUpdate(platform, current_version);
    
    if (response.has_update && response.update) {
      const shouldUpdate = await showUpdateDialog(response.update);
      
      if (shouldUpdate) {
        // 开始下载更新
        uni.showLoading({ title: '下载更新中...' });
        
        try {
          const updateFile = await downloadUpdate(platform, response.update.version);
          uni.hideLoading();
          
          // 如果是APP环境，调用plus.runtime.install安装更新
          // #ifdef APP-PLUS
          plus.runtime.install(URL.createObjectURL(updateFile), {
            force: false
          }, () => {
            uni.showModal({
              title: '更新完成',
              content: '应用更新完成，请重启应用',
              showCancel: false,
              success: () => {
                plus.runtime.restart();
              }
            });
          }, (error) => {
            uni.showModal({
              title: '更新失败',
              content: `安装更新失败：${error.message}`,
              showCancel: false
            });
          });
          // #endif
          
          // 如果是H5环境，直接刷新页面
          // #ifdef H5
          window.location.reload();
          // #endif
        } catch (error) {
          uni.hideLoading();
          uni.showModal({
            title: '更新失败',
            content: '下载更新失败，请稍后重试',
            showCancel: false
          });
        }
      }
    }
  } catch (error) {
    console.error('检查更新失败:', error);
    // 检查更新失败时不显示错误提示，静默失败
  }
};

/**
 * 获取所有应用更新记录
 * @returns {Promise<AppUpdatesResponse>} 更新记录列表
 */
export const getAppUpdates = async (): Promise<AppUpdatesResponse> => {
  try {
    const response = await request({
      method: 'GET',
      url: '/api/app-updates'
    }) as AppUpdatesResponse;
    return response;
  } catch (error) {
    console.error('获取应用更新记录失败:', error);
    throw error;
  }
};

/**
 * 检查应用更新
 * @param platform 平台类型
 * @param current_version 当前版本号
 * @returns {Promise<CheckUpdateResponse>} 更新检查结果
 */
export const checkUpdate = async (
  platform: string,
  current_version: string
): Promise<CheckUpdateResponse> => {
  try {
    const response = await request({
      method: 'GET',
      url: '/api/app-updates/check',
      params: {
        platform,
        current_version
      }
    }) as CheckUpdateResponse;
    return response;
  } catch (error) {
    console.error('检查应用更新失败:', error);
    throw error;
  }
};

/**
 * 下载应用更新
 * @param platform 平台类型
 * @param version 版本号
 * @returns {Promise<Blob>} 更新文件
 */
export const downloadUpdate = async (
  platform: string,
  version: string
): Promise<Blob> => {
  try {
    const response = await request({
      method: 'GET',
      url: `/api/app-updates/download/${platform}/${version}`,
      responseType: 'blob'
    }) as Blob;
    return response;
  } catch (error) {
    console.error('下载应用更新失败:', error);
    throw error;
  }
};

export default {
  getAppUpdates,
  checkUpdate,
  downloadUpdate,
  checkAndPromptUpdate
}; 