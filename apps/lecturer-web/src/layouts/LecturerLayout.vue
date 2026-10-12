<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { ArrowDown, House, SwitchButton } from '@element-plus/icons-vue'

import { ROLE_LABELS, type UserRole } from '@yanmai/shared'
import { YmLogo } from '@yanmai/ui'

import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const menus = [{ path: '/workspace', title: '我的工作台', icon: House }]
const activePath = computed(() => route.path)
const currentTitle = computed(() => (route.meta.title as string | undefined) ?? '')
const roleLabel = computed(() => {
  const r = auth.profile?.role
  return r ? ROLE_LABELS[r as UserRole] : '-'
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
    <el-aside width="220px" class="layout__aside">
      <div class="layout__brand">
        <YmLogo :size="34" subtitle="讲师工作台" inverse />
      </div>

      <el-menu :default-active="activePath" router class="layout__menu">
        <el-menu-item v-for="m in menus" :key="m.path" :index="m.path">
          <el-icon><component :is="m.icon" /></el-icon>
          <span>{{ m.title }}</span>
        </el-menu-item>
      </el-menu>

      <div class="layout__aside-foot">
        <span>Phase 1 · 工程骨架</span>
        <span class="layout__dim">建课 / 直播 / 教务将在 Phase 3 起交付</span>
      </div>
    </el-aside>

    <el-container>
      <el-header height="58px" class="layout__header">
        <div class="layout__crumb">
          <span class="layout__crumb-root">讲师工作台</span>
          <template v-if="currentTitle">
            <span class="layout__crumb-sep">/</span>
            <span class="layout__crumb-cur">{{ currentTitle }}</span>
          </template>
        </div>

        <el-dropdown trigger="click">
          <span class="layout__user">
            <span class="layout__avatar">{{ auth.displayName.slice(0, 1) || '师' }}</span>
            <span class="layout__user-info">
              <b>{{ auth.displayName || '未登录' }}</b>
              <el-tag type="warning" size="small" effect="light">{{ roleLabel }}</el-tag>
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
        <router-view />
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
  background: #c9a227;
  color: #fff;
}

.layout__aside-foot {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 14px 20px 18px;
  border-top: 1px solid #22304a;
  font-size: 11.5px;
  color: #8c9bb8;
}

.layout__dim {
  color: #6b7b98;
  line-height: 1.6;
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

.layout__crumb-root,
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
}

.layout__user:hover {
  background: var(--ym-accent-light);
}

.layout__avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: var(--ym-accent);
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
