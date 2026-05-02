# 领域模型

## 核心实体

### User

用户实体，由 Ent ORM 管理，对应数据库 `users` 表。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | int | PK, AutoIncrement | 用户 ID |
| `username` | string | Unique, NotEmpty, MaxLen(50) | 用户名 |
| `email` | string | Unique, NotEmpty, MaxLen(255) | 邮箱 |
| `password` | string | NotEmpty, Sensitive | 密码（哈希存储） |
| `nickname` | string | Optional, MaxLen(100) | 昵称 |
| `avatar` | string | Optional, MaxLen(500) | 头像 URL |
| `phone` | string | Optional, MaxLen(20) | 手机号 |
| `is_active` | bool | Default(true) | 是否激活 |
| `created_at` | time.Time | AutoCreate | 创建时间（mixin.TimeMixin） |
| `updated_at` | time.Time | AutoUpdate | 更新时间（mixin.TimeMixin） |

**索引**：

- `username` 唯一索引
- `email` 唯一索引

**Ent Schema 定义**：`internal/ent/schema/user.go`

## 核心接口

### Handler 层接口

```go
// internal/handler/user/handler.go
type Handler interface {
    UserCreate(c *fiber.Ctx) error
    UserGetByID(c *fiber.Ctx) error
    UserGetByUsername(c *fiber.Ctx) error
    UserGetByEmail(c *fiber.Ctx) error
    UserUpdate(c *fiber.Ctx) error
    UserDelete(c *fiber.Ctx) error
    UserList(c *fiber.Ctx) error
}
```

### Service 层接口

```go
// internal/service/user/service.go
type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
    GetUserByID(ctx context.Context, id int) (*UserResponse, error)
    GetUserByUsername(ctx context.Context, username string) (*UserResponse, error)
    GetUserByEmail(ctx context.Context, email string) (*UserResponse, error)
    UpdateUser(ctx context.Context, id int, req UpdateUserRequest) (*UserResponse, error)
    DeleteUser(ctx context.Context, id int) error
    ListUsers(ctx context.Context, page, pageSize int) (*ListUsersResponse, error)
}
```

### Repo 层接口

```go
// internal/repo/user/repo.go
type UserRepo interface {
    Create(ctx context.Context, user *CreateUserParams) (*ent.User, error)
    GetByID(ctx context.Context, id int) (*ent.User, error)
    GetByUsername(ctx context.Context, username string) (*ent.User, error)
    GetByEmail(ctx context.Context, email string) (*ent.User, error)
    Update(ctx context.Context, id int, params *UpdateUserParams) (*ent.User, error)
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, page, pageSize int) ([]*ent.User, int, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    ExistsByUsername(ctx context.Context, username string) (bool, error)
}
```

## 错误码体系

使用 `errno.Errno` 接口统一错误表示：

```go
type Errno interface {
    Error() string
    GetHTTPStatus() int
    GetCode() int
    GetMessage() string
    GetData() any
}
```

### 通用错误码（1000-1099）

| Code | 名称 | HTTP 状态码 | 说明 |
|------|------|-------------|------|
| 0 | OK | 200 | 成功 |
| 1000 | InternalServerError | 500 | 服务器内部错误 |
| 1001 | Deprecated | 410 | API 已废弃 |
| 1002 | RequestParserError | 400 | 请求解析错误 |
| 1003 | RequestValidateError | 400 | 参数校验失败 |
| 1004 | RequestConversionError | 400 | 请求转换错误 |
| 1005 | DataConversionError | 500 | 数据转换错误 |

### 认证/授权错误码（1100-1199）

| Code | 名称 | HTTP 状态码 | 说明 |
|------|------|-------------|------|
| 1101 | TokenInvalidError | 401 | Token 无效 |
| 1102 | TokenParserError | 401 | Token 解析错误 |
| 1103 | PermissionDeniedError | 403 | 权限不足 |
| 1104 | RateLimitExceededError | 429 | 超出限流 |

### 文件错误码（1200-1299）

| Code | 名称 | HTTP 状态码 | 说明 |
|------|------|-------------|------|
| 1201 | RetrievingFileError | 404 | 文件获取错误 |

### 数据库错误码（1300+）

| Code | 名称 | HTTP 状态码 | 说明 |
|------|------|-------------|------|
| 1301 | DatabaseError | 500 | 数据库错误 |
| 1302 | DatabaseConversionError | 500 | 数据库转换错误 |

### 业务错误码（2000+）

| Code | 名称 | HTTP 状态码 | 说明 |
|------|------|-------------|------|
| 2001 | UserNotFoundError | 404 | 用户不存在 |
| 2002 | UserAlreadyExistsError | 409 | 用户已存在 |
| 2003 | UserCreateFailedError | 500 | 创建用户失败 |
| 2004 | UserUpdateFailedError | 500 | 更新用户失败 |
| 2005 | UserDeleteFailedError | 500 | 删除用户失败 |
| 2006 | InvalidUserIDError | 400 | 无效的用户 ID |
| 2007 | InvalidEmailError | 400 | 无效的邮箱格式 |
| 2008 | InvalidUsernameError | 400 | 无效的用户名格式 |

### 分页错误码（3000+）

| Code | 名称 | HTTP 状态码 | 说明 |
|------|------|-------------|------|
| 3001 | InvalidPageError | 400 | 无效的页码 |
| 3002 | InvalidPageSizeError | 400 | 无效的每页数量 |

## 错误响应格式

所有 API 错误响应格式统一：

```json
{
  "code": 2001,
  "message": "User not found",
  "data": null
}
```

## 数据流

### 创建用户流程

```
POST /api/users
  → Recovery（Panic 恢复）
  → RateLimit（限流检查）
  → Trace（创建 Span）
  → Metrics（记录请求指标）
  → Logger（记录请求日志）
  → UserHandler.UserCreate()
    → validator 参数校验
    → UserService.CreateUser()
      → 检查用户名/邮箱唯一性
      → UserRepo.Create()
        → Ent ORM 执行 INSERT
      → 返回 UserResponse
    → errno.Decode(data, err) 构造响应
  → JSON Response
```
