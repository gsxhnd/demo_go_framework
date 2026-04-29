# Grafana v2 监控栈

ClickHouse + OTel Collector + Grafana OSS，用于存储和可视化 OpenTelemetry 数据。

## 目录结构

```
grafana.v2/
├── clickhouse/
│   ├── docker-compose.yml        # ClickHouse 服务
│   ├── init.sql                  # OTel 数据表 DDL（手动执行）
│   ├── .env.example
│   └── .env
├── otel-collector/
│   ├── docker-compose.yml        # OTel Collector 服务
│   ├── config.yaml               # Collector 配置（输出到 ClickHouse）
│   ├── .env.example
│   └── .env
└── grafana/
    ├── docker-compose.yml        # Grafana OSS 13.0.1
    ├── config.ini                # Grafana 全局配置
    ├── .env.example
    ├── .env
    └── provisioning/
        └── datasources/
            └── datasources.yaml  # ClickHouse 数据源
```

## 数据流

```
Go App
  │
  │  OTLP gRPC/HTTP
  ▼
OTel Collector (v2)
  │
  │  clickhouse exporter (TCP 9000)
  ▼
ClickHouse
  │
  │  grafana-clickhouse-datasource
  ▼
Grafana v2
```

## 快速开始

### 1. 启动 ClickHouse

```bash
cd clickhouse
docker compose up -d
```

### 2. 初始化 OTel 数据表

ClickHouse 启动后手动执行建表脚本：

```bash
docker exec -i demo-clickhouse \
  clickhouse-client --multiquery < init.sql
```

验证建表结果：

```bash
docker exec -it demo-clickhouse \
  clickhouse-client --query "SHOW TABLES FROM otel"
```

预期输出：

```
otel_logs
otel_metrics_exponential_histogram
otel_metrics_gauge
otel_metrics_histogram
otel_metrics_summary
otel_metrics_sum
otel_traces
```

### 3. 启动 OTel Collector

```bash
cd otel-collector
docker compose up -d
```

### 4. 启动 Grafana v2

```bash
cd grafana
docker compose up -d
```

访问 http://localhost:3001，ClickHouse 数据源已自动配置。

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| ClickHouse | 8123 | HTTP 接口（查询、健康检查） |
| ClickHouse | 9000 | Native TCP 接口（clickhouse-client） |
| OTel Collector | 4319 | OTLP gRPC 接收 |
| OTel Collector | 4320 | OTLP HTTP 接收 |
| OTel Collector | 8889 | Collector 自身 metrics |
| OTel Collector | 13134 | Health check |
| Grafana v2 | 3001 | Web UI |

> 端口与 grafana.v1 的 OTel Collector（4317/4318/8888/13133）错开，两套可同时运行。

## 数据表说明

所有表位于 `otel` 数据库，schema 与 [OTel ClickHouse Exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/clickhouseexporter) 保持一致。

| 表名 | 数据类型 | TTL |
|------|----------|-----|
| `otel_metrics_gauge` | Gauge 指标 | 90 天 |
| `otel_metrics_sum` | Sum/Counter 指标 | 90 天 |
| `otel_metrics_histogram` | Histogram 指标 | 90 天 |
| `otel_metrics_exponential_histogram` | 指数 Histogram 指标 | 90 天 |
| `otel_metrics_summary` | Summary 指标 | 90 天 |
| `otel_logs` | 日志 | 30 天 |
| `otel_traces` | 追踪 Span | 30 天 |

## 应用接入

将应用的 OTLP exporter 指向 OTel Collector v2 的端口：

```yaml
# config.yaml 示例
otel:
  endpoint: "http://localhost:4319"   # gRPC
  # endpoint: "http://localhost:4320" # HTTP
```

## 停止服务

```bash
# 停止 Grafana v2
cd grafana && docker compose down

# 停止 OTel Collector
cd otel-collector && docker compose down

# 停止 ClickHouse
cd clickhouse && docker compose down

# 停止并删除数据卷（加 -v 参数）
docker compose down -v
```
