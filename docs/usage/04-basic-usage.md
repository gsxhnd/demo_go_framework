# 基础使用

## 概述

本文档介绍 go_sample_code 的核心 API 使用方法。当前已实现用户管理 CRUD 接口和健康检查接口。

> **注意**：用户管理路由当前未在 `RegisterHooks` 中注册，需手动添加后方可使用。仅 `/api/health` 路由默认可用。

## 健康检查

### 检查服务状态

```bash
curl http://localhost:8080/api/health
```

响应：

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "status": "healthy",
    "timestamp": "2026-05-02T10:00:00Z"
  }
}
```

## 用户管理

所有用户管理 API 的基础路径为 `/api/users`。所有响应格式统一：

```json
{
  "code": 0,
  "message": "OK",
  "data": { ... }
}
```

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123",
    "nickname": "Alice"
  }'
```

响应 (201)：

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "nickname": "Alice",
    "avatar": "",
    "phone": "",
    "is_active": true,
    "created_at": "2026-05-02T10:00:00Z",
    "updated_at": "2026-05-02T10:00:00Z"
  }
}
```

### 查询用户

**按 ID 查询**：

```bash
curl http://localhost:8080/api/users/1
```

**按用户名查询**：

```bash
curl http://localhost:8080/api/users/username/alice
```

**按邮箱查询**：

```bash
curl http://localhost:8080/api/users/email/alice@example.com
```

### 更新用户

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "Alice Updated",
    "is_active": false
  }'
```

### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/users/1
```

响应：

```json
{
  "code": 0,
  "message": "OK",
  "data": null
}
```

### 分页列表

```bash
curl "http://localhost:8080/api/users?page=1&page_size=10"
```

响应：

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "list": [ ... ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

## 错误响应

所有错误响应格式统一：

```json
{
  "code": 2001,
  "message": "User not found",
  "data": null
}
```

常见错误码：

| HTTP 状态码 | code | 说明 |
|-------------|------|------|
| 400 | 1003 | 参数校验失败 |
| 404 | 2001 | 用户不存在 |
| 409 | 2002 | 用户已存在 |
| 429 | 1104 | 超出限流 |
| 500 | 1000 | 服务器内部错误 |

完整错误码参考 → [API 参考](./07-api-reference.md)
