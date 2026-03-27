# Go Fiber API

基于 Go Fiber 的 Web API 项目，支持从单体模板到微服务架构的渐进式演进。

## 项目定位

**第一阶段：API Web Framework Template**

提供开箱即用的 Web API 基础框架，包含：
- HTTP 路由与中间件（Fiber）
- 依赖注入（fx）
- 结构化日志（zap + OpenTelemetry）
- 健康检查端点

**第二阶段：微服务集成**

在第一阶段基础上扩展微服务能力：
- 服务注册与发现
- 分布式链路追踪
- 负载均衡与熔断
- 消息队列集成

## 项目结构

```
go_sample_code/
├── cmd/
│   └── server/         # 服务入口
├── internal/           # 内部代码
│   ├── handler/        # HTTP 处理器
│   ├── middleware/     # 中间件
│   ├── service/        # 业务逻辑
│   └── errno/          # 错误码
├── pkg/                # 公共包
│   └── logger/         # 结构化日志
└── README.md
```

## 快速开始

### 运行服务

```bash
go run cmd/server/main.go
```

### 测试健康检查

```bash
curl http://localhost:8080/api/health
```

## 技术栈

- **Web 框架**: [Fiber](https://github.com/gofiber/fiber) - 高性能 Web 框架
- **依赖注入**: [uber-go/fx](https://github.com/uber-go/fx) - 依赖注入框架
- **日志**: [zap](https://github.com/uber-go/zap) + OpenTelemetry
- **配置**: 支持 YAML/JSON 配置文件

## 未来规划

### 第二阶段待办

- [ ] 集成 gRPC 服务
- [ ] 添加服务注册中心（Consul/etcd）
- [ ] 实现分布式链路追踪（Jaeger）
- [ ] 添加消息队列支持（Kafka/RabbitMQ）
- [ ] 实现熔断和重试机制

## 开发指南

### 添加新的 Handler

1. 在 `internal/handler/` 创建新的 handler 包
2. 实现 handler 接口和构造函数
3. 在 `cmd/server/main.go` 的 `RegisterHooks` 中注册路由

### 中间件使用

项目已内置中间件：
- `Recovery` -  panic 恢复
- `Logger` - 请求日志记录

在 `RegisterHooks` 的 `OnStart` 中按顺序添加中间件。

## 贡献指南

欢迎提交 Issue 和 PR！

请确保：
- 代码通过 `go fmt` 格式化
- 添加单元测试
- 更新相关文档

## License

MIT License
