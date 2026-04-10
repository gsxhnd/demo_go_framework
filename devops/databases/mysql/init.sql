-- =============================================================================
-- MySQL 初始化脚本
-- =============================================================================

-- 创建应用数据库
CREATE DATABASE IF NOT EXISTS demo_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'demo_user'@'%';

-- 设置密码
ALTER USER 'demo_user'@'%' IDENTIFIED BY 'demo_password';

-- 授予权限
GRANT ALL PRIVILEGES ON demo_db.* TO 'demo_user'@'%';

-- 刷新权限
FLUSH PRIVILEGES;

-- 使用数据库
USE demo_db;

-- =============================================================================
-- 应用表结构（示例）
-- =============================================================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar VARCHAR(500),
    phone VARCHAR(20),
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试用户（密码为 'password123' 的 bcrypt hash）
INSERT INTO users (email, password, nickname, is_active) VALUES
    ('admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin', 1),
    ('user@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User', 1)
ON DUPLICATE KEY UPDATE email = email;
