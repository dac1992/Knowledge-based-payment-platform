import { STORAGE_KEYS, storage, type AuthTokens, type UserProfile } from '@yanmai/shared'

import { ApiClient } from './http'

/**
 * Web / 桌面端（localStorage 可用）的会话与客户端装配。
 *
 * 管理后台、讲师端、PC 端、桌面端都直接复用这里，避免各端重复实现令牌读写与 401 处理。
 */

/** 读取 Access Token。 */
export function getAccessToken(): string | undefined {
  return storage.get<string>(STORAGE_KEYS.accessToken)
}

/** 读取 Refresh Token。 */
export function getRefreshToken(): string | undefined {
  return storage.get<string>(STORAGE_KEYS.refreshToken)
}

/** 读取本地缓存的用户资料。 */
export function getProfile(): UserProfile | undefined {
  return storage.get<UserProfile>(STORAGE_KEYS.profile)
}

/** 持久化登录态。 */
export function saveSession(tokens: AuthTokens): void {
  storage.set(STORAGE_KEYS.accessToken, tokens.accessToken)
  storage.set(STORAGE_KEYS.refreshToken, tokens.refreshToken)
  storage.set(STORAGE_KEYS.profile, tokens.user)
}

/** 仅更新本地缓存的资料（令牌不变）。 */
export function saveProfile(profile: UserProfile): void {
  storage.set(STORAGE_KEYS.profile, profile)
}

/** 清除登录态。 */
export function clearSession(): void {
  storage.remove(STORAGE_KEYS.accessToken)
  storage.remove(STORAGE_KEYS.refreshToken)
  storage.remove(STORAGE_KEYS.profile)
}

type UnauthorizedHandler = (error?: unknown) => void

let unauthorizedHandler: UnauthorizedHandler | null = null

/** 注册「登录态彻底失效」处理器（通常绑定到路由跳转登录页）。 */
export function setUnauthorizedHandler(handler: UnauthorizedHandler | null): void {
  unauthorizedHandler = handler
}

/** 触发「登录态彻底失效」处理。 */
export function emitUnauthorized(error?: unknown): void {
  unauthorizedHandler?.(error)
}

/** createWebApiClient 参数。 */
export interface CreateWebApiClientOptions {
  /** 接口根地址，例如 /api/v1 */
  baseURL: string
  timeout?: number
}

/**
 * 创建已接好会话管理的接口客户端：
 * 自动携带令牌、401 自动续期、续期失败清空会话并回调。
 */
export function createWebApiClient(options: CreateWebApiClientOptions): ApiClient {
  return new ApiClient({
    baseURL: options.baseURL,
    timeout: options.timeout ?? 15000,
    getAccessToken,
    getRefreshToken,
    onTokensRefreshed: (tokens) => saveSession(tokens),
    onUnauthorized: (error) => {
      clearSession()
      emitUnauthorized(error)
    },
  })
}
