# =============================================================================
# DevOps 文档
# =============================================================================
# Demo Go Framework 开发环境配置指南
# =============================================================================

## 目录结构

```
devops/
├── database/
│   └── docker-compose.yml       # 数据库服务（MySQL, PostgreSQL, Redis）
│
├── grafana.v1/
│   ├── docker-compose.yml       # Grafana v1 监控栈（Prometheus, Grafana, Tempo, Loki, OTel）
│   ├── grafana/
│   │   ├── config.ini
│   │   ├── provisioning/
│   │   │   ├── datasources/datasources.yaml
│   │   │   └── dashboards/dashboards.yaml
│   │   └── dashboards/
│   │       └── http-metrics.json
│   ├── prometheus/
│   │   ├── config.yaml
│   │   └── rules/alerts.yml
│   ├── tempo/
│   │   └── config.yaml
│   ├── loki/
│   │   └── config.yaml
│   └── otel-collector/
│       └── config.yaml
│
├── grafana.v2/
│   ├── docker-compose.yml       # Grafana v2 监控栈（ClickHouse + Grafana OSS）
│   ├── clickhouse/
│   │   └── init.sql             # ClickHouse 初始化脚本
│   └── grafana/
│       ├── config.ini
│       └── provisioning/
│           └── datasources/datasources.yaml
│
├── .env.example                 # 环境变量配置示例
├── .env                         # 环境变量配置（从 .env.example 复制）
├── config/
│   └── config.template.yaml     # 应用配置模板
├── databases/
│   ├── mysql/
│   │   └── init.sql             # MySQL 初始化脚本
│   └── postgres/
│       └── init.sql             # PostgreSQL 初始化脚本
└── README.md
```

## 快速开始

### 1. 环境准备

每个子目录下都有独立的 `.env.example`，首次使用时复制为 `.env`：

```bash
cp devops/database/.env.example devops/database/.env
cp devops/grafana.v1/.env.example devops/grafana.v1/.env
cp devops/grafana.v2/clickhouse/.env.example devops/grafana.v2/clickhouse/.env
cp devops/grafana.v2/grafana/.env.example devops/grafana.v2/grafana/.env
```

### 2. 启动数据库

```bash
cd devops/database
docker compose up -d
```

### 3. 启动 Grafana v1 监控栈

Prometheus + Tempo + Loki + OTel Collector + Grafana 10.3.3

```bash
cd devops/grafana.v1
docker compose up -d
```

### 4. 启动 Grafana v2 监控栈

ClickHouse + Grafana OSS 13.0.1

```bash
# 先启动 ClickHouse
cd devops/grafana.v2/clickhouse
docker compose up -d

# 初始化 OTel 数据表（仅首次）
docker exec -i ${COMPOSE_PROJECT_NAME:-demo}-clickhouse \
  clickhouse-client --multiquery < init.sql

# 再启动 Grafana v2
cd devops/grafana.v2/grafana
docker compose up -d
```

### 5. 停止服务

```bash
# 停止数据库
cd devops/database && docker compose down

# 停止 Grafana v1
cd devops/grafana.v1 && docker compose down

# 停止 Grafana v2
cd devops/grafana.v2/grafana && docker compose down
cd devops/grafana.v2/clickhouse && docker compose down

# 停止并删除数据卷（加 -v 参数）
docker compose down -v
```

### 2. 启动数据库

```bash
cd devops/database
docker compose --env-file ../.env up -d

# 查看服务状态
docker compose --env-file ../.env ps
```

### 3. 启动 Grafana v1 监控栈

Prometheus + Tempo + Loki + OTel Collector + Grafana 10.3.3

```bash
cd devops/grafana.v1
docker compose --env-file ../.env up -d
```

### 4. 启动 Grafana v2 监控栈

ClickHouse + Grafana OSS 13.0.1

```bash
cd devops/grafana.v2
docker compose --env-file ../.env up -d
```

### 5. 停止服务

```bash
# 停止数据库
cd devops/database
docker compose --env-file ../.env down

# 停止 Grafana v1
cd devops/grafana.v1
docker compose --env-file ../.env down

# 停止 Grafana v2
cd devops/grafana.v2
docker compose --env-file ../.env down

# 停止并删除数据卷（加 -v 参数）
docker compose --env-file ../.env down -v
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Go App | 8080 | 应用服务 |
| MySQL | 3306 | MySQL 数据库 |
| PostgreSQL | 5432 | PostgreSQL 数据库 |
| Redis | 6379 | Redis 缓存 |
| Grafana v1 | 3000 | 可视化看板（Prometheus 数据源） |
| Grafana v2 | 3001 | 可视化看板（ClickHouse 数据源） |
| Prometheus | 9090 | Metrics 存储 |
| Tempo | 3200 | 分布式追踪 |
| Loki | 3100 | 日志存储 |
| OTel Collector | 4317/4318 | OTLP 接收 |
| ClickHouse | 8123/9000 | 列式数据库（HTTP/TCP） |

## 默认凭证

### Grafana (v1 & v2)
- URL: http://localhost:3000 (v1) / http://localhost:3001 (v2)
- 用户名: (匿名访问已启用)
- 角色: Admin

### 数据库
- MySQL: demo_user / demo_password
- PostgreSQL: demo_user / demo_password

## 配置说明

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| COMPOSE_PROJECT_NAME | demo | 项目名称 |
| MYSQL_PORT | 3306 | MySQL 端口 |
| POSTGRES_PORT | 5432 | PostgreSQL 端口 |
| REDIS_PORT | 6379 | Redis 端口 |
| GRAFANA_PORT | 3000 | Grafana v1 端口 |
| GRAFANA_V2_PORT | 3001 | Grafana v2 端口 |
| PROMETHEUS_PORT | 9090 | Prometheus 端口 |
| LOKI_PORT | 3100 | Loki 端口 |
| TEMPO_API_PORT | 3200 | Tempo API 端口 |
| CLICKHOUSE_HTTP_PORT | 8123 | ClickHouse HTTP 端口 |
| CLICKHOUSE_TCP_PORT | 9000 | ClickHouse TCP 端口 |

### 应用配置

复制 `config/config.template.yaml` 为 `config/config.yaml` 并修改相应配置：

```bash
cp config/config.template.yaml config/config.yaml
```

## 数据源说明

### Grafana v1 数据源

| 数据源 | 说明 |
|--------|------|
| Prometheus | 存储指标数据，被 Grafana 用于查询 metrics |
| Tempo | 存储分布式追踪数据，支持 OTLP 协议 |
| Loki | 存储日志数据，支持从日志中提取 trace_id 跳转到 Tempo |

### Grafana v2 数据源

| 数据源 | 说明 |
|--------|------|
| ClickHouse | 列式数据库，用于存储和查询大规模日志、指标数据 |

### OpenTelemetry Collector
- 统一接收应用的 traces、metrics、logs
- 处理后转发到对应的后端存储
- 提供健康检查端点

## 健康检查

各服务健康检查端点：

| 服务 | 端点 |
|------|------|
| Go App | http://localhost:8080/api/health |
| OTel Collector | http://localhost:13133/health |
| Prometheus | http://localhost:9090/-/healthy |
| Grafana v1 | http://localhost:3000/api/health |
| Grafana v2 | http://localhost:3001/api/health |
| Tempo | http://localhost:3200/ready |
| Loki | http://localhost:3100/ready |
| ClickHouse | http://localhost:8123/ping |

## 清理

```bash
# 停止数据库并删除数据卷
cd devops/database && docker compose --env-file ../.env down -v

# 停止 Grafana v1 并删除数据卷
cd devops/grafana.v1 && docker compose --env-file ../.env down -v

# 停止 Grafana v2 并删除数据卷
cd devops/grafana.v2 && docker compose --env-file ../.env down -v

# 删除所有未使用的镜像
docker image prune -f

# 完全清理（包括未使用的卷和网络）
docker system prune -af --volumes
```

## 常见问题

### Q: 服务启动失败怎么办？
```bash
# 查看详细日志（以数据库为例）
cd devops/database
docker compose --env-file ../.env logs [service-name]

# 检查端口占用
netstat -tlnp | grep [port]
```

### Q: 如何查看 OTel Collector 配置？
```bash
# 进入容器
docker exec -it demo-otel-collector sh

# 查看配置
cat /etc/otel-collector-config.yaml
```

### Q: 如何重新加载 Prometheus 配置？
```bash
curl -X POST http://localhost:9090/-/reload
```

### Q: 如何重置 Grafana v1？
```bash
cd devops/grafana.v1
docker compose --env-file ../.env down
docker volume rm demo-grafana-data
docker compose --env-file ../.env up -d
```

### Q: 如何重置 Grafana v2？
```bash
cd devops/grafana.v2
docker compose --env-file ../.env down
docker volume rm demo-grafana-v2-data
docker compose --env-file ../.env up -d
```

### Q: Grafana v2 的 ClickHouse 插件安装失败？
Grafana v2 通过 `GF_INSTALL_PLUGINS` 环境变量自动安装 `grafana-clickhouse-datasource` 插件。
如果网络不通，可以手动下载插件并挂载到容器的 `/var/lib/grafana/plugins` 目录。
