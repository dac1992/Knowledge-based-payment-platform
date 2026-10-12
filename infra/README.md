# infra 部署基础设施

| 文件 | 用途 |
| --- | --- |
| `nginx.conf` | 生产 Nginx 模板：HTTPS 终止、静态资源托管、把 `/api` 与 `/health` 反代到 Go 服务 |
| `db-init.sql` | 云 PostgreSQL 建库建账号参考 SQL 与连接串规范 |

## 环境划分

| 环境 | 数据库 | 说明 |
| --- | --- | --- |
| local | 云开发库 `yanmai_dev` | 本机无 Docker，直连云开发库；配置在 `services/api/.env.local`（不入 Git） |
| dev | 云开发库 | 前后端联调 |
| staging | 云预发布库 `yanmai_staging` | 与生产同构，用于验收 |
| prod | 云生产库 `yanmai_prod` | 主从 + 自动备份（保留 30 天）+ 白名单访问 |

## 部署顺序（Phase 14 使用）

1. 建库建账号：执行 `db-init.sql`
2. 配置后端：写入 `services/api/.env`（`DB_AUTO_MIGRATE=false`，迁移改由发布流程执行）
3. 发布后端：`cd services/api && make build`，把 `bin/yanmai-api.exe`（Linux 下为 `yanmai-api`）上传并托管为系统服务
4. 构建前端：`pnpm build`，把 `apps/*/dist` 上传到 `__WEB_ROOT__`
5. 配置 Nginx：替换 `nginx.conf` 中的占位符并 reload

## 安全约定

- 所有密钥（数据库密码、JWT Secret、支付密钥）只存环境变量，**绝不入库、不进 Git**
- 生产库开启自动备份与白名单访问
- JWT Secret 用 `openssl rand -hex 32` 生成
