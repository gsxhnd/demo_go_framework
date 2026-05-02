# Repo 模块

> 数据访问层，封装 Ent ORM 操作，提供类型安全的数据库访问接口。

## 设计决策

### 为什么需要这个模块

Repo 层将数据访问逻辑与业务逻辑分离，使得：

- 数据库操作可被 mock，便于 Service 层单元测试
- 数据库变更（如切换 ORM 或数据库类型）仅影响 Repo 层
- 复杂查询逻辑集中管理，避免分散在各 Service 方法中

### 为什么这么设计

- **选择了**：每个实体一个 Repo 接口，方法按 CRUD 操作组织
- **而不是**：使用通用的 Repository 模式或直接在 Service 中使用 Ent
- **原因**：接口驱动便于 mock 测试，方法粒度与业务操作对齐

## 关键类型与接口

### userrepo.UserRepo

- **定义位置**：`internal/repo/user/repo.go`
- **用途**：用户数据访问接口

**方法**：

| 方法 | 说明 |
|------|------|
| `Create(ctx, params) (*ent.User, error)` | 创建用户记录 |
| `GetByID(ctx, id) (*ent.User, error)` | 按 ID 查询 |
| `GetByUsername(ctx, username) (*ent.User, error)` | 按用户名查询 |
| `GetByEmail(ctx, email) (*ent.User, error)` | 按邮箱查询 |
| `Update(ctx, id, params) (*ent.User, error)` | 更新用户 |
| `Delete(ctx, id) error` | 删除用户 |
| `List(ctx, page, pageSize) ([]*ent.User, int, error)` | 分页列表（返回数据 + 总数） |
| `ExistsByEmail(ctx, email) (bool, error)` | 邮箱是否存在 |
| `ExistsByUsername(ctx, username) (bool, error)` | 用户名是否存在 |

## 模块结构

```text
internal/repo/user/
├── repo.go                     # UserRepo 接口与实现
├── user_create.go              # 创建用户
├── user_delete.go              # 删除用户
├── user_update.go              # 更新用户
├── user_get_by_id.go           # 按 ID 查询
├── user_get_by_username.go     # 按用户名查询
├── user_get_by_email.go        # 按邮箱查询
├── user_list.go                # 分页列表
├── user_exists_by_email.go     # 邮箱存在性检查
└── user_exists_by_username.go  # 用户名存在性检查
```

| 文件 | 职责 |
|------|------|
| `repo.go` | 定义接口、实现结构体（持有 `*ent.Client`）、构造函数 |
| `user_*.go` | 各数据操作方法，封装 Ent 查询构建器 |

## 与其他模块的关系

### 依赖

- **ent**：Ent ORM Client，执行数据库操作

### 被依赖

- **service**：Service 层调用 Repo 层

### 依赖关系图

```text
service
  ↑ (调用)
repo
  ↑ (调用)
ent (Ent ORM Client)
```

## 注意事项

- Repo 层不应包含业务逻辑，只做数据存取
- 返回类型使用 Ent 生成的实体类型（`*ent.User`），由 Service 层转换为响应类型
- 唯一性检查方法（`ExistsByEmail`、`ExistsByUsername`）用于 Service 层的业务校验
- 分页方法返回 `([]*ent.User, int, error)`，其中 `int` 为总记录数
