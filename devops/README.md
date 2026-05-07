# DevOps 文档

Demo Go Framework 开发环境配置指南。所有配置已内置固定值，开箱即用，无需额外配置。

## 目录结构

```
devops/
├── database/
│   └── docker-compose.yml       # 数据库服务（MySQL, PostgreSQL, Redis）
├── monitor.v1.grafana/
│   ├── docker-compose.yml       # LGTM 后端：Prometheus + Tempo + Loki + OTel（无 Grafana UI）
│   ├── README.md
│   ├── loki/config.yaml
│   ├── otel-collector/config.yaml
│   ├── prometheus/
│   │   ├── config.yaml
│   │   └── rules/alerts.yml
│   └── tempo/config.yaml
├── monitor.v2.clickhouse/
│   ├── README.md                # ClickHouse + OTel 详细步骤
│   ├── clickhouse/
│   │   ├── docker-compose.yml
│   │   └── init.sql
│   └── otel-collector/
│       ├── docker-compose.yml
│       └── config.yaml
├── monitor.v3.influxdb/
│   ├── README.md                # InfluxDB 2 + OTel（metrics/traces/logs）
│   ├── docker-compose.yml
│   └── otel-collector/
│       └── config.yaml
├── monitor.grafana.panel/
│   ├── docker-compose.yml       # Grafana OSS 13（v1 + v2 数据源）
│   ├── README.md
│   └── grafana/
│       ├── config.ini
│       ├── dashboards/
│       └── provisioning/        # 含 v1 / v2 / v3 数据源
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

### 2. 启动监控后端（二选一或按需）

**方案 A — monitor.v1.grafana**（Prometheus + Tempo + Loki + OTel Collector）

```bash
cd devops/monitor.v1.grafana
docker compose up -d
```

**方案 B — monitor.v2.clickhouse**（ClickHouse + OTel Collector）

```bash
cd devops/monitor.v2.clickhouse/clickhouse
docker compose up -d
docker exec -i demo-clickhouse clickhouse-client --multiquery < init.sql
cd ../otel-collector
docker compose up -d
```

详见 `monitor.v2.clickhouse/README.md`。

**方案 C — monitor.v3.influxdb**（InfluxDB 2 + OTel Collector）

```bash
cd devops/monitor.v3.influxdb
docker compose up -d
```

详见 `monitor.v3.influxdb/README.md`。应用 OTLP 请指向 **4319（gRPC）/ 4320（HTTP）**，与 v1/v2 的 4317/4318 区分。

### 3. 启动 Grafana 面板（统一 UI）

```bash
cd devops/monitor.grafana.panel
docker compose up -d
```

- Grafana: <http://localhost:3000>（匿名 Admin）

预置数据源：Prometheus、Tempo、Loki（依赖方案 A）、ClickHouse（依赖方案 B）、InfluxDB（依赖方案 C）。未启动的后端在 Explore 中会失败，属正常。

### 4. 停止服务

```bash
cd devops/database && docker compose down

cd devops/monitor.v1.grafana && docker compose down
cd devops/monitor.v3.influxdb && docker compose down

cd devops/monitor.grafana.panel && docker compose down
cd devops/monitor.v2.clickhouse/otel-collector && docker compose down
cd ../clickhouse && docker compose down

# 停止并删除数据卷（加 -v）
docker compose down -v
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Go App | 8080 | 应用服务 |
| MySQL | 3306 | MySQL 数据库 |
| PostgreSQL | 5432 | PostgreSQL 数据库 |
| Redis | 6379 | Redis 缓存 |
| Grafana | 3000 | `monitor.grafana.panel` |
| Prometheus | 9090 | 仅 v1 后端 |
| Tempo | 3200 | 仅 v1 后端 |
| Loki | 3100 | 仅 v1 后端 |
| OTel Collector | 4317/4318 | v1 与 v2 **互斥**占用，勿同时启动两套 Collector |
| OTel Collector (v3) | 4319/4320 | `monitor.v3.influxdb`，与 v1/v2 端口错开 |
| OTel Collector metrics | 8888 | v1/v2 Collector 自身 metrics |
| OTel Collector metrics (v3) | 8889 | v3 Collector 自身 metrics |
| OTel Collector health | 13133 | v1/v2 Health check |
| OTel Collector health (v3) | 13134 | v3 Health check |
| InfluxDB | 8086 | 仅 v3 |
| Prometheus scrape | 9464 | OTel Prometheus exporter（仅 v1） |
| ClickHouse HTTP | 8123 | 仅 v2 |
| ClickHouse TCP | 9000 | 仅 v2 |

> **注意**：两套 OTel Collector 默认共用 **4317/4318**；Grafana 单独在 `monitor.grafana.panel`，与后端通过 Docker 网络 `demo-network` 互联。

## 默认凭证

### Grafana

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

OTel 追踪配置（指向当前运行的 Collector，v1 或 v2 二选一）：

```yaml
trace:
  endpoint: localhost:4317
```

## 监控方案对比

| 维度 | monitor.v1.grafana | monitor.v2.clickhouse | monitor.v3.influxdb |
|------|---------------------|----------------------|---------------------|
| 数据存储 | Prometheus + Tempo + Loki | ClickHouse 统一存储 | InfluxDB 2 统一 bucket |
| 组件数量 | 4 个后端服务 | 2 个后端服务 + 面板 | InfluxDB + Collector + 面板 |
| Grafana | `monitor.grafana.panel`（OSS 13，含 ClickHouse 插件） | 同上 | 同上 |
| 数据源插件 | 内置 Prometheus / Tempo / Loki | 内置 + grafana-clickhouse-datasource | 内置 InfluxDB（Flux） |
| 适用场景 | 经典 LGTM 栈 | OTel ClickHouse exporter | Influx 生态与 Flux 查询 |

## 健康检查

| 服务 | 端点 |
|------|------|
| Go App | <http://localhost:8080/api/health> |
| OTel Collector | <http://localhost:13133/health> |
| Prometheus | <http://localhost:9090/-/healthy> |
| Grafana | <http://localhost:3000/api/health> |
| Tempo | <http://localhost:3200/ready> |
