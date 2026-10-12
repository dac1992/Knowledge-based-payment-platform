<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Lock, User } from '@element-plus/icons-vue'

import { YmLogo } from '@yanmai/ui'

import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const formRef = ref<FormInstance>()
const form = ref({ login: 'lecturer', password: 'Lecturer@123456' })

const rules: FormRules = {
  login: [{ required: true, message: '请输入登录名或手机号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  const ok = await auth.login({ ...form.value })
  if (!ok) {
    ElMessage.error(auth.error || '登录失败')
    return
  }

  ElMessage.success(`欢迎回来，${auth.displayName}`)
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/workspace'
  await router.replace(redirect)
}
</script>

<template>
  <div class="login">
    <div class="login__card">
      <YmLogo :size="40" subtitle="讲师工作台" />

      <h1>欢迎回来</h1>
      <p class="login__sub">登录后进入讲师工作台，管理课程、直播与教务</p>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        size="large"
        @submit.prevent="submit"
      >
        <el-form-item label="登录名 / 手机号" prop="login">
          <el-input
            v-model="form.login"
            placeholder="请输入登录名或手机号"
            :prefix-icon="User"
            clearable
          />
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

        <el-button type="primary" class="login__submit" :loading="auth.loading" @click="submit">
          登录
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
            请确认后端服务已启动，并在 <code>.env.local</code> 中配置好云 PostgreSQL 连接串。
          </span>
        </template>
      </el-alert>

      <p class="login__hint">联调账号：<code>lecturer / Lecturer@123456</code></p>
    </div>
  </div>
</template>

<style scoped>
.login {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #0f1a2e;
  padding: 32px;
}

.login__card {
  width: 100%;
  max-width: 400px;
  padding: 34px 34px 26px;
  border-radius: 14px;
  background: #fff;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.28);
}

.login__card h1 {
  margin: 26px 0 6px;
  font-size: 22px;
  font-weight: 500;
  color: var(--ym-ink);
}

.login__sub {
  margin: 0 0 24px;
  font-size: 13px;
  color: var(--ym-ink-3);
}

.login__submit {
  width: 100%;
  margin-top: 4px;
}

.login__alert {
  margin-top: 16px;
}

.login__alert-body {
  font-size: 12.5px;
  line-height: 1.7;
}

.login__hint {
  margin: 20px 0 0;
  font-size: 12px;
  color: var(--ym-ink-4);
  text-align: center;
}

.login__hint code,
.login__alert-body code {
  padding: 1px 5px;
  border-radius: 4px;
  background: var(--ym-line-light);
  font-family: var(--ym-font-mono);
}
</style>
