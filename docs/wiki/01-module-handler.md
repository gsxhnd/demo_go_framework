# Handler 模块

> HTTP 请求处理层，负责请求参数解析、参数校验、调用 Service 层和构造 JSON 响应。

## 设计决策

### 为什么需要这个模块

Handler 层是 HTTP 请求的入口，负责将 HTTP 协议相关的逻辑（请求解析、响应构造）与业务逻辑分离。这使得 Service 层可以独立于 HTTP 协议进行测试和复用。

### 为什么这么设计

- **选择了**：每个 Handler 实现一个接口，每个接口方法对应一个 HTTP 端点
- **而不是**：使用通用路由分发器或反射路由
- **原因**：接口驱动便于 mock 测试，方法命名与端点一一对应，代码清晰可读

## 关键类型与接口

### healthhandler.Handler

- **定义位置**：`internal/handler/health/handler.go`
- **用途**：健康检查处理器
- **方法**：`Check(c *fiber.Ctx) error`

### userhandler.Handler

- **定义位置**：`internal/handler/user/handler.go`
- **用途**：用户管理 HTTP 处理器
- **方法**：
  - `UserCreate(c *fiber.Ctx) error`
  - `UserGetByID(c *fiber.Ctx) error`
  - `UserGetByUsername(c *fiber.Ctx) error`
  - `UserGetByEmail(c *fiber.Ctx) error`
  - `UserUpdate(c *fiber.Ctx) error`
  - `UserDelete(c *fiber.Ctx) error`
  - `UserList(c *fiber.Ctx) error`

### 请求/响应类型

- **定义位置**：`internal/handler/user/types.go`
- **用途**：定义各端点的请求体和响应体结构
- **包含**：`CreateUserRequest`, `UpdateUserRequest`, `UserResponse`, `ListUsersResponse` 等

## 模块结构

```text
internal/handler/
├── health/
│   ├── handler.go          # Health Handler 接口与实现
│   ├── check.go            # GET /api/health 端点
│   └── handler_test.go     # 单元测试
│
└── user/
    ├── handler.go           # User Handler 接口与实现
    ├── types.go             # 请求/响应类型定义
    ├── validator.go         # 自定义校验规则
    ├── user_create.go       # POST /api/users
    ├── user_delete.go       # DELETE /api/users/:id
    ├── user_update.go       # PUT /api/users/:id
    ├── user_get_by_id.go    # GET /api/users/:id
    ├── user_get_by_username.go # GET /api/users/username/:username
    ├── user_get_by_email.go    # GET /api/users/email/:email
    ├── user_list.go         # GET /api/users
    └── handler_test.go      # 单元测试
```

| 文件 | 职责 |
|------|------|
| `handler.go` | 定义 Handler 接口和实现结构体，通过构造函数注入 Service、Logger、Tracer、Validator |
| `types.go` | 请求体和响应体结构定义，含 `json` / `query` / `params` 标签 |
| `validator.go` | 自定义校验规则（如 UpdateUserRequest 至少传一个字段） |
| `user_*.go` | 各端点实现：解析参数 → 校验 → 调用 Service → `errno.Decode` 构造响应 |

## 与其他模块的关系

### 依赖

- **userservice**：调用 Service 层执行业务逻辑
- **validator**：参数校验
- **logger**：记录日志
- **trace**：创建子 Span
- **errno**：统一错误响应格式

### 被依赖

- **cmd/server**：在 `RegisterHooks` 中注册路由

### 依赖关系图

```text
cmd/server
  ↑ (注册路由)
handler
  ↑ (调用)
  ├── service (业务逻辑)
  ├── validator (参数校验)
  ├── logger
  ├── trace
  └── errno (错误转换)
```

## 注意事项

- Handler 层不应包含业务逻辑，只做协议转换
- 所有错误通过 `errno.Decode(data, err)` 统一转换为 JSON 响应
- 每个端点方法需要自己解析路径参数、查询参数和请求体
- 测试时使用 `fiber.App` + `httptest.NewRequest` 模拟 HTTP 请求
- 用户管理路由当前未在 RegisterHooks 中注册，需手动添加
