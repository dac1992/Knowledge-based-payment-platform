<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { ArrowDown, Odometer, SwitchButton, User } from '@element-plus/icons-vue'

import { ROLE_LABELS, ROLE_TAG_TYPES } from '@yanmai/shared'
import { YmLogo } from '@yanmai/ui'

import { useAuthStore } from '@/stores/auth'
import type { UserRole } from '@yanmai/shared'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const menus = [
  { path: '/dashboard', title: '工作台', icon: Odometer },
  { path: '/users', title: '用户管理', icon: User },
]

const activePath = computed(() => route.path)
const currentTitle = computed(() => (route.meta.title as string | undefined) ?? '')
const roleLabel = computed(() => {
  const r = auth.profile?.role
  return r ? ROLE_LABELS[r as UserRole] : '-'
})
const roleTagType = computed(() => {
  const r = auth.profile?.role
  return r ? ROLE_TAG_TYPES[r as UserRole] : 'primary'
})

async function handleLogout(): Promise<void> {
  try {
    await ElMessageBox.confirm('确认退出登录吗？', '退出登录', {
      type: 'warning',
      confirmButtonText: '退出',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  await auth.logout()
  await router.replace({ name: 'login' })
}
</script>

<template>
  <el-container class="layout">
    <el-aside width="232px" class="layout__aside">
      <div class="layout__brand">
        <YmLogo :size="34" subtitle="管理后台" inverse />
      </div>

      <el-menu :default-active="activePath" router class="layout__menu" :collapse="false">
        <el-menu-item v-for="m in menus" :key="m.path" :index="m.path">
          <el-icon><component :is="m.icon" /></el-icon>
          <span>{{ m.title }}</span>
        </el-menu-item>
      </el-menu>

      <div class="layout__aside-foot">
        <span class="layout__phase">Phase 1 · 工程骨架</span>
        <span class="layout__version">API v0.1.0</span>
      </div>
    </el-aside>

    <el-container>
      <el-header height="58px" class="layout__header">
        <div class="layout__crumb">
          <span class="layout__crumb-root">管理后台</span>
          <template v-if="currentTitle">
            <span class="layout__crumb-sep">/</span>
            <span class="layout__crumb-cur">{{ currentTitle }}</span>
          </template>
        </div>

        <el-dropdown trigger="click">
          <span class="layout__user">
            <span class="layout__avatar">{{ auth.displayName.slice(0, 1) || '研' }}</span>
            <span class="layout__user-info">
              <b>{{ auth.displayName || '未登录' }}</b>
              <el-tag :type="roleTagType" size="small" effect="light">{{ roleLabel }}</el-tag>
            </span>
            <el-icon class="layout__caret"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :icon="SwitchButton" @click="handleLogout"
                >退出登录</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>

      <el-main class="layout__main">
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  height: 100vh;
}

.layout__aside {
  display: flex;
  flex-direction: column;
  background: #0f1a2e;
  border-right: 1px solid #22304a;
}

.layout__brand {
  padding: 18px 18px 16px;
  border-bottom: 1px solid #22304a;
}

.layout__menu {
  flex: 1;
  border-right: none;
  background: transparent;
  padding: 12px 10px;
}

.layout__menu :deep(.el-menu-item) {
  height: 42px;
  line-height: 42px;
  margin-bottom: 4px;
  border-radius: 8px;
  color: #aebbd0;
  font-size: 13.5px;
}

.layout__menu :deep(.el-menu-item:hover) {
  background: #17233a;
  color: #fff;
}

.layout__menu :deep(.el-menu-item.is-active) {
  background: var(--ym-primary);
  color: #fff;
}

.layout__aside-foot {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 14px 20px 18px;
  border-top: 1px solid #22304a;
  font-size: 11.5px;
  color: #6b7b98;
}

.layout__phase {
  color: #8c9bb8;
}

.layout__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid var(--ym-line);
  padding: 0 22px;
}

.layout__crumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13.5px;
}

.layout__crumb-root {
  color: var(--ym-ink-4);
}

.layout__crumb-sep {
  color: var(--ym-ink-4);
}

.layout__crumb-cur {
  color: var(--ym-ink);
  font-weight: 500;
}

.layout__user {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 5px 9px;
  border-radius: 8px;
  cursor: pointer;
  outline: none;
  transition: background 0.15s ease;
}

.layout__user:hover {
  background: var(--ym-primary-light);
}

.layout__avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: var(--ym-primary);
  color: #fff;
  font-size: 13px;
}

.layout__user-info {
  display: flex;
  align-items: center;
  gap: 7px;
}

.layout__user-info b {
  font-size: 13.5px;
  font-weight: 500;
  color: var(--ym-ink);
}

.layout__caret {
  color: var(--ym-ink-4);
}

.layout__main {
  background: var(--ym-bg);
  padding: 22px 24px 40px;
  overflow-y: auto;
}
</style>
