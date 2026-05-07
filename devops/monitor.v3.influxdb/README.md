# monitor.v3.influxdb

InfluxDB 2.x + OpenTelemetry Collector（contrib），将应用的 **metrics / traces / logs** 经 OTLP 写入同一 bucket `otel`（measurement：`spans`、`logs`、各指标名）。

Grafana 使用 `../monitor.grafana.panel`，预置 **InfluxDB (OTel v3)** 数据源（Flux，指向本栈的 `influxdb:8086`）。

### 为何需要 `influxdb/config.yml`（尤其 2.9）

官方镜像的 `entrypoint.sh` 会用 `dasel` 从配置里读取并改写 **`http-bind-address`、`tls-cert`、`tls-key`** 以做首次 `DOCKER_INFLUXDB_INIT_MODE=setup`。从 **InfluxDB 2.9** 起，镜像内置的默认 `config.yml` 不再包含这些键，`dasel` 报错后 init 失败，日志里会出现 **“cleaning bolt and engine files to prevent conflicts on retry”**，容器随即退出。

本目录的 `influxdb/config.yml` 补齐上述字段；Compose 将其挂到 **`/etc/influxdb2/config.yml`**。挂载目标**不要**使用 `:ro`，否则 entrypoint 对配置的 `chown` 会失败。

从 **2.7 数据卷升级到 2.9** 前请先备份；若 init 已反复失败，bolt 可能被清理，需 `docker compose down -v` 后重新初始化（会丢本地演示数据）。

## 启动

```bash
cd devops/monitor.v3.influxdb
docker compose up -d
```

## 应用 OTLP 端点

与 v1/v2 的 Collector **不要同时占用主机 OTLP 端口**。本栈映射：

| 协议 | 主机地址 |
|------|----------|
| gRPC | `localhost:4319` |
| HTTP | `localhost:4320` |

示例（`config.local.yaml`）：

```yaml
trace:
  endpoint: localhost:4319
```

若应用使用 OTLP HTTP，请使用 `http://localhost:4320`（按 SDK 要求配置）。

## InfluxDB UI / Token

- UI: <http://localhost:8086>
- 用户: `admin` / `demo_password_change_me`
- Org: `demo`，Bucket: `otel`
- Admin token（与 Compose / Collector / Grafana 预置一致）: `demo_influx_admin_token_for_local_dev_only_12345`

## Grafana Explore（Flux）

数据源选 **InfluxDB (OTel v3)**，示例：

```flux
from(bucket: "otel")
  |> range(start: -1h)
  |> filter(fn: (r) => r["_measurement"] == "spans" or r["_measurement"] == "logs")
  |> limit(n: 20)
```

指标 measurement 一般为 OTel 指标名（`telegraf-prometheus-v1` 模式）。

## 健康检查

- Collector: <http://localhost:13134/health>
- InfluxDB: <http://localhost:8086/ready>

## 停止

```bash
docker compose down
```

删除 InfluxDB 数据卷：`docker compose down -v`。
