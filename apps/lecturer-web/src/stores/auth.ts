import { ApiError, clearSession, getAccessToken, getProfile, saveProfile, saveSession } from '@yanmai/api-client'
import type { LoginPayload, UserProfile } from '@yanmai/shared'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { api } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | undefined>(getAccessToken())
  const profile = ref<UserProfile | null>(getProfile() ?? null)
  const loading = ref(false)
  const error = ref('')

  const isAuthenticated = computed(() => Boolean(accessToken.value && profile.value))
  const isLecturer = computed(() => profile.value?.role === 'lecturer' || profile.value?.role === 'admin')
  const displayName = computed(() => profile.value?.nickname || profile.value?.username || '')

  function applySession(token: string, user: UserProfile): void {
    accessToken.value = token
    profile.value = user
  }

  function reset(): void {
    accessToken.value = undefined
    profile.value = null
    error.value = ''
  }

  async function login(payload: LoginPayload): Promise<boolean> {
    loading.value = true
    error.value = ''
    try {
      const tokens = await api.auth.login(payload)
      saveSession(tokens)
      applySession(tokens.accessToken, tokens.user)
      return true
    } catch (err) {
      error.value = err instanceof ApiError ? err.message : '登录失败，请稍后重试'
      return false
    } finally {
      loading.value = false
    }
  }

  async function refreshProfile(): Promise<void> {
    const user = await api.user.me()
    profile.value = user
    saveProfile(user)
  }

  async function logout(): Promise<void> {
    try {
      await api.auth.logout()
    } catch {
      // 无状态登出
    }
    clearSession()
    reset()
  }

  return {
    accessToken,
    profile,
    loading,
    error,
    isAuthenticated,
    isLecturer,
    displayName,
    login,
    logout,
    refreshProfile,
    reset,
  }
})
