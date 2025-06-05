
# Logger 日志库功能文档

## 简介

该日志库是一个轻量、高度可配置的日志框架，支持按名称多实例管理、群组配置继承、日志文件轮转、控制台输出颜色等功能，适用于中大型 Go 项目中的日志体系建设。

## 主要功能

### 1. 日志实例管理

- 每个日志实例由唯一名称（`name`）标识。
- 可动态创建多个日志配置，分别用于不同模块或服务。

```go
logger.Configure("auth-service", logger.WithFilePath("logs/auth"))
logger.Info("auth-service", "用户 %s 登录成功", username)
```

### 2. 群组配置继承（Group Inheritance）

- 支持将多个实例归类到同一 `Group`。
- 子实例可自动继承群组配置（如日志路径、最大文件大小等），可按需覆盖。

```go
logger.Configure("service-A", logger.WithGroup("common"))
logger.Configure("service-B", logger.WithGroup("common"))

// 统一配置组
logger.Configure("common",
    logger.WithFilePath("logs/shared"),
    logger.WithMaxSize(50),
    logger.WithRotateCount(5),
)
```

### 3. 日志输出控制

- 支持控制是否输出到控制台或文件：

```go
logger.Configure("debug", logger.WithConsoleOutput(true), logger.WithFileOutput(false))
```

- 支持最大文件大小控制（单位：MB）：

```go
logger.WithMaxSize(100) // 100MB
```

### 4. 日志轮转与命名规则

- 支持基于文件大小的日志轮转。
- 支持自定义文件名格式，使用占位符自动生成：

| 占位符      | 含义             |
|-------------|------------------|
| `{name}`    | 实例名           |
| `{group}`   | 分组名           |
| `{index}`   | 当前轮转编号     |
| `{date}`    | 日期（YYYYMMDD） |
| `{time}`    | 时间（HHMMSS）   |
| `{pid}`     | 进程ID           |
| `{ts}`      | 时间戳（秒）     |

例如：

```go
logger.WithFilePattern("{group}/{name}-{date}-{index}.log")
```

### 5. 控制台颜色输出

- 默认日志支持颜色输出，增强可读性。
- 支持以下颜色：

| 颜色      | 示例               |
|-----------|--------------------|
| 红色      | 错误日志           |
| 黄色      | 警告日志           |
| 绿色      | 成功操作           |
| 蓝色      | 普通信息           |

### 6. 自定义颜色配置（可选）

- 可设置不同实例或级别对应的颜色配置：

```go
logger.Configure("payment",
    logger.WithConsoleColor("green"), // 自定义控制台输出为绿色
)
```

> 支持颜色代码或 ANSI 转义序列，如 `"[31m"` 代表红色。

### 7. 异步日志队列（如已启用）

- 内部可启用异步日志队列以缓解写入延迟。
- 日志写入操作将提交到缓冲通道，异步落地。

```go
logger.WithAsyncQueue(10000) // 设置队列长度为 10000 条
```

## 使用示例

```go
logger.Init("logs")

logger.Configure("user-api",
    logger.WithGroup("backend"),
    logger.WithFilePattern("{name}-{date}-{index}.log"),
    logger.WithConsoleColor("blue"),
)

logger.Info("user-api", "创建用户: %s", "Alice")
```

## 未来可扩展方向

- 支持日志级别过滤（DEBUG, INFO, WARN, ERROR）
- JSON 日志输出（支持结构化）
- 支持钩子函数（Hook）
- 接入远程日志服务（如 ELK）
