# 配置指南

## 概述

项目支持 YAML 配置文件和命令行参数两种配置方式。

## 启动命令

```bash
# 使用默认配置
go run cmd/server/main.go

# 指定配置文件
go run cmd/server/main.go -c config/config.local.yaml
```

## 配置文件结构

```yaml
# 监听地址
common:
  listen: ":8080"           # 默认 ":8080"

# 数据库配置
database:
  relational:
    driver: "postgres"     # postgres 或 mysql
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
      conn_max_lifetime: 3600  # 秒

  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    pool_size: 100

  health_check_timeout: 5  # 秒

# 日志配置
logger:
  level: "info"            # debug, info, warn, error
  encoding: "json"         # json 或 console
  output_paths:            # 日志输出路径
    - "stdout"

# 链路追踪配置
trace:
  enabled: true
  endpoint: "localhost:4317"  # OTLP gRPC 地址
  service_name: "go_sample_code"
  sampling_ratio: 1.0      # 采样比例 0.0-1.0

# 指标配置
metrics:
  enabled: true
  endpoint: "localhost:4317"  # OTLP gRPC 地址
  service_name: "go_sample_code"
```

## 配置默认值

如未提供配置文件，使用以下默认值：

| 配置项 | 默认值 |
|--------|--------|
| `common.listen` | `:8080` |
| `database.relational.driver` | `postgres` |
| `database.relational.postgres.host` | `localhost` |
| `database.relational.postgres.port` | `5432` |
| `database.relational.postgres.user` | `postgres` |
| `database.relational.postgres.password` | `postgres` |
| `database.relational.postgres.dbname` | `demo` |
| `database.redis.addr` | `localhost:6379` |
| `logger.level` | `info` |
| `logger.encoding` | `json` |
| `trace.enabled` | `true` |
| `metrics.enabled` | `true` |

## 配置加载流程

```
main.go
    │
    ├─ flag.Parse()           解析命令行参数
    │
    ├─ NewAppConfig()         加载配置文件
    │       │
    │       ├─ viper.SetConfigFile()  设置配置文件路径
    │       ├─ viper.ReadInConfig()   读取配置文件
    │       └─ config.Unmarshal()     反序列化到结构体
    │
    └─ fx.Provide()           注入配置到各组件
```

## 配置文件模板

参考 `config/config.template.yaml`：

```yaml
common:
  listen: ":8080"

database:
  relational:
    driver: "postgres"
    postgres:
      host: "localhost"
      port: 5432
      user: "postgres"
      password: "postgres"
      dbname: "demo"
      max_open_conns: 100
      max_idle_conns: 10
      conn_max_lifetime: 3600
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
  health_check_timeout: 5

logger:
  level: "info"
  encoding: "json"
  output_paths:
    - "stdout"

trace:
  enabled: true
  endpoint: "localhost:4317"
  service_name: "go_sample_code"
  sampling_ratio: 1.0

metrics:
  enabled: true
  endpoint: "localhost:4317"
  service_name: "go_sample_code"
```

## 环境变量

当前版本不支持环境变量配置，如需支持可通过 Viper 的 `AutomaticEnv()` 功能扩展。

## 数据库配置

### PostgreSQL

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

## Redis 配置

```yaml
database:
  redis:
    addr: "localhost:6379"
    password: ""       # 无密码则留空
    db: 0             # DB 编号
    pool_size: 100    # 连接池大小
```

## 可观测性配置

### 链路追踪 (Trace)

```yaml
trace:
  enabled: true
  endpoint: "localhost:4317"      # OTLP Collector 地址
  service_name: "go_sample_code"
  sampling_ratio: 1.0              # 1.0 = 100% 采样
```

### 指标 (Metrics)

```yaml
metrics:
  enabled: true
  endpoint: "localhost:4317"      # OTLP Collector 地址
  service_name: "go_sample_code"
```
