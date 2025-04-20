<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
import { ref } from 'vue'
import appUpdates from '@/api/app-updates'

// 版本和平台信息
const version = ref('')
const platform = ref('')

// 获取系统信息和版本信息
const getSystemInfo = () => {
  uni.getSystemInfo({
    success: (res) => {
      platform.value = res.platform.toLowerCase()
      console.log('当前平台：', platform.value)
    },
    fail: (err) => {
      console.error('获取系统信息失败：', err)
    }
  })

  // #ifdef APP-PLUS
  try {
    const runtimeVersion = plus.runtime.version
    if (runtimeVersion) {
      version.value = runtimeVersion
      console.log('当前版本：', version.value)
    } else {
      version.value = '1.0.0' // 如果获取失败，使用默认版本号
      console.warn('无法获取版本号，使用默认版本：1.0.0')
    }
  } catch (error) {
    console.error('获取版本信息失败：', error)
    version.value = '1.0.0' // 发生错误时使用默认版本号
  }
  // #endif

  // #ifdef H5
  version.value = '1.0.0' // H5环境下使用默认版本号
  // #endif
}

// 检查应用更新
const checkUpdate = async () => {
  // 只在 APP 环境下检查更新
  // #ifdef APP-PLUS
  try {
    const result = await appUpdates.checkUpdate(platform.value, version.value)
    
    if (result.needs_update && result.update_info) {
      // 显示更新提示
      uni.showModal({
        title: '发现新版本',
        content: `最新版本：${result.update_info.version}\n${result.update_info.update_notes}`,
        confirmText: '立即更新',
        cancelText: result.update_info.force_update ? '退出应用' : '暂不更新',
        success: async (res) => {
          if (res.confirm) {
            // 用户点击更新
            try {
              await appUpdates.downloadUpdate(platform.value, result.update_info!.version)
            } catch (error) {
              uni.showToast({
                title: '更新失败，请稍后重试',
                icon: 'none'
              })
            }
          } else if (result.update_info?.force_update) {
            // 强制更新，用户点击取消则退出应用
            plus.runtime.quit()
          }
        }
      })
    }
  } catch (error) {
    console.error('检查更新失败：', error)
  }
  // #endif
}

onLaunch(() => {
  console.log("App Launch");
  getSystemInfo() // 在应用启动时获取系统信息
  checkUpdate() // 检查应用更新
});

onShow(() => {
  console.log("App Show");
});

onHide(() => {
  console.log("App Hide");
});

// 导出版本信息，供其他组件使用
const getAppInfo = () => {
  return {
    version: version.value,
    platform: platform.value
  }
}

defineExpose({
  getAppInfo
})
</script>
<style>
	/* 全局样式 */
	page {
		background-color: #000000;
		color: #FFFFFF;
		font-family: -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Helvetica, Segoe UI, Arial, Roboto, 'PingFang SC', 'miui', 'Hiragino Sans GB', 'Microsoft Yahei', sans-serif;
	}
	
	/* 自定义tabBar样式 */
	.uni-tabbar {
		height: 80px !important;
		padding-bottom: 3px !important;
	}
	
	.uni-tabbar__item {
		padding-top: 2px !important;
	}
	
	.uni-tabbar__icon {
		width: 32px !important;
		height: 40px !important;
		margin-bottom: 4px !important;
		object-fit: contain !important;
	}
	
	.uni-tabbar__label {
		font-size: 14px !important;
		line-height: 1.2 !important;
		margin-top: 2px !important;
	}
</style>
