# 进阶使用

## 概述

本文档介绍 go_sample_code 的高级功能，包括可观测性集成、中间件配置、权限控制和新模块开发。

## 可观测性

### 启动可观测性栈

```bash
# 基础版
docker compose -f devops/grafana.v1/docker-compose.yml up -d

# ClickHouse 版（更高性能）
docker compose -f devops/grafana.v2/docker-compose.yml up -d
```

启动后可访问：

| 服务 | 地址 | 用途 |
|------|------|------|
| Grafana | <http://localhost:3000> | 仪表盘和可视化 |
| Prometheus | <http://localhost:9090> | 指标查询 |
| Tempo | <http://localhost:3100> | 链路追踪查询 |
| Loki | <http://localhost:3100> | 日志查询 |

### 链路追踪

所有请求自动生成 Trace ID 和 Span ID，通过 OTLP gRPC 发送到 Collector。

在 HTTP 响应头中可获取 `traceparent` 用于关联。

### 指标采集

自动采集以下指标：

| 指标 | 类型 | 说明 |
|------|------|------|
| `http_requests_total` | Counter | 请求总数 |
| `http_request_duration_seconds` | Histogram | 请求耗时分布 |
| `http_requests_in_flight` | Gauge | 当前并发请求数 |

标签：`method`、`path`、`status_code`、`error`

## 中间件配置

### 限流配置

默认限流：20 req/s，桶容量 50。在 `internal/middleware/rate_limit.go` 中修改：

```go
config := RateLimitConfig{
    Rate:     20,          // 每秒令牌数
    Burst:    50,          // 桶容量
    CleanupInterval: 5 * time.Minute, // 清理间隔
}
```

### 启用认证中间件

在 `cmd/server/main.go` 的 `RegisterHooks` 中添加：

```go
app.Use(middleware.Auth(jwxService))
```

认证中间件验证请求头 `Authorization: Bearer <token>` 中的 JWT Token。

### 启用权限中间件

在 `cmd/server/main.go` 的 `RegisterHooks` 中添加：

```go
app.Use(middleware.RBAC(rbacService))
```

权限策略配置在 `pkg/rbac/model.conf` 和 `pkg/rbac/policy.csv` 中。

## 新增业务模块

### 1. 定义 Ent Schema

在 `internal/ent/schema/` 中创建实体定义文件：

```go
// internal/ent/schema/product.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type Product struct {
    ent.Schema
}

func (Product) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").NotEmpty(),
        field.Float("price").Positive(),
    }
}
```

### 2. 生成 Ent 代码

```bash
ent generate ./internal/ent/schema
```

### 3. 创建 Repo 层

在 `internal/repo/product/` 中创建数据访问接口和实现。

### 4. 创建 Service 层

在 `internal/service/product/` 中创建业务逻辑接口和实现。

### 5. 创建 Handler 层

在 `internal/handler/product/` 中创建 HTTP 处理器。

### 6. 注册路由

在 `cmd/server/main.go` 的 `RegisterHooks` 中注册新路由。

## 自定义中间件

在 `internal/middleware/` 中创建新文件：

```go
func NewCustomMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 前置处理
        start := time.Now()

        err := c.Next() // 执行后续中间件/Handler

        // 后置处理
        latency := time.Since(start)
        // ...

        return err
    }
}
```

然后在 `RegisterHooks` 中注册：

```go
app.Use(middleware.NewCustomMiddleware())
```

中间件按 `app.Use()` 顺序执行（洋葱模型）。

## 性能调优

### 连接池配置

```yaml
database:
  relational:
    postgres:
      max_open_conns: 200    # 最大连接数
      max_idle_conns: 50     # 最大空闲连接数
      conn_max_lifetime: 1800 # 连接最大存活时间（秒）
  redis:
    pool_size: 200           # Redis 连接池大小
```

### 日志级别

生产环境建议使用 `info` 级别，开发调试使用 `debug` 级别：

```yaml
logger:
  level: "info"
```

### 采样率

高流量场景下可降低 trace 采样率：

```yaml
trace:
  sampling_ratio: 0.1  # 仅采样 10% 的请求
```
