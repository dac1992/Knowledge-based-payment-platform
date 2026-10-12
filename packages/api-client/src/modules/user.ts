import type { ListUsersQuery, PageResult, UpdateProfilePayload, UserProfile } from '@yanmai/shared'

import type { ApiClient } from '../http'

/** 用户模块接口。 */
export class UserApi {
  constructor(private readonly client: ApiClient) {}

  /** 当前登录用户资料。 */
  me(): Promise<UserProfile> {
    return this.client.request<UserProfile>({ url: '/users/me', method: 'GET' })
  }

  /** 修改当前登录用户资料。 */
  updateMe(payload: UpdateProfilePayload): Promise<UserProfile> {
    return this.client.request<UserProfile>({ url: '/users/me', method: 'PATCH', data: payload })
  }

  /** 用户列表（需 admin 角色）。 */
  list(query: ListUsersQuery = {}): Promise<PageResult<UserProfile>> {
    return this.client.request<PageResult<UserProfile>>({
      url: '/admin/users',
      method: 'GET',
      params: query,
    })
  }

  /** 用户详情（需 admin 角色）。 */
  detail(id: number): Promise<UserProfile> {
    return this.client.request<UserProfile>({ url: `/admin/users/${id}`, method: 'GET' })
  }
}
