<script setup lang="ts">
import { MODULES } from '@yanmai/shared'

const courseModule = MODULES.find((m) => m.key === 'course')
</script>

<template>
  <view class="ym-page course">
    <view class="course__head">
      <text class="course__title">课程中心</text>
      <text class="ym-sub">分类浏览 · 课程详情 · 搜索筛选 · 讲师主页</text>
    </view>

    <view class="course__placeholder">
      <view class="course__placeholder-mark">03</view>
      <text class="course__placeholder-title">该模块将在 Phase 3 交付</text>
      <text class="course__placeholder-desc">{{ courseModule?.summary }}</text>
    </view>

    <view class="ym-card">
      <text class="ym-title">规划能力</text>
      <view class="course__list">
        <view class="course__item">
          <text class="course__dot" />
          <text>分类体系：学科分类树、标签体系、专题聚合</text>
        </view>
        <view class="course__item">
          <text class="course__dot" />
          <text>课程模型：SPU / SKU 分离，支持单课、专栏、训练营</text>
        </view>
        <view class="course__item">
          <text class="course__dot" />
          <text>内容组织：章 - 节 - 课时，支持视频 / 图文 / 附件</text>
        </view>
        <view class="course__item">
          <text class="course__dot" />
          <text>讲师建课：草稿 → 提交审核 → 上架，素材存 COS</text>
        </view>
      </view>
    </view>

    <view class="ym-card">
      <text class="ym-title">当前状态</text>
      <text class="ym-sub">
        Phase 1 仅交付工程骨架：登录态、路由、请求封装已打通。课程数据模型将在 Phase 3
        随数据库迁移一并落地。
      </text>
    </view>
  </view>
</template>

<style scoped>
.course__head {
  padding: 8rpx 8rpx 28rpx;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.course__title {
  font-size: 40rpx;
  font-weight: 600;
  color: #1a1f2b;
}

.course__placeholder {
  background: #ffffff;
  border: 1rpx dashed #c6d5f8;
  border-radius: 20rpx;
  padding: 56rpx 30rpx;
  margin-bottom: 24rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
}

.course__placeholder-mark {
  font-size: 56rpx;
  font-weight: 600;
  color: #c6d5f8;
}

.course__placeholder-title {
  font-size: 30rpx;
  color: #1a1f2b;
}

.course__placeholder-desc {
  font-size: 24rpx;
  color: #6b7688;
  text-align: center;
  line-height: 1.6;
}

.course__list {
  margin-top: 18rpx;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.course__item {
  display: flex;
  align-items: flex-start;
  gap: 14rpx;
  font-size: 26rpx;
  color: #3d4759;
  line-height: 1.6;
}

.course__dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  background: #1a56db;
  margin-top: 14rpx;
  flex-shrink: 0;
}
</style>
