# 📘 qwq-server

[![Go Version](https://img.shields.io/github/go-mod/go-version/yourname/qwqserver)](https://github.com/yushulinfengxl/qwq-server)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/license/MIT)

qwq-server 是一个模块化、高性能、极简可扩展的 Go 后端框架，基于 Gin 构建，强调解耦设计、清晰目录结构、易维护性，适用于中大型项目的后端开发。

<p align="center">
  <img src="./resources/docs/logo.png" alt="qwqserver logo" width="180">
</p>

[//]: # (![QWQServer Logo]&#40;./logo.png&#41;)

> 注：当前缺乏专业前端支持，诚邀前端开发者共同构建后台管理界面，欢迎提交意见和 PR！


## 🔍 速览

* 👉 项目演示地址（敬请期待）: [http://www.iqwq.com](http://www.iqwq.com)
* 👉 接口文档 Swagger：[http://localhost:2333/swagger/index.html](http://localhost:2333/swagger/index.html)
* 👉 前端仓库（待）：[https://github.com/yushulinfengxl/qwqserver-web](https://github.com/yushulinfengxl/qwqserver-web)
* 👉 文档入口：[`/docs`](/docs)

---

## ⚙️ 技术栈

| 类别       | 技术                  | 描述                   |
| -------- |---------------------| -------------------- |
| 后端语言     | Go                  | 高性能后端开发语言            |
| 框架       | Gin                 | 轻量级 HTTP Web 框架      |
| 数据库      | MySQL/Postgres/SQLite | 多数据库支持               |
| 缓存       | Redis               | 高性能缓存/消息队列           |
| 安全认证     | JWT                 | 支持访问令牌与刷新令牌，内置无感刷新机制 |
| 配置加载     | Viper（内置封装）         | 动态配置、环境变量加载          |
| 日志记录     | Zap/Slog（待）         | 高性能结构化日志记录           |
| 热重载      | Air（待）              | 本地开发自动重载             |
| 文档生成     | swaggo/swag         | Swagger UI 自动生成      |
| Docker部署 | Docker Compose      | 容器一键部署支持             |
| 前端支持     | Nuxt/Next（待）        | 可选前端集成支持             |

---

## 🧩 功能模块

### 模块结构
- **cmd/**：应用启动入口（main.go）
- **internal/app/**：组合所有服务逻辑
- **internal/config/**：配置加载
- **internal/common/**：公共
- **internal/base/**：基础服务
- **internal/server/**：Gin/gRPC 启动逻辑
- **internal/middleware/**：中间件
- **internal/handler/**：HTTP 控制器
- **internal/model/**：数据模型
- **internal/repository/**：数据库操作
- **internal/service/**：业务逻辑
- **pkg/**：可复用库


### ✅ 用户模块

* 用户注册、登录、登出
* JWT 鉴权（Access/Refresh 分离）
* 登录失败保护机制、密码哈希升级
* 修改密码 / 更新信息

### ✅ 权限系统

* 基于位掩码的权限控制（RBAC 模拟）
* 角色 - 权限映射（待）

### ✅ 配置模块

* 支持 YAML/ENV 动态配置文件
* 多环境切换（dev/test/prod）

### ✅ 日志模块

* 控制台彩色输出
* 支持日志等级（DEBUG/INFO/WARN/ERROR）
* 可扩展文件输出（待）

### ✅ 安全机制

* 图形验证码功能
* 支持 CORS 跨域请求配置
* 支持 CSRF/XSS 安全防护（待）
* 邮箱验证（支持 QQ/Gmail 等）

### ✅ 中间件

* 日志拦截器
* 恶意请求速率限制器（待）
* 鉴权中间件

### ✅ 其他

* OpenAPI 规范文档
* 自动数据库迁移
* 可选 Redis 缓存接入
* Postman 请求集合（待）
* 插件系统（待）

---

## 💻 本地开发

### 安装依赖

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go mod tidy
```

### 配置文件

编辑 `configs/config.yaml`：

```yaml
listen:
  port: 2333
  log_level: debug
  max_concurrent: 10

database:
  driver: mysql
  dsn: root:password@tcp(localhost:3306)/qwq?charset=utf8mb4&parseTime=True&loc=Local

redis:
  addr: localhost:6379
  password: ""

admin_user:
  username: admin
  password: admin
  email: admin@example.com

jwt:
  secret: "your-secret"
  access_expire: 3600
  refresh_expire: 86400
```

### 启动服务

```bash
go run ./cmd/server/main.go
```

或使用 Air：

```bash
go install github.com/air-verse/air@latest
air -c ./configs/.air.toml
```

访问接口：[http://localhost:2333/ping](http://localhost:2333/ping)

---

## 🐳 Docker 部署（MySQL）

编辑 `configs/config.yaml` 和 `docker-compose.yaml`：

```yaml
APP_HOST: "0.0.0.0"
database:
  driver: mysql
  dsn: root:password@tcp(mysql:3306)/qwq?charset=utf8mb4&parseTime=True&loc=Local
```

```yaml
environment:
  - MYSQL_ROOT_PASSWORD=yourpassword
```

启动：

```bash
docker-compose up -d
```

---

## 🗂 接口文档

* Swagger 接口文档：[http://localhost:2333/swagger/index.html](http://localhost:2333/swagger/index.html)
* Postman 文档（待）
* README.md 文档入口：`docs/`

---

## 🛣️ Roadmap

* ✅ 用户注册/登录
* ✅ JWT 管理机制
* ✅ RBAC 权限模型（简化版）
* ✅ 接口分层架构
* ✅ 数据库迁移 + GORM 封装
* ⏳ 插件系统（待）
* ⏳ 管理后台前端（待）
* ⏳ 邮箱验证 + 通知机制（待）
* ⏳ WebSocket / 消息推送（待）
* ⏳ 多租户支持（待）

---

## 🤝 官方社区

欢迎加入微信官方社区，共同探讨后端架构与实践：

* 微信号：qwqcnxl
* QQ 群：5897746
* 邮箱：[xiaolin@iqwq.com](mailto:xiaolin@iqwq.com)

⚠️ 请遵守开源社群规范，禁止非法内容。

---

## 🙏 特别鸣谢

感谢每一位开源贡献者对 qwqserver 的支持，期待更多开发者一起协作完善！

---

## 📊 代码统计（GitHub Actions 自动更新）

| 语言     | 文件数 | 代码行数 | 注释行 | 空白行 | 占比   |
| ------ | --- | ---- | --- | --- | ---- |
| Go     | 80  | 3000 | 600 | 650 | 91%  |
| YAML   | 3   | 200  | 20  | 30  | 6%   |
| Docker | 1   | 30   | 5   | 10  | 1.5% |
| 其他     | 3   | 50   | 3   | 10  | 1.5% |

> 统计排除：docs/, go.mod, LICENSE 等非核心代码。

---

[//]: # (## 📄 许可证)

## ⚖️ 使用声明与道德说明

本项目遵循 [MIT License](https://opensource.org/license/MIT) 协议开源，任何人均可自由使用、复制、修改和发布本项目源代码。

> **本软件按“原样”提供（AS IS），作者不对因使用本项目造成的任何损失或后果承担任何责任。**

同时，作为本项目的作者，我们郑重声明：

- 本项目的初衷是服务于 **合法、正当、正向** 的技术研究与产品开发；
- **严禁将本项目用于任何形式的非法活动**，包括但不限于诈骗、赌博、色情、病毒传播、数据窃取、隐私监控、压迫性技术等；
- 对任何违反法律法规或违背基本人类价值的使用方式，我们**坚决反对**，并保留公开谴责、移除协作、技术封禁等权利。

🙏 请自觉遵守法律法规和道德底线，守护开源生态的纯净与开放。
