# 常见问题

## 项目相关

### Q: 这是什么项目？

A: go_sample_code 是一个基于 Go Fiber 的 Web API 基础框架，提供开箱即用的分层架构、中间件体系和可观测性支持，可作为新项目的脚手架。

### Q: 模块名为什么是 `go_sample_code`？

A: 这是 Go module 的名称，与目录名 `demo_go_framework` 不一致。模块名在 `go.mod` 中定义，导入时使用 `go_sample_code`。

### Q: 用户管理 API 可以用吗？

A: 用户管理 API 的代码已完整实现，但路由尚未在 `RegisterHooks` 中注册。需要在 `cmd/server/main.go` 中手动添加路由注册。当前仅 `/api/health` 可用。

## 配置相关

### Q: 如何切换数据库？

A: 修改配置文件中 `database.relational.driver` 字段为 `"postgres"` 或 `"mysql"`，并配置对应的连接信息。

### Q: 支持环境变量配置吗？

A: 当前版本不支持。所有配置通过 YAML 文件和 `-c` 命令行参数管理。如需环境变量支持，可在 `cmd/server/config.go` 中集成 Viper 的 `AutomaticEnv()` 功能。

### Q: 如何修改限流参数？

A: 修改 `internal/middleware/rate_limit.go` 中的 `RateLimitConfig` 结构体中的 `Rate`（每秒令牌数）和 `Burst`（桶容量）。

## 开发相关

### Q: 如何新增一个实体？

A: 1) 在 `internal/ent/schema/` 中创建 Schema 定义；2) 运行 `ent generate ./internal/ent/schema`；3) 按 Handler → Service → Repo 分层创建对应代码；4) 在 `RegisterHooks` 中注册路由。

### Q: 测试如何运行？

A: 所有测试均为单元测试，不依赖数据库/Redis：

```bash
# 运行全部测试
go test ./...

# 运行单个包的测试
go test ./internal/middleware/...
```

### Q: 为什么 Auth/RBAC 中间件未启用？

A: Auth 和 RBAC 中间件已实现但有意未注册。项目目前处于 MVP 阶段，认证和授权策略需要根据具体业务需求定制后才能启用。

## 部署相关

### Q: 生产环境推荐什么配置？

A: 建议使用 PostgreSQL + Redis，启用 trace 和 metrics（`enabled: true`），`logger.level` 设为 `info`，`trace.sampling_ratio` 根据流量调整（高流量 0.1-0.5）。

### Q: 如何实现零停机部署？

A: 项目支持优雅关闭（30s 超时），关闭顺序为 HTTP Server → Redis → DB → MeterProvider → TracerProvider → Logger。部署时先启动新实例，待健康检查通过后再关闭旧实例即可。

### Q: 支持 Kubernetes 部署吗？

A: 项目未内置 Kubernetes 配置，但可作为标准 Go 应用部署到 K8s。需要自行编写 Dockerfile、Deployment、Service 等资源配置。

## 可观测性相关

### Q: 可观测性栈有哪些版本？

A: 两个版本：

- `grafana.v1`：Grafana + Prometheus + Tempo + Loki + OTel Collector
- `grafana.v2`：Grafana + ClickHouse + OTel Collector（更高性能）

### Q: 如何查看请求链路？

A: 在 Grafana (<http://localhost:3000>) 的 Tempo 数据源中，使用 Trace ID 查询。Trace ID 可在响应头 `traceparent` 或日志中获取。
