---
name: go-service-layer
description: 在 Go Fiber 项目中创建标准化的 service 层实现。使用场景：需要添加新的业务服务、实现 CRUD 操作、创建服务层代码时。参考 internal/service/user/ 作为实现示例。
---

# Go Service Layer Implementation

## 项目结构

```
internal/service/{name}/
├── {name}.go           # 接口定义 + DTO
├── {name}_impl.go      # 服务实现
└── {name}_test.go      # 单元测试
```

## 实现步骤

### 1. 创建接口和 DTO (`{name}.go`)

```go
package user

import "context"

// Entity 定义实体
type User struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

// Request/Response DTO
type CreateRequest struct {
    Name string `json:"name"`
}

type ListResponse struct {
    Users []*User `json:"users"`
    Total int64  `json:"total"`
}

// Service 接口
type Service interface {
    Create(ctx context.Context, req *CreateRequest) (*User, error)
    List(ctx context.Context) (*ListResponse, error)
}
```

### 2. 创建 Tracer (`tracer.go`)

```go
package user

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    otel_trace "go.opentelemetry.io/otel/trace"
)

const scopeName = "user.service"

type Tracer struct {
    tracer otel_trace.Tracer
}

func NewTracer(serviceName string) *Tracer {
    tracer := otel.Tracer(scopeName)
    return &Tracer{tracer: tracer}
}

func (t *Tracer) TraceCreate(ctx context.Context, name string) (context.Context, otel_trace.Span) {
    ctx, span := t.tracer.Start(ctx, "user.Create")
    span.SetAttributes(attribute.String("user.name", name))
    return ctx, span
}
```

### 3. 创建服务实现 (`{name}_impl.go`)

```go
package user

import (
    "context"
    "go_sample_code/internal/errno"
    "go_sample_code/pkg/logger"
    "go.uber.org/zap"
)

type ServiceImpl struct {
    repo   Repository
    log    logger.Logger
    tracer *Tracer
}

type Repository interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id int64) (*User, error)
}

func NewService(repo Repository, log logger.Logger, tracer *Tracer) Service {
    return &ServiceImpl{repo: repo, log: log, tracer: tracer}
}

func (s *ServiceImpl) Create(ctx context.Context, req *CreateRequest) (*User, error) {
    ctx, span := s.tracer.TraceCreate(ctx, req.Name)
    defer span.End()

    s.log.InfoCtx(ctx, "creating user", zap.String("name", req.Name))

    user := &User{Name: req.Name}
    created, err := s.repo.Create(ctx, user)
    if err != nil {
        s.log.ErrorCtx(ctx, "failed to create user", zap.Error(err))
        span.RecordError(err)
        return nil, err
    }

    s.log.InfoCtx(ctx, "user created", zap.Int64("id", created.ID))
    return created, nil
}
```

### 4. 创建单元测试 (`{name}_test.go`)

```go
package user

import (
    "context"
    "testing"
    "go_sample_code/internal/errno"
    "go_sample_code/pkg/logger"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func newTestService(t *testing.T) Service {
    log, _ := logger.NewLogger(logger.DefaultConfig())
    repo := NewMockRepository()
    tracer := NewTracer("test")
    return NewService(repo, log, tracer)
}

func TestService_Create(t *testing.T) {
    svc := newTestService(t)
    ctx := context.Background()

    user, err := svc.Create(ctx, &CreateRequest{Name: "test"})
    require.NoError(t, err)
    assert.Equal(t, "test", user.Name)

    // 错误处理测试
    _, err = svc.Create(ctx, &CreateRequest{Name: ""})
    assert.Error(t, err)
    var wrapper *errno.ErrnoWrapper
    assert.ErrorAs(t, err, &wrapper)
}
```

### 5. 注册 Fx 模块

在 `internal/service/provider.go` 中添加：

```go
var Module = fx.Module(
    "service",
    fx.Provide(
        user.NewRepository,
        user.NewTracer,
        user.NewService,
    ),
)
```

## 错误处理规范

使用 `errno.ErrnoWrapper` 包装错误：

```go
// 定义错误码 (internal/errno/business.go)
var UserNotFoundError = errno errno{HTTPStatus: 404, Code: 2001, Message: "User not found"}

// 返回错误
return nil, &errno.ErrnoWrapper{Errno: errno.UserNotFoundError}
```

## 日志规范

```go
s.log.InfoCtx(ctx, "operation description", zap.String("key", value))
s.log.DebugCtx(ctx, "debug info", zap.Int("count", n))
s.log.ErrorCtx(ctx, "error occurred", zap.Error(err))
s.log.WarnCtx(ctx, "validation failed", zap.Error(err))
```

## Tracing 规范

```go
ctx, span := s.tracer.TraceMethodName(ctx, param)
defer span.End()
// ... 业务逻辑 ...
if err != nil {
    span.RecordError(err)
    return nil, err
}
```

## 参考实现

- `internal/service/user/` - 完整的用户服务实现
- `internal/errno/` - 错误码定义
- `pkg/logger/` - 日志接口
- `cmd/server/main.go` - Fx 依赖注入配置
