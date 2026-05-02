# Database 模块

> 数据库与 Redis 客户端管理，提供连接初始化、健康检查和优雅关闭。

## 设计决策

### 为什么需要这个模块

将数据库连接管理与业务代码分离，统一管理连接池配置、健康检查和生命周期。支持 PostgreSQL 和 MySQL 双驱动切换。

### 为什么这么设计

- **选择了**：通过 YAML 配置驱动的 `driver` 字段选择数据库类型
- **而不是**：编译时固定一种数据库驱动
- **原因**：运行时切换数据库类型，适应不同部署环境

## 关键类型与接口

### database.DatabaseConfig

- **定义位置**：`internal/database/config.go`
- **用途**：数据库配置结构，包含关系型数据库和 Redis 配置

```go
type DatabaseConfig struct {
    Relational RelationalConfig `yaml:"relational"`
    Redis      RedisConfig      `yaml:"redis"`
    HealthCheckTimeout int      `yaml:"health_check_timeout"`
}
```

### database.HealthChecker

- **定义位置**：`internal/database/health.go`
- **用途**：健康检查接口，检查 DB 和 Redis 连通性

### 关键函数

| 函数 | 文件 | 说明 |
|------|------|------|
| `NewEntClient(cfg, log)` | `relational.go` | 创建 `*sql.DB` 和 `*ent.Client` |
| `NewRedisClient(cfg)` | `redis.go` | 创建 `*redis.Client` |
| `NewHealthChecker(db, redis, driver, log, timeout)` | `health.go` | 创建健康检查器 |
| `CloseEntClient(db, entClient, log)` | `relational.go` | 优雅关闭 Ent 和 SQL 连接 |
| `CloseRedisClient(client, log)` | `redis.go` | 优雅关闭 Redis 连接 |

## 模块结构

```text
internal/database/
├── config.go          # 数据库配置结构体 + 默认值 + 验证
├── relational.go      # sql.DB + Ent Client 初始化
├── redis.go           # Redis Client 初始化
├── health.go          # 健康检查接口
└── *_test.go          # 单元测试
```

| 文件 | 职责 |
|------|------|
| `config.go` | 定义 `DatabaseConfig`、`RelationalConfig`、`RedisConfig`，含 `ApplyDefaults()` 和 `Validate()` |
| `relational.go` | 根据 driver 类型创建对应的 `sql.DB`，创建 Ent Client，自动迁移 schema |
| `redis.go` | 创建 Redis 客户端，含连接池配置 |
| `health.go` | 实现 `HealthChecker` 接口，Ping DB 和 Redis |

## 与其他模块的关系

### 依赖

- **ent**：Ent ORM Client
- **pkg/logger**：日志接口
- **github.com/lib/pq**：PostgreSQL 驱动
- **github.com/go-sql-driver/mysql**：MySQL 驱动
- **github.com/redis/go-redis/v9**：Redis 客户端

### 被依赖

- **cmd/server**：通过 fx.Provide 注入 DB 和 Redis 客户端

## 注意事项

- 连接池参数配置在 `RelationalConfig` 中（`max_open_conns`、`max_idle_conns`、`conn_max_lifetime`）
- Redis 连接池通过 `pool_size` 配置
- 优雅关闭顺序：先关 HTTP Server，再关 Redis，最后关 DB
- Ent schema 迁移在 `NewEntClient` 中自动执行（`client.Schema.Create(ctx)`）
- 健康检查超时通过 `health_check_timeout` 配置（秒），默认 5s
