# 故障排查

## 启动问题

### 服务启动失败：数据库连接错误

**现象**：启动后立即退出，日志显示数据库连接错误。

**解决方案**：

1. 检查数据库是否已启动：

```bash
docker compose -f devops/database/docker-compose.yml ps
```

1. 检查配置文件中的数据库连接信息是否正确（host、port、user、password、dbname）。
2. 如未使用配置文件，默认连接 `localhost:5432`（PostgreSQL）。

### 服务启动失败：Redis 连接错误

**现象**：日志显示 Redis 连接错误。

**解决方案**：

1. 确保 Redis 已启动：

```bash
docker compose -f devops/database/docker-compose.yml ps
```

1. 检查配置文件中的 Redis 地址是否正确。

### 端口已被占用

**现象**：`bind: address already in use`

**解决方案**：

1. 更改监听端口：

```yaml
common:
  listen: ":8081"
```

1. 或杀掉占用进程：

```bash
lsof -i :8080
kill -9 <PID>
```

## API 问题

### 请求返回 404

**现象**：调用用户管理 API 返回 404。

**原因**：用户管理路由当前未在 `RegisterHooks` 中注册，仅 `/api/health` 可用。

**解决方案**：在 `cmd/server/main.go` 的 `RegisterHooks` 中添加用户路由注册。

### 请求返回 429（限流）

**现象**：频繁请求后返回 429 Too Many Requests。

**原因**：触发了限流机制（默认 20 req/s）。

**解决方案**：

- 降低请求频率
- 或修改限流配置（`internal/middleware/rate_limit.go`）

### 请求返回 400（参数校验失败）

**现象**：创建用户时返回 `code: 1003`。

**解决方案**：

- 检查必填字段是否已传（`username`、`email`、`password`）
- 检查字段格式是否符合要求（邮箱格式、字符长度等）
- 检查用户名和邮箱是否已被占用

### 请求返回 500

**现象**：服务器返回 500 Internal Server Error。

**解决方案**：

1. 检查数据库是否正常连接
2. 查看服务日志定位具体错误
3. 检查 Ent schema 是否已迁移（自动执行）

## 可观测性问题

### Grafana 无法访问

**现象**：<http://localhost:3000> 无法访问。

**解决方案**：

```bash
# 检查可观测性栈是否已启动
docker compose -f devops/monitor.grafana.panel/docker-compose.yml ps

# 如未启动，启动 Grafana 面板及（按需）v1 后端
docker compose -f devops/monitor.grafana.panel/docker-compose.yml up -d
docker compose -f devops/monitor.v1.grafana/docker-compose.yml up -d
```

### 无 Trace 数据

**现象**：Grafana Tempo 中无链路追踪数据。

**解决方案**：

1. 确保 trace 已启用（`trace.enabled: true`）
2. 确保 OTel Collector 已启动
3. 检查 `trace.endpoint` 配置是否正确

### 无 Metrics 数据

**现象**：Prometheus 中无指标数据。

**解决方案**：

1. 确保 metrics 已启用（`metrics.enabled: true`）
2. 确保 OTel Collector 已启动
3. 检查 `metrics.endpoint` 配置是否正确

## 编译问题

### Ent 代码生成失败

**现象**：运行 `ent generate` 报错。

**解决方案**：

1. 确保已安装 ent CLI：

```bash
go install entgo.io/ent/cmd/ent@latest
```

1. 确保在项目根目录下运行：

```bash
cd <项目根目录>
ent generate ./internal/ent/schema
```

### Go 版本不兼容

**现象**：编译时提示 Go 版本不匹配。

**解决方案**：项目要求 Go 1.25.0+，升级 Go 版本：

```bash
go version
# 如版本过低，安装最新版本
```
