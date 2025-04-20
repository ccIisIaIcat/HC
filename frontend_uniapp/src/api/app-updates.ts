import request from '@/utils/request';
import { apiConfig } from '@/config';

// 应用更新信息接口
interface AppUpdateInfo {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: null | string;
  version: string;
  platform: string;
  update_url: string;
  force_update: boolean;
  update_notes: string;
  release_date: string;
  file_size: number;
  md5: string;
}

// 获取所有更新响应接口
interface GetAllUpdatesResponse {
  data: AppUpdateInfo[];
}

// 检查更新响应接口
interface CheckUpdateResponse {
  needs_update: boolean;
  update_info?: AppUpdateInfo;
}

export default {
  /**
   * 获取所有应用更新信息
   */
  getAllUpdates(): Promise<GetAllUpdatesResponse> {
    return request.get('/api/app-updates');
  },

  /**
   * 检查应用更新
   * @param platform 平台（android/ios）
   * @param currentVersion 当前版本号
   */
  checkUpdate(platform: string, currentVersion: string): Promise<CheckUpdateResponse> {
    return request.get('/api/app-updates/check', {
      data: {
        platform,
        current_version: currentVersion
      }
    });
  },

  /**
   * 下载更新
   * @param platform 平台（android/ios）
   * @param version 版本号
   * @returns 下载链接
   */
  getDownloadUrl(platform: string, version: string): string {
    return `/api/app-updates/download/${platform}/${version}`;
  },

  /**
   * 执行下载更新
   * @param platform 平台（android/ios）
   * @param version 版本号
   */
  async downloadUpdate(platform: string, version: string): Promise<void> {
    // #ifdef APP-PLUS
    try {
      const downloadUrl = this.getDownloadUrl(platform, version);
      const fullUrl = apiConfig.baseURL + downloadUrl;
      
      return new Promise((resolve, reject) => {
        if (platform === 'android') {
          // 使用系统浏览器下载
          plus.runtime.openURL(fullUrl, (err) => {
            if (err) {
              uni.showModal({
                title: '打开下载失败',
                content: '请手动复制链接在浏览器中打开下载：' + fullUrl,
                showCancel: true,
                confirmText: '复制链接',
                success: (res) => {
                  if (res.confirm) {
                    uni.setClipboardData({
                      data: fullUrl,
                      success: () => {
                        uni.showToast({
                          title: '链接已复制',
                          icon: 'success'
                        });
                      }
                    });
                  }
                }
              });
              reject(err);
            } else {
              uni.showToast({
                title: '正在打开浏览器下载',
                icon: 'none'
              });
              resolve();
            }
          });
        } else if (platform === 'ios') {
          // iOS 跳转到 App Store
          plus.runtime.openURL(fullUrl);
          resolve();
        }
      });
    } catch (error: any) {
      const errorMsg = error.message || '未知错误';
      console.error('更新失败：', errorMsg);
      uni.showModal({
        title: '更新失败',
        content: errorMsg,
        showCancel: false
      });
      throw error;
    }
    // #endif
  }
}; 