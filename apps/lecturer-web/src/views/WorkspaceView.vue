<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'

import { ApiError } from '@yanmai/api-client'
import { MODULES, ROLE_LABELS, formatDateTime, maskPhone, type UserRole } from '@yanmai/shared'
import { YmPageHeader, YmStatCard } from '@yanmai/ui'

import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const loading = ref(false)
const healthText = ref('未连接')
const healthOk = ref(false)
const errorText = ref('')

const lecturerModules = computed(() => MODULES.filter((m) => ['course', 'live', 'finance'].includes(m.key)))
const roleLabel = computed(() => {
  const r = auth.profile?.role
  return r ? ROLE_LABELS[r as UserRole] : '-'
})

async function load(): Promise<void> {
  loading.value = true
  errorText.value = ''
  try {
    const health = await api.system.health()
    healthOk.value = health.status === 'ok'
    healthText.value = health.status === 'ok' ? '正常' : '降级'
  } catch (err) {
    healthOk.value = false
    healthText.value = '未连接'
    errorText.value = err instanceof ApiError ? err.message : '无法连接后端服务'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="ws">
    <YmPageHeader title="我的工作台" description="Phase 1 骨架已打通登录态与接口调用，讲师业务能力将在 Phase 3 起逐步开放">
      <template #actions>
        <el-button :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
      </template>
    </YmPageHeader>

    <el-alert
      v-if="errorText"
      class="ws__alert"
      type="warning"
      :closable="false"
      show-icon
      title="后端服务暂不可用"
    >
      <template #default>
        <span class="ws__alert-body">{{ errorText }}</span>
      </template>
    </el-alert>

    <section class="ws__stats">
      <YmStatCard label="当前账号" :value="auth.displayName || '-'" :hint="`角色：${roleLabel}`" />
      <YmStatCard
        label="账号资料"
        :value="maskPhone(auth.profile?.phone)"
        :hint="auth.profile?.email || '未填写邮箱'"
      />
      <YmStatCard
        label="后端服务"
        :value="healthText"
        :tone="healthOk ? 'up' : 'warn'"
        hint="接口连通性探测结果"
        :loading="loading"
      />
      <YmStatCard
        label="最近登录"
        :value="formatDateTime(auth.profile?.lastLoginAt) === '-' ? '首次登录' : '已登录'"
        :hint="formatDateTime(auth.profile?.lastLoginAt)"
      />
    </section>

    <div class="ym-card ws__panel">
      <h3 class="ym-section-title">讲师相关能力规划</h3>
      <el-table :data="lecturerModules" size="small">
        <el-table-column prop="key" label="模块" width="120">
          <template #default="{ row }">
            <span class="ym-mono">{{ row.key }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="130" />
        <el-table-column prop="summary" label="规划能力" min-width="320" show-overflow-tooltip />
        <el-table-column label="阶段" width="90" align="center">
          <template #default="{ row }">Phase {{ row.phase }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.delivered ? 'success' : 'info'" size="small" effect="light">
              {{ row.delivered ? '已交付' : '待开发' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="ym-card ws__panel">
      <h3 class="ym-section-title">当前阶段说明</h3>
      <p class="ws__p">
        Phase 1 交付的是<b>工程骨架</b>：统一账号体系（JWT 双 Token）、模块化单体后端、
        三端前端外壳、共享组件与接口 SDK、CI 流水线。
      </p>
      <p class="ws__p">
        讲师端真正可用的建课、直播开播、教务与收益结算，按计划在
        <b>Phase 3（课程内容中心）</b>与 <b>Phase 9（直播）</b>、<b>Phase 12（财务结算）</b> 交付。
      </p>
    </div>
  </div>
</template>

<style scoped>
.ws__alert {
  margin-bottom: 18px;
}

.ws__alert-body {
  font-size: 12.5px;
}

.ws__stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 14px;
  margin-bottom: 18px;
}

.ws__panel {
  padding: 18px 20px 20px;
  margin-bottom: 18px;
}

.ws__p {
  margin: 0 0 10px;
  font-size: 13px;
  line-height: 1.9;
  color: var(--ym-ink-2);
}

.ws__p b {
  font-weight: 500;
  color: var(--ym-ink);
}
</style>
