# Errno 模块

> 统一错误码体系，定义错误接口和标准错误响应格式。

## 设计决策

### 为什么需要这个模块

统一的错误码体系使得：

- 前后端通过错误码而非字符串匹配进行错误处理
- 错误码按范围分段，快速定位错误来源
- 所有 API 返回一致的 `{code, message, data}` 格式

### 为什么这么设计

- **选择了**：`Errno` 接口 + 预定义错误变量 + `Decode(data, err)` 函数
- **而不是**：使用自定义 error 类型 + HTTP middleware 统一处理
- **原因**：显式编码，每个预定义错误携带 HTTP 状态码和业务错误码，在 Handler 层统一转换

## 关键类型与接口

### Errno 接口

```go
type Errno interface {
    Error() string
    GetHTTPStatus() int
    GetCode() int
    GetMessage() string
    GetData() any
}
```

### Decode 函数

```go
func Decode(data any, err error) Errno
```

- 如果 `err == nil`，返回 `{code: 0, message: "OK", data: data}`
- 如果 `err != nil`，尝试提取 `Errno` 接口，失败则返回 `InternalServerError`

### 错误码分段

| 范围 | 类别 | 示例 |
|------|------|------|
| `0` | 成功 | `OK` |
| `1000-1099` | 通用错误 | `InternalServerError`, `RequestParserError`, `RequestValidateError` |
| `1100-1199` | 认证/授权 | `TokenInvalidError`, `PermissionDeniedError`, `RateLimitExceededError` |
| `1200-1299` | 文件相关 | `RetrievingFileError` |
| `1300+` | 数据库 | `DatabaseError`, `DatabaseConversionError` |
| `2000+` | 业务错误 | `UserNotFoundError`, `UserAlreadyExistsError` |
| `3000+` | 分页错误 | `InvalidPageError`, `InvalidPageSizeError` |

## 模块结构

```text
internal/errno/
├── errno.go          # Errno 接口定义 + Decode 函数
├── code.go           # 通用/认证/文件/数据库错误码定义
├── business.go       # 业务错误码定义（用户、分页等）
└── errno_test.go     # 单元测试
```

| 文件 | 职责 |
|------|------|
| `errno.go` | 定义 `Errno` 接口和 `errno` 结构体，实现 `Decode` 函数 |
| `code.go` | 定义基础设施层错误码（1000-1399） |
| `business.go` | 定义业务层错误码（2000+） |

## 与其他模块的关系

### 依赖

- 无内部模块依赖（基础层）

### 被依赖

- **handler**：构造 HTTP 错误响应
- **service**：返回业务错误

## 注意事项

- `Errno` 接口的 `Error()` 方法返回 `Message`，因此可直接作为 `error` 使用
- 所有预定义错误变量统一使用 `var` 声明，不可变
- 新增业务错误码时，在 `business.go` 中添加对应变量
- 新增基础设施错误码时，在 `code.go` 中添加对应变量
- 错误响应格式：`{code: <int>, message: "<string>", data: <any>}`
