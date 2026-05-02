# DevOps 文档

Demo Go Framework 开发环境配置指南。所有配置已内置固定值，开箱即用，无需额外配置。

## 目录结构

```
devops/
├── database/
│   └── docker-compose.yml       # 数据库服务（MySQL, PostgreSQL, Redis）
├── grafana.v1/
│   ├── docker-compose.yml       # Grafana v1 监控栈（Prometheus + Tempo + Loki + OTel）
│   ├── grafana/
│   │   ├── config.ini
│   │   ├── dashboards/http-metrics.json
│   │   └── provisioning/
│   ├── loki/config.yaml
│   ├── otel-collector/config.yaml
│   ├── prometheus/
│   │   ├── config.yaml
│   │   └── rules/alerts.yml
│   └── tempo/config.yaml
├── grafana.v2/
│   ├── README.md                # Grafana v2 详细文档
│   ├── clickhouse/
│   │   ├── docker-compose.yml
│   │   └── init.sql
│   ├── grafana/
│   │   ├── docker-compose.yml
│   │   ├── config.ini
│   │   ├── otel_dashboard.json
│   │   └── provisioning/
│   └── otel-collector/
│       ├── config.yaml
│       └── docker-compose.yml
├── config/
│   └── config.template.yaml     # 应用配置模板
├── databases/
│   ├── mysql/init.sql
│   └── postgres/init.sql
└── README.md
```

## 快速开始

> 所有服务端口和凭证已内置固定值，直接启动即可，无需复制 `.env` 文件。

### 1. 启动数据库

```bash
cd devops/database
docker compose up -d
```

### 2. 启动监控栈（二选一）

`grafana.v1` 和 `grafana.v2` 是两套**不同的监控实现方案**，只需启动其中一套：

**方案 A — Grafana v1**（Prometheus + Tempo + Loki + OTel Collector）

```bash
cd devops/grafana.v1
docker compose up -d
```

- Grafana: <http://localhost:3000（匿名> Admin 访问）

**方案 B — Grafana v2**（ClickHouse + OTel Collector + Grafana OSS 13）

```bash
# 先启动 ClickHouse
cd devops/grafana.v2/clickhouse
docker compose up -d

# 初始化 OTel 数据表（仅首次）
docker exec -i demo-clickhouse clickhouse-client --multiquery < init.sql

# 启动 OTel Collector
cd devops/grafana.v2/otel-collector
docker compose up -d

# 启动 Grafana v2
cd devops/grafana.v2/grafana
docker compose up -d
```

- Grafana v2: <http://localhost:3000（匿名> Admin 访问）

### 3. 停止服务

```bash
# 停止数据库
cd devops/database && docker compose down

# 停止 Grafana v1
cd devops/grafana.v1 && docker compose down

# 停止 Grafana v2
cd devops/grafana.v2/grafana && docker compose down
cd devops/grafana.v2/otel-collector && docker compose down
cd devops/grafana.v2/clickhouse && docker compose down

# 停止并删除数据卷（加 -v 参数）
docker compose down -v
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Go App | 8080 | 应用服务 |
| MySQL | 3306 | MySQL 数据库 |
| PostgreSQL | 5432 | PostgreSQL 数据库 |
| Redis | 6379 | Redis 缓存 |
| Grafana (v1 / v2) | 3000 | 可视化看板（两者互斥，不冲突） |
| Prometheus | 9090 | Metrics 存储（仅 v1） |
| Tempo | 3200 | 分布式追踪（仅 v1） |
| Loki | 3100 | 日志存储（仅 v1） |
| OTel Collector | 4317/4318 | OTLP gRPC / HTTP 接收 |
| OTel Collector metrics | 8888 | Collector 自身 metrics |
| OTel Collector health | 13133 | Health check |
| Prometheus scrape | 9464 | OTel Prometheus exporter（仅 v1） |
| ClickHouse HTTP | 8123 | HTTP 查询接口（仅 v2） |
| ClickHouse TCP | 9000 | Native TCP 接口（仅 v2） |

> **注意**：grafana.v1 和 grafana.v2 共用端口（OTel: 4317/4318, Grafana: 3000），这是有意设计——两套方案只需启动其中一套，不会同时运行。

## 默认凭证

### Grafana (v1 & v2)

- URL: <http://localhost:3000>
- 匿名访问已启用，角色为 Admin
- 手动登录: `admin` / `admin`

### 数据库

- MySQL: `demo_user` / `demo_password`，库 `demo_db`
- PostgreSQL: `demo_user` / `demo_password`，库 `demo_db`

## 应用配置

启动基础设施后，在 `config/config.local.yaml` 中配置应用连接信息（或直接使用默认值）：

```yaml
database:
  relational:
    driver: postgres  # 或 mysql
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

OTel 追踪配置（v1 和 v2 使用相同的 OTLP 端口）：

```yaml
trace:
  endpoint: localhost:4317
```

## 监控方案对比

| 维度 | Grafana v1 | Grafana v2 |
|------|-----------|-----------|
| 数据存储 | Prometheus + Tempo + Loki | ClickHouse 统一存储 |
| 组件数量 | 5 个服务 | 3 个服务 |
| Grafana 版本 | 10.3.3 | 13.0.1 (OSS) |
| 数据源插件 | 内置 (Prometheus / Tempo / Loki) | grafana-clickhouse-datasource |
| 适用场景 | 经典 Grafana 技术栈学习 | 列式存储方案学习 |

## 健康检查

| 服务 | 端点 |
|------|------|
| Go App | <http://localhost:8080/api/health> |
| OTel Collector | <http://localhost:13133/health> |
| Prometheus | <http://localhost:9090/-/healthy> |
| Grafana | <http://localhost:3000/api/health> |
| Tempo | <http://localhost:3200/ready> |
