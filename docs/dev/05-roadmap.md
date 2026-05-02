# 开发路线图

## 当前状态

项目处于 **MVP 已交付** 阶段。核心分层架构、中间件体系、用户管理 CRUD 示例和可观测性栈已就绪。

## 阶段划分

### Phase 1 — 基础框架（已完成 ✅）

- [x] 分层架构（Handler / Service / Repo）
- [x] uber-go/fx 依赖注入
- [x] Fiber HTTP 框架集成
- [x] Ent ORM + PostgreSQL / MySQL 双驱动
- [x] Redis 集成
- [x] YAML 配置文件加载
- [x] 统一错误码体系（errno）
- [x] 参数校验框架（go-playground/validator）

### Phase 2 — 用户管理示例（已完成 ✅）

- [x] 用户 CRUD API（创建、查询、更新、删除、列表）
- [x] 多条件查询（按 ID / 用户名 / 邮箱）
- [x] 分页支持
- [x] 用户名/邮箱唯一性校验

### Phase 3 — 中间件体系（已完成 ✅）

- [x] Recovery（Panic 恢复）
- [x] RateLimit（令牌桶限流）
- [x] Trace（OpenTelemetry 链路追踪）
- [x] Metrics（HTTP 指标采集）
- [x] Logger（请求日志）
- [x] Auth（JWT 认证中间件，已实现）
- [x] RBAC（Casbin 权限中间件，已实现）

### Phase 4 — 路由注册与集成（已完成 ✅）

- [x] 在 `RegisterHooks` 中注册用户管理路由
- [x] 启用 Auth 中间件并配置路由级认证策略
- [x] 启用 RBAC 中间件并配置权限策略
- [x] 配置 Casbin 模型和策略文件

### Phase 5 — 生产就绪（规划中）

- [ ] CI/CD 流水线（GitHub Actions）
- [ ] Docker 镜像构建与发布
- [ ] API 文档自动生成（Swagger/OpenAPI）
- [ ] 数据库迁移管理（golang-migrate）
- [ ] 结构化配置验证
- [ ] 环境变量配置支持
- [ ] 代码 lint 配置（.golangci.yml）
- [ ] 集成测试

### Phase 6 — 高级特性（远期规划）

- [ ] WebSocket 支持
- [ ] 多实例部署与分布式支持
- [ ] 分布式限流（基于 Redis）
- [ ] 消息队列集成

## 里程碑

| 里程碑 | 目标日期 | 状态 |
|--------|----------|------|
| M1: 基础框架可用 | - | ✅ 已完成 |
| M2: 用户管理 API 可用 | - | ✅ 已完成 |
| M3: 中间件体系完整 | - | ✅ 已完成 |
| M4: 全路由注册 + Auth/RBAC 启用 | - | ✅ 已完成 |
| M5: 生产就绪 | TBD | 🔲 待规划 |
