import axios, {
  type AxiosAdapter,
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from 'axios'

import type { ApiResponse, AuthTokens } from '@yanmai/shared'

import { AuthApi } from './modules/auth'
import { SystemApi } from './modules/system'
import { UserApi } from './modules/user'

/** 业务错误（响应信封中 code !== 0，或网络/HTTP 层失败）。 */
export class ApiError extends Error {
  /** 业务错误码；网络层失败时为 0 */
  readonly code: number
  /** HTTP 状态码 */
  readonly status?: number
  /** 链路 ID，便于对照后端日志 */
  readonly traceId?: string
  /** 原始响应体 */
  readonly payload?: unknown

  constructor(params: {
    code: number
    message: string
    status?: number
    traceId?: string
    payload?: unknown
  }) {
    super(params.message)
    this.name = 'ApiError'
    this.code = params.code
    this.status = params.status
    this.traceId = params.traceId
    this.payload = params.payload
  }

  /** 是否为「登录态失效」类错误。 */
  get isAuthError(): boolean {
    return this.status === 401 || [20005, 20006, 20007, 20008, 10003].includes(this.code)
  }
}

/** 客户端构造参数。 */
export interface ApiClientOptions {
  /** 接口根地址，例如 http://127.0.0.1:8080/api/v1 */
  baseURL: string
  timeout?: number
  /** 读取当前 Access Token */
  getAccessToken?: () => string | undefined
  /** 读取当前 Refresh Token */
  getRefreshToken?: () => string | undefined
  /** 令牌刷新成功回调（用于写回本地存储） */
  onTokensRefreshed?: (tokens: AuthTokens) => void
  /** 刷新失败 / 未登录时回调（用于跳转登录页） */
  onUnauthorized?: (error: ApiError) => void
  /** 自定义适配器，便于测试注入 */
  adapter?: AxiosAdapter
}

type RetriableConfig = InternalAxiosRequestConfig & { _retried?: boolean; _skipAuth?: boolean }

/**
 * 研脉学术知识付费平台接口客户端。
 *
 * 统一职责：
 *   - 自动携带 Bearer Token
 *   - 拆解 `{ code, message, data, traceId }` 信封，只有 data 会返回给调用方
 *   - 401 时用 Refresh Token 自动续期并重放一次请求（并发请求共享同一次续期）
 *
 * 契约来源：后端 swaggo 注解生成的 OpenAPI（见 scripts/generate-openapi-types.mjs）。
 */
export class ApiClient {
  readonly auth: AuthApi
  readonly user: UserApi
  readonly system: SystemApi

  /** 底层 axios 实例，供业务侧绕过封装直接调用时使用 */
  readonly http: AxiosInstance

  private readonly options: ApiClientOptions
  private refreshing: Promise<string> | null = null

  constructor(options: ApiClientOptions) {
    this.options = options
    this.http = axios.create({
      baseURL: options.baseURL,
      timeout: options.timeout ?? 15000,
      adapter: options.adapter,
      headers: { 'Content-Type': 'application/json' },
    })

    this.http.interceptors.request.use(this.attachToken)
    this.http.interceptors.response.use(this.unwrap, this.handleError)

    this.auth = new AuthApi(this)
    this.user = new UserApi(this)
    this.system = new SystemApi(this)
  }

  /** 发起请求并返回已拆包的 data。 */
  async request<T>(config: AxiosRequestConfig): Promise<T> {
    // 响应拦截器已把 response.data 替换为信封里的业务数据
    const res = await this.http.request<ApiResponse<T>>(config)
    return res.data as unknown as T
  }

  /** 主动触发一次令牌续期（并发调用只会真正刷新一次）。 */
  refreshTokens(): Promise<string> {
    if (!this.refreshing) {
      this.refreshing = this.doRefresh().finally(() => {
        this.refreshing = null
      })
    }
    return this.refreshing
  }

  private doRefresh = async (): Promise<string> => {
    const refreshToken = this.options.getRefreshToken?.()
    if (!refreshToken) {
      throw new ApiError({ code: 20005, message: '本地没有 Refresh Token，需要重新登录' })
    }

    try {
      const tokens = await this.auth.refresh(refreshToken)
      this.options.onTokensRefreshed?.(tokens)
      return tokens.accessToken
    } catch (err) {
      const apiErr =
        err instanceof ApiError ? err : new ApiError({ code: 10006, message: '令牌续期失败' })
      this.options.onUnauthorized?.(apiErr)
      throw apiErr
    }
  }

  private attachToken = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token = this.options.getAccessToken?.()
    if (token && !(config as RetriableConfig)._skipAuth) {
      config.headers.set?.('Authorization', `Bearer ${token}`)
    }
    return config
  }

  /**
   * 拆解统一响应信封：code !== 0 抛 ApiError，否则把 data 写回 response.data。
   * 保留 AxiosResponse 形态是为了让 axios 的拦截器类型保持成立。
   */
  private unwrap = (response: AxiosResponse): AxiosResponse => {
    const body = response.data as ApiResponse<unknown> | undefined

    // 非信封结构（例如静态文件、第三方接口）原样返回
    if (!body || typeof body !== 'object' || !('code' in body)) {
      return response
    }

    if (body.code !== 0) {
      throw new ApiError({
        code: body.code,
        message: body.message || '请求失败',
        status: response.status,
        traceId: body.traceId,
        payload: body,
      })
    }

    response.data = body.data ?? null
    return response
  }

  private handleError = async (error: unknown): Promise<AxiosResponse> => {
    const err = error as {
      response?: { status?: number; data?: ApiResponse<unknown> }
      config?: RetriableConfig
      message?: string
      code?: string
    }

    // 网络层失败（超时 / 断网 / 跨域被拦）
    if (!err.response) {
      throw new ApiError({
        code: 0,
        message:
          err.code === 'ECONNABORTED'
            ? '请求超时，请检查后端服务是否已启动'
            : '无法连接后端服务，请确认 API 已启动且地址配置正确',
      })
    }

    const { status, data } = err.response
    const apiErr = new ApiError({
      code: data?.code ?? 0,
      message: data?.message || `请求失败（HTTP ${status}）`,
      status,
      traceId: data?.traceId,
      payload: data,
    })

    const config = err.config
    const canRetry =
      apiErr.isAuthError &&
      status === 401 &&
      config !== undefined &&
      !config._retried &&
      !config._skipAuth &&
      Boolean(this.options.getRefreshToken?.())

    if (!canRetry || !config) {
      if (apiErr.isAuthError) this.options.onUnauthorized?.(apiErr)
      throw apiErr
    }

    config._retried = true
    try {
      const token = await this.refreshTokens()
      config.headers?.set?.('Authorization', `Bearer ${token}`)
      return await this.http.request(config)
    } catch (refreshErr) {
      if (refreshErr instanceof ApiError) throw refreshErr
      throw apiErr
    }
  }
}

/** 创建接口客户端。 */
export function createApiClient(options: ApiClientOptions): ApiClient {
  return new ApiClient(options)
}
