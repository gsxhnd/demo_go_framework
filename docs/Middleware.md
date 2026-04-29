# 中间件文档

## 概述

项目实现了多个 Fiber 中间件，按配置顺序依次执行。

## 中间件顺序

```
Request → Recovery → RateLimit → Trace → Metrics → Logger → Handler → Response
```

## 已启用中间件

### 1. Recovery

**文件**: `internal/middleware/recovery.go`

**功能**: 捕获 panic 异常，返回 500 错误响应。

**响应**:

```json
{
  "code": 1000,
  "message": "Internal Server Error",
  "data": null
}
```

---

### 2. RateLimit

**文件**: `internal/middleware/rate_limit.go`

**功能**: 基于令牌桶算法的请求限流。

**配置** (代码中):

| 配置项 | 默认值 |
|--------|--------|
| 每秒令牌数 | 20 |
| 桶容量 | 50 |
| 清理间隔 | 5 分钟 |

**特性**:

- 按 IP 地址限流
- `/api/health` 路径跳过限流
- 使用 `golang.org/x/time/rate`

**响应** (超出限流):

```json
{
  "code": 1103,
  "message": "Rate limit exceeded",
  "data": null
}
```

**HTTP 状态码**: `429 Too Many Requests`

---

### 3. Trace

**文件**: `internal/middleware/trace.go`

**功能**: OpenTelemetry 链路追踪。

**特性**:

- 创建根 span
- 注入 trace context 到请求上下文
- 支持 W3C Trace Context 标准
- 使用 OTLP exporter 发送数据

---

### 4. Metrics

**文件**: `internal/middleware/metrics.go`

**功能**: HTTP 请求指标采集。

**采集指标**:

| 指标名 | 类型 | 描述 |
|--------|------|------|
| `http_requests_total` | Counter | 请求总数 |
| `http_request_duration_seconds` | Histogram | 请求耗时 |
| `http_requests_in_flight` | Gauge | 当前处理中请求数 |

**标签**:

- `method`: HTTP 方法 (GET, POST, etc.)
- `path`: 请求路径
- `status_code`: HTTP 状态码
- `error`: 是否有错误 (true/false)

---

### 5. Logger

**文件**: `internal/middleware/logger.go`

**功能**: 请求/响应日志记录。

**日志字段**:

- `trace_id`: 链路 ID
- `span_id`: Span ID
- `method`: HTTP 方法
- `path`: 请求路径
- `status`: HTTP 状态码
- `latency`: 请求耗时
- `client_ip`: 客户端 IP
- `user_agent`: User Agent

---

## 未启用中间件

以下中间件已实现但未在 `RegisterHooks` 中注册：

### 6. Auth

**文件**: `internal/middleware/auth.go`

**功能**: JWT Bearer Token 认证。

**Header**: `Authorization: Bearer <token>`

**验证**: 调用 `pkg/jwx` 验证 Token 有效性。

**启用方式**: 在 `cmd/server/main.go` 的 `RegisterHooks` 中添加。

---

### 7. RBAC

**文件**: `internal/middleware/rbac.go`

**功能**: 基于 Casbin 的 RBAC/ABAC 权限控制。

**特性**:

- 支持 RBAC (基于角色的访问控制)
- 支持 ABAC (基于属性的访问控制)
- 从请求上下文中提取用户信息

**启用方式**: 在 `cmd/server/main.go` 的 `RegisterHooks` 中添加。

---

## 自定义中间件

### 添加新中间件

1. 在 `internal/middleware/` 目录创建文件，如 `custom.go`

2. 实现 Fiber Middleware 接口：

```go
func NewCustomMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 前置处理
        start := time.Now()
        
        // 继续执行后续中间件/Handler
        err := c.Next()
        
        // 后置处理
        latency := time.Since(start)
        // 记录日志等
        
        return err
    }
}
```

1. 在 `cmd/server/main.go` 的 `RegisterHooks` 中注册：

```go
app.Use(custom.NewCustomMiddleware())
```

### 中间件执行顺序

中间件按照 `app.Use()` 调用顺序执行，后添加的先执行（洋葱模型）：

```go
app.Use(Recovery())  // 1. 最后执行
app.Use(RateLimit()) // 2. 
app.Use(Trace())     // 3.
app.Use(Metrics())   // 4. 最先执行
```

---

## 测试

限流中间件测试: `internal/middleware/rate_limit_test.go`

指标中间件测试: `internal/middleware/metrics_test.go`

运行测试:

```bash
go test ./internal/middleware/...
```
