# go_sample_code 开发文档

## 概述

基于 Go Fiber 的 Web API 基础框架，提供开箱即用的分层架构、中间件体系和可观测性支持。

本目录是 go_sample_code 的开发文档，包含产品定义、架构设计、领域模型、技术栈说明、开发路线图和待决问题。

## 阅读顺序

| 文件 | 内容 |
|------|------|
| `01-product-scope.md` | 产品定位、需求、MVP 范围与验收标准 |
| `02-architecture.md` | 系统分层、模块职责、依赖关系 |
| `03-domain-model.md` | 核心实体、接口定义、错误码体系 |
| `04-tech-stack.md` | 语言版本、依赖管理、工具链、常用命令 |
| `05-roadmap.md` | 开发路线图、阶段划分、里程碑 |
| `06-open-questions.md` | 未决设计问题与决策记录 |

**建议阅读顺序**：`01` → `02` → `03` → `04` → `05` → `06`

## 文档规则

- `docs/dev/` 是当前唯一有效的开发文档源
- 产品范围以 `01-product-scope.md` 为准
- 架构边界以 `02-architecture.md` 与 `03-domain-model.md` 为准
- 开发顺序以 `05-roadmap.md` 为准
- 未决问题统一记录在 `06-open-questions.md`
- 文档默认为 draft，代码落地后再更新为已验证描述

## 设计原则

- 先定义边界，再定义接口，再定义实现细节
- Handler → Service → Repo 分层，接口驱动，依赖注入
- 模块间通过明确的接口通信，禁止循环依赖
- 所有错误通过统一的 `errno.Errno` 接口返回
