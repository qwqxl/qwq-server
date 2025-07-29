# 🔐 AuthX — 企业级 Go JWT 鉴权框架

> 作者：xiaolin
> 描述：`AuthX` 是一个基于 Go 构建的模块化、可扩展、状态驱动的企业级 JWT 鉴权工具包，支持单点登录（SSO）、单端登录（SLO）、无感刷新、多端隔离、生命周期钩子、Redis 会话存储，适用于微服务与中大型系统的用户登录场景。

---

## 🚀 核心特性

* ✅ AccessToken + RefreshToken 双层机制
* 🧠 Redis 中心态管理，实现会话状态控制
* 🧍‍♂️ SSO（单点登录）与 SLO（单端登录）灵活切换
* 🔄 自动刷新 AccessToken（无感刷新）
* 📱 多端设备隔离（支持 pc/web、mobile/app 等）
* 🔧 生命周期钩子（登录、退出、刷新、踢出）
* 🧩 插拔式中间件，支持 Gin/Echo/Fiber 等框架
* 🔐 自定义 Claims 支持扩展字段
* 📦 支持黑名单、OpenTelemetry 链路追踪、i18n 等扩展能力

---

## 📦 模块结构

```txt
authx/
├── authx.go          // 框架主入口与调度
├── config.go         // 配置结构定义
├── token.go          // JWT 签发与验证
├── claims.go         // 自定义 Claims 结构
├── store.go          // Redis 状态存储逻辑
├── middleware.go     // HTTP 中间件封装
├── lifecycle.go      // 生命周期钩子管理
├── errors.go         // 自定义错误封装
├── model.go          // 输入输出结构体
└── blacklist.go      // Token 黑名单机制（可选扩展）
```

## 使用 jwt + redis

```go
import (
    "github.com/golang-jwt/jwt/v5"
    "github.com/redis/go-redis/v9"
)

```

---

## 🧾 自定义 Claims 结构

```go
type CustomClaims struct {
  UserID       string `json:"uid"`
  DeviceSign   string `json:"did"`   // 设备标识：pc/web、mobile/app
  PlatformSign string `json:"plat"`  // 平台标识：如 qwq、tv、draw
  jwt.RegisteredClaims
}
```

---

## 🗃 Redis Key 设计规范

```
authx:token:{user_id}:{platform}:{device_sign}
```

示例：

```
authx:token:10086:qwq:pc/web
```

> 所有 token 均按用户 + 平台 + 设备隔离，自动 TTL 控制，易于集中踢出。

---

## 🔑 AccessToken 与 RefreshToken

| 类型           | 描述                      | 推荐有效期 |
| ------------ | ----------------------- | ----- |
| AccessToken  | 短效 JWT，用于用户快速鉴权         | 15分钟  |
| RefreshToken | 长效 JWT，用于刷新 AccessToken | 7天    |

---

## 📋 配置结构体

```go
type Config struct {
  JWTSecret    string
  AccessTTL    time.Duration
  RefreshTTL   time.Duration
  EnableSSO    bool
  RedisPrefix  string
  RedisClient  *redis.Client
  Hooks        LifecycleHooks
  UseBlacklist bool  // 启用黑名单模式（可选）
}
```

---

## 🪝 生命周期钩子机制

```go
type HookFunc func(userID, platform, device string) error

type LifecycleHooks struct {
  OnLogin   HookFunc
  OnLogout  HookFunc
  OnRefresh HookFunc
  OnKick    HookFunc
}
```

使用示例：

```go
authx.UseHooks(authx.LifecycleHooks{
  OnLogin: func(uid, plat, dev string) error {
    log.Printf("[LOGIN] %s on %s / %s", uid, plat, dev)
    return nil
  },
})
```

---

## 🧪 使用示例

### 初始化 AuthX：

```go
auth := authx.New(&authx.Config{
  JWTSecret:   "your-secret",
  AccessTTL:   15 * time.Minute,
  RefreshTTL:  7 * 24 * time.Hour,
  RedisPrefix: "authx:token",
  EnableSSO:   true,
  RedisClient: redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
})
```

### 登录：

```go
pair, err := auth.Login(context.TODO(), &authx.LoginInput{
  UserID:     "10086",
  Platform:   "qwq",
  DeviceSign: "pc/web",
})
```

### 中间件（以 Gin 为例）：

```go
r := gin.Default()
r.Use(authx.Middleware())

r.GET("/me", func(c *gin.Context) {
  claims := authx.GetClaims(c)
  c.JSON(200, gin.H{
    "user_id": claims.UserID,
    "platform": claims.PlatformSign,
  })
})
```

### 刷新 Token：

```go
pair, err := auth.Refresh(ctx, refreshToken)
```

---

## 🛡 黑名单机制（可选）

> Token 黑名单用于实现“主动失效”，避免前端持有的 token 长期有效。

```go
err := authx.BlacklistToken(tokenString)
```

---

## 📈 可观测性扩展建议

| 能力               | 说明                          |
| ---------------- | --------------------------- |
| 🔒 RSA/ES256 签发  | 提升安全性，推荐用于前后端分离系统           |
| 🚫 Token 黑名单     | Redis 实现黑名单踢出策略             |
| ☁️ 多租户支持         | Redis key 增加 `tenant_id` 前缀 |
| 🔍 OpenTelemetry | 全链路跟踪 JWT 的签发、验证、刷新         |
| 🌐 i18n 国际化      | 多语言错误输出支持                   |

---

## ✅ 总结

`AuthX` 是一款为企业级项目量身打造的 Go 鉴权框架，通过标准化的 JWT 管理 + Redis 状态中心 + 生命周期控制，提供了安全、灵活、可观测的用户认证体验。
