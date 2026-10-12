<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Refresh, Search } from '@element-plus/icons-vue'

import { ApiError } from '@yanmai/api-client'
import {
  ROLE_LABELS,
  ROLE_TAG_TYPES,
  STATUS_LABELS,
  formatDateTime,
  maskPhone,
  type ListUsersQuery,
  type UserProfile,
  type UserRole,
  type UserStatus,
} from '@yanmai/shared'
import { YmEmptyState, YmPageHeader } from '@yanmai/ui'

import { api } from '@/api/client'

const loading = ref(false)
const errorText = ref('')
const rows = ref<UserProfile[]>([])
const total = ref(0)

const query = reactive<ListUsersQuery>({
  page: 1,
  pageSize: 10,
  keyword: '',
  role: undefined,
  status: undefined,
})

const roleOptions: UserRole[] = ['admin', 'lecturer', 'student']
const statusOptions: UserStatus[] = [1, 0]

async function load(): Promise<void> {
  loading.value = true
  errorText.value = ''
  try {
    const res = await api.user.list({ ...query })
    rows.value = res.list
    total.value = res.total
  } catch (err) {
    errorText.value = err instanceof ApiError ? err.message : '加载用户列表失败'
    rows.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function onSearch(): void {
  query.page = 1
  void load()
}

function onReset(): void {
  query.keyword = ''
  query.role = undefined
  query.status = undefined
  query.page = 1
  void load()
}

function onPageChange(page: number): void {
  query.page = page
  void load()
}

function onSizeChange(size: number): void {
  query.pageSize = size
  query.page = 1
  void load()
}

onMounted(load)
</script>

<template>
  <div class="users">
    <YmPageHeader title="用户管理" description="平台账号列表，支持关键字、角色与状态筛选（仅管理员可见）">
      <template #actions>
        <el-button :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
      </template>
    </YmPageHeader>

    <div class="ym-card users__toolbar">
      <el-form :model="query" inline @submit.prevent="onSearch">
        <el-form-item label="关键字">
          <el-input
            v-model="query.keyword"
            placeholder="登录名 / 昵称 / 手机号"
            clearable
            style="width: 220px"
            @keyup.enter="onSearch"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="query.role" placeholder="全部" clearable style="width: 140px">
            <el-option v-for="r in roleOptions" :key="r" :label="ROLE_LABELS[r]" :value="r" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width: 120px">
            <el-option v-for="s in statusOptions" :key="s" :label="STATUS_LABELS[s]" :value="s" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="onSearch">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-alert
      v-if="errorText"
      class="users__alert"
      type="warning"
      :closable="false"
      show-icon
      :title="errorText"
    />

    <div class="ym-card users__table">
      <el-table v-loading="loading" :data="rows" size="default" :empty-text="' '">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="登录名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column label="手机号" min-width="140">
          <template #default="{ row }">{{ maskPhone(row.phone) }}</template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180">
          <template #default="{ row }">{{ row.email || '-' }}</template>
        </el-table-column>
        <el-table-column label="角色" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="ROLE_TAG_TYPES[row.role as UserRole]" size="small" effect="light">
              {{ ROLE_LABELS[row.role as UserRole] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small" effect="light">
              {{ STATUS_LABELS[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最近登录" width="170">
          <template #default="{ row }">{{ formatDateTime(row.lastLoginAt) }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template>
        </el-table-column>

        <template #empty>
          <YmEmptyState
            size="sm"
            :title="errorText ? '数据加载失败' : '暂无用户数据'"
            :description="errorText ? '请确认后端服务与数据库连接正常后重试。' : '可以尝试调整筛选条件。'"
          />
        </template>
      </el-table>

      <div v-if="total > 0" class="users__pager">
        <el-pagination
          :current-page="query.page"
          :page-size="query.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @current-change="onPageChange"
          @size-change="onSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.users__toolbar {
  padding: 16px 18px 2px;
  margin-bottom: 16px;
}

.users__alert {
  margin-bottom: 16px;
}

.users__table {
  padding: 6px 8px 8px;
}

.users__pager {
  display: flex;
  justify-content: flex-end;
  padding: 14px 12px 8px;
}
</style>
