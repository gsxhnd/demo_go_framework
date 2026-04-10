-- =============================================================================
-- PostgreSQL 初始化脚本
-- =============================================================================

-- 创建用户
DO
$do$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE  rolname = 'demo_user') THEN
        CREATE ROLE demo_user WITH LOGIN PASSWORD 'demo_password';
    END IF;
END
$do$;

-- 创建数据库
SELECT 'CREATE DATABASE demo_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'demo_db')\gexec

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE demo_db TO demo_user;

-- 连接到数据库并创建表
\c demo_db

-- =============================================================================
-- 应用表结构（示例）
-- =============================================================================

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar VARCHAR(500),
    phone VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT users_email_key UNIQUE (email)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- 插入测试用户（密码为 'password123' 的 bcrypt hash）
INSERT INTO users (email, password, nickname, is_active)
VALUES
    ('admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin', TRUE),
    ('user@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User', TRUE)
ON CONFLICT (email) DO NOTHING;
