# 快速开始

> 1 分钟上手 go_sample_code

## 前置条件

- Go 1.25.0+
- Docker 和 Docker Compose（用于启动数据库和 Redis）

## 最小示例

```bash
# 1. 克隆项目
git clone <仓库地址>
cd go_sample_code

# 2. 启动数据库和 Redis
docker compose -f devops/database/docker-compose.yml up -d

# 3. 启动服务
go run cmd/server/main.go -c config/config.local.yaml
```

预期输出：

```
{"level":"info","msg":"server started on :8080"}
```

## 验证

```bash
curl http://localhost:8080/api/health
```

预期响应：

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "status": "healthy",
    "timestamp": "2026-05-02T10:00:00Z"
  }
}
```

## 下一步

- 了解更多安装方式 → [安装指南](./02-installation.md)
- 了解配置选项 → [配置说明](./03-configuration.md)
- 开始使用 API → [基础使用](./04-basic-usage.md)
