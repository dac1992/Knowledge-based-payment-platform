import type { HealthStatus, ModuleMeta } from '@yanmai/shared'

import type { ApiClient } from '../http'

/**
 * 系统与元信息接口。
 *
 * 注意：健康检查与模块元信息挂在服务根路径（不带 /api/v1 前缀），
 * 因此这里需要显式覆盖 baseURL。
 */
export class SystemApi {
  private readonly root: string

  constructor(private readonly client: ApiClient) {
    this.root = client.http.defaults.baseURL?.replace(/\/api\/v1\/?$/, '') ?? ''
  }

  /** 综合健康检查。 */
  health(): Promise<HealthStatus> {
    return this.client.request<HealthStatus>({ url: '/health', baseURL: this.root, method: 'GET' })
  }

  /** 已挂载的业务模块清单。 */
  modules(): Promise<ModuleMeta> {
    return this.client.request<ModuleMeta>({ url: '/meta/modules', method: 'GET' })
  }
}
