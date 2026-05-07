# 系统架构

## 架构概述

go_sample_code 采用分层架构设计，通过 uber-go/fx 依赖注入框架实现各层之间的解耦。所有层（Handler、Service、Repo）均基于接口抽象，构造函数注入。

## 分层设计

```
┌─────────────────────────────────────────────────────────────┐
│                      Handler Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Health    │  │    User     │  │    Future Domains   │ │
│  │   Handler   │  │   Handler   │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│  ┌─────────────────────────┐  ┌─────────────────────────┐  │
│  │      User Service       │  │   Future Domain Services │ │
│  └─────────────────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       Repo Layer                             │
│  ┌─────────────────────────┐  ┌─────────────────────────┐  │
│  │       User Repo         │  │   Future Domain Repos    │ │
│  └─────────────────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Ent ORM / Database                      │
│  ┌───────────┐  ┌───────────┐  ┌─────────────────────────┐  │
│  │ PostgreSQL│  │   MySQL   │  │        Redis            │ │
│  └───────────┘  └───────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 模块职责

| 模块 | 职责 | 关键文件 |
|------|------|----------|
| `cmd/server/` | 应用入口，FX 依赖注入组装，路由注册 | `main.go`, `config.go` |
| `internal/handler/` | HTTP 请求解析、参数校验、响应构造 | `handler/user/handler.go` |
| `internal/service/` | 业务逻辑处理、业务规则校验、错误转换 | `service/user/service.go` |
| `internal/repo/` | 数据访问，Ent ORM 操作封装 | `repo/user/repo.go` |
| `internal/middleware/` | Fiber 中间件：Recovery、RateLimit、Trace、Metrics、Logger、Auth、RBAC | `middleware/*.go` |
| `internal/database/` | 数据库和 Redis 客户端初始化、健康检查 | `database/relational.go`, `redis.go` |
| `internal/ent/` | Ent ORM 生成代码（勿直接编辑） | `ent/client.go` |
| `internal/ent/schema/` | 实体定义（可编辑） | `schema/user.go` |
| `internal/errno/` | 统一错误码定义与错误响应格式 | `errno/errno.go`, `code.go` |
| `pkg/logger/` | 日志接口与 zap 实现 | `logger/logger.go` |
| `pkg/trace/` | OpenTelemetry TracerProvider | `trace/trace.go` |
| `pkg/metrics/` | HTTP 指标采集（OTLP 导出） | `metrics/metrics.go` |
| `pkg/validator/` | 参数校验封装（go-playground/validator） | `validator/validator.go` |
| `pkg/jwx/` | JWT/JWK 操作 | `jwx/jwx.go` |
| `pkg/rbac/` | Casbin RBAC/ABAC 权限控制 | `rbac/rbac.go` |
| `config/` | YAML 配置文件模板 | `config.local.yaml` |
| `devops/` | Docker Compose 部署配置 | `database/`, `monitor.v1.grafana/`, `monitor.v2.clickhouse/`, `monitor.grafana.panel/` |

## 依赖关系

```text
cmd/server (入口层，组装所有依赖)
  ↑
  ├── internal/handler (依赖 service, validator, logger, trace)
  │     ↑
  │     └── internal/service (依赖 repo, errno)
  │           ↑
  │           └── internal/repo (依赖 ent, database)
  │
  ├── internal/middleware (依赖 pkg/*)
  ├── internal/database (依赖 ent, redis)
  ├── internal/errno (无内部依赖)
  │
  └── pkg/* (公共包，可被外部引用，不依赖 internal/)
```

**依赖规则**：

- `pkg/` 是公共包，不依赖 `internal/` 中的任何模块
- `internal/` 是内部包，不对外暴露
- Handler → Service → Repo 单向依赖，上层依赖下层
- 所有层通过接口通信，禁止循环依赖

## 依赖注入

使用 `uber-go/fx` 实现依赖注入：

```
fx.New()
  ├── fx.Supply(cfgPath)           提供配置文件路径
  ├── fx.Provide(...)               注册所有构造函数
  │   ├── NewAppConfig             配置加载
  │   ├── NewLogger                日志实例
  │   ├── newEntClients            Ent + sql.DB
  │   ├── database.NewRedisClient   Redis 客户端
  │   ├── newTracerProvider         Trace Provider
  │   ├── newMeterProvider          Metrics Provider
  │   ├── NewFiberApp              Fiber 应用实例
  │   ├── NewValidator             参数校验器
  │   ├── userrepo.NewUserRepo     User Repo
  │   ├── userservice.NewUserService User Service
  │   ├── healthhandler.NewHandler  Health Handler
  │   └── userhandler.NewHandler    User Handler
  │
  └── fx.Invoke(RegisterHooks)     注册路由和中间件
```

## 中间件链

请求处理顺序（洋葱模型）：

```
Request → Recovery → RateLimit → Trace → Metrics → Logger → Handler → Response
```

| 顺序 | 中间件 | 职责 |
|:----:|--------|------|
| 1 | Recovery | Panic 恢复，返回 500 错误 |
| 2 | RateLimit | 令牌桶限流，默认 20 req/s，跳过 `/api/health` |
| 3 | Trace | OpenTelemetry 链路追踪，创建 root span |
| 4 | Metrics | HTTP 请求指标采集（Counter、Histogram、Gauge） |
| 5 | Logger | 请求日志，含 trace_id、latency、status 等字段 |

**未启用中间件**：Auth（JWT 认证）、RBAC（Casbin 权限控制）已实现但未在 RegisterHooks 中注册。

## 运行时模型

- 单进程、单实例运行
- Fiber 使用 fasthttp 作为底层 HTTP 引擎
- 未启用 Prefork 模式
- 启动超时 30s，优雅关闭超时 30s
- 关闭顺序：HTTP Server → Redis → Ent/sql.DB → MeterProvider → TracerProvider → Logger

## 错误处理策略

- 所有错误通过 `errno.Errno` 接口统一表示
- `errno.Decode(data, err)` 函数统一转换为 `{code, message, data}` JSON 响应
- 错误码分段：
  - `0`：OK
  - `1000-1099`：通用错误（Internal、参数解析、参数校验等）
  - `1100-1199`：认证/授权错误（Token无效、权限不足、限流等）
  - `1200-1299`：文件相关错误
  - `1300+`：数据库错误
  - `2000+`：业务错误（用户不存在、用户已存在等）
  - `3000+`：分页相关错误
- 底层错误向上传播，在 Handler 层统一转换为 HTTP 响应
