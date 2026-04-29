# API 接口文档

## 概述

| 基础 URL | 说明 |
|----------|------|
| `/api` | API 前缀 |

## 健康检查

### GET /api/health

健康检查接口。

**请求**

```
GET /api/health
```

**响应**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "status": "healthy",
    "timestamp": "2026-04-29T10:00:00Z"
  }
}
```

---

## 用户管理

### POST /api/users

创建用户。

**请求**

```json
{
  "username": "string (required, 1-50 chars)",
  "email": "string (required, valid email, 1-255 chars)",
  "password": "string (required, min 6 chars)",
  "nickname": "string (optional, max 100 chars)",
  "avatar": "string (optional, max 500 chars)",
  "phone": "string (optional, max 20 chars)"
}
```

**响应 (201)**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "id": 1,
    "username": "test",
    "email": "test@example.com",
    "nickname": "Test User",
    "avatar": "",
    "phone": "",
    "is_active": true,
    "created_at": "2026-04-29T10:00:00Z",
    "updated_at": "2026-04-29T10:00:00Z"
  }
}
```

**错误码**

| code | 说明 |
|------|------|
| 1001 | 参数解析错误 |
| 1002 | 参数校验失败 |
| 2001 | 用户已存在 |
| 1303 | 数据库错误 |

---

### GET /api/users/:id

根据 ID 获取用户。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 用户 ID |

**响应 (200)**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "id": 1,
    "username": "test",
    "email": "test@example.com",
    "nickname": "Test User",
    "avatar": "",
    "phone": "",
    "is_active": true,
    "created_at": "2026-04-29T10:00:00Z",
    "updated_at": "2026-04-29T10:00:00Z"
  }
}
```

**错误码**

| code | 说明 |
|------|------|
| 1001 | 参数解析错误 |
| 2002 | 用户不存在 |

---

### GET /api/users/username/:username

根据用户名获取用户。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| username | string | 用户名 |

**响应 (200)**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "id": 1,
    "username": "test",
    "email": "test@example.com",
    "nickname": "Test User",
    "avatar": "",
    "phone": "",
    "is_active": true,
    "created_at": "2026-04-29T10:00:00Z",
    "updated_at": "2026-04-29T10:00:00Z"
  }
}
```

---

### GET /api/users/email/:email

根据邮箱获取用户。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| email | string | 邮箱 |

**响应 (200)**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "id": 1,
    "username": "test",
    "email": "test@example.com",
    "nickname": "Test User",
    "avatar": "",
    "phone": "",
    "is_active": true,
    "created_at": "2026-04-29T10:00:00Z",
    "updated_at": "2026-04-29T10:00:00Z"
  }
}
```

---

### PUT /api/users/:id

更新用户。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 用户 ID |

**请求**

```json
{
  "email": "string (optional, valid email)",
  "nickname": "string (optional, max 100 chars)",
  "avatar": "string (optional, max 500 chars)",
  "phone": "string (optional, max 20 chars)",
  "is_active": "boolean (optional)"
}
```

**响应 (200)**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "id": 1,
    "username": "test",
    "email": "test@example.com",
    "nickname": "Updated Name",
    "avatar": "",
    "phone": "",
    "is_active": true,
    "created_at": "2026-04-29T10:00:00Z",
    "updated_at": "2026-04-29T11:00:00Z"
  }
}
```

---

### DELETE /api/users/:id

删除用户。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 用户 ID |

**响应 (200)**

```json
{
  "code": 0,
  "message": "OK",
  "data": null
}
```

---

### GET /api/users

分页获取用户列表。

**查询参数**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 10 | 每页数量 |

**响应 (200)**

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "test",
        "email": "test@example.com",
        "nickname": "Test User",
        "avatar": "",
        "phone": "",
        "is_active": true,
        "created_at": "2026-04-29T10:00:00Z",
        "updated_at": "2026-04-29T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 错误响应格式

所有 API 错误响应格式统一：

```json
{
  "code": <错误码>,
  "message": "<错误消息>",
  "data": null
}
```

## 错误码表

| 错误码 | HTTP 状态 | 说明 |
|--------|-----------|------|
| 0 | 200 | 成功 |
| 1000 | 500 | 内部错误 |
| 1001 | 400 | 请求参数解析错误 |
| 1002 | 400 | 请求参数校验失败 |
| 1003 | 410 | 已废弃 |
| 1100 | 401 | 认证失败 |
| 1101 | 401 | Token 无效 |
| 1102 | 403 | 权限不足 |
| 1103 | 429 | 请求过于频繁 |
| 1200 | 500 | 文件错误 |
| 1300 | 500 | 数据库错误 |
| 1301 | 404 | 记录不存在 |
| 1302 | 409 | 记录冲突 |
| 1303 | 409 | 记录已存在 |
| 2000 | 400 | 业务错误 |
| 2001 | 409 | 用户已存在 |
| 2002 | 404 | 用户不存在 |
