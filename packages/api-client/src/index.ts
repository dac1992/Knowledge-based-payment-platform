/**
 * @yanmai/api-client —— 类型安全的接口 SDK。
 *
 * ## 契约来源
 * 后端通过 swaggo 注解产出 OpenAPI 3 文档（services/api/docs/swagger.json），
 * 再由 `pnpm --filter @yanmai/api-client gen:openapi` 生成 TS 类型。
 *
 * Phase 1 的接口集很小且稳定，SDK 采用手写封装（路径与后端路由一一对应），
 * 同时保留 OpenAPI 生成脚本；Phase 2 接口规模上来后切换为「注解 → 生成」流程。
 *
 * ## 使用方式
 * Web / 桌面端（localStorage 可用）：
 * ```ts
 * import { createWebApiClient, setUnauthorizedHandler } from '@yanmai/api-client'
 * export const api = createWebApiClient({ baseURL: '/api/v1' })
 * ```
 * 小程序端因缺少 XMLHttpRequest，使用 uni.request 自实现传输层。
 */
export { ApiClient, ApiError, createApiClient } from './http'
export type { ApiClientOptions } from './http'

export {
  clearSession,
  createWebApiClient,
  emitUnauthorized,
  getAccessToken,
  getProfile,
  getRefreshToken,
  saveProfile,
  saveSession,
  setUnauthorizedHandler,
} from './web'
export type { CreateWebApiClientOptions } from './web'

export { AuthApi } from './modules/auth'
export { UserApi } from './modules/user'
export { SystemApi } from './modules/system'
