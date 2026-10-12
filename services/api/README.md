# 研脉学术知识付费平台 · API 服务（services/api）

Go 1.26 + Gin + GORM + golang-migrate 的模块化单体后端。

## 快速开始

```bash
# 1. 配置数据库连接（必需）
cp .env.example .env.local
#    编辑 .env.local，把 DATABASE_URL 换成云 PostgreSQL 连接串

# 2. 启动（首次启动会自动建表 + 写入种子账号）
go run ./cmd/server

# 3. 验证
curl http://127.0.0.1:8080/health
```

未配置 `DATABASE_URL` 时服务会**明确报错退出**并列出缺失项，不会静默降级。

## 目录结构

```
services/api/
├── cmd/server/main.go          程序入口（迁移 → 连接 → 种子 → 路由 → 优雅关闭）
├── internal/
│   ├── config/                 Viper 配置加载与校验
│   ├── database/               PostgreSQL 连接、迁移、健康探测、GORM 日志桥接
│   ├── middleware/             链路 ID / 访问日志 / panic 恢复 / CORS / 限流 / JWT 鉴权
│   ├── modules/                业务模块（每域一个包，六件套分层）
│   │   ├── module.go           Module 契约与 Registry
│   │   ├── placeholder.go      未实现模块的统一占位实现
│   │   ├── auth/               认证（已实现）
│   │   ├── user/               用户（已实现）
│   │   └── course|order|...    10 个占位模块
│   └── router/                 全局路由聚合与健康检查
├── migrations/                 golang-migrate SQL 脚本（go:embed 进二进制）
└── pkg/                        公共库：errs / httpx / response / jwt / logger / hash / queue
```

## 架构约定

- **模块化单体**：一个二进制，按业务域严格分包；模块只操作自己的表，跨模块通过 Service 接口协作（消费方定义接口，依赖倒置）。
- **分层**：每个模块内部 `handler → service → repository`，`dto` 承载请求/响应，`model` 承载 GORM 结构，`router` 负责路由声明。
- **统一响应**：`{ code, message, data, traceId }`；错误码集中在 `pkg/errs`，HTTP 状态码映射集中在 `pkg/httpx`，业务层不感知 HTTP。
- **迁移**：纯 SQL、可回滚、embed 进二进制；启动时按 `DB_AUTO_MIGRATE` 决定是否执行。

## 接口一览（Phase 1）

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/health` | 综合健康检查（含数据库探测） | — |
| GET | `/health/live` | 存活探针 | — |
| GET | `/health/ready` | 就绪探针（数据库故障返回 503） | — |
| GET | `/api/v1/meta/modules` | 已挂载模块清单 | — |
| POST | `/api/v1/auth/register` | 注册 | — |
| POST | `/api/v1/auth/login` | 登录 | — |
| POST | `/api/v1/auth/refresh` | 刷新令牌 | — |
| GET | `/api/v1/auth/me` | 当前登录账号 | Bearer |
| POST | `/api/v1/auth/logout` | 登出（客户端丢弃令牌） | Bearer |
| GET | `/api/v1/users/me` | 个人资料 | Bearer |
| PATCH | `/api/v1/users/me` | 修改个人资料 | Bearer |
| GET | `/api/v1/admin/users` | 用户列表（分页） | Bearer + admin |
| GET | `/api/v1/admin/users/:id` | 用户详情 | Bearer + admin |
| GET | `/api/v1/<module>` | 未交付模块占位，返回 501 | — |

## 种子账号（`DB_SEED=true` 时写入）

| 登录名 | 密码 | 角色 |
| --- | --- | --- |
| admin | Admin@123456 | 管理员 |
| lecturer | Lecturer@123456 | 讲师 |
| student | Student@123456 | 学员 |

## 常用命令

```bash
make run      # 启动服务
make test     # 测试 + 覆盖率
make lint     # go vet + gofmt 校验
make verify   # lint + test
make build    # 编译到 bin/
```

## 环境变量

见 `.env.example`。所有密钥只存环境变量，绝不入库、不进 Git。
