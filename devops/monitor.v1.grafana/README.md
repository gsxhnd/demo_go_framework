# monitor.v1.grafana（Prometheus + Tempo + Loki + OTel Collector）

经典 LGTM 后端，**不含 Grafana**。面板请单独启动 **`../monitor.grafana.panel`**。

```bash
docker compose up -d
```

默认与 `monitor.v2.clickhouse` 的 OTel Collector **共用 4317/4318**，请勿同时在本机启动两套 Collector，或自行修改端口映射。
