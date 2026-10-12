/**
 * 后端统一响应信封（见 services/api/pkg/response）。
 */
export interface ApiResponse<T = unknown> {
  /** 业务错误码，0 表示成功 */
  code: number
  message: string
  data?: T
  traceId?: string
}

/** 分页响应数据。 */
export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

/** 角色。Phase 2 升级为 RBAC 后此处会扩展为权限点集合。 */
export type UserRole = 'admin' | 'lecturer' | 'student'

/** 账号状态：1 正常，0 禁用。 */
export type UserStatus = 0 | 1

/** 用户资料（对应后端 user.Profile）。 */
export interface UserProfile {
  id: number
  username: string
  phone: string
  email: string
  nickname: string
  avatar: string
  role: UserRole
  status: UserStatus
  lastLoginAt?: string
  createdAt: string
}

/** 登录态（对应后端 auth.TokenResponse）。 */
export interface AuthTokens {
  tokenType: string
  accessToken: string
  expiresIn: number
  refreshToken: string
  user: UserProfile
}

/** 登录请求。 */
export interface LoginPayload {
  /** 登录名或手机号 */
  login: string
  password: string
}

/** 注册请求。 */
export interface RegisterPayload {
  username: string
  password: string
  phone?: string
  email?: string
  nickname?: string
}

/** 修改个人资料请求。 */
export interface UpdateProfilePayload {
  nickname?: string
  avatar?: string
  email?: string
  phone?: string
}

/** 用户列表查询参数。 */
export interface ListUsersQuery {
  page?: number
  pageSize?: number
  keyword?: string
  role?: UserRole
  status?: UserStatus
}

/** 健康检查结果。 */
export interface HealthStatus {
  status: 'ok' | 'degraded'
  service: string
  env: string
  version: string
  time: string
  uptime: string
  database: {
    status: 'up' | 'down'
    latency: string
    error?: string
  }
}

/** 模块元信息（对应 GET /api/v1/meta/modules）。 */
export interface ModuleMeta {
  service: string
  version: string
  env: string
  modules: string[]
  queue: boolean
  time: string
}
