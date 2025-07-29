
# Logger

一个功能强大且可配置的 Go 日志库，支持异步输出、日志轮转、颜色控制台输出、分组配置等特性。

## ✨ 特性

- 支持日志级别（DEBUG / INFO / WARN / ERROR）过滤输出
- 支持控制台彩色输出（可自定义颜色）
- 支持异步日志写入（带队列大小控制）
- 支持日志文件轮转（按文件大小 + 限制数量）
- 支持配置分组及继承
- 支持路径/文件名模式占位符解析
- 支持多实例日志器，适合微服务日志隔离

## 📦 安装（未上传，存在/pkg/logger目录）

```bash
go get github.com/your/repo/logger
```

## 🚀 快速开始

```go
package main

import "logger"

func main() {
    logger.Init("./logs/") // 初始化日志根目录

    log := logger.Get("app") // 获取默认配置的日志器

    log.Debug("This is a debug message")
    log.Info("This is an info message")
    log.Warn("This is a warning")
    log.Error("This is an error")

    log.Close() // 关闭日志（确保异步日志刷新完成）
}
```

## ⚙️ 自定义配置

```go
logger.Init("./logs/")

logger.Configure("custom", 
    logger.WithGroup("web"),
    logger.WithLogLevel(logger.INFO),
    logger.WithConsoleOutput(true),
    logger.WithAsync(true, 500),
    logger.WithMaxSize(10),
    logger.WithRotateCount(3),
    logger.WithFilePath("{group}/{year}/{month}/{day}"),
    logger.WithFilePattern("{name}_{hh}{ii}{ss}-{i}.log"),
)

log := logger.Get("custom")
log.Info("Custom logger initialized")
```

## 🎨 控制台颜色配置

自定义控制台输出颜色（ANSI 颜色码）：

```go
logger.ConfigureGroup("colorful",
    logger.WithColors(map[logger.Level]string{
        logger.DEBUG: "\033[35m", // 紫色
        logger.INFO:  "\033[36m", // 青色
        logger.WARN:  "\033[33m", // 黄色
        logger.ERROR: "\033[31m", // 红色
    }),
)
```

## 📚 使用案例

```go
    //初始化日志
	logger.Init("./logs")
	defer logger.Close()

	配置数据库日志组
	logger.ConfigureGroup("db",
		logger.WithAsync(false, 0),
		logger.WithRotateCount(5),
		logger.WithFilePath("{group}/{name}/{year}/{month}/{day}"),
		logger.WithFilePattern("{group}_{name}_{yyyy}{mm}{dd}_{hh}{ii}{ss}_{index}.log"),
		logger.WithMaxSize(50),
	)
	
	// 配置服务器日志组
	logger.ConfigureGroup("server",
		logger.WithAsync(false, 0),
		logger.WithRotateCount(5),
		logger.WithFilePath("{group}/{name}/{year}/{month}/{day}"),
		logger.WithFilePattern("{group}_{name}_{yyyy}{mm}{dd}_{hh}{ii}{ss}_{index}.log"),
		logger.WithMaxSize(50),
	)
	
	logger.ConfigureGroup("app",
		logger.WithAsync(false, 0),
		logger.WithRotateCount(5),
		logger.WithFilePath("{group}/{name}/{year}/{month}/{day}"),
		logger.WithFilePattern("{group}_{name}_{yyyy}{mm}{dd}_{hh}{ii}{ss}_{index}.log"),
		logger.WithMaxSize(50),
	)
	
	logger.Configure(config.DBDriverName(),
		logger.WithGroup("db"),
		logger.WithConsoleOutput(true),
		//logger.WithAsync(false, 0), // 关闭异步，方便测试
	)
	
	logger.Configure(config.CacheDriverName(),
		logger.WithGroup("db"),
		logger.WithConsoleOutput(true),
	)
	
	logger.Configure("http",
		logger.WithGroup("server"),
		logger.WithConsoleOutput(true),
	)
	
	logger.Configure("auth",
		logger.WithGroup("server"),
		logger.WithConsoleOutput(true),
	)
	
	logger.Configure("dev",
		logger.WithGroup("app"),
		logger.WithConsoleOutput(true),
	)
```

## 🔁 文件轮转机制

- 单个日志文件超过 `MaxFileSizeMB`（MB）将轮转至下一个文件
- 最多保留 `RotateCount` 个轮转文件
- 当所有文件满时，覆盖第一个文件

## 🧰 支持的文件名/路径占位符

| 占位符       | 含义         |
|--------------|--------------|
| `{name}`     | 日志器名称    |
| `{group}`    | 分组名        |
| `{year}`     | 年 (4位)      |
| `{month}`    | 月 (2位)      |
| `{day}`      | 日 (2位)      |
| `{hour}`     | 小时 (2位)    |
| `{min}`      | 分钟 (2位)    |
| `{sec}`      | 秒 (2位)      |
| `{i}` / `{index}` | 文件序号    |
| `{weekday}`  | 星期几        |
| `{yyyy}` `{mm}` `{dd}` `{hh}` `{ii}` `{ss}` | 简写时间格式支持 |

## 🧪 日志级别控制

可通过 `WithLogLevel` 设置日志级别：

- `DEBUG`（默认）
- `INFO`
- `WARN`
- `ERROR`

```go
logger.WithLogLevel(logger.WARN)
```

将只输出 `WARN` 和 `ERROR` 等级的日志。

## 🔚 关闭日志器

为了确保异步日志正常写入，建议在程序退出前手动关闭：

```go
log.Close()
```

或统一关闭所有日志器（待实现）：

```go
logger.Close() // TODO: 支持全局关闭所有日志器
```

## 📁 日志目录结构示例

```text
./logs/
  └── app/
      └── 2025/
          └── 05/
              └── 30/
                  ├── app_140530-0.log
                  ├── app_140530-1.log
```

## 📜 License

MIT License

---

由 Go 构建，致力于简洁、高性能日志处理。
