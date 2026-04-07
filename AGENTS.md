# AGENTS.md

This file provides guidance to Qoder (qoder.com) when working with code in this repository.

## 项目概述

基于 Go Fiber 的 Web API 项目，使用 fx 进行依赖注入，支持从单体到微服务的渐进式演进。

## 常用命令

```bash
# 运行服务
go run cmd/server/main.go -c config.yaml

# 运行单个测试
go test ./internal/handler/health/...

# 运行所有测试
go test ./...

# 运行测试（带详细输出）
go test -v ./...

# 格式化代码
go fmt ./...

# 构建
go build ./...

# 依赖下载
go mod download
```

## 架构设计

### 依赖注入 (fx)

入口点在 `cmd/server/main.go`，使用 fx 管理生命周期：
- `fx.Provide` - 提供依赖（Logger、Fiber App、Handler）
- `fx.Invoke(RegisterHooks)` - 在 `RegisterHooks` 中注册路由和中间件
- 生命周期钩子：`OnStart`（启动）、`OnStop`（关闭）

添加新 Handler 的流程：
1. 在 `internal/handler/` 创建新包，实现 Handler 接口
2. 在 `main.go` 的 `fx.Provide` 中添加构造函数
3. 在 `RegisterHooks` 的 `OnStart` 中注册路由

### 错误处理 (internal/errno)

统一错误码系统，核心函数 `errno.Decode(data, err)`：
- 成功：返回 `OK` 及 data
- 失败：提取 wrapped errno 或返回 `InternalServerError`

预定义错误码（见 `internal/errno/code.go`）：
- 1xxx: 通用错误（1000 InternalServerError、1002 RequestParserError）
- 11xx: 认证授权（1101 TokenInvalidError、1103 PermissionDeniedError）
- 12xx: 文件资源（1201 RetrievingFileError）
- 13xx: 数据库（1301 DatabaseError）

Handler 中使用：
```go
decoded := errno.Decode("payload", nil)  // 成功
decoded := errno.Decode(nil, RequestValidateError.WithData("bad input"))  // 失败
return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
```

### 中间件 (internal/middleware)

已内置两个中间件，按注册顺序执行：
- `Recovery(log)` - panic 恢复，记录错误，返回 500
- `Logger(log)` - 请求日志（method、path、status、latency、ip）

### 日志 (pkg/logger)

基于 zap，支持：
- console/file 两种输出模式
- OpenTelemetry 集成（通过 `otel_enable` 配置）
- Context 支持（自动提取 trace_id、span_id）

### JWT (pkg/jwx)

支持两种 Token 验证：
- `SignSelfToken`/`ValidateSelfToken` - 自签发 JWT（HS256）
- `ValidateOauthToken` - 第三方 OAuth JWK 验证

### 链路追踪 (pkg/trace)

OpenTelemetry 配置，提供 `TracerProvider` 和 `Tracer`。

## 项目结构

```
cmd/server/          # 服务入口
  main.go            # fx 启动、路由注册
  config.go          # Config 结构定义

internal/
  handler/           # HTTP 处理器（按功能模块划分）
    health/          # 示例：健康检查
  middleware/        # 中间件（Recovery、Logger）
  errno/             # 错误码定义和 Decode 工具
  service/          # 业务逻辑层（未来扩展）

pkg/                 # 公共包
  logger/           # 结构化日志
  jwx/               # JWT 工具
  trace/             # OpenTelemetry 追踪
```

## 测试框架

使用 `testify`（assert/require）。Handler 测试示例见 `internal/handler/health/handler_test.go`。

## 配置文件

通过 `-c` 标志指定，默认 `config.yaml`。Logger 和 Trace 配置支持：
- `output`: console/file
- `otel_enable`: true/false
- `otel_endpoint`: OTLP 端点
- `otel_service_name`: 服务名称
