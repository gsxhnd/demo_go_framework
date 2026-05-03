# Go Fiber API — 学习与演示项目 🎓

> ⚠️ **这是一个 Demo 项目，主要用于学习和实践 Go Web 后端开发。**
> 项目整合了 Go 生态中常见的主流库和框架，通过一个完整的用户管理 CRUD 示例，展示了分层架构、依赖注入、中间件链、可观测性等工程实践。**不推荐直接用于生产环境。**

## 🎯 学习目标

本项目旨在帮助你理解和掌握以下 Go 后端开发中的核心概念：

| 主题 | 涉及内容 |
|------|---------|
| **项目结构** | 分层架构（Handler → Service → Repo），`internal/pkg` 目录规范 |
| **依赖注入** | `uber-go/fx` 的构造函数注入与生命周期管理 |
| **ORM 实践** | Ent ORM 的 Schema 定义、代码生成、CRUD 操作 |
| **中间件设计** | Recovery / Logger / Trace / RateLimit / Auth / RBAC 中间件链 |
| **错误处理** | 统一错误码体系（`errno`），HTTP 状态码映射 |
| **可观测性** | OpenTelemetry 分布式追踪 + 结构化日志 + Prometheus 指标 |
| **参数校验** | `go-playground/validator` 与自定义校验器 |
| **配置管理** | YAML 配置文件加载与默认值策略 |
| **测试实践** | 基于 `testify` 的分层单元测试（mock 下层接口） |

适合人群：正在学习 Go Web 开发，希望了解工程化项目组织方式的中级开发者。

## 📁 项目结构

```
go_sample_code/
├── cmd/
│   └── server/         # 服务入口
├── config/             # 配置文件
├── devops/             # Docker 编排 & 监控套件（开箱即用）
│   ├── database/       # 数据库 + Redis
│   ├── grafana.v1/     # 经典监控栈（Prometheus / Tempo / Loki）
│   └── grafana.v2/     # 现代监控栈（ClickHouse + Grafana OSS）
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
│   ├── metrics/        # 指标采集
│   ├── rbac/           # 角色权限控制
│   ├── trace/          # 分布式追踪
│   └── validator/      # 参数校验
└── README.md
```

## 🛠 技术栈

本项目刻意整合了多种 Go 生态主流库，帮助你一次性了解它们的基本用法与配合方式：

- **Web 框架**: [Fiber](https://github.com/gofiber/fiber) — 高性能 HTTP 框架（类 Express API）
- **ORM**: [Ent](https://entgo.io/) — 类型安全的 Go ORM，支持 MySQL / PostgreSQL
- **依赖注入**: [uber-go/fx](https://github.com/uber-go/fx) — 构造函数注入与生命周期管理
- **日志**: [zap](https://github.com/uber-go/zap) + OpenTelemetry — 结构化日志 + 链路追踪
- **参数校验**: [go-playground/validator](https://github.com/go-playground/validator) — 结构体标签校验
- **缓存**: [go-redis](https://github.com/redis/go-redis/v9) — Redis 客户端
- **JWT**: [lestrrat-go/jwx](https://github.com/lestrrat-go/jwx) — JWT 令牌签发与解析
- **配置**: YAML 配置文件，`-c` 参数指定路径
- **链路追踪**: OpenTelemetry（对接 Grafana + Tempo / ClickHouse）
- **容器化**: Docker Compose 编排（数据库 + 两套可观测性方案）

## 🚀 快速开始

### 前置条件

- Go 1.25+
- Docker（用于启动数据库和 Redis）

### 1. 启动基础设施

```bash
# 启动 PostgreSQL + Redis（MySQL 可选）
docker compose -f devops/database/docker-compose.yml up -d

# （可选）启动可观测性套件 — 二选一：
# 方案 A: 经典 Grafana 栈（Prometheus + Tempo + Loki）
docker compose -f devops/grafana.v1/docker-compose.yml up -d
# 方案 B: ClickHouse 方案（现代化列式存储）
# 详见 devops/grafana.v2/README.md
```

### 2. 配置文件

编辑 `config/config.local.yaml` 或直接使用默认值。程序启动时通过 `-c` 指定配置文件：

```yaml
relational:
  driver: postgres  # 支持 postgres 或 mysql
  postgres:
    host: localhost
    port: 5432
    user: demo_user
    password: demo_password
    dbname: demo_db
  mysql:
    host: localhost
    port: 3306
    user: demo_user
    password: demo_password
    dbname: demo_db

redis:
  addr: localhost:6379
```

### 3. 运行服务

```bash
go run cmd/server/main.go -c config/config.local.yaml
```

### 4. 验证服务

```bash
# 健康检查
curl http://localhost:8080/api/health
# → {"code":0,"message":"OK","data":{"database":"up","redis":"up"}}
```

### API 端点

以下为用户管理的完整 CRUD 接口（演示分层架构中各层的协作方式）：

| 方法   | 路径                              | 描述             |
|--------|-----------------------------------|------------------|
| GET    | `/api/health`                     | 健康检查（DB + Redis） |
| POST   | `/api/user`                       | 创建用户         |
| GET    | `/api/user/:id`                   | 根据 ID 获取用户 |
| GET    | `/api/user/username/:username`   | 根据用户名获取   |
| GET    | `/api/user/email/:email`          | 根据邮箱获取     |
| PUT    | `/api/user/:id`                   | 更新用户         |
| DELETE | `/api/user/:id`                   | 删除用户         |
| GET    | `/api/users`                      | 分页获取用户列表  |

### 接口测试

```bash
# 创建用户
curl -X POST http://localhost:8080/api/user \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# 获取用户列表（分页）
curl "http://localhost:8080/api/users?page=1&page_size=10"

# 根据 ID 获取
curl http://localhost:8080/api/user/1

# 更新用户
curl -X PUT http://localhost:8080/api/user/1 \
  -H "Content-Type: application/json" \
  -d '{"username":"updated"}'

# 删除用户
curl -X DELETE http://localhost:8080/api/user/1
```

## 📖 开发指南

### 核心设计理念

本项目采用 **分层架构 + 依赖注入** 的设计，每一层都通过接口定义，便于测试和解耦：

```
HTTP Request → Handler（参数校验 & 响应）→ Service（业务逻辑）→ Repo（数据访问）
                    ↑                            ↑                     ↑
               Fiber + validator            fx 注入接口           Ent ORM
```

### 如何添加新功能

以「新增一个实体 CRUD」为例，体验完整的分层开发流程：

1. **定义 Schema** — 在 `internal/ent/schema/` 创建实体定义
2. **生成代码** — 运行 `ent generate ./internal/ent/schema`
3. **实现 Repo 层** — 在 `internal/repo/` 创建数据访问接口与实现
4. **实现 Service 层** — 在 `internal/service/` 编写业务逻辑
5. **实现 Handler 层** — 在 `internal/handler/` 处理 HTTP 请求
6. **注册路由** — 在 `cmd/server/main.go` 的 `RegisterHooks` 中添加路由，并用 `fx.Provide` 注册构造函数

### 中间件链

请求经过以下中间件（按顺序），是学习中间件设计的良好范例：

```plain
HTTP 请求
  │  traceparent 头（如果有）
  ▼
Recovery → RateLimit → Trace → Metrics → Logger → 业务 Handler
```

具体流程：

```plain
Trace 中间件 → 创建/恢复 span，注入 ctx
  ▼
Logger 中间件 → InfoCtx(ctx) → otelzap bridge
  │                ├─ TraceId = span.TraceID()   ← 写入 otel_logs.TraceId
  │                ├─ SpanId  = span.SpanID()    ← 写入 otel_logs.SpanId
  │                ├─ method/path/status/...     ← 写入 otel_logs.LogAttributes
  │                └─ Body = ""                  ← 写入 otel_logs.Body
  └─ stdout JSON 也包含 trace_id/span_id
```

> ⚠️ 注意：Auth 和 RBAC 中间件已实现但**尚未挂载**到中间件链，你可以将其作为练习自行接入。

### Ent 代码生成

```bash
# 安装 Ent CLI
go install entgo.io/ent/cmd/ent@latest

# 修改 internal/ent/schema/ 下的定义后，重新生成
ent generate --feature sql/modifier ./internal/ent/schema
```

> 💡 **提示**：生成的代码在 `internal/ent/` 目录下，**请勿手动编辑**——只修改 `schema/` 目录下的定义文件。

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试（带详细输出）
go test -v ./internal/middleware/...
go test -v ./internal/handler/...
```

测试策略：所有测试均为**单元测试**，无需数据库/Redis。通过 mock 下层接口来隔离测试目标层。

### API 文档 (Swagger)

```bash
# 生成 Swagger API 文档
swag init -d ./cmd/server -g main.go --outputTypes json --pdl 3
```

```bash
# 安装 swag CLI（如未安装）
go install github.com/swaggo/swag/cmd/swag@latest
```

## 📚 延伸学习

在理解本 Demo 的基础上，建议进一步探索：

- 接入 Auth + RBAC 中间件，实现认证与授权
- 添加更多业务实体（如订单、商品），练习复杂关联查询
- 切换并对比两套监控方案：`grafana.v1`（经典 LGTM 栈）vs `grafana.v2`（ClickHouse）
- 在 Grafana 中构建自定义 Dashboard，可视化链路追踪和性能指标
- 编写集成测试，覆盖完整的请求-响应链路
- 尝试替换 Fiber 为 `net/http` + `chi` 路由，对比框架差异

> 📖 更多基础设施配置细节见 [`devops/README.md`](devops/README.md)。

## 📄 License

MIT License — 请自由使用、修改和学习本项目代码。
