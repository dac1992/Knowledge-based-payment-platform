<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowRight, CircleCheck, InfoFilled, Lock, User } from '@element-plus/icons-vue'

import { DELIVERED_MODULE_COUNT, MODULES } from '@yanmai/shared'
import { YmLogo } from '@yanmai/ui'

import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const formRef = ref<FormInstance>()
const form = ref({ login: 'admin', password: 'Admin@123456' })

const rules: FormRules = {
  login: [
    { required: true, message: '请输入登录名或手机号', trigger: 'blur' },
    { min: 3, max: 64, message: '长度需在 3-64 个字符之间', trigger: 'blur' },
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const seedAccounts = [
  { login: 'admin', password: 'Admin@123456', role: '管理员' },
  { login: 'lecturer', password: 'Lecturer@123456', role: '讲师' },
  { login: 'student', password: 'Student@123456', role: '学员' },
]

const totalModules = computed(() => MODULES.length)

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  const ok = await auth.login({ ...form.value })
  if (!ok) {
    ElMessage.error(auth.error || '登录失败')
    return
  }

  ElMessage.success(`欢迎回来，${auth.displayName}`)
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard'
  await router.replace(redirect)
}

function useSeed(login: string, password: string): void {
  form.value = { login, password }
}
</script>

<template>
  <div class="login">
    <aside class="login__brand">
      <YmLogo :size="42" subtitle="面向学术科研人群的知识付费平台" inverse />

      <div class="login__headline">
        <h1>让学术成长<br />有迹可循</h1>
        <p>课程 · 科研训练营 · 专家直播 · 学术社群</p>
      </div>

      <ul class="login__facts">
        <li>
          <el-icon><CircleCheck /></el-icon>
          <span>六端矩阵：管理后台 / 讲师端 / 小程序 / PC / APP / 桌面端</span>
        </li>
        <li>
          <el-icon><CircleCheck /></el-icon>
          <span>Phase 1 骨架已交付，当前进度 {{ DELIVERED_MODULE_COUNT }} / {{ totalModules }} 个业务模块</span>
        </li>
        <li>
          <el-icon><CircleCheck /></el-icon>
          <span>统一账号体系：手机号 + 微信 UnionID，JWT 双 Token 跨端互通</span>
        </li>
      </ul>

      <div class="login__foot">研脉学术（北京）智能科技有限公司</div>
    </aside>

    <main class="login__panel">
      <div class="login__form-wrap">
        <h2>登录管理后台</h2>
        <p class="login__sub">请使用平台账号登录，登录后进入运营工作台</p>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          size="large"
          @submit.prevent="submit"
        >
          <el-form-item label="登录名 / 手机号" prop="login">
            <el-input v-model="form.login" placeholder="请输入登录名或手机号" :prefix-icon="User" clearable />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
              @keyup.enter="submit"
            />
          </el-form-item>

          <el-button
            type="primary"
            class="login__submit"
            :loading="auth.loading"
            @click="submit"
          >
            登录
            <el-icon class="login__submit-icon"><ArrowRight /></el-icon>
          </el-button>
        </el-form>

        <el-alert
          v-if="auth.error"
          class="login__alert"
          type="error"
          :closable="false"
          show-icon
          :title="auth.error"
        >
          <template #default>
            <span class="login__alert-body">
              请确认后端服务已启动（<code>cd services/api &amp;&amp; make run</code>），并在
              <code>.env.local</code> 中配置好云 PostgreSQL 连接串。
            </span>
          </template>
        </el-alert>

        <div class="login__seed">
          <div class="login__seed-head">
            <el-icon><InfoFilled /></el-icon>
            <span>本地联调种子账号（点击即可填入）</span>
          </div>
          <div class="login__seed-list">
            <button
              v-for="acc in seedAccounts"
              :key="acc.login"
              type="button"
              class="login__seed-item"
              @click="useSeed(acc.login, acc.password)"
            >
              <b>{{ acc.role }}</b>
              <span>{{ acc.login }} / {{ acc.password }}</span>
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.login {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 1fr);
  min-height: 100vh;
}

.login__brand {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 44px 52px;
  background: #0f1a2e;
  color: #c7d2e4;
}

.login__headline h1 {
  margin: 46px 0 10px;
  font-size: 34px;
  line-height: 1.35;
  font-weight: 500;
  color: #fff;
}

.login__headline p {
  margin: 0;
  font-size: 14px;
  color: #8c9bb8;
}

.login__facts {
  list-style: none;
  margin: 40px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.login__facts li {
  display: flex;
  align-items: flex-start;
  gap: 9px;
  font-size: 13.5px;
  color: #aebbd0;
}

.login__facts .el-icon {
  margin-top: 4px;
  color: #2e6bff;
  flex: none;
}

.login__foot {
  font-size: 12px;
  color: #6b7b98;
}

.login__panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 32px;
  background: #fff;
}

.login__form-wrap {
  width: 100%;
  max-width: 380px;
}

.login__form-wrap h2 {
  margin: 0;
  font-size: 23px;
  font-weight: 500;
  color: var(--ym-ink);
}

.login__sub {
  margin: 8px 0 26px;
  font-size: 13.5px;
  color: var(--ym-ink-3);
}

.login__submit {
  width: 100%;
  margin-top: 6px;
}

.login__submit-icon {
  margin-left: 6px;
}

.login__alert {
  margin-top: 18px;
}

.login__alert-body {
  font-size: 12.5px;
  line-height: 1.7;
}

.login__alert-body code {
  padding: 1px 5px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.06);
  font-family: var(--ym-font-mono);
  font-size: 12px;
}

.login__seed {
  margin-top: 26px;
  padding-top: 20px;
  border-top: 1px dashed var(--ym-line);
}

.login__seed-head {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: var(--ym-ink-4);
  margin-bottom: 10px;
}

.login__seed-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.login__seed-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 7px 11px;
  border: 1px solid var(--ym-line);
  border-radius: 8px;
  background: #fbfcfe;
  font-family: inherit;
  font-size: 12.5px;
  color: var(--ym-ink-2);
  cursor: pointer;
  transition: all 0.15s ease;
}

.login__seed-item:hover {
  border-color: var(--ym-primary);
  background: var(--ym-primary-light);
}

.login__seed-item b {
  font-weight: 500;
  color: var(--ym-primary);
}

@media (max-width: 960px) {
  .login {
    grid-template-columns: 1fr;
  }

  .login__brand {
    display: none;
  }
}
</style>
