# 产品范围

## 项目简介

go_sample_code 是一个基于 Go Fiber 的 Web API 基础框架，提供开箱即用的分层架构、中间件体系和可观测性支持，可作为新项目的脚手架快速启动。

## 项目定位

面向 Go 后端开发者，提供一个结构清晰、功能完备的 Web API 项目模板。与同类框架模板的差异：

- **分层架构**：Handler → Service → Repo 三层分离，接口驱动
- **依赖注入**：基于 uber-go/fx，构造函数注入，易于测试和替换
- **可观测性内置**：OpenTelemetry (Trace + Metrics + Log) 开箱即用
- **Ent ORM**：类型安全的 ORM，代码生成避免手写 SQL
- **最小依赖原则**：只引入必要的外部依赖，保持轻量

## 目标

- 提供标准化的 Go Web API 项目结构
- 内置用户管理 CRUD 作为业务示例
- 集成完整的中间件体系（Recovery、RateLimit、Trace、Metrics、Logger、Auth、RBAC）
- 支持 PostgreSQL / MySQL 双数据库驱动
- 提供生产可用的可观测性栈（Grafana + Prometheus + Tempo + Loki）

## 目标用户

- **Go 后端开发者**：需要一个成熟的项目模板快速启动新项目
- **团队技术负责人**：需要统一团队的项目结构和编码规范

## 功能需求

### 用户管理

- 用户创建（POST /api/users）
- 按 ID 查询用户（GET /api/users/:id）
- 按用户名查询用户（GET /api/users/username/:username）
- 按邮箱查询用户（GET /api/users/email/:email）
- 更新用户信息（PUT /api/users/:id）
- 删除用户（DELETE /api/users/:id）
- 分页列表（GET /api/users）

### 健康检查

- 服务健康检查（GET /api/health）

### 中间件体系

- Panic 恢复（Recovery）
- 请求限流（RateLimit，令牌桶算法）
- 链路追踪（OpenTelemetry Trace）
- 指标采集（HTTP Metrics）
- 请求日志（Logger）
- 认证中间件（Auth，已实现未启用）
- 权限控制（RBAC/ABAC，已实现未启用）

## 非功能性需求

- **性能**：Fiber 高性能 HTTP 框架，单机支持万级 QPS
- **安全**：JWT 认证支持、密码哈希存储（Ent Sensitive 字段）、参数校验
- **兼容性**：PostgreSQL 14+ / MySQL 8.0+，Go 1.25+
- **可观测性**：OpenTelemetry 标准，支持 OTLP 导出

## 入口模式

- **HTTP API**：通过 Fiber 提供 RESTful API 服务
- **CLI**：通过 `-c` 参数指定配置文件路径

## MVP 范围

### MVP 包含

- 分层架构（Handler / Service / Repo）
- 用户管理 CRUD API
- 健康检查接口
- 核心中间件：Recovery、RateLimit、Trace、Metrics、Logger
- 统一错误码体系
- 参数校验框架
- YAML 配置文件支持
- PostgreSQL / MySQL 双驱动支持
- Redis 集成
- Docker Compose 开发环境

### MVP 不包含

- 认证/授权中间件的启用和路由集成
- 文件上传/下载
- WebSocket 支持
- gRPC 服务
- 数据库迁移工具集成
- CI/CD 流水线

### MVP 约束

- 仅支持单实例部署，不支持分布式
- 仅 `/api/health` 路由已注册，用户相关路由未在 RegisterHooks 中注册

## 验收标准

- 服务启动后 `/api/health` 返回 200
- 配置文件中指定的数据库和 Redis 连接正常
- 中间件按顺序正确执行
- 错误响应格式统一为 `{code, message, data}`

## 延期项

- Auth/RBAC 中间件的路由级集成和策略配置
- 多实例部署与分布式支持
- 自动生成 API 文档（Swagger/OpenAPI）
- 数据库迁移管理（如 golang-migrate）
- CI/CD 流水线配置
