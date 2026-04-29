# 部署指南

## 环境要求

| 依赖 | 版本 | 说明 |
|------|------|------|
| Go | 1.25.0+ | 编译运行 |
| PostgreSQL | 14+ | 生产环境推荐 |
| MySQL | 8.0+ | 可选 |
| Redis | 6+ | 缓存/限流 |

## 开发环境

### 1. 启动数据库

```bash
docker compose -f devops/database/docker-compose.yml up -d
```

### 2. 配置

编辑 `config/config.local.yaml` 或使用默认配置。

### 3. 运行服务

```bash
go run cmd/server/main.go -c config/config.local.yaml
```

### 4. 验证

```bash
curl http://localhost:8080/api/health
```

---

## 生产部署

### 编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# macOS
GOOS=darwin GOARCH=amd64 go build -o server ./cmd/server
```

### Docker 部署

创建 `Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
COPY config/ ./config/

EXPOSE 8080
CMD ["./server", "-c", "config/config.local.yaml"]
```

构建并运行:

```bash
docker build -t go_sample_code .
docker run -p 8080:8080 go_sample_code
```

### Docker Compose 完整部署

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./config:/app/config

  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: demo
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

---

## 可观测性部署

### 快速启动 (开发用)

```bash
docker compose -f devops/grafana.v1/docker-compose.yml up -d
```

启动以下服务:

- Grafana (<http://localhost:3000>)
- Prometheus (<http://localhost:9090>)
- Tempo (<http://localhost:3100>)
- Loki (<http://localhost:3100>)
- OTel Collector

### ClickHouse 版本

```bash
docker compose -f devops/grafana.v2/docker-compose.yml up -d
```

---

## 数据库初始化

### PostgreSQL

```sql
CREATE DATABASE demo;
```

### Ent 代码生成

如修改实体定义，需重新生成代码：

```bash
# 安装 ent 工具
go install entgo.io/ent/cmd/ent@latest

# 重新生成
ent generate ./internal/ent/schema
```

---

## 环境变量 (可选)

如需支持环境变量，可在 `cmd/server/config.go` 中配置 Viper：

```go
viper.AutomaticEnv()
viper.SetEnvPrefix("APP")
```

---

## 目录结构建议

生产环境建议目录结构：

```
/opt/go_sample_code/
├── config/
│   └── config.yaml
├── data/
│   ├── postgres/
│   └── redis/
├── logs/
└── server
```

---

## 健康检查

### 应用健康

```bash
curl http://localhost:8080/api/health
```

### 数据库健康

在 `config.yaml` 中配置 `health_check_timeout` 后自动检查。

### Docker 健康检查

```yaml
healthcheck:
  test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/api/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

---

## 日志配置

推荐生产环境使用 JSON 格式日志：

```yaml
logger:
  level: "info"
  encoding: "json"
  output_paths:
    - "/app/logs/app.log"
    - "stdout"
```

---

## 安全建议

1. **不要提交敏感配置** - 使用环境变量或密钥管理服务
2. **生产环境关闭 Debug 日志** - `level: "info"`
3. **配置 TLS** - 使用 Nginx/Caddy 反向代理
4. **限流配置** - 根据实际需求调整
5. **启用 Auth/RBAC** - 在生产环境中启用认证和授权
