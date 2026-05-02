# 安装指南

## 系统要求

| 项目 | 要求 |
|------|------|
| 操作系统 | macOS / Linux / Windows |
| 运行时 | Go 1.25.0+ |
| 数据库 | PostgreSQL 14+ 或 MySQL 8.0+ |
| 缓存 | Redis 6+ |
| 磁盘空间 | 最小 100MB |

## 安装方式

### 方式一：从源码编译（推荐）

```bash
git clone <仓库地址>
cd go_sample_code

# 编译
go build -o server ./cmd/server

# 运行
./server -c config/config.local.yaml
```

### 方式二：直接运行

```bash
git clone <仓库地址>
cd go_sample_code

# 直接运行（无需编译）
go run cmd/server/main.go -c config/config.local.yaml
```

### 方式三：跨平台编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o server-linux ./cmd/server

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o server-darwin-arm64 ./cmd/server

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o server-darwin-amd64 ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o server.exe ./cmd/server
```

### 方式四：Docker

创建 `Dockerfile`：

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

构建并运行：

```bash
docker build -t go_sample_code .
docker run -p 8080:8080 go_sample_code
```

## 依赖服务

### 数据库和 Redis

```bash
# 使用 Docker Compose 启动
docker compose -f devops/database/docker-compose.yml up -d
```

启动的服务：

| 服务 | 端口 | 说明 |
|------|------|------|
| PostgreSQL | 5432 | 关系型数据库 |
| Redis | 6379 | 缓存 |

### 可观测性栈（可选）

```bash
# 基础版（Grafana + Prometheus + Tempo + Loki）
docker compose -f devops/grafana.v1/docker-compose.yml up -d

# ClickHouse 版
docker compose -f devops/grafana.v2/docker-compose.yml up -d
```

## 安装验证

```bash
# 检查服务是否正常
curl http://localhost:8080/api/health
```

预期输出：

```json
{"code":0,"message":"OK","data":{"status":"healthy","timestamp":"..."}}
```

## 下一步

安装完成后，请阅读 [配置说明](./03-configuration.md) 进行初始配置。
