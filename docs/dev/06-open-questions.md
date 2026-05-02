# 待决问题

## 架构设计

### 路由注册策略

**问题**：当前仅 `/api/health` 路由已注册。用户管理路由（`POST /api/users` 等）虽已实现但未在 `RegisterHooks` 中注册。

**影响**：用户管理 API 无法通过 HTTP 访问。

**建议方案**：在 `RegisterHooks` 中添加用户路由组注册。

### Auth/RBAC 中间件启用时机

**问题**：Auth 和 RBAC 中间件已实现但未在中间件链中启用。

**影响**：当前所有接口无需认证即可访问。

**待决策**：

- Auth 中间件应放在哪个位置？RateLimit 之后还是 Logger 之后？
- 哪些路由需要认证？哪些可以公开访问？
- RBAC 策略如何配置（model.conf / policy.csv）？

### 配置文件扩展

**问题**：当前不支持环境变量覆盖配置，仅支持 YAML 文件。

**建议**：集成 Viper 的 `AutomaticEnv()` 功能，支持环境变量覆盖。

### 数据库迁移策略

**问题**：当前使用 Ent 的 `AutoMigrate` 进行数据库迁移，缺乏版本化迁移管理。

**建议**：集成 `golang-migrate` 或 Ent 的 `atlas` 迁移工具，实现版本化迁移。

## 技术选型

### 限流策略

**问题**：当前使用内存令牌桶（`golang.org/x/time/rate`）按 IP 限流，单实例有效。

**待决策**：是否需要支持基于 Redis 的分布式限流？

### API 文档生成

**问题**：当前无自动生成的 API 文档。

**建议**：集成 `swaggo/swag` 或 `go-swagger` 自动生成 OpenAPI 文档。

### 代码质量

**问题**：无 `.golangci.yml` 配置文件，无 CI lint 检查。

**建议**：添加 golangci-lint 配置并在 CI 中集成。

## 业务逻辑

### 用户认证流程

**问题**：Auth 中间件已实现 JWT 验证，但缺少登录接口（签发 Token）和注册流程。

**待实现**：

- POST /api/auth/login（用户名/密码 → JWT Token）
- POST /api/auth/register（注册 + 自动签发 Token）
- POST /api/auth/refresh（刷新 Token）

### 密码策略

**问题**：当前密码仅标记为 `Sensitive`（序列化时隐藏），无密码强度校验。

**待决策**：是否需要密码复杂度要求（最小长度、特殊字符等）？

## 运维相关

### 日志轮转

**问题**：已引入 `timberjack` 但未确认是否在默认配置中启用。

### 健康检查深度

**问题**：当前 `/api/health` 仅返回静态 healthy 状态。

**建议**：增加 DB 和 Redis 连通性检查，返回详细健康状态。

### 优雅关闭超时

**问题**：当前优雅关闭各组件使用独立的 context 超时，缺乏统一的关闭超时管理。
