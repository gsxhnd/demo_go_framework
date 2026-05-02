# Service 模块

> 业务逻辑处理层，负责业务规则校验、流程编排和错误转换。

## 设计决策

### 为什么需要这个模块

Service 层是业务逻辑的核心所在。将业务逻辑从 Handler 层分离可以：

- 使业务逻辑独立于 HTTP 协议，便于测试和复用
- 支持未来可能的多种入口（HTTP API、gRPC、CLI）共享同一套业务逻辑

### 为什么这么设计

- **选择了**：每个 Service 实现一个接口，方法按业务操作组织
- **而不是**：使用贫血模型（业务逻辑放在 Handler 或 Repo 中）
- **原因**：保持 Handler 薄、Service 厚、Repo 纯粹的分层原则

## 关键类型与接口

### userservice.UserService

- **定义位置**：`internal/service/user/service.go`
- **用途**：用户业务逻辑服务

**方法**：

| 方法 | 说明 |
|------|------|
| `CreateUser(ctx, req) (*UserResponse, error)` | 创建用户，含唯一性校验 |
| `GetUserByID(ctx, id) (*UserResponse, error)` | 按 ID 查询 |
| `GetUserByUsername(ctx, username) (*UserResponse, error)` | 按用户名查询 |
| `GetUserByEmail(ctx, email) (*UserResponse, error)` | 按邮箱查询 |
| `UpdateUser(ctx, id, req) (*UserResponse, error)` | 更新用户信息 |
| `DeleteUser(ctx, id) error` | 删除用户 |
| `ListUsers(ctx, page, pageSize) (*ListUsersResponse, error)` | 分页列表 |

## 模块结构

```text
internal/service/user/
├── service.go              # UserService 接口与实现
├── user_create.go          # 创建用户业务逻辑
├── user_delete.go          # 删除用户业务逻辑
├── user_update.go          # 更新用户业务逻辑
├── user_get_by_id.go       # 按 ID 查询
├── user_get_by_username.go # 按用户名查询
├── user_get_by_email.go    # 按邮箱查询
└── user_list.go            # 分页列表
```

| 文件 | 职责 |
|------|------|
| `service.go` | 定义接口、实现结构体和构造函数 |
| `user_create.go` | 创建用户：检查唯一性 → 调用 Repo → 返回 UserResponse |
| `user_*.go` | 各业务方法实现 |

## 与其他模块的关系

### 依赖

- **userrepo**：调用 Repo 层进行数据操作
- **errno**：业务错误定义与转换

### 被依赖

- **handler**：Handler 层调用 Service 层

### 依赖关系图

```text
handler
  ↑ (调用)
service
  ↑ (调用)
  ├── repo (数据访问)
  └── errno (错误定义)
```

## 注意事项

- Service 层负责业务规则校验（如唯一性检查），不负责参数格式校验（由 Handler + Validator 负责）
- 所有业务错误使用 `errno` 中定义的变量（如 `UserNotFoundError`、`UserAlreadyExistsError`）
- Service 方法的参数使用专用请求类型（如 `CreateUserRequest`），而非直接使用 Ent 实体
- 返回值使用专用响应类型（如 `UserResponse`），隐藏内部实体细节
