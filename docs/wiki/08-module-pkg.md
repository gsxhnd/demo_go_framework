# Pkg 模块

> 公共工具包集合，包含 logger、trace、metrics、validator、jwx、rbac 六个子包。

## 设计决策

### 为什么放在 pkg/

Go 约定：`internal/` 中的包不可被外部项目引用，`pkg/` 中的包可以被外部引用。将可复用的工具包放在 `pkg/` 使其可以被其他项目导入。

### 为什么拆分为多个子包

- **选择了**：按功能域拆分为独立的子包
- **而不是**：一个大的 `pkg/common` 包含所有工具
- **原因**：单一职责，每个子包可独立演进和测试，按需导入减少编译依赖

## 子包总览

| 子包 | 路径 | 用途 |
|------|------|------|
| logger | `pkg/logger/` | 日志接口 + zap 实现 |
| trace | `pkg/trace/` | OpenTelemetry TracerProvider |
| metrics | `pkg/metrics/` | HTTP 指标采集与 OTLP 导出 |
| validator | `pkg/validator/` | 参数校验封装（go-playground/validator） |
| jwx | `pkg/jwx/` | JWT/JWK 操作 |
| rbac | `pkg/rbac/` | Casbin RBAC/ABAC 权限控制 |

## 关键类型与接口

### pkg/logger

| 类型/函数 | 说明 |
|-----------|------|
| `logger.Logger` | 日志接口（`Info`, `Error`, `Warn`, `Debug`, `Shutdown` 等） |
| `logger.NewLogger(cfg)` | 创建 zap 日志实例 |
| `logger.LoggerConfig` | 日志配置（Level、Encoding、Output、OtelServiceName） |

### pkg/trace

| 类型/函数 | 说明 |
|-----------|------|
| `trace.Tracer` | 类型别名，使用 OpenTelemetry Tracer |
| `trace.NewTracerProvider(cfg)` | 创建 OTLP TracerProvider |
| `trace.NewInMemoryProvider()` | 创建内存 TracerProvider（测试用） |
| `trace.TraceConfig` | 链路追踪配置 |

### pkg/metrics

| 类型/函数 | 说明 |
|-----------|------|
| `metrics.HTTPRecorder` | HTTP 指标采集器（Counter、Histogram、Gauge） |
| `metrics.NewHTTPRecorder(provider)` | 创建 HTTPRecorder |
| `metrics.NewMeterProvider(cfg)` | 创建 OTLP MeterProvider |
| `metrics.MetricsConfig` | 指标配置 |

### pkg/validator

| 类型/函数 | 说明 |
|-----------|------|
| `validator.Validate` | 校验器，包装 `go-playground/validator` |
| `validator.New()` | 创建校验器实例 |
| `validator.ValidationError` | 校验错误类型 |

### pkg/jwx

| 类型/函数 | 说明 |
|-----------|------|
| `jwx.ParseToken(token)` | 解析 JWT Token |
| `jwx.ValidateToken(token)` | 验证 Token 有效性 |

### pkg/rbac

| 类型/函数 | 说明 |
|-----------|------|
| `rbac.NewRBACService(cfg)` | 创建 RBAC 服务 |
| `rbac.NewABACService(cfg)` | 创建 ABAC 服务 |
| `rbac.Middleware` | 权限中间件 |
| `rbac.Config` | 权限配置 |

## 模块结构

```text
pkg/
├── logger/
│   ├── logger.go          # Logger 接口 + zap 实现
│   ├── config.go          # 日志配置结构
│   └── README.md
│
├── trace/
│   └── trace.go           # TracerProvider 创建 + 内存 Provider（测试用）
│
├── metrics/
│   ├── metrics.go         # MeterProvider 创建
│   ├── config.go          # 指标配置
│   ├── http.go            # HTTP 指标采集器
│   └── *_test.go
│
├── validator/
│   ├── validator.go       # 校验器封装 + params 标签解析
│   ├── errors.go          # 校验错误类型
│   └── *_test.go
│
├── jwx/
│   └── jwx.go             # JWT/JWK 解析与验证
│
└── rbac/
    ├── rbac.go            # RBAC 服务
    ├── abac.go            # ABAC 服务
    ├── service.go         # 权限服务统一入口
    ├── middleware.go       # 权限中间件
    ├── config.go          # 权限配置
    ├── errors.go          # 权限错误定义
    ├── model.conf         # Casbin RBAC 模型
    ├── abac_model.conf    # Casbin ABAC 模型
    ├── policy.csv         # RBAC 策略文件
    └── *_test.go
```

## 依赖关系

```text
pkg/ 各子包之间独立，无相互依赖

logger ──── 无内部依赖
trace  ──── 无内部依赖
metrics ─── 无内部依赖
validator ─ 无内部依赖
jwx    ──── 无内部依赖
rbac   ──── 无内部依赖
```

## 注意事项

- `pkg/` 中的包不应依赖 `internal/` 中的任何模块
- `validator.Validate` 扩展了 Fiber 的路径参数解析（`params` tag），因 Fiber 原生不支持
- `trace.NewInMemoryProvider()` 用于单元测试，避免真正导出 trace 数据
- `metrics.HTTPRecorder` 采集的指标通过 OTLP gRPC 导出到 Collector
- `rbac` 子包包含 Casbin 模型文件和策略文件，部署时需确保文件路径正确
- `pkg/` 的包设计为可被外部项目引用，API 变更需考虑向后兼容
