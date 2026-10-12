/**
 * 契约一致性校验层。
 *
 * `./generated/schema.d.ts` 由后端 OpenAPI 文档生成，是接口的「真相」；
 * `@yanmai/shared` 中的类型则是前端消费时使用的形态。两者一旦漂移，
 * 通常只会在运行时才暴露。本文件把漂移提前到编译期。
 *
 * ## 断言方向
 * 检查「前端声明的字段」是否仍全部存在于契约中：
 * - 后端**删除或改名**字段 → 编译立即失败（最危险的不兼容变更）✅
 * - 后端**新增**字段 → 不报错（兼容变更，前端按需取用）✅
 *
 * 之所以不做反向断言（要求契约字段全部出现在前端类型中），是因为 swag
 * 未标注 required 的字段在契约里一律是可选的，反向断言会产生大量无意义告警。
 */
import type {
  AuthTokens,
  ListUsersQuery,
  LoginPayload,
  RegisterPayload,
  UpdateProfilePayload,
  UserProfile,
} from '@yanmai/shared'

import type { components } from './generated/schema'

/** OpenAPI 的组件集合。 */
export type Schemas = components['schemas']

/**
 * swag 生成的模型名 → 前端可读别名。
 *
 * 键名规则：Go 包路径以 `_` 连接，跨包引用带完整模块前缀。后端若调整包结构
 * 导致键名变化，只需修改此处，其余业务代码不受影响。
 */
export type ContractUserProfile = Schemas['internal_modules_user.Profile']
export type ContractTokenResponse = Schemas['internal_modules_auth.TokenResponse']
export type ContractLoginRequest = Schemas['internal_modules_auth.LoginRequest']
export type ContractRegisterRequest = Schemas['internal_modules_auth.RegisterRequest']
export type ContractUpdateProfileRequest = Schemas['internal_modules_user.UpdateProfileRequest']
export type ContractEnvelope = Schemas['github_com_yanmai_api_pkg_httpx.Envelope']
export type ContractPageEnvelope = Schemas['github_com_yanmai_api_pkg_httpx.PageEnvelope']

/** 分页查询参数在契约中被内联进 paths，此处保留前端类型作为参照。 */
export type ContractListUsersQuery = ListUsersQuery

// ---------------------------------------------------------------------------
// 类型级断言
// ---------------------------------------------------------------------------

/** 仅当 T 为 true 时通过编译。 */
type Expect<T extends true> = T

/** 不做联合类型分发：用元组包裹两侧。 */
type Extends<A, B> = [A] extends [B] ? true : false

/**
 * 编译期断言。以下每一项都必须为 true，否则 `tsc` 直接报错。
 *
 * 为什么写成 `declare const` 而不是纯类型别名：TypeScript 对**未被引用**的
 * 类型别名会跳过泛型约束求值——实测把契约里的字段删掉后 `pnpm typecheck`
 * 依然通过，断言形同虚设。改成值级声明后约束会在实例化时强制求值。
 * `declare` 只产生类型，不产生任何运行时代码。
 *
 * 断言失败的处理方式：同步修改 `@yanmai/shared` 中对应的类型定义，并检查所有消费方。
 */
export declare const contractAsserts: [
  Expect<Extends<keyof UserProfile, keyof ContractUserProfile>>,
  Expect<Extends<keyof AuthTokens, keyof ContractTokenResponse>>,
  Expect<Extends<keyof LoginPayload, keyof ContractLoginRequest>>,
  Expect<Extends<keyof RegisterPayload, keyof ContractRegisterRequest>>,
  Expect<Extends<keyof UpdateProfilePayload, keyof ContractUpdateProfileRequest>>,
]
