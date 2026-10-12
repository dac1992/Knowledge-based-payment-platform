import { createWebApiClient } from '@yanmai/api-client'

/**
 * 全局唯一的接口客户端。
 * 令牌读写与 401 处理由 @yanmai/api-client 的 Web 适配层统一实现。
 */
export const api = createWebApiClient({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
})
