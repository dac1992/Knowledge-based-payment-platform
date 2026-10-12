# 数据库初始化参考 SQL（腾讯云 PostgreSQL 控制台执行）

Phase 1 的表结构由 golang-migrate 在服务启动时自动创建（`DB_AUTO_MIGRATE=true`），
本文件只保留**建库与建账号**部分——这部分需要实例管理员权限，不适合放进迁移脚本。

## 1. 创建开发库与生产库

```sql
-- 本地/联调
CREATE DATABASE yanmai_dev  ENCODING 'UTF8' LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0;
-- 预发布
CREATE DATABASE yanmai_staging ENCODING 'UTF8' LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0;
-- 生产
CREATE DATABASE yanmai_prod ENCODING 'UTF8' LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0;
```

## 2. 创建业务账号（最小权限）

```sql
CREATE ROLE yanmai_app LOGIN PASSWORD '替换为强密码';

GRANT CONNECT ON DATABASE yanmai_dev TO yanmai_app;
GRANT CONNECT ON DATABASE yanmai_staging TO yanmai_app;
GRANT CONNECT ON DATABASE yanmai_prod TO yanmai_app;

-- 切到目标库后执行（示例：yanmai_dev）
\c yanmai_dev
GRANT USAGE, CREATE ON SCHEMA public TO yanmai_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO yanmai_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO yanmai_app;

-- 生产账号建议只给 DML，DDL 交由发布流程使用独立账号
-- \c yanmai_prod
-- GRANT USAGE ON SCHEMA public TO yanmai_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO yanmai_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO yanmai_app;
```

## 3. 连接串格式

```
postgres://yanmai_app:密码@主机:5432/yanmai_dev?sslmode=require
```

- 腾讯云 / 云托管实例通常强制 SSL，需保留 `sslmode=require`
- 阿里云 / Supabase 同理，按其控制台给出的 URI 为准
- 连接串只放 `services/api/.env.local`（已 gitignore），**不要**写进代码或提交到仓库

## 4. 白名单

云数据库默认只允许白名单 IP 访问。本地开发需要把**本机公网出口 IP** 加入白名单；
家庭宽带出口 IP 会变化，长期开发建议：
- 使用固定公网 IP 的跳板机，或
- 临时放行 `0.0.0.0/0`（仅联调期，完成后立即收回）
