<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection, Refresh, VideoPlay } from '@element-plus/icons-vue'

import { ApiError } from '@yanmai/api-client'
import {
  CLIENTS,
  DELIVERED_MODULE_COUNT,
  MODULES,
  ROLE_LABELS,
  formatDateTime,
  type HealthStatus,
  type ModuleMeta,
  type UserRole,
} from '@yanmai/shared'
import { YmEmptyState, YmPageHeader, YmStatCard } from '@yanmai/ui'

import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const loading = ref(true)
const health = ref<HealthStatus | null>(null)
const meta = ref<ModuleMeta | null>(null)
const errorText = ref('')

const totalModules = MODULES.length
const roleLabel = computed(() => {
  const r = auth.profile?.role
  return r ? ROLE_LABELS[r as UserRole] : '-'
})

const backendTone = computed<'up' | 'down' | 'warn'>(() => {
  if (!health.value) return 'warn'
  if (health.value.status === 'ok') return 'up'
  return 'warn'
})

const backendText = computed(() => {
  if (!health.value) return '未连接'
  return health.value.status === 'ok' ? '正常' : '降级'
})

async function load(): Promise<void> {
  loading.value = true
  errorText.value = ''
  try {
    const [h, m] = await Promise.all([api.system.health(), api.system.modules()])
    health.value = h
    meta.value = m
  } catch (err) {
    errorText.value =
      err instanceof ApiError ? err.message : '无法获取后端信息，请确认 API 服务已启动'
  } finally {
    loading.value = false
  }
}

/** 探活单个业务模块：命中占位路由会返回 501，据此确认路由已挂载。 */
async function probe(key: string): Promise<void> {
  try {
    await api.request({ url: `/${key}`, method: 'GET' })
    ElMessage.success(`${key} 模块已实现`)
  } catch (err) {
    const message = err instanceof ApiError ? err.message : '探测失败'
    ElMessage({ type: 'info', message, duration: 4500, showClose: true })
  }
}

onMounted(load)
</script>

<template>
  <div class="dash">
    <YmPageHeader title="运营工作台" description="Phase 1 工程骨架运行状态与业务模块交付进度">
      <template #actions>
        <el-button :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
      </template>
    </YmPageHeader>

    <el-alert
      v-if="errorText"
      class="dash__alert"
      type="warning"
      :closable="false"
      show-icon
      title="后端服务暂不可用"
    >
      <template #default>
        <p class="dash__alert-p">{{ errorText }}</p>
        <p class="dash__alert-p">
          启动方式：<code>cd services/api &amp;&amp; make run</code>；并确认
          <code>.env.local</code> 中的 <code>DATABASE_URL</code> 已指向云 PostgreSQL。
        </p>
      </template>
    </el-alert>

    <section class="dash__stats">
      <YmStatCard
        label="后端服务"
        :value="backendText"
        :tone="backendTone"
        :hint="health ? `${health.service} · ${health.env} · v${health.version}` : '等待探测'"
        :loading="loading"
      />
      <YmStatCard
        label="已挂载模块"
        :value="meta ? meta.modules.length : '-'"
        hint="后端 internal/modules 下注册的业务域"
        :loading="loading"
      />
      <YmStatCard
        label="已交付模块"
        :value="`${DELIVERED_MODULE_COUNT} / ${totalModules}`"
        hint="Phase 1 已可用的业务域"
        :loading="loading"
      />
      <YmStatCard
        label="当前登录角色"
        :value="roleLabel"
        :hint="auth.profile ? `账号 ${auth.profile.username}` : '未登录'"
        :loading="loading"
      />
    </section>

    <el-row :gutter="18" class="dash__row">
      <el-col :xs="24" :lg="10">
        <div class="ym-card dash__panel">
          <h3 class="ym-section-title">服务健康</h3>
          <el-descriptions v-if="health" :column="1" size="small" border>
            <el-descriptions-item label="整体状态">
              <el-tag
                :type="health.status === 'ok' ? 'success' : 'warning'"
                size="small"
                effect="light"
              >
                {{ health.status === 'ok' ? '正常' : '降级' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="数据库">
              <el-tag
                :type="health.database.status === 'up' ? 'success' : 'danger'"
                size="small"
                effect="light"
              >
                {{ health.database.status === 'up' ? '已连通' : '不可用' }}
              </el-tag>
              <span class="dash__latency">探测耗时 {{ health.database.latency }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="运行环境">{{ health.env }}</el-descriptions-item>
            <el-descriptions-item label="服务版本">v{{ health.version }}</el-descriptions-item>
            <el-descriptions-item label="当前时间">{{
              formatDateTime(health.time)
            }}</el-descriptions-item>
            <el-descriptions-item label="队列状态">
              <el-tag :type="meta?.queue ? 'success' : 'info'" size="small" effect="light">
                {{ meta?.queue ? '已启用' : '未启用（缺 REDIS_ADDR）' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
          <YmEmptyState
            v-else
            size="sm"
            title="尚未获取到后端信息"
            description="确认 API 服务已启动后点击右上角「刷新」。"
          />

          <p v-if="health?.database.error" class="dash__db-error">
            <el-icon><Connection /></el-icon>
            {{ health.database.error }}
          </p>
        </div>
      </el-col>

      <el-col :xs="24" :lg="14">
        <div class="ym-card dash__panel">
          <h3 class="ym-section-title">六端矩阵</h3>
          <el-table :data="CLIENTS" size="small" :show-header="true">
            <el-table-column prop="name" label="端" min-width="130" />
            <el-table-column prop="stack" label="技术栈" min-width="200" />
            <el-table-column label="交付批次" width="100" align="center">
              <template #default="{ row }">第 {{ row.batch }} 批</template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.ready ? 'success' : 'info'" size="small" effect="light">
                  {{ row.ready ? '已就绪' : '后置' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>

    <div class="ym-card dash__panel">
      <h3 class="ym-section-title">业务模块交付进度</h3>
      <el-table :data="MODULES" size="small">
        <el-table-column prop="key" label="模块" width="120">
          <template #default="{ row }">
            <span class="ym-mono">{{ row.key }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="130" />
        <el-table-column prop="summary" label="规划能力" min-width="300" show-overflow-tooltip />
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
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" :icon="VideoPlay" @click="probe(row.key)">
              探活
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <p class="dash__tip">
        点击「探活」会请求
        <code>/api/v1/&lt;模块&gt;</code>：已交付模块返回正常数据，待开发模块由后端统一返回 501
        并说明计划阶段。
      </p>
    </div>
  </div>
</template>

<style scoped>
.dash__alert {
  margin-bottom: 18px;
}

.dash__alert-p {
  margin: 2px 0;
  font-size: 12.5px;
  line-height: 1.8;
}

.dash__alert-p code {
  padding: 1px 5px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.06);
  font-family: var(--ym-font-mono);
  font-size: 12px;
}

.dash__stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 14px;
  margin-bottom: 18px;
}

.dash__row {
  margin-bottom: 18px;
}

.dash__panel {
  padding: 18px 20px 20px;
  height: 100%;
}

.dash__latency {
  margin-left: 8px;
  font-size: 12px;
  color: var(--ym-ink-4);
}

.dash__db-error {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin: 12px 0 0;
  padding: 9px 11px;
  border-radius: 8px;
  background: var(--ym-warn-bg);
  color: var(--ym-warn);
  font-size: 12px;
  line-height: 1.7;
  word-break: break-all;
}

.dash__tip {
  margin: 12px 0 0;
  font-size: 12px;
  color: var(--ym-ink-4);
  line-height: 1.8;
}

.dash__tip code {
  padding: 1px 5px;
  border-radius: 4px;
  background: var(--ym-line-light);
  font-family: var(--ym-font-mono);
}
</style>
