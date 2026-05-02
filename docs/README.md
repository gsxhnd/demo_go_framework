# go_sample_code 文档中心

基于 Go Fiber 的 Web API 基础框架，提供开箱即用的分层架构、中间件体系和可观测性支持。

## 快速导航

| 如果你是... | 推荐阅读 |
|-------------|----------|
| 普通用户 | → [使用指南 (usage)](./usage/) |
| 开发者 | → [开发文档 (dev)](./dev/) + [代码描述 (wiki)](./wiki/) |

## 文档分类

### [dev](./dev/) — 开发文档

面向开发者的产品定义、架构设计和技术文档。

| 编号 | 文件 | 说明 |
|:----:|------|------|
| 00 | [00-readme.md](./dev/00-readme.md) | 开发文档导读与阅读顺序 |
| 01 | [01-product-scope.md](./dev/01-product-scope.md) | 产品定位、需求、MVP 范围与验收标准 |
| 02 | [02-architecture.md](./dev/02-architecture.md) | 系统分层、模块职责、依赖关系 |
| 03 | [03-domain-model.md](./dev/03-domain-model.md) | 核心实体、接口定义、错误码体系 |
| 04 | [04-tech-stack.md](./dev/04-tech-stack.md) | 语言版本、依赖管理、工具链、常用命令 |
| 05 | [05-roadmap.md](./dev/05-roadmap.md) | 开发路线图、阶段划分、里程碑 |
| 06 | [06-open-questions.md](./dev/06-open-questions.md) | 未决设计问题与决策记录 |

### [wiki](./wiki/) — 代码描述

面向开发者的模块级代码说明与设计决策文档。

| 编号 | 文件 | 说明 |
|:----:|------|------|
| 00 | [00-readme.md](./wiki/00-readme.md) | 模块总览与命名规则 |
| 01 | [01-module-handler.md](./wiki/01-module-handler.md) | HTTP 请求处理层 |
| 02 | [02-module-service.md](./wiki/02-module-service.md) | 业务逻辑层 |
| 03 | [03-module-repo.md](./wiki/03-module-repo.md) | 数据访问层 |
| 04 | [04-module-middleware.md](./wiki/04-module-middleware.md) | 中间件链 |
| 05 | [05-module-database.md](./wiki/05-module-database.md) | 数据库与 Redis 客户端 |
| 06 | [06-module-ent.md](./wiki/06-module-ent.md) | Ent ORM 实体与代码生成 |
| 07 | [07-module-errno.md](./wiki/07-module-errno.md) | 统一错误码体系 |
| 08 | [08-module-pkg.md](./wiki/08-module-pkg.md) | 公共工具包 |

### [usage](./usage/) — 使用指南

面向最终用户的操作手册和参考文档。

| 编号 | 文件 | 说明 |
|:----:|------|------|
| 00 | [00-readme.md](./usage/00-readme.md) | 使用指南导读 |
| 01 | [01-getting-started.md](./usage/01-getting-started.md) | 1 分钟快速上手 |
| 02 | [02-installation.md](./usage/02-installation.md) | 安装方式与系统要求 |
| 03 | [03-configuration.md](./usage/03-configuration.md) | 配置文件与配置项说明 |
| 04 | [04-basic-usage.md](./usage/04-basic-usage.md) | 核心 API 基础使用 |
| 05 | [05-advanced-usage.md](./usage/05-advanced-usage.md) | 高级功能与进阶技巧 |
| 06 | [06-cli-reference.md](./usage/06-cli-reference.md) | CLI 命令行参考 |
| 07 | [07-api-reference.md](./usage/07-api-reference.md) | API 接口完整参考 |
| 08 | [08-troubleshooting.md](./usage/08-troubleshooting.md) | 常见错误与解决方案 |
| 09 | [09-faq.md](./usage/09-faq.md) | 常见问题解答 |

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
