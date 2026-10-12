/**
 * 本文件由脚本自动生成，请勿手动修改。
 *
 * 来源：services/api/docs/swagger.json
 * 生成：cd services/api && make sdk
 */

export interface paths {
    "/admin/users": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * 用户列表（分页）
         * @description 仅管理员可访问，支持关键字、角色、状态筛选
         */
        get: {
            parameters: {
                query?: {
                    /** @description 登录名 / 昵称 / 手机号 模糊匹配 */
                    keyword?: string;
                    /** @description 页码，从 1 开始 */
                    page?: number;
                    /** @description 每页条数，1-100 */
                    pageSize?: number;
                    /** @description 角色 */
                    role?: "admin" | "lecturer" | "student";
                    /** @description 状态 */
                    status?: 0 | 1;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["github_com_yanmai_api_pkg_httpx.PageEnvelope"];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/admin/users/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** 用户详情 */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description 用户 ID */
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["internal_modules_user.Profile"];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/auth/login": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * 登录
         * @description 支持使用登录名或手机号登录，返回 Access + Refresh 双 Token
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description 登录参数 */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["internal_modules_auth.LoginRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["internal_modules_auth.TokenResponse"];
                        };
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description 账号或密码不正确 */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description 账号被禁用 */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/auth/logout": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * 登出
         * @description Phase 1 为无状态登出，由客户端丢弃令牌；服务端黑名单在 Phase 2 引入 Redis 后实现
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/auth/me": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** 当前登录账号 */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/auth/refresh": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * 刷新令牌
         * @description 用 Refresh Token 换取新的令牌对；会回源校验账号仍存在且可用
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description 刷新参数 */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["internal_modules_auth.RefreshRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["internal_modules_auth.TokenResponse"];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/auth/register": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * 注册
         * @description 创建学员账号并直接返回登录态（Access + Refresh 双 Token）
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description 注册参数 */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["internal_modules_auth.RegisterRequest"];
                };
            };
            responses: {
                /** @description Created */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["internal_modules_auth.TokenResponse"];
                        };
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users/me": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** 当前用户资料 */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["internal_modules_user.Profile"];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /** 修改个人资料 */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description 待修改字段，未传的字段保持不变 */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["internal_modules_user.UpdateProfileRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"] & {
                            data?: components["schemas"]["internal_modules_user.Profile"];
                        };
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
                /** @description 手机号已被占用 */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["github_com_yanmai_api_pkg_httpx.Envelope"];
                    };
                };
            };
        };
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        "github_com_yanmai_api_internal_modules_user.Profile": {
            avatar?: string;
            createdAt?: string;
            email?: string;
            id?: number;
            lastLoginAt?: string;
            nickname?: string;
            phone?: string;
            role?: string;
            status?: number;
            username?: string;
        };
        "github_com_yanmai_api_pkg_httpx.Envelope": {
            code?: number;
            data?: unknown;
            message?: string;
            traceId?: string;
        };
        "github_com_yanmai_api_pkg_httpx.PageEnvelope": {
            list?: unknown;
            page?: number;
            pageSize?: number;
            total?: number;
        };
        "internal_modules_auth.LoginRequest": {
            login: string;
            password: string;
        };
        "internal_modules_auth.RefreshRequest": {
            refreshToken: string;
        };
        "internal_modules_auth.RegisterRequest": {
            email?: string;
            nickname?: string;
            password: string;
            phone?: string;
            username: string;
        };
        "internal_modules_auth.TokenResponse": {
            accessToken?: string;
            expiresIn?: number;
            refreshToken?: string;
            tokenType?: string;
            user?: components["schemas"]["github_com_yanmai_api_internal_modules_user.Profile"];
        };
        "internal_modules_user.Profile": {
            avatar?: string;
            createdAt?: string;
            email?: string;
            id?: number;
            lastLoginAt?: string;
            nickname?: string;
            phone?: string;
            role?: string;
            status?: number;
            username?: string;
        };
        "internal_modules_user.UpdateProfileRequest": {
            avatar?: string;
            email?: string;
            nickname?: string;
            phone?: string;
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export type operations = Record<string, never>;

