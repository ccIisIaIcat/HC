<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
import { ref } from 'vue'

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
    version.value = plus.runtime.version
    console.log('当前版本：', version.value)
  } catch (error) {
    console.error('获取版本信息失败：', error)
  }
  // #endif

  // #ifdef H5
  version.value = '1.0.0' // H5环境下使用默认版本号
  // #endif
}

onLaunch(() => {
  console.log("App Launch");
  getSystemInfo() // 在应用启动时获取系统信息
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
