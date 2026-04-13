# =============================================================================
# DevOps 文档
# =============================================================================
# Demo Go Framework 开发环境配置指南
# =============================================================================

## 目录结构

```
devops/
├── docker-compose.database.yml  # 数据库服务（MySQL, PostgreSQL, Redis）
├── docker-compose.monitoring.yml # 监控栈（Prometheus, Grafana, Tempo, Loki, OTel）
├── docker-compose.all.yml       # 一键启动 DevOps 服务（数据库 + 监控栈）
├── .env.example                 # 环境变量配置示例
│
├── config/
│   └── config.template.yaml     # 应用配置模板
│
├── otel-collector/
│   └── config.yaml              # OpenTelemetry Collector 配置
│
├── prometheus/
│   ├── config.yaml              # Prometheus 配置
│   └── rules/
│       └── alerts.yml           # Prometheus 告警规则
│
├── grafana/
│   ├── config.ini               # Grafana 全局配置
│   ├── provisioning/
│   │   ├── datasources/
│   │   │   └── datasources.yaml # Grafana 数据源配置
│   │   └── dashboards/
│   │       └── dashboards.yaml   # Grafana Dashboard 配置
│   └── dashboards/
│       └── http-metrics.json     # HTTP Metrics Dashboard
│
├── tempo/
│   └── config.yaml              # Tempo 配置
│
├── loki/
│   └── config.yaml              # Loki 配置
│
└── databases/
    ├── mysql/
    │   └── init.sql              # MySQL 初始化脚本
    └── postgres/
        └── init.sql              # PostgreSQL 初始化脚本
```

## 快速开始

### 1. 环境准备

```bash
# 复制环境变量文件
cd devops
cp .env.example .env

# 编辑 .env 文件，根据需要修改配置
vim .env
```

### 2. 一键启动所有服务

```bash
# 启动所有 DevOps 服务（数据库 + 监控栈）
docker compose -f docker-compose.all.yml up -d

# 查看服务状态
docker compose -f docker-compose.all.yml ps

# 查看日志
docker compose -f docker-compose.all.yml logs -f
```

### 3. 分步启动

#### 启动数据库
```bash
docker compose -f docker-compose.database.yml up -d
```

#### 启动监控栈
```bash
docker compose -f docker-compose.monitoring.yml up -d
```

### 4. 停止服务

```bash
# 停止所有服务
docker compose -f docker-compose.all.yml down

# 停止并删除数据卷
docker compose -f docker-compose.all.yml down -v
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Go App | 8080 | 应用服务 |
| MySQL | 3306 | MySQL 数据库 |
| PostgreSQL | 5432 | PostgreSQL 数据库 |
| Redis | 6379 | Redis 缓存 |
| Grafana | 3000 | 可视化看板 |
| Prometheus | 9090 | Metrics 存储 |
| Tempo | 3200 | 分布式追踪 |
| Loki | 3100 | 日志存储 |
| OTel Collector | 4317/4318 | OTLP 接收 |

## 默认凭证

### Grafana
- URL: http://localhost:3000
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
| GRAFANA_PORT | 3000 | Grafana 端口 |
| PROMETHEUS_PORT | 9090 | Prometheus 端口 |
| LOKI_PORT | 3100 | Loki 端口 |
| TEMPO_API_PORT | 3200 | Tempo API 端口 |

### 应用配置

复制 `config/config.template.yaml` 为 `config/config.yaml` 并修改相应配置：

```bash
cp config/config.template.yaml config/config.yaml
```

## 数据源说明

### Prometheus
- 存储指标数据
- 被 Grafana 用于查询 metrics
- 抓取 OTel Collector 暴露的 metrics

### Tempo
- 存储分布式追踪数据
- 支持 OTLP 协议接收 traces
- 与 Grafana 集成用于追踪可视化

### Loki
- 存储日志数据
- 支持从日志中提取 trace_id 跳转到 Tempo
- Grafana 原生支持

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
| Grafana | http://localhost:3000/api/health |
| Tempo | http://localhost:3200/ready |
| Loki | http://localhost:3100/ready |

## 清理

```bash
# 停止所有服务并删除容器
docker-compose -f docker-compose.all.yml down

# 停止所有服务，删除容器和数据卷
docker-compose -f docker-compose.all.yml down -v

# 删除所有未使用的镜像
docker image prune -f

# 完全清理（包括未使用的卷和网络）
docker system prune -af --volumes
```

## 常见问题

### Q: 服务启动失败怎么办？
```bash
# 查看详细日志
docker compose -f docker-compose.all.yml logs [service-name]

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

### Q: 如何重置 Grafana？
```bash
# 停止服务
docker compose -f docker-compose.monitoring.yml down

# 删除数据卷
docker volume rm demo-grafana-data

# 重新启动
docker compose -f docker-compose.monitoring.yml up -d
```
