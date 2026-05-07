# monitor.v2.clickhouse（ClickHouse + OTel Collector）

将 OTLP 数据写入 ClickHouse。Grafana UI 请使用 **`../monitor.grafana.panel`**（预置 ClickHouse 数据源与 v2 面板）。

## 目录结构

```
monitor.v2.clickhouse/
├── clickhouse/
│   ├── docker-compose.yml
│   └── init.sql
└── otel-collector/
    ├── docker-compose.yml
    └── config.yaml
```

## 数据流

```
Go App → OTel Collector → ClickHouse → Grafana（monitor.grafana.panel）
```

## 快速开始

### 1. 启动 ClickHouse

```bash
cd clickhouse
docker compose up -d
```

### 2. 初始化 OTel 表

```bash
docker exec -i demo-clickhouse clickhouse-client --multiquery < init.sql
```

### 3. 启动 OTel Collector

```bash
cd ../otel-collector
docker compose up -d
```

### 4. 启动 Grafana 面板

```bash
cd ../../monitor.grafana.panel
docker compose up -d
```

## 端口

| 服务           | 端口                    |
|----------------|-------------------------|
| ClickHouse     | 8123, 9000              |
| OTel Collector | 4317, 4318, 8888, 13133 |

## 停止

```bash
cd otel-collector && docker compose down
cd ../clickhouse && docker compose down
```
