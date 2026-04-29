# 代码结构

## 目录结构

```
go_sample_code/
├── cmd/                          # 应用入口
│   └── server/
│       ├── main.go               # 主程序入口
│       ├── config.go             # 配置定义
│       └── README.md
│
├── internal/                     # 内部包（不对外暴露）
│   ├── database/                 # 数据库相关
│   │   ├── config.go             # 数据库配置结构
│   │   ├── relational.go         # 关系型数据库客户端
│   │   ├── redis.go              # Redis 客户端
│   │   ├── health.go             # 健康检查接口
│   │   └── *_test.go
│   │
│   ├── ent/                      # Ent ORM 生成代码（勿直接编辑）
│   │   ├── client.go             # Ent 客户端
│   │   ├── ent.go                # Ent 主入口
│   │   ├── mutation.go            # 变更记录
│   │   ├── user.go               # User 实体
│   │   ├── user_create.go        # User 创建操作
│   │   ├── user_update.go        # User 更新操作
│   │   ├── user_delete.go        # User 删除操作
│   │   ├── user_query.go         # User 查询
│   │   ├── user/                 # User 查询 API
│   │   ├── hook/                 # Ent Hooks
│   │   ├── predicate/            # Ent 谓词
│   │   ├── migrate/              # 迁移工具
│   │   └── schema/               # 实体定义（可编辑）
│   │       ├── user.go          # User 实体定义
│   │       └── mixin/
│   │           └── time.go       # 时间戳混入
│   │
│   ├── errno/                    # 错误码
│   │   ├── errno.go              # 错误接口定义
│   │   ├── code.go               # 错误码定义
│   │   ├── business.go           # 业务错误码
│   │   └── errno_test.go
│   │
│   ├── handler/                  # HTTP 处理器
│   │   ├── health/
│   │   │   ├── handler.go        # 健康检查处理器接口/实现
│   │   │   ├── check.go          # 健康检查端点
│   │   │   └── handler_test.go
│   │   │
│   │   └── user/
│   │       ├── handler.go        # 用户处理器接口/实现
│   │       ├── types.go          # 请求/响应类型
│   │       ├── user_create.go    # 创建用户
│   │       ├── user_delete.go    # 删除用户
│   │       ├── user_update.go    # 更新用户
│   │       ├── user_get_by_id.go # 按 ID 获取
│   │       ├── user_get_by_username.go  # 按用户名获取
│   │       ├── user_get_by_email.go     # 按邮箱获取
│   │       ├── user_list.go      # 用户列表
│   │       ├── validator.go      # 请求校验
│   │       └── handler_test.go
│   │
│   ├── middleware/               # 中间件
│   │   ├── recovery.go           #  Panic 恢复
│   │   ├── rate_limit.go         #  限流
│   │   ├── trace.go              #  链路追踪
│   │   ├── metrics.go            #  指标采集
│   │   ├── logger.go             #  请求日志
│   │   ├── auth.go               #  认证（未启用）
│   │   ├── rbac.go               #  权限（未启用）
│   │   ├── metrics_test.go
│   │   └── rate_limit_test.go
│   │
│   ├── repo/                     # 数据访问层
│   │   └── user/
│   │       ├── repo.go           # 用户仓储接口/实现
│   │       ├── user_create.go    # 创建用户
│   │       ├── user_delete.go    # 删除用户
│   │       ├── user_update.go    # 更新用户
│   │       ├── user_get_by_id.go # 按 ID 获取
│   │       ├── user_get_by_username.go  # 按用户名获取
│   │       ├── user_get_by_email.go     # 按邮箱获取
│   │       ├── user_list.go      # 用户列表
│   │       ├── user_exists_by_email.go   # 邮箱是否存在
│   │       └── user_exists_by_username.go # 用户名是否存在
│   │
│   └── service/                  # 业务逻辑层
│       └── user/
│           ├── service.go        # 用户服务接口/实现
│           ├── user_create.go    # 创建用户
│           ├── user_delete.go    # 删除用户
│           ├── user_update.go    # 更新用户
│           ├── user_get_by_id.go # 按 ID 获取
│           ├── user_get_by_username.go  # 按用户名获取
│           ├── user_get_by_email.go     # 按邮箱获取
│           └── user_list.go      # 用户列表
│
├── pkg/                          # 公共包
│   ├── logger/                   # 日志
│   │   ├── logger.go             # Logger 接口 + zap 实现
│   │   ├── config.go             # 日志配置
│   │   └── README.md
│   │
│   ├── trace/                    # 链路追踪
│   │   └── trace.go              # OTel TracerProvider
│   │
│   ├── metrics/                  # 指标
│   │   ├── metrics.go            # MeterProvider
│   │   ├── config.go            # 指标配置
│   │   ├── http.go              # HTTP 指标
│   │   └── *_test.go
│   │
│   ├── validator/                # 参数校验
│   │   ├── validator.go         # 校验器封装
│   │   └── errors.go            # 校验错误
│   │
│   ├── jwx/                      # JWT
│   │   └── jwx.go               # JWT/JWK 操作
│   │
│   └── rbac/                     # 权限控制
│       ├── rbac.go              # RBAC 服务
│       ├── abac.go              # ABAC 服务
│       ├── service.go          # 权限服务
│       ├── middleware.go        # 权限中间件
│       ├── config.go           # 权限配置
│       ├── errors.go           # 权限错误
│       ├── model.conf          # Casbin 模型
│       ├── abac_model.conf     # ABAC 模型
│       └── policy.csv          # 权限策略
│
├── config/                       # 配置文件
│   ├── config.local.yaml        # 本地开发配置
│   └── config.template.yaml     # 配置模板
│
├── devops/                       # 部署配置
│   ├── database/
│   │   └── docker-compose.yml   # 数据库服务
│   ├── grafana.v1/              # 可观测性栈 v1
│   │   ├── docker-compose.yml
│   │   ├── grafana/
│   │   ├── prometheus/
│   │   ├── tempo/
│   │   └── loki/
│   └── grafana.v2/              # 可观测性栈 v2 (ClickHouse)
│
├── docs/                         # 文档
├── go.mod
├── go.sum
├── README.md
├── AGENTS.md
└── LICENSE
```

## 分层职责

### Handler 层 (`internal/handler/`)

负责 HTTP 请求/响应处理。

**职责**：

- 解析请求参数
- 参数校验
- 调用 Service 层
- 构造响应

**接口定义模式**：

```go
// 接口定义
type Handler interface {
    UserCreate(c *fiber.Ctx) error
    UserGetByID(c *fiber.Ctx) error
    // ...
}

// 实现结构
type handler struct {
    userService userservice.UserService
    log         *zap.Logger
    tracer      trace.Tracer
    validate    *validator.Validate
}
```

### Service 层 (`internal/service/`)

负责业务逻辑处理。

**职责**：

- 业务规则校验
- 业务逻辑执行
- 调用 Repo 层
- 错误转换

**接口定义模式**：

```go
type UserService interface {
    CreateUser(ctx context.Context, req *userrepo.CreateUserRequest) (*UserResponse, errno.Errno)
    GetUserByID(ctx context.Context, id int) (*UserResponse, errno.Errno)
    // ...
}

type userService struct {
    userRepo userrepo.UserRepo
    log      *zap.Logger
    tracer   trace.Tracer
}
```

### Repo 层 (`internal/repo/`)

负责数据访问。

**职责**：

- 数据库操作
- Ent ORM 调用
- 数据转换

**接口定义模式**：

```go
type UserRepo interface {
    UserCreate(ctx context.Context, req *CreateUserRequest) (*ent.User, error)
    UserGetByID(ctx context.Context, id int) (*ent.User, error)
    // ...
}

type userRepo struct {
    client * ent.Client
    log    *zap.Logger
    tracer trace.Tracer
}
```

## Ent 实体定义

实体定义位于 `internal/ent/schema/`，编辑后需重新生成代码。

```bash
ent generate ./internal/ent/schema
```

### User 实体

```go
// internal/ent/schema/user.go
type User struct {
    ent.Schema
}

func (User) Mixin() []ent.Mixin {
    return []ent.Mixin{mixin.TimeMixin{}}
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").Unique().NotEmpty().MaxLen(50),
        field.String("email").Unique().NotEmpty().MaxLen(255),
        field.String("password").NotEmpty().Sensitive(),
        field.String("nickname").MaxLen(100).Optional(),
        field.String("avatar").MaxLen(500).Optional(),
        field.String("phone").MaxLen(20).Optional(),
        field.Bool("is_active").Default(true),
    }
}
```

### 时间戳混入

```go
// internal/ent/schema/mixin/time.go
type TimeMixin struct {
    ent.Schema
}

func (TimeMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Time("created_at").Default(time.Now).Immutable(),
        field.Time("updated_at").UpdateDefault(time.Now),
    }
}
```

## 包依赖关系

```
cmd/server/main.go
         │
         ▼
┌─────────────────────────────────────────┐
│  pkg/logger      - 日志接口             │
│  pkg/trace       - 追踪接口             │
│  pkg/metrics     - 指标接口             │
│  pkg/validator   - 校验接口             │
│  pkg/jwx         - JWT 接口             │
│  pkg/rbac        - 权限接口             │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  internal/database  - DB/Redis 客户端   │
│  internal/ent       - ORM 客户端        │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  internal/repo/user   - 数据访问       │
│  internal/service/user - 业务逻辑       │
│  internal/handler/user - HTTP 处理      │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  internal/middleware  - 中间件          │
└─────────────────────────────────────────┘
```
