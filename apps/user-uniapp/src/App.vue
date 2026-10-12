<script setup lang="ts">
import { onLaunch, onShow } from '@dcloudio/uni-app'

onLaunch(() => {
  // Phase 1 骨架：仅打印启动日志，后续在此恢复登录态、初始化埋点
  console.log('[研脉学术] App Launch')
})

onShow(() => {
  console.log('[研脉学术] App Show')
})
</script>

<style>
page {
  background-color: #f4f6fa;
  color: #1a1f2b;
  font-size: 28rpx;
  line-height: 1.7;
}

.ym-page {
  min-height: 100vh;
  padding: 32rpx 28rpx 60rpx;
  box-sizing: border-box;
}

.ym-card {
  background: #ffffff;
  border: 1rpx solid #e4e8ef;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
}

.ym-title {
  font-size: 34rpx;
  font-weight: 600;
  color: #1a1f2b;
}

.ym-sub {
  font-size: 24rpx;
  color: #6b7688;
}

.ym-muted {
  color: #98a2b3;
}
</style>
