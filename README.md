# Go Fiber API

基于 Go Fiber 的 Web API 项目，提供开箱即用的 Web 基础框架，支持用户管理功能。

## 项目结构

```
go_sample_code/
├── cmd/
│   └── server/         # 服务入口
├── internal/
│   ├── database/       # 数据库配置 (MySQL/PostgreSQL/Redis)
│   ├── ent/            # Ent ORM 代码生成
│   ├── errno/          # 错误码定义
│   ├── handler/        # HTTP 处理器
│   │   ├── health/     # 健康检查
│   │   └── user/       # 用户管理
│   ├── middleware/     # 中间件
│   ├── repo/           # 数据访问层
│   │   └── user/       # 用户仓储
│   └── service/        # 业务逻辑层
│       └── user/       # 用户服务
├── pkg/
│   ├── jwx/            # JWT 令牌处理
│   ├── logger/         # 结构化日志
│   └── trace/          # 分布式追踪
├── config.yaml         # 配置文件
└── README.md
```

## 技术栈

- **Web 框架**: [Fiber](https://github.com/gofiber/fiber) - 高性能 Web 框架
- **ORM**: [Ent](https://entgo.io/) - Go ORM 框架，支持 MySQL/PostgreSQL
- **依赖注入**: [uber-go/fx](https://github.com/uber-go/fx) - 依赖注入框架
- **日志**: [zap](https://github.com/uber-go/zap) + OpenTelemetry
- **验证**: [go-playground/validator](https://github.com/go-playground/validator) - 参数校验
- **缓存**: [go-redis](https://github.com/redis/go-redis/v9) - Redis 客户端
- **JWT**: [lestrrat-go/jwx](https://github.com/lestrrat-go/jwx) - JWT 处理
- **配置**: YAML 配置文件支持

## 快速开始

### 配置

编辑 `config.yaml` 或使用默认配置（localhost 默认值）：

```yaml
relational:
  driver: postgres  # postgres 或 mysql
  postgres:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    dbname: demo
  mysql:
    host: localhost
    port: 3306
    user: root
    password: root
    dbname: demo

redis:
  addr: localhost:6379
```

### 运行服务

```bash
go run cmd/server/main.go
```

### API 端点

| 方法   | 路径                              | 描述             |
|--------|-----------------------------------|------------------|
| GET    | `/api/health`                     | 健康检查         |
| POST   | `/api/user`                       | 创建用户         |
| GET    | `/api/user/:id`                   | 根据 ID 获取用户 |
| GET    | `/api/user/username/:username`   | 根据用户名获取   |
| GET    | `/api/user/email/:email`          | 根据邮箱获取     |
| PUT    | `/api/user/:id`                   | 更新用户         |
| DELETE | `/api/user/:id`                   | 删除用户         |
| GET    | `/api/users`                      | 分页获取用户列表 |

### 测试示例

```bash
# 健康检查
curl http://localhost:8080/api/health

# 创建用户
curl -X POST http://localhost:8080/api/user \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# 获取用户列表
curl "http://localhost:8080/api/users?page=1&page_size=10"
```

## 开发指南

### 添加新的 Handler

1. 在 `internal/handler/` 创建新的 handler 包
2. 实现 handler 接口和构造函数
3. 在 `cmd/server/main.go` 的 `fx.Provide` 中注册
4. 在 `RegisterHooks` 中添加路由

### 中间件使用

项目已内置中间件：

- `Recovery` - panic 恢复
- `Logger` - 请求日志记录
- `Trace`

```plain
HTTP 请求
  │  traceparent 头（如果有）
  ▼
Trace 中间件 → 创建/恢复 span，注入 ctx
  ▼
Logger 中间件 → InfoCtx(ctx) → otelzap bridge
  │                ├─ TraceId = span.TraceID()   ← 写入 otel_logs.TraceId
  │                              ├─ SpanId  = span.SpanID()    ← 写入 otel_logs.SpanId
  │                              ├─ method/path/status/...     ← 写入 otel_logs.LogAttributes
  │                              └─ Body = ""                  ← 写入 otel_logs.Body
  └─ stdout JSON 也包含 trace_id/span_id（traceExtract 追加）
```

### 数据库迁移

使用 Ent 进行数据库迁移：

```bash
# 安装 Ent CLI
go install entgo.io/ent/cmd/ent@latest

# 生成代码
ent generate ./internal/ent/schema

# 创建迁移
ent migrate up ./internal/ent
```

### 运行测试

```bash
go test ./...
```

## 贡献指南

欢迎提交 Issue 和 PR！

请确保：

- 代码通过 `go fmt` 格式化
- 添加单元测试
- 更新相关文档

## License

MIT License
