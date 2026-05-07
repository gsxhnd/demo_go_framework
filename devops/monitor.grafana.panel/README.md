# Grafana 面板（v1 + v2 数据源）

独立启动的 Grafana OSS，预置 **Prometheus / Tempo / Loki**（`monitor.v1.grafana`）、**ClickHouse**（`monitor.v2.clickhouse`）与 **InfluxDB**（`monitor.v3.influxdb`）数据源。所有栈需使用同一 Docker 网络 **`demo-network`**，以便容器名 `prometheus`、`tempo`、`loki`、`clickhouse`、`influxdb` 可解析。

## 启动

```bash
cd devops/monitor.grafana.panel
docker compose up -d
```

访问 <http://localhost:3000>。

## 与后端栈的配合

1. 按需启动 `../monitor.v1.grafana`、`../monitor.v2.clickhouse`（ClickHouse + OTel）、`../monitor.v3.influxdb`（InfluxDB + OTel）。
2. 再启动本目录的 Compose。若某后端未运行，对应数据源在 Explore 中会报错，属预期行为。

> **端口**：v1/v2 的 Collector 默认 **4317/4318**（二者互斥）；v3 使用 **4319/4320**，可与 v1 或 v2 之一并行（应用只连一套 Collector）。

## 预置 Dashboard

- `grafana/dashboards/http-metrics.json` — Prometheus（v1）
- `grafana/dashboards/otel_dashboard.json` — ClickHouse OTel（v2）

## 停止

```bash
docker compose down
```

删除 Grafana 持久化数据卷：`docker compose down -v`。
