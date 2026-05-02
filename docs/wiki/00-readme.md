# go_sample_code 代码描述

## 概述

本目录包含 go_sample_code 各模块的代码描述和设计决策文档。每个文档解释对应模块的职责、设计理由、关键类型和注意事项。

## 模块总览

| 编号 | 模块 | 职责 | 文档 |
|:----:|------|------|------|
| 01 | handler | HTTP 请求处理，参数解析与响应构造 | [01-module-handler.md](./01-module-handler.md) |
| 02 | service | 业务逻辑处理，规则校验与错误转换 | [02-module-service.md](./02-module-service.md) |
| 03 | repo | 数据访问层，Ent ORM 操作封装 | [03-module-repo.md](./03-module-repo.md) |
| 04 | middleware | Fiber 中间件链 | [04-module-middleware.md](./04-module-middleware.md) |
| 05 | database | 数据库与 Redis 客户端管理 | [05-module-database.md](./05-module-database.md) |
| 06 | ent | Ent ORM 实体与代码生成 | [06-module-ent.md](./06-module-ent.md) |
| 07 | errno | 统一错误码体系 | [07-module-errno.md](./07-module-errno.md) |
| 08 | pkg | 公共工具包（logger、trace、metrics、validator、jwx、rbac） | [08-module-pkg.md](./08-module-pkg.md) |

## 命名规则

- 文件名格式：`{{编号}}-module-{{模块名}}.md`
- 编号按模块依赖层级排序：基础层（errno、pkg）→ 数据层（database、ent）→ 业务层（repo、service）→ 接口层（handler、middleware）
- 模块名使用小写英文 + 连字符

## 文档维护

- 新增模块时，在本目录下创建对应的文档文件
- 模块重构后，更新对应文档的内容
- 模块删除后，将对应文档标记为已废弃或删除
