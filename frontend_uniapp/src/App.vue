<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
onLaunch(() => {
  console.log("App Launch");
});
onShow(() => {
  console.log("App Show");
});
onHide(() => {
  console.log("App Hide");
});
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
