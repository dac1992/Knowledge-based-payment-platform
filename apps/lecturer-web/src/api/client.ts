import { createWebApiClient } from '@yanmai/api-client'

/** 全局唯一的接口客户端。 */
export const api = createWebApiClient({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
})
