<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
import 'uno.css';

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
/* 全局样式 - 使用 UnoCSS 替代 */
</style>
