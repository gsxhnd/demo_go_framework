# Ent 模块

> Ent ORM 实体定义与代码生成，提供类型安全的数据库访问层。

## 设计决策

### 为什么选择 Ent

Ent 是 Go 语言的类型安全 ORM，通过代码生成提供编译时类型检查。相比手写 SQL 或使用其他 ORM：

- **类型安全**：查询构建器在编译时检查字段和类型
- **代码生成**：减少样板代码，实体变更后重新生成即可
- **图遍历**：支持自然的关联查询（Edges）
- **迁移管理**：内置 schema 迁移支持

### 为什么代码生成不自动化

- **选择了**：手动运行 `ent generate ./internal/ent/schema` 而非 `go:generate`
- **原因**：代码生成需要 `ent` CLI 工具预先安装，手动触发更可控

## 关键类型与接口

### 目录结构

```text
internal/ent/
├── client.go              # Ent Client（自动生成，勿编辑）
├── ent.go                 # Ent 核心（自动生成，勿编辑）
├── mutation.go            # Mutation 构建器（自动生成，勿编辑）
├── runtime.go             # 运行时（自动生成，勿编辑）
├── tx.go                  # 事务支持（自动生成，勿编辑）
├── user.go                # User 实体（自动生成，勿编辑）
├── user_create.go         # User 创建构建器（自动生成，勿编辑）
├── user_update.go         # User 更新构建器（自动生成，勿编辑）
├── user_delete.go         # User 删除构建器（自动生成，勿编辑）
├── user_query.go          # User 查询构建器（自动生成，勿编辑）
├── enttest/               # 测试工具
├── hook/                  # Ent Hooks
├── migrate/               # Schema 迁移
├── predicate/             # 查询谓词
├── runtime/               # 运行时配置
├── schema/                # 实体定义（可编辑！）
│   ├── user.go            # User Schema 定义
│   └── mixin/             # Mixin（如 TimeMixin）
│       └── time.go        # 时间戳混入（created_at, updated_at）
└── user/                  # User 查询辅助
    ├── user.go            # 字段常量
    └── where.go           # 查询条件构建器
```

### User Schema（`internal/ent/schema/user.go`）

```go
type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").Unique().NotEmpty().MaxLen(50),
        field.String("email").Unique().NotEmpty().MaxLen(255),
        field.String("password").NotEmpty().Sensitive(),
        field.String("nickname").MaxLen(100).Optional(),
        field.String("avatar").MaxLen(500).Optional(),
        field.String("phone").MaxLen(20).Optional(),
        field.Bool("is_active").Default(true),
    }
}
```

### TimeMixin（`internal/ent/schema/mixin/time.go`）

为实体自动添加 `created_at` 和 `updated_at` 时间戳字段。

## 与其他模块的关系

### 依赖

- 数据库驱动（`github.com/lib/pq` / `github.com/go-sql-driver/mysql`）

### 被依赖

- **database**：创建 Ent Client
- **repo**：使用 Ent Client 执行查询

## 注意事项

- `internal/ent/schema/` 是唯一可编辑的目录，其他文件由 `ent generate` 自动生成
- 修改 schema 后必须运行 `ent generate ./internal/ent/schema` 重新生成
- `password` 字段使用 `Sensitive()` 标记，序列化时自动隐藏
- 唯一性约束通过 `field.Unique()` 和 `index.Fields().Unique()` 双重保证
- 代码生成需要先安装 ent CLI：`go install entgo.io/ent/cmd/ent@latest`
