# 配置说明

## 概述

项目使用 YAML 格式的配置文件，通过 `-c` 命令行参数指定路径。如未提供配置文件，使用硬编码默认值。

## 配置文件

### 启动时指定配置

```bash
go run cmd/server/main.go -c config/config.local.yaml
```

### 不指定配置（使用默认值）

```bash
go run cmd/server/main.go
```

此时使用硬编码默认值：PostgreSQL `localhost:5432`，Redis `localhost:6379`。

## 配置文件结构

```yaml
# 通用配置
common:
  listen: ":8080"           # 监听地址，默认 ":8080"

# 数据库配置
database:
  relational:
    driver: "postgres"      # postgres 或 mysql
    postgres:
      host: "localhost"
      port: 5432
      user: "postgres"
      password: "postgres"
      dbname: "demo"
      max_open_conns: 100
      max_idle_conns: 10
      conn_max_lifetime: 3600  # 秒
    mysql:
      host: "localhost"
      port: 3306
      user: "root"
      password: "root"
      dbname: "demo"
      max_open_conns: 100
      max_idle_conns: 10
      conn_max_lifetime: 3600

  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    pool_size: 100

  health_check_timeout: 5   # 健康检查超时（秒）

# 日志配置
logger:
  level: "info"             # debug / info / warn / error
  encoding: "json"          # json 或 console
  output_paths:
    - "stdout"

# 链路追踪配置
trace:
  enabled: true
  endpoint: "localhost:4317"   # OTLP gRPC 地址
  service_name: "go_sample_code"
  sampling_ratio: 1.0          # 采样比例 0.0-1.0

# 指标配置
metrics:
  enabled: true
  endpoint: "localhost:4317"
  service_name: "go_sample_code"
```

## 默认值一览

| 配置项 | 默认值 |
|--------|--------|
| `common.listen` | `:8080` |
| `database.relational.driver` | `postgres` |
| `database.relational.postgres.host` | `localhost` |
| `database.relational.postgres.port` | `5432` |
| `database.relational.postgres.user` | `postgres` |
| `database.relational.postgres.password` | `postgres` |
| `database.relational.postgres.dbname` | `demo` |
| `database.relational.postgres.max_open_conns` | `100` |
| `database.relational.postgres.max_idle_conns` | `10` |
| `database.relational.postgres.conn_max_lifetime` | `3600` |
| `database.redis.addr` | `localhost:6379` |
| `database.redis.pool_size` | `100` |
| `database.health_check_timeout` | `5` |
| `logger.level` | `info` |
| `logger.encoding` | `json` |
| `trace.enabled` | `true` |
| `trace.service_name` | `demo-go-framework` |
| `trace.sampling_ratio` | `1.0` |
| `metrics.enabled` | `true` |

## 数据库切换

### PostgreSQL（默认）

```yaml
database:
  relational:
    driver: "postgres"
    postgres:
      host: "localhost"
      port: 5432
      user: "postgres"
      password: "postgres"
      dbname: "demo"
```

### MySQL

```yaml
database:
  relational:
    driver: "mysql"
    mysql:
      host: "localhost"
      port: 3306
      user: "root"
      password: "root"
      dbname: "demo"
```

## 配置模板

完整配置模板参考 `config/config.template.yaml`，本地开发配置参考 `config/config.local.yaml`。

## 环境变量

当前版本不支持环境变量覆盖配置。如需支持，可在代码中集成 Viper 的 `AutomaticEnv()` 功能。
