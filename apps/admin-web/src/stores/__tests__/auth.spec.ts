import { ApiError } from '@yanmai/api-client'
import type { AuthTokens, UserProfile } from '@yanmai/shared'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { api } from '@/api/client'

import { useAuthStore } from '../auth'

vi.mock('@/api/client', () => ({
  api: {
    auth: { login: vi.fn(), logout: vi.fn() },
    user: { me: vi.fn() },
  },
}))

const adminProfile: UserProfile = {
  id: 1,
  username: 'admin',
  phone: '13800000001',
  email: '',
  nickname: '平台管理员',
  avatar: '',
  role: 'admin',
  status: 1,
  createdAt: '2026-10-12T10:00:00Z',
}

const tokens: AuthTokens = {
  tokenType: 'Bearer',
  accessToken: 'access-token-value',
  refreshToken: 'refresh-token-value',
  expiresIn: 7200,
  user: adminProfile,
}

beforeEach(() => {
  setActivePinia(createPinia())
  localStorage.clear()
  vi.clearAllMocks()
})

describe('useAuthStore', () => {
  it('初始状态为未登录', () => {
    const store = useAuthStore()

    expect(store.isAuthenticated).toBe(false)
    expect(store.isAdmin).toBe(false)
    expect(store.displayName).toBe('')
  })

  it('登录成功后写入状态并持久化令牌', async () => {
    vi.mocked(api.auth.login).mockResolvedValue(tokens)
    const store = useAuthStore()

    await expect(store.login({ login: 'admin', password: 'Admin@123456' })).resolves.toBe(true)

    expect(store.isAuthenticated).toBe(true)
    expect(store.isAdmin).toBe(true)
    expect(store.displayName).toBe('平台管理员')
    expect(store.error).toBe('')

    expect(localStorage.getItem('yanmai:accessToken')).toBe('"access-token-value"')
    expect(localStorage.getItem('yanmai:refreshToken')).toBe('"refresh-token-value"')
    expect(JSON.parse(localStorage.getItem('yanmai:profile') as string).username).toBe('admin')
  })

  it('登录失败时记录错误且不写入登录态', async () => {
    vi.mocked(api.auth.login).mockRejectedValue(
      new ApiError({ code: 20003, message: '账号或密码不正确' }),
    )
    const store = useAuthStore()

    await expect(store.login({ login: 'admin', password: 'wrong' })).resolves.toBe(false)

    expect(store.isAuthenticated).toBe(false)
    expect(store.error).toBe('账号或密码不正确')
    expect(localStorage.getItem('yanmai:accessToken')).toBeNull()
  })

  it('登录过程中的 loading 标记会被复位', async () => {
    vi.mocked(api.auth.login).mockRejectedValue(
      new ApiError({ code: 0, message: '无法连接后端服务' }),
    )
    const store = useAuthStore()

    await store.login({ login: 'admin', password: 'x' })

    expect(store.loading).toBe(false)
  })

  it('登出会清理本地状态，即便后端调用失败', async () => {
    vi.mocked(api.auth.login).mockResolvedValue(tokens)
    vi.mocked(api.auth.logout).mockRejectedValue(new ApiError({ code: 0, message: '服务不可用' }))

    const store = useAuthStore()
    await store.login({ login: 'admin', password: 'Admin@123456' })
    expect(store.isAuthenticated).toBe(true)

    await store.logout()

    expect(store.isAuthenticated).toBe(false)
    expect(store.profile).toBeNull()
    expect(localStorage.getItem('yanmai:accessToken')).toBeNull()
  })

  it('仅 admin 角色被识别为管理员', async () => {
    const lecturerTokens: AuthTokens = {
      ...tokens,
      user: { ...adminProfile, id: 2, username: 'lecturer', role: 'lecturer' },
    }
    vi.mocked(api.auth.login).mockResolvedValue(lecturerTokens)

    const store = useAuthStore()
    await store.login({ login: 'lecturer', password: 'Lecturer@123456' })

    expect(store.isAuthenticated).toBe(true)
    expect(store.isAdmin).toBe(false)
    expect(store.role).toBe('lecturer')
  })

  it('refreshProfile 会回写本地缓存', async () => {
    vi.mocked(api.auth.login).mockResolvedValue(tokens)
    vi.mocked(api.user.me).mockResolvedValue({ ...adminProfile, nickname: '新昵称' })

    const store = useAuthStore()
    await store.login({ login: 'admin', password: 'Admin@123456' })
    await store.refreshProfile()

    expect(store.displayName).toBe('新昵称')
    expect(JSON.parse(localStorage.getItem('yanmai:profile') as string).nickname).toBe('新昵称')
  })
})
