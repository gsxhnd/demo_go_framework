# 技术栈

## 运行时与语言

| 项目 | 版本/说明 |
|------|-----------|
| Go | 1.25.0 |
| 模块名 | `go_sample_code` |

## 核心依赖

### Web 框架

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/gofiber/fiber/v2` | v2.52.12 | 高性能 HTTP Web 框架（基于 fasthttp） |

### ORM

| 依赖 | 版本 | 用途 |
|------|------|------|
| `entgo.io/ent` | v0.14.6 | 类型安全的 Go ORM，代码生成 |

### 依赖注入

| 依赖 | 版本 | 用途 |
|------|------|------|
| `go.uber.org/fx` | v1.24.0 | 依赖注入框架，构造函数注入 |

### 日志

| 依赖 | 版本 | 用途 |
|------|------|------|
| `go.uber.org/zap` | v1.27.1 | 高性能结构化日志 |
| `github.com/DeRuina/timberjack` | v1.4.1 | 日志轮转 |

### 可观测性

| 依赖 | 版本 | 用途 |
|------|------|------|
| `go.opentelemetry.io/otel` | v1.43.0 | OpenTelemetry SDK |
| `go.opentelemetry.io/otel/trace` | v1.43.0 | 链路追踪 API |
| `go.opentelemetry.io/otel/metric` | v1.43.0 | 指标 API |
| `go.opentelemetry.io/contrib/bridges/otelzap` | v0.18.0 | OTel + Zap 桥接 |

### 数据库驱动

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/lib/pq` | v1.10.9 | PostgreSQL 驱动 |
| `github.com/go-sql-driver/mysql` | v1.9.3 | MySQL 驱动 |

### 缓存

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/redis/go-redis/v9` | v9.18.0 | Redis 客户端 |

### 认证与授权

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/lestrrat-go/jwx/v3` | v3.0.13 | JWT/JWK 操作 |
| `github.com/casbin/casbin/v2` | v2.92.0 | RBAC/ABAC 权限控制 |

### 参数校验

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/go-playground/validator/v10` | v10.30.2 | 结构体标签校验 |

### 工具

| 依赖 | 版本 | 用途 |
|------|------|------|
| `gopkg.in/yaml.v3` | v3.0.1 | YAML 配置解析 |
| `golang.org/x/time` | v0.15.0 | 令牌桶限流算法 |
| `github.com/stretchr/testify` | v1.11.1 | 测试断言库 |

## 构建与工具链

### 常用命令

```bash
# 运行服务
go run cmd/server/main.go -c config/config.local.yaml

# 运行测试 — 单个包或全部
go test ./internal/middleware/...
go test ./...

# Ent 代码生成（修改 internal/ent/schema/ 后运行）
go install entgo.io/ent/cmd/ent@latest
ent generate ./internal/ent/schema

# 格式化代码
go fmt ./...

# 编译二进制
go build -o server ./cmd/server
```

### 代码生成

- **Ent ORM**：修改 `internal/ent/schema/*.go` 后手动运行 `ent generate ./internal/ent/schema`
- 无 `go:generate` 指令，纯手动触发
- 生成的代码位于 `internal/ent/`，禁止直接编辑

### 配置管理

- YAML 格式配置文件，使用 `yaml` struct tags
- 通过 `-c` 命令行参数指定配置文件路径
- 无配置文件时使用硬编码默认值
- 不支持环境变量覆盖（可通过 Viper 扩展）

## 开发环境

### 本地开发

```bash
# 启动数据库和 Redis
docker compose -f devops/database/docker-compose.yml up -d

# 启动可观测性栈（可选）
docker compose -f devops/monitor.v1.grafana/docker-compose.yml up -d
docker compose -f devops/monitor.grafana.panel/docker-compose.yml up -d

# 启动服务
go run cmd/server/main.go -c config/config.local.yaml
```

### Docker Compose 服务

| 目录 | 服务 | 端口 |
|------|------|------|
| `devops/database/` | PostgreSQL + Redis | 5432, 6379 |
| `devops/monitor.v1.grafana/` | Prometheus + Tempo + Loki + OTel Collector | 9090, 3100, 3200, 4317 |
| `devops/monitor.v2.clickhouse/` | ClickHouse + OTel Collector | 8123, 9000, 4317 |
| `devops/monitor.grafana.panel/` | Grafana OSS 13（v1+v2 数据源） | 3000 |

## 测试策略

- 所有测试均为单元测试级别，不依赖 DB/Redis
- 使用 `stretchr/testify` 进行断言
- Mock 模式：Handler 测试 mock Service 接口，Service 测试 mock Repo 接口
- 内存 Tracer：`trace.NewInMemoryProvider()` 用于测试
- 测试文件位置：`internal/*/` 和 `test/` 目录
