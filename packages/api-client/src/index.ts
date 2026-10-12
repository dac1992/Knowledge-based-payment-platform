/**
 * @yanmai/api-client —— 类型安全的接口 SDK。
 *
 * ## 契约来源
 * 后端 swaggo 注解 → `services/api/docs/swagger.json`（Swagger 2.0）
 * → swagger2openapi 转 OpenAPI 3.0 → openapi-typescript 生成 `src/generated/schema.d.ts`。
 * 一条命令完成：`cd services/api && make sdk`。
 *
 * 请求封装仍为手写（路径与后端路由一一对应，便于加拦截与类型推断），
 * 但 `src/contract.ts` 会用生成的契约对 `@yanmai/shared` 的类型做编译期断言：
 * 后端一旦删除或改名字段，前端在 `pnpm typecheck` 阶段即报错。
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

// 契约类型与一致性断言（断言在 typecheck 时生效，见 contract.ts）
export type * from './contract'
