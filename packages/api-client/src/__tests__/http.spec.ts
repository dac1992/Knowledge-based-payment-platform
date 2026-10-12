import type { AuthTokens, UserProfile } from '@yanmai/shared'
import type { AxiosAdapter, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { describe, expect, it, vi } from 'vitest'

import { ApiError, createApiClient } from '../http'

/* -------------------------------------------------------------------------- */
/* 测试脚手架：用自定义 axios adapter 精确编排响应，不依赖任何真实网络        */
/* -------------------------------------------------------------------------- */

interface Scenario {
  /** HTTP 状态码，默认 200 */
  status?: number
  /** 响应体 */
  body?: unknown
}

type ScenarioHandler = (
  config: InternalAxiosRequestConfig,
  index: number,
) => Scenario | Promise<Scenario>

/**
 * 创建脚本化 adapter。
 *
 * 注意：自定义 adapter 需要自己承担 axios 内建 adapter 的 `settle` 职责——
 * 非 2xx 必须 reject 一个带 `response` 字段的错误，否则 axios 会当成成功响应。
 */
function createAdapter(handler: ScenarioHandler) {
  const calls: InternalAxiosRequestConfig[] = []

  const adapter: AxiosAdapter = async (config) => {
    const index = calls.push(config) - 1
    const scenario = await handler(config, index)
    const status = scenario.status ?? 200

    const response = {
      data: scenario.body,
      status,
      statusText: `${status}`,
      headers: {},
      config,
    } as AxiosResponse

    if (status >= 200 && status < 300) return response

    const error = new Error(`Request failed with status code ${status}`) as Error & {
      response?: AxiosResponse
      config?: unknown
      isAxiosError?: boolean
    }
    error.response = response
    error.config = config
    error.isAxiosError = true
    return Promise.reject(error)
  }

  return { adapter, calls }
}

/** 成功信封。 */
function ok<T>(data: T, message = 'ok'): Scenario {
  return { status: 200, body: { code: 0, message, data, traceId: 'trace-ok' } }
}

/** 业务失败信封（默认 HTTP 200 + 非 0 业务码，与后端 Fail 行为一致）。 */
function fail(code: number, message: string, status = 200): Scenario {
  return { status, body: { code, message, data: null, traceId: 'trace-fail' } }
}

/** 读取请求头（AxiosHeaders 与普通对象都能取到）。 */
function readHeader(
  config: InternalAxiosRequestConfig | undefined,
  name: string,
): string | undefined {
  const headers = config?.headers as unknown as
    { get?: (key: string) => unknown; [key: string]: unknown } | undefined
  if (!headers) return undefined

  const viaGet = typeof headers.get === 'function' ? headers.get(name) : undefined
  const value = viaGet ?? headers[name]
  return typeof value === 'string' ? value : undefined
}

const profile: UserProfile = {
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

const refreshedTokens: AuthTokens = {
  tokenType: 'Bearer',
  accessToken: 'new-access',
  refreshToken: 'new-refresh',
  expiresIn: 7200,
  user: profile,
}

/* -------------------------------------------------------------------------- */

describe('ApiClient 响应信封拆解', () => {
  it('成功响应只把 data 交给调用方', async () => {
    const { adapter } = createAdapter(() => ok(profile))
    const client = createApiClient({ baseURL: '/api/v1', adapter })

    await expect(client.user.me()).resolves.toEqual(profile)
  })

  it('data 为空时返回 null 而不是整个信封', async () => {
    const { adapter } = createAdapter(() => ({
      status: 200,
      body: { code: 0, message: 'ok', data: null },
    }))
    const client = createApiClient({ baseURL: '/api/v1', adapter })

    await expect(client.user.me()).resolves.toBeNull()
  })

  it('业务码非 0 时抛出 ApiError 并保留 code / traceId / status', async () => {
    const { adapter } = createAdapter(() => fail(20003, '账号或密码不正确'))
    const client = createApiClient({ baseURL: '/api/v1', adapter })

    const error = await client.auth
      .login({ login: 'admin', password: 'wrong' })
      .then(() => null)
      .catch((err: unknown) => err)

    expect(error).toBeInstanceOf(ApiError)
    const apiError = error as ApiError
    expect(apiError.code).toBe(20003)
    expect(apiError.message).toBe('账号或密码不正确')
    expect(apiError.traceId).toBe('trace-fail')
    expect(apiError.status).toBe(200)
  })

  it('非信封结构（静态资源等）原样返回', async () => {
    const { adapter } = createAdapter(() => ({ status: 200, body: 'plain-text' }))
    const client = createApiClient({ baseURL: '/api/v1', adapter })

    await expect(client.user.me()).resolves.toBe('plain-text')
  })
})

describe('ApiClient 令牌注入', () => {
  it('已登录时自动携带 Authorization', async () => {
    const { adapter, calls } = createAdapter(() => ok(profile))
    const client = createApiClient({
      baseURL: '/api/v1',
      adapter,
      getAccessToken: () => 'access-1',
    })

    await client.user.me()

    expect(readHeader(calls[0], 'Authorization')).toBe('Bearer access-1')
  })

  it('未登录时不携带 Authorization', async () => {
    const { adapter, calls } = createAdapter(() => ok(profile))
    const client = createApiClient({ baseURL: '/api/v1', adapter })

    await client.user.me()

    expect(readHeader(calls[0], 'Authorization')).toBeUndefined()
  })
})

describe('ApiClient 401 自动续期', () => {
  it('令牌过期时续期一次并重放原请求', async () => {
    const onTokensRefreshed = vi.fn()
    // 模拟真实的会话存储：续期成功后 Access Token 会被写回，重放时才带得上新令牌。
    let accessToken = 'stale-access'

    const { adapter, calls } = createAdapter((config) => {
      if (config.url === '/auth/refresh') return ok(refreshedTokens)
      if (readHeader(config, 'Authorization') === 'Bearer new-access') return ok(profile)
      return fail(20007, '访问令牌已过期', 401)
    })

    const client = createApiClient({
      baseURL: '/api/v1',
      adapter,
      getAccessToken: () => accessToken,
      getRefreshToken: () => 'refresh-1',
      onTokensRefreshed: (tokens) => {
        onTokensRefreshed(tokens)
        accessToken = tokens.accessToken
      },
    })

    await expect(client.user.me()).resolves.toEqual(profile)

    // 1) 原始请求 2) 续期请求 3) 重放请求
    expect(calls.map((c) => c.url)).toEqual(['/users/me', '/auth/refresh', '/users/me'])
    expect(onTokensRefreshed).toHaveBeenCalledTimes(1)
    expect(onTokensRefreshed).toHaveBeenCalledWith(refreshedTokens)
    expect(readHeader(calls[2], 'Authorization')).toBe('Bearer new-access')
  })

  it('续期请求本身不带旧令牌，避免自我递归', async () => {
    let accessToken = 'stale-access'

    const { adapter, calls } = createAdapter((config) => {
      if (config.url === '/auth/refresh') return ok(refreshedTokens)
      if (readHeader(config, 'Authorization') === 'Bearer new-access') return ok(profile)
      return fail(20007, '访问令牌已过期', 401)
    })

    const client = createApiClient({
      baseURL: '/api/v1',
      adapter,
      getAccessToken: () => accessToken,
      getRefreshToken: () => 'refresh-1',
      onTokensRefreshed: (tokens) => {
        accessToken = tokens.accessToken
      },
    })

    await client.user.me()

    const refreshCall = calls.find((c) => c.url === '/auth/refresh')
    expect(readHeader(refreshCall, 'Authorization')).toBeUndefined()
  })

  it('并发请求同时过期时只续期一次', async () => {
    let refreshCount = 0
    let meCount = 0

    const { adapter } = createAdapter((config) => {
      if (config.url === '/auth/refresh') {
        refreshCount += 1
        return ok(refreshedTokens)
      }
      meCount += 1
      return meCount <= 2 ? fail(20007, '访问令牌已过期', 401) : ok(profile)
    })

    const client = createApiClient({
      baseURL: '/api/v1',
      adapter,
      getAccessToken: () => 'stale-access',
      getRefreshToken: () => 'refresh-1',
    })

    const [a, b] = await Promise.all([client.user.me(), client.user.me()])

    expect(a).toEqual(profile)
    expect(b).toEqual(profile)
    expect(refreshCount).toBe(1)
  })

  it('续期失败时回调 onUnauthorized 并抛出错误', async () => {
    const onUnauthorized = vi.fn()
    const { adapter } = createAdapter((config) =>
      config.url === '/auth/refresh'
        ? fail(20006, '刷新令牌无效', 401)
        : fail(20007, '访问令牌已过期', 401),
    )

    const client = createApiClient({
      baseURL: '/api/v1',
      adapter,
      getAccessToken: () => 'stale-access',
      getRefreshToken: () => 'refresh-1',
      onUnauthorized,
    })

    await expect(client.user.me()).rejects.toBeInstanceOf(ApiError)
    expect(onUnauthorized).toHaveBeenCalled()
  })

  it('没有 Refresh Token 时不尝试续期，直接回调 onUnauthorized', async () => {
    const onUnauthorized = vi.fn()
    const { adapter, calls } = createAdapter(() => fail(20007, '访问令牌已过期', 401))

    const client = createApiClient({
      baseURL: '/api/v1',
      adapter,
      getAccessToken: () => 'stale-access',
      onUnauthorized,
    })

    await expect(client.user.me()).rejects.toBeInstanceOf(ApiError)

    expect(calls.map((c) => c.url)).toEqual(['/users/me'])
    expect(onUnauthorized).toHaveBeenCalled()
  })
})

describe('ApiClient 网络层失败', () => {
  it('连接不上后端时抛出 code 0 的 ApiError', async () => {
    const adapter: AxiosAdapter = () =>
      Promise.reject(Object.assign(new Error('Network Error'), { code: 'ECONNREFUSED' }))

    const client = createApiClient({ baseURL: '/api/v1', adapter })

    const error = await client.user.me().then(
      () => null,
      (err: unknown) => err,
    )

    expect(error).toBeInstanceOf(ApiError)
    expect((error as ApiError).code).toBe(0)
    expect((error as ApiError).message).toContain('无法连接后端服务')
  })

  it('请求超时时给出可读提示', async () => {
    const adapter: AxiosAdapter = () =>
      Promise.reject(
        Object.assign(new Error('timeout of 15000ms exceeded'), { code: 'ECONNABORTED' }),
      )

    const client = createApiClient({ baseURL: '/api/v1', adapter })

    const error = (await client.user.me().then(
      () => null,
      (err: unknown) => err,
    )) as ApiError

    expect(error.code).toBe(0)
    expect(error.message).toContain('请求超时')
  })
})

describe('ApiError', () => {
  it('isAuthError 能识别令牌类错误', () => {
    expect(new ApiError({ code: 20007, message: '过期', status: 401 }).isAuthError).toBe(true)
    expect(new ApiError({ code: 0, message: '网络失败' }).isAuthError).toBe(false)
    expect(new ApiError({ code: 10000, message: '参数错误', status: 400 }).isAuthError).toBe(false)
  })
})
