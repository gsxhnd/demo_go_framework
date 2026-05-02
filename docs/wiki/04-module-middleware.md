# Middleware 模块

> Fiber 中间件链，提供 Recovery、RateLimit、Trace、Metrics、Logger、Auth、RBAC 等横切关注点处理。

## 设计决策

### 为什么需要这个模块

中间件处理横切关注点（日志、追踪、限流等），避免在每个 Handler 中重复实现。Fiber 的中间件模型基于洋葱模型，支持请求前后处理。

### 为什么这么设计

- **选择了**：每个中间件一个独立文件，在 `RegisterHooks` 中按顺序注册
- **而不是**：使用 Fiber 内置中间件或路由组级中间件
- **原因**：自定义中间件可以完全控制行为，与项目的 Logger、Tracer、Metrics 体系深度集成

## 关键类型与接口

### 中间件执行顺序

```
Request → Recovery → RateLimit → Trace → Metrics → Logger → Handler → Response
```

### 已启用中间件

| 顺序 | 中间件 | 文件 | 功能 |
|:----:|--------|------|------|
| 1 | Recovery | `recovery.go` | Panic 恢复，返回 `{code: 1000, message: "Internal Server Error"}` |
| 2 | RateLimit | `rate_limit.go` | 基于令牌桶的 IP 限流，默认 20 req/s，跳过 `/api/health` |
| 3 | Trace | `trace.go` | 创建 OpenTelemetry root span，注入 trace context |
| 4 | Metrics | `metrics.go` | 采集 `http_requests_total`、`http_request_duration_seconds`、`http_requests_in_flight` |
| 5 | Logger | `logger.go` | 记录请求日志，字段含 trace_id、method、path、status、latency |

### 未启用中间件

| 中间件 | 文件 | 功能 |
|--------|------|------|
| Auth | `auth.go` | JWT Bearer Token 认证，验证 Header `Authorization: Bearer <token>` |
| RBAC | `rbac.go` | 基于 Casbin 的 RBAC/ABAC 权限控制 |

## 模块结构

```text
internal/middleware/
├── recovery.go          # Panic 恢复中间件
├── rate_limit.go        # 令牌桶限流中间件
├── trace.go             # OpenTelemetry 链路追踪中间件
├── metrics.go           # HTTP 指标采集中间件
├── logger.go            # 请求日志中间件
├── auth.go              # JWT 认证中间件（未启用）
├── rbac.go              # Casbin 权限中间件（未启用）
├── rate_limit_test.go   # 限流单元测试
└── metrics_test.go      # 指标单元测试
```

| 文件 | 职责 |
|------|------|
| `recovery.go` | 捕获 panic，返回统一错误响应 |
| `rate_limit.go` | 令牌桶算法，按 IP 存储 limiter，定期清理过期条目 |
| `trace.go` | 从请求头提取或创建新 trace context，设置到 Fiber ctx |
| `metrics.go` | 记录请求计数、耗时分布、并发数，带 method/path/status/error 标签 |
| `logger.go` | 请求完成后记录结构化日志 |
| `auth.go` | 验证 Bearer Token 有效性（调用 `pkg/jwx`） |
| `rbac.go` | 从 ctx 提取用户，调用 Casbin 执行权限检查 |

## 与其他模块的关系

### 依赖

- **pkg/logger**：日志接口
- **pkg/trace**：TracerProvider
- **pkg/metrics**：HTTPRecorder
- **pkg/jwx**：JWT 验证（Auth 中间件）
- **pkg/rbac**：权限检查（RBAC 中间件）
- **golang.org/x/time/rate**：令牌桶实现

### 被依赖

- **cmd/server**：在 `RegisterHooks` 中注册

### 依赖关系图

```text
cmd/server
  ↑ (注册)
middleware
  ↑ (依赖)
  ├── pkg/logger
  ├── pkg/trace
  ├── pkg/metrics
  ├── pkg/jwx (auth)
  └── pkg/rbac (rbac)
```

## 注意事项

- 中间件按 `app.Use()` 调用顺序执行，洋葱模型：先注册的先执行前置逻辑，后执行后置逻辑
- RateLimit 使用内存存储，单实例有效，多实例需改为 Redis 存储
- RateLimit 跳过 `/api/health` 路径，避免健康检查被限流
- Auth 和 RBAC 中间件已实现但未注册，需在 `RegisterHooks` 中添加并配置路由策略
- Metrics 中间件使用 OTLP exporter 导出到 OTel Collector
- 自定义中间件需在 `internal/middleware/` 中创建，并在 `RegisterHooks` 中注册
