# CLI 参考

## 概述

go_sample_code 通过命令行参数控制启动行为。

## 命令格式

```bash
./server [选项]
```

## 选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-c <path>` | 指定 YAML 配置文件路径 | `config.yaml` |

## 使用示例

### 使用默认配置启动

```bash
./server
```

程序将尝试读取当前目录下的 `config.yaml`，如文件不存在则使用硬编码默认值：PostgreSQL `localhost:5432`、Redis `localhost:6379`、监听 `:8080`。

### 指定配置文件

```bash
./server -c config/config.local.yaml
```

### 开发模式

```bash
# 使用 go run
go run cmd/server/main.go -c config/config.local.yaml

# 编译后运行
go build -o server ./cmd/server
./server -c config/config.local.yaml
```

### 生产模式

```bash
# 编译 Linux 二进制
GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# 部署并运行
./server -c /opt/go_sample_code/config/config.yaml
```

## 环境变量

当前版本不支持环境变量配置。所有配置通过 YAML 文件和命令行参数 `-c` 管理。
