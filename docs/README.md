# 项目文档

基于 Go Fiber 的 Web API 项目，提供开箱即用的 Web 基础框架。

## 目录结构

| 文档 | 描述 |
|------|------|
| [Architecture](Architecture.md) | 系统架构设计 |
| [CodeStructure](CodeStructure.md) | 代码结构详解 |
| [API](API.md) | API 接口文档 |
| [Configuration](Configuration.md) | 配置指南 |
| [Middleware](Middleware.md) | 中间件说明 |
| [Deployment](Deployment.md) | 部署指南 |

---

## 项目概览

| 项目信息 | 说明 |
|----------|------|
| **模块名** | `go_sample_code` |
| **Go 版本** | 1.25.0 |
| **Web 框架** | Fiber v2 |
| **ORM** | Ent |
| **依赖注入** | Uber FX |
| **数据库** | PostgreSQL / MySQL |
| **缓存** | Redis |
| **可观测性** | OpenTelemetry (Trace, Metrics, Log) |

## 技术栈

- **Web 框架**: [Fiber](https://github.com/gofiber/fiber) - 高性能 Web 框架
- **ORM**: [Ent](https://entgo.io/) - Go ORM 框架
- **依赖注入**: [uber-go/fx](https://github.com/uber-go/fx)
- **日志**: [zap](https://github.com/uber-go/zap) + OpenTelemetry
- **验证**: [go-playground/validator](https://github.com/go-playground/validator)
- **缓存**: [go-redis](https://github.com/redis/go-redis/v9)
- **JWT**: [lestrrat-go/jwx](https://github.com/lestrrat-go/jwx)
- **权限**: [Casbin](https://github.com/casbin/casbin/v2)
