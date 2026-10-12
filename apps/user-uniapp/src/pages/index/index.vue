<script setup lang="ts">
import { onMounted, ref } from 'vue'

import {
  CLIENTS,
  DELIVERED_MODULE_COUNT,
  MODULES,
  type HealthStatus,
  type ModuleMeta,
} from '@yanmai/shared'

import { UniApiError, api } from '@/api/request'

const loading = ref(true)
const errorText = ref('')
const health = ref<HealthStatus | null>(null)
const meta = ref<ModuleMeta | null>(null)

const totalModules = MODULES.length
const upcoming = MODULES.filter((m) => !m.delivered).slice(0, 4)

async function load(): Promise<void> {
  loading.value = true
  errorText.value = ''
  try {
    const [h, m] = await Promise.all([api.health(), api.modules()])
    health.value = h
    meta.value = m
  } catch (err) {
    errorText.value = err instanceof UniApiError ? err.message : '无法连接后端服务'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <view class="ym-page home">
    <view class="home__hero">
      <view class="home__mark">研</view>
      <view class="home__hero-text">
        <text class="home__brand">研脉学术</text>
        <text class="home__slogan">面向学术科研人群的知识付费与在线学习平台</text>
      </view>
    </view>

    <view class="ym-card">
      <view class="home__row">
        <text class="ym-title">服务状态</text>
        <text
          class="home__badge"
          :class="health ? (health.status === 'ok' ? 'is-ok' : 'is-warn') : 'is-off'"
        >
          {{ health ? (health.status === 'ok' ? '正常' : '降级') : '未连接' }}
        </text>
      </view>

      <view v-if="loading" class="ym-sub">正在探测后端服务…</view>

      <view v-else-if="errorText" class="home__error">
        <text class="home__error-title">{{ errorText }}</text>
        <text class="home__error-tip">
          请确认 API 服务已启动（cd services/api &amp;&amp; make run），并已配置云 PostgreSQL
          连接串。
        </text>
      </view>

      <view v-else class="home__kv">
        <view class="home__kv-row">
          <text class="ym-muted">服务</text>
          <text>{{ health?.service }} · {{ health?.env }} · v{{ health?.version }}</text>
        </view>
        <view class="home__kv-row">
          <text class="ym-muted">数据库</text>
          <text>{{ health?.database.status === 'up' ? '已连通' : '不可用' }}</text>
        </view>
        <view class="home__kv-row">
          <text class="ym-muted">已挂载模块</text>
          <text>{{ meta?.modules.length ?? '-' }} 个</text>
        </view>
        <view class="home__kv-row">
          <text class="ym-muted">队列</text>
          <text>{{ meta?.queue ? '已启用' : '未启用' }}</text>
        </view>
      </view>
    </view>

    <view class="ym-card">
      <text class="ym-title">交付进度</text>
      <view class="home__progress">
        <view class="home__progress-bar">
          <view
            class="home__progress-inner"
            :style="{ width: `${Math.round((DELIVERED_MODULE_COUNT / totalModules) * 100)}%` }"
          />
        </view>
        <text class="ym-sub"
          >{{ DELIVERED_MODULE_COUNT }} / {{ totalModules }} 个业务模块已交付</text
        >
      </view>
    </view>

    <view class="ym-card">
      <text class="ym-title">即将上线</text>
      <view v-for="m in upcoming" :key="m.key" class="home__upcoming">
        <view class="home__upcoming-head">
          <text class="home__upcoming-name">{{ m.name }}</text>
          <text class="home__upcoming-phase">Phase {{ m.phase }}</text>
        </view>
        <text class="ym-sub">{{ m.summary }}</text>
      </view>
    </view>

    <view class="ym-card">
      <text class="ym-title">六端矩阵</text>
      <view v-for="c in CLIENTS" :key="c.key" class="home__client">
        <view class="home__client-main">
          <text class="home__client-name">{{ c.name }}</text>
          <text class="ym-sub">{{ c.stack }}</text>
        </view>
        <text class="home__client-tag" :class="{ 'is-ready': c.ready }">
          {{ c.ready ? '已就绪' : '后置' }}
        </text>
      </view>
    </view>

    <view class="home__foot">研脉学术（北京）智能科技有限公司</view>
  </view>
</template>

<style scoped>
.home__hero {
  display: flex;
  align-items: center;
  gap: 22rpx;
  padding: 8rpx 8rpx 36rpx;
}

.home__mark {
  width: 96rpx;
  height: 96rpx;
  border-radius: 24rpx;
  background: #1a56db;
  color: #ffffff;
  font-size: 44rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.home__hero-text {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.home__brand {
  font-size: 40rpx;
  font-weight: 600;
  color: #1a1f2b;
}

.home__slogan {
  font-size: 24rpx;
  color: #6b7688;
  line-height: 1.5;
}

.home__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18rpx;
}

.home__badge {
  padding: 4rpx 18rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
}

.home__badge.is-ok {
  background: #e9f8ef;
  color: #16a34a;
}

.home__badge.is-warn {
  background: #fef5e7;
  color: #d97706;
}

.home__badge.is-off {
  background: #eef1f6;
  color: #6b7688;
}

.home__kv {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.home__kv-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 26rpx;
}

.home__error {
  background: #fef5e7;
  border-radius: 14rpx;
  padding: 20rpx 22rpx;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.home__error-title {
  font-size: 26rpx;
  color: #d97706;
}

.home__error-tip {
  font-size: 23rpx;
  color: #b45309;
  line-height: 1.6;
}

.home__progress {
  margin-top: 16rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.home__progress-bar {
  height: 14rpx;
  border-radius: 999rpx;
  background: #eef1f6;
  overflow: hidden;
}

.home__progress-inner {
  height: 100%;
  border-radius: 999rpx;
  background: #1a56db;
}

.home__upcoming {
  margin-top: 22rpx;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.home__upcoming-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.home__upcoming-name {
  font-size: 28rpx;
  color: #1a1f2b;
}

.home__upcoming-phase {
  font-size: 22rpx;
  color: #c9a227;
  background: #fbf4df;
  padding: 2rpx 14rpx;
  border-radius: 999rpx;
}

.home__client {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
  margin-top: 20rpx;
}

.home__client-main {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  min-width: 0;
}

.home__client-name {
  font-size: 27rpx;
  color: #1a1f2b;
}

.home__client-tag {
  font-size: 22rpx;
  color: #98a2b3;
  background: #eef1f6;
  padding: 3rpx 16rpx;
  border-radius: 999rpx;
  flex-shrink: 0;
}

.home__client-tag.is-ready {
  color: #16a34a;
  background: #e9f8ef;
}

.home__foot {
  text-align: center;
  font-size: 22rpx;
  color: #98a2b3;
  padding: 20rpx 0 10rpx;
}
</style>
