# qwqserver

[![Go Version](https://img.shields.io/github/go-mod/go-version/yourname/qwqserver)](https://github.com/yushulinfengxl/qwq-server)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

一个基于 Go 语言构建的低性能后端。

## 功能特性

- ✅ **RESTful API**：支持标准的 HTTP 请求处理
- 🚀 **高性能路由**：基于 [Gin](https://github.com/gin-gonic/gin) 框架
- 🔒 **身份验证**：JWT 中间件
- 📦 **模块化设计**：清晰的代码分层（Handler/Service/Repository）
- 📊 **数据存储**：支持 MySQL/PostgreSQL + GORM/SQLx

### qwqserver 实现了以下 2 类功能：

- 用户管理： 支持用户注册、用户登录、获取用户列表、获取用户详情、更新用户信息、修改用户密码等操作；
- 博客管理： 支持创建博客、获取博客列表、获取博客详情、更新博客内容、删除博客等操作。

### 模块结构

```
qwqserver/
├── cmd/                 # 主程序入口（一个子目录对应一个可执行文件）
│   └── server/          # 主服务入口（例如：main.go）
│       └── main.go
├── internal/            # 私有代码（禁止被其他项目导入）
│   ├── app/             # 应用核心逻辑初始化（连接组件）
│   ├── config/          # 配置解析（环境变量、YAML/TOML等）
│   ├── server/          # HTTP/gRPC 服务启动和路由配置
│   ├── handler/         # HTTP 请求处理器（按业务拆分）
│   ├── middleware/      # 中间件（认证、日志、限流等）
│   ├── model/           # 数据模型定义（结构体）
│   ├── repository/      # 数据库操作层（ORM/SQL）
│   └── service/         # 业务逻辑层（解耦 handler 和 repository）
├── pkg/                 # 可公开复用的工具库（其他项目可导入）
│   ├── logger/          # 日志模块（Zap/Slog 封装）
│   ├── util/            # 通用工具（加密、验证等）
│   │   ├── cache/       # 缓存操作（Redis/Memcached）
│   │   ├── error/       # 错误处理（自定义错误类型）
│   │   ├── validator/   # 验证器（Gin-Validator/Go-Validator）
│   │   └── security/    # 安全工具（JWT/OAuth2）
│   └── database/        # 数据库连接池和扩展方法
├── api/                 # API 协议定义（Protobuf/OpenAPI）
├── configs/             # 配置文件模板（YAML/TOML/ENV）
├── deployments/         # 部署配置（Dockerfile, k8s, compose）
├── scripts/             # 辅助脚本（部署、代码生成）
├── test/                # 集成测试和测试数据
├── web/                 # 前端资源（可选，如静态文件）
├── go.mod               # Go 模块定义
├── go.sum               # 依赖校验
└── README.md            # 项目文档
```

## 快速开始

### 前置条件

- Go 1.21+
- MySQL/Redis（或 Docker 容器）

### 安装与运行

```bash
# 克隆项目
git clone https://github.com/yushulinfengxl/qwqserver.git
cd qwqserver

# 安装依赖
go mod download

# 复制配置文件模板
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件（按需修改数据库等配置）
vim configs/config.yaml

# 构造项目
go mod tidy

# 构建项目Linux
go build -o bin/qwqserver cmd/server/main.go

# 构建项目Windows
go build -o bin/qwqserver.exe cmd/server/main.go

# 启动服务
go run cmd/server/main.go

# 访问 http://localhost:8080/healthcheck