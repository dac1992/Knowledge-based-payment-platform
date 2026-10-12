# 研脉学术知识付费平台

面向学术科研人群（研究生 / 高校教师 / 科研院所研究员）的知识付费与在线学习平台。
核心业务形态：**学术课程 + 科研训练营 + 专家直播 + 学术社群**。

> 当前进度：**Phase 1 — 工程骨架（Walking Skeleton）已交付**。
> 详见 `docs/开发计划与技术选型.html` 与 `docs/Phase1-验收报告.html`。

---

## 一、仓库结构

```
研脉学术知识付费平台/
├── apps/                        各端前端应用
│   ├── admin-web/               管理后台      Vue3 + Vite + Element Plus
│   ├── lecturer-web/            讲师工作台    Vue3 + Vite + Element Plus
│   └── user-uniapp/             小程序/H5/APP uni-app (Vue3)
├── services/
│   └── api/                     Go + Gin 模块化单体
│       ├── cmd/server/          程序入口
│       ├── internal/
│       │   ├── config/          Viper 配置加载与校验
│       │   ├── database/        PostgreSQL 连接 / 迁移 / 健康探测
│       │   ├── middleware/      链路 ID / 日志 / 恢复 / CORS / 限流 / 鉴权
│       │   ├── modules/         12 个业务域（每域六件套分层）
│       │   └── router/          全局路由聚合
│       ├── migrations/          golang-migrate SQL（embed 进二进制）
│       ├── pkg/                 errs / httpx / response / jwt / logger / hash / queue
│       └── scripts/smoke.sh     端到端冒烟测试
├── packages/                    前端共享包
│   ├── ui/                      @yanmai/ui         设计系统组件库（Vue3）
│   ├── api-client/              @yanmai/api-client 类型安全接口 SDK
│   ├── shared/                  @yanmai/shared     类型 / 常量 / 工具
│   └── config/                  @yanmai/config     ESLint / Prettier / TS 共享配置
├── infra/                       nginx 模板 / 数据库初始化 SQL / 部署说明
├── docs/                        产品功能文档 / UI 设计规范 / 开发计划 / 验收报告
└── .github/workflows/           CI：api.yml（Go）/ web.yml（前端）
```

## 二、环境要求

| 组件 | 版本 | 本机实测 |
| --- | --- | --- |
| Go | 1.26+ | 1.26.3 ⚠️ `windows/386`（32 位，建议升级 64 位） |
| Node.js | 20+ | 22.22.2 ✓ |
| pnpm | 9+ | 10.30.3 ✓ |
| PostgreSQL | 16（云托管） | 云实例，本机无需安装 Docker |

## 三、快速开始

### 1. 准备数据库

需要一个可访问的 PostgreSQL 实例（腾讯云 / Supabase 均可）。建库与账号参考 `infra/db-init.sql`。

```bash
cd services/api
cp .env.example .env.local
# 编辑 .env.local，把 DATABASE_URL 换成你的真实连接串
# 格式：postgres://用户:密码@主机:5432/yanmai_dev?sslmode=require
```

> 未配置 `DATABASE_URL` 时服务会**明确报错退出**并列出缺失项，不会静默降级。

### 2. 启动后端

```bash
cd services/api
make run          # 首次启动自动建表 + 写入种子账号
```

验证：

```bash
curl http://127.0.0.1:8080/health          # 综合健康检查
curl http://127.0.0.1:8080/api/v1/meta/modules   # 模块清单
make smoke        # 端到端冒烟（登录 / 鉴权 / RBAC / 刷新令牌）
```

### 3. 启动前端

```bash
pnpm install      # 首次

pnpm dev:admin    # 管理后台   http://127.0.0.1:5173
pnpm dev:lecturer # 讲师工作台 http://127.0.0.1:5174
pnpm dev:mp       # 微信小程序（产物 dist/dev/mp-weixin，用微信开发者工具打开）
pnpm dev:h5       # uni-app H5
```

前端已配置代理，`/api` 与 `/health` 会转发到 `127.0.0.1:8080`，无需额外跨域配置。

### 4. 种子账号（`DB_SEED=true` 时写入）

| 登录名 | 密码 | 角色 |
| --- | --- | --- |
| admin | `Admin@123456` | 平台管理员 |
| lecturer | `Lecturer@123456` | 讲师 |
| student | `Student@123456` | 学员 |

> 仅用于本地联调。生产环境必须保持 `DB_SEED=false`。

## 四、常用命令

```bash
# 根目录
pnpm build        # turbo 全量构建（含类型检查）
pnpm lint         # 全量 ESLint
pnpm test         # 全量单元测试
pnpm typecheck    # 全量类型检查

# 后端
pnpm api:run      # 启动 API
pnpm api:test     # 测试 + 覆盖率
pnpm api:lint     # go vet + gofmt 校验
cd services/api && make help     # 查看全部 Makefile 目标
```

## 五、Phase 1 交付范围

已实现：

- **Monorepo**：pnpm workspace + Turborepo，apps / services / packages 就位
- **Go 后端骨架**：Gin 路由、12 个业务模块（auth / user 完整实现，其余为占位）、统一响应与错误码、zap 结构化日志、panic 恢复、链路 ID、CORS、按 IP 限流
- **认证闭环**：注册 / 登录 / JWT 双 Token / 刷新 / 鉴权中间件 / RBAC 角色守卫
- **数据库**：GORM + golang-migrate（SQL embed 进二进制）+ `users` 表迁移 + 可回滚 down 脚本 + 种子数据
- **管理后台**：登录页 → 工作台 → 用户管理（分页 / 筛选）→ 401 自动跳登录
- **讲师工作台**：登录页 → 工作台（品牌色区分，与后台共享组件）
- **uni-app**：可编译微信小程序与 H5，含首页与课程占位页
- **共享包**：`@yanmai/ui` / `@yanmai/api-client` / `@yanmai/shared` / `@yanmai/config`
- **工程规范**：ESLint 9 Flat Config + Prettier + TypeScript strict
- **CI**：Go 质量门禁 + **真实 PostgreSQL 的端到端冒烟**；前端 lint / test / build
- **文档**：本 README、`services/api/README.md`、`infra/README.md`、`docs/Phase1-验收报告.html`

后置到后续阶段：

- 微信登录 / 短信验证码登录、RBAC 细粒度权限点、令牌黑名单（Phase 2）
- swaggo 注解 → OpenAPI → 前端 SDK 自动生成（当前 SDK 为手写，路径与后端一一对应）
- Meilisearch 搜索、asynq 任务消费者（当前仅生产者侧，缺 `REDIS_ADDR` 时自动禁用）
- Husky + commitlint（需先 `git init` 后 `pnpm install` 才会激活）

## 六、工程约定

- **模块化单体**：模块只操作自己的表；跨模块调用必须通过对方 Service 接口（消费方定义接口，依赖倒置）
- **统一响应**：`{ code, message, data, traceId }`；业务错误码集中在 `pkg/errs`，HTTP 状态映射集中在 `pkg/httpx`
- **配置**：所有密钥只放环境变量（`.env.local`，已 gitignore）；日志中连接串自动脱敏
- **迁移**：纯 SQL、可回滚；本地启动自动执行，生产改由发布流程执行
- **提交规范**：建议使用 Conventional Commits

## 七、团队

研脉学术（北京）智能科技有限公司 · 内部资料
