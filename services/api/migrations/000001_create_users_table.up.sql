-- 000001 创建用户表
-- Phase 1 骨架版本：仅承载认证所需字段。
-- Phase 2 将在此基础上扩展 RBAC 角色权限、多端身份绑定（微信 UnionID / 小程序 openid）与手机号验证码登录。

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL,
    phone         VARCHAR(20)  NOT NULL DEFAULT '',
    email         VARCHAR(128) NOT NULL DEFAULT '',
    password_hash VARCHAR(255) NOT NULL,
    nickname      VARCHAR(64)  NOT NULL DEFAULT '',
    avatar        VARCHAR(512) NOT NULL DEFAULT '',
    role          VARCHAR(32)  NOT NULL DEFAULT 'student',
    status        SMALLINT     NOT NULL DEFAULT 1,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

COMMENT ON TABLE  users               IS '用户表';
COMMENT ON COLUMN users.username      IS '登录名，全局唯一';
COMMENT ON COLUMN users.phone         IS '手机号，Phase 2 用于短信登录';
COMMENT ON COLUMN users.email         IS '邮箱';
COMMENT ON COLUMN users.password_hash IS 'bcrypt 密码哈希，禁止存明文';
COMMENT ON COLUMN users.role          IS '角色：admin / lecturer / student，Phase 2 升级为 RBAC 模型';
COMMENT ON COLUMN users.status        IS '状态：1 正常，0 禁用';
COMMENT ON COLUMN users.last_login_at IS '最近登录时间';
COMMENT ON COLUMN users.deleted_at    IS '软删除时间，非空表示已删除';

-- 未删除记录的唯一性约束（部分索引），软删除后可重新注册同名账号
CREATE UNIQUE INDEX IF NOT EXISTS uk_users_username
    ON users (username) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uk_users_phone
    ON users (phone) WHERE deleted_at IS NULL AND phone <> '';

CREATE INDEX IF NOT EXISTS idx_users_role
    ON users (role);

CREATE INDEX IF NOT EXISTS idx_users_deleted
    ON users (deleted_at);
