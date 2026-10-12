import { STORAGE_KEYS, storage, type ApiResponse, type AuthTokens, type HealthStatus, type ModuleMeta, type UserProfile } from '@yanmai/shared'

/**
 * uni-app 请求封装。
 *
 * 为什么不直接用 @yanmai/api-client？
 * 小程序环境没有 XMLHttpRequest，axios 不可用。这里用 uni.request 实现
 * 与 Web 端完全一致的「信封拆解 + 业务错误码 + Bearer 令牌」语义。
 * Phase 2 会把传输层抽象为可注入适配器，让 api-client 在 Web 与小程序间共用。
 */

const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

/** 业务错误。 */
export class UniApiError extends Error {
  readonly code: number
  readonly traceId?: string

  constructor(params: { code: number; message: string; traceId?: string }) {
    super(params.message)
    this.name = 'UniApiError'
    this.code = params.code
    this.traceId = params.traceId
  }
}

/** 我们对外暴露的 HTTP 方法集合。 */
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

/**
 * uni 的类型定义只收录了 GET/POST/PUT/DELETE/OPTIONS/HEAD/TRACE/CONNECT，
 * 但小程序运行时（wx.request）实际支持任意方法名，包括 PATCH。
 * 这里在类型边界做一次收窄断言，避免我们自己的 API 面被第三方类型定义限制住。
 */
type UniRequestMethod = Parameters<typeof uni.request>[0]['method']

const toUniMethod = (method: HttpMethod): UniRequestMethod => method as UniRequestMethod

export interface RequestOptions {
  url: string
  method?: HttpMethod
  data?: unknown
  /** 是否携带 Access Token，默认 true */
  auth?: boolean
  /** 覆盖根地址（健康检查挂在服务根路径） */
  baseURL?: string
}

/** 发起请求并返回已拆包的业务数据。 */
export function request<T>(options: RequestOptions): Promise<T> {
  const { url, method = 'GET', data, auth = true, baseURL } = options

  const header: Record<string, string> = { 'Content-Type': 'application/json' }
  if (auth) {
    const token = storage.get<string>(STORAGE_KEYS.accessToken)
    if (token) header.Authorization = `Bearer ${token}`
  }

  return new Promise<T>((resolve, reject) => {
    uni.request({
      url: `${baseURL ?? BASE_URL}${url}`,
      method: toUniMethod(method),
      data: data as string | Record<string, unknown> | ArrayBuffer | undefined,
      header,
      timeout: 15000,
      success: (res) => {
        const body = res.data as ApiResponse<T> | undefined
        if (!body || typeof body !== 'object' || !('code' in body)) {
          resolve(body as unknown as T)
          return
        }
        if (body.code !== 0) {
          reject(new UniApiError({ code: body.code, message: body.message || '请求失败', traceId: body.traceId }))
          return
        }
        resolve(body.data as T)
      },
      fail: (err) => {
        reject(
          new UniApiError({
            code: 0,
            message: err.errMsg?.includes('timeout') ? '请求超时，请检查网络' : '无法连接后端服务',
          }),
        )
      },
    })
  })
}

/** 与后端接口一一对应的薄封装。 */
export const api = {
  health: () => request<HealthStatus>({ url: '/health', baseURL: '', auth: false }),
  modules: () => request<ModuleMeta>({ url: '/meta/modules', auth: false }),
  login: (payload: { login: string; password: string }) =>
    request<AuthTokens>({ url: '/auth/login', method: 'POST', data: payload, auth: false }),
  profile: () => request<UserProfile>({ url: '/users/me' }),
}
