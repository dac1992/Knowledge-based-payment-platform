import type { AuthTokens, LoginPayload, RegisterPayload, UserProfile } from '@yanmai/shared'

import type { ApiClient } from '../http'

/** 认证模块接口。 */
export class AuthApi {
  constructor(private readonly client: ApiClient) {}

  /** 注册（注册成功即返回登录态）。 */
  register(payload: RegisterPayload): Promise<AuthTokens> {
    return this.client.request<AuthTokens>({ url: '/auth/register', method: 'POST', data: payload })
  }

  /** 登录。 */
  login(payload: LoginPayload): Promise<AuthTokens> {
    return this.client.request<AuthTokens>({ url: '/auth/login', method: 'POST', data: payload })
  }

  /** 刷新令牌（内部使用 _skipAuth，避免续期请求自身再次触发续期）。 */
  refresh(refreshToken: string): Promise<AuthTokens> {
    return this.client.request<AuthTokens>({
      url: '/auth/refresh',
      method: 'POST',
      data: { refreshToken },
      _skipAuth: true,
    } as never)
  }

  /** 当前登录账号资料。 */
  me(): Promise<UserProfile> {
    return this.client.request<UserProfile>({ url: '/auth/me', method: 'GET' })
  }

  /** 登出（无状态，客户端需自行清除令牌）。 */
  logout(): Promise<{ revoked: boolean; hint: string }> {
    return this.client.request({ url: '/auth/logout', method: 'POST' })
  }
}
