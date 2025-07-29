# 🔐 **Go AuthX — 企业级 JWT 鉴权工具包**

> 作者：xiaolin  
> 描述：**AuthX** 是一个基于 Go 实现的模块化、可扩展、支持 Redis 状态校验的 JWT 鉴权工具包。支持多平台、单点登录、单端登录、Token 自动刷新、生命周期钩子等中大型系统常用能力。

---

## 🚀 为什么选择 AuthX？

- 支持 **多平台多设备隔离登录**
- 支持 **单点登录 / 单端登录策略切换**
- 采用 **Redis 状态同步管理**，支持注销与踢下线
- 提供 **Token 无感刷新机制**
- 内置 **生命周期钩子**
- 易于集成到 Gin、Echo、Fiber 等主流框架
- 支持自定义 `Claims`，适配多业务系统

---

## 📦 项目结构（优化命名）

```text
authx/
├── authx.go          // 对外主接口，组织调度器
├── config.go         // 配置结构体定义
├── token.go          // JWT Token 签发与解析逻辑
├── store.go          // 状态存储模块（Redis 实现）
├── claims.go         // 自定义 Claims 定义与处理
├── middleware.go     // Gin/Fiber 等 HTTP 中间件封装
├── lifecycle.go      // 生命周期钩子定义与触发
├── model.go          // 通用输入输出模型
└── errors.go         // 错误定义与封装
```

---

## 📌 自定义 Claims 结构

```go
type CustomClaims struct {
  UserID     string `json:"uid"`
  DeviceSign string `json:"did"`   // 设备标识：如 pc/web、mobile/app
  Platform   string `json:"plat"`  // 平台标识：如 qwq、tv、draw
  jwt.RegisteredClaims
}
```

---

## 🔑 Redis Key 命名规则

```
authx:token:{user_id}:{platform}:{device_sign}
```

### 示例：

```
authx:token:10086:qwq:pc/web
```

> ✅ 多平台隔离  
> ✅ 多设备隔离  
> ✅ 结构化统一命名便于后期扩展

---

## ✅ 功能能力总览

| 功能                 | 说明                                              |
|----------------------|---------------------------------------------------|
| Access Token         | 短期有效（15分钟），用户认证主凭证               |
| Refresh Token        | 长期有效（7天），用于无感刷新 AccessToken        |
| 单点登录（SSO）       | 同一用户只能在一个平台登录                      |
| 单端登录（SLO）       | 同一平台设备仅允许一端在线                      |
| Redis 状态管理       | Token 状态存储在 Redis，支持主动失效             |
| 生命周期钩子         | 登录 / 注销 / 刷新 / 踢下线均可触发自定义事件   |
| 中间件封装           | 提供统一鉴权中间件，便于与 Gin、Fiber 等集成    |

---

## 🔄 Token 无感刷新机制

```go
func (a *AuthX) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
```

1. 校验 RefreshToken 合法性与 Redis 状态一致性
2. 颁发新 AccessToken 与 RefreshToken
3. 更新 Redis 内容
4. 支持触发 `OnRefresh` 钩子事件

---

## 🔁 生命周期钩子定义

```go
type HookFunc func(userID, platform, device string) error

type LifecycleHooks struct {
  OnLogin   HookFunc
  OnLogout  HookFunc
  OnRefresh HookFunc
  OnKick    HookFunc
}
```

注册方式：

```go
authx.UseHooks(authx.LifecycleHooks{
  OnLogin: func(uid, plat, dev string) error {
    log.Printf("[LOGIN] %s @ %s / %s", uid, plat, dev)
    return nil
  },
})
```

---

## 🧪 快速使用示例

### 初始化：

```go
auth := authx.New(&authx.Config{
  JWTSecret:     "your-secret",
  AccessTTL:     15 * time.Minute,
  RefreshTTL:    7 * 24 * time.Hour,
  RedisPrefix:   "authx:token",
  EnableSSO:     true,
  RedisClient:   redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
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

返回：

```json
{
  "access_token": "xxx",
  "refresh_token": "yyy",
  "expires_in": 900
}
```

### 鉴权中间件（Gin 示例）：

```go
r := gin.Default()
r.Use(authx.Middleware())

r.GET("/profile", func(c *gin.Context) {
  claims := authx.GetClaims(c)
  // ...
})
```

---

## ⚙️ 配置结构体

```go
type Config struct {
  JWTSecret   string
  AccessTTL   time.Duration
  RefreshTTL  time.Duration
  EnableSSO   bool
  RedisPrefix string
  RedisClient *redis.Client
  Hooks       LifecycleHooks
}
```

---

## 📈 可扩展建议

| 建议方向             | 内容                                                         |
|----------------------|--------------------------------------------------------------|
| 支持非对称加密 JWT   | 使用 RSA 或 ES256 公私钥对                                   |
| 黑名单机制           | 配置 Redis BlackList Key 做拉黑判断                         |
| WebSocket 通知支持   | 登出/踢出时通知客户端断连                                   |
| 多语言错误返回       | 使用 i18n 封装错误消息                                       |
| 多租户支持           | Redis Key 增加 tenant_id 维度支持 SaaS 应用                 |

---

## ✅ 总结

**AuthX** 是一款现代化的 Go 鉴权工具包，专为复杂应用场景设计。相比传统无状态 JWT，它结合 Redis 存储状态，兼顾「高性能」与「可控性」，让你轻松实现：

- Token 登录鉴权
- 多平台设备隔离登录
- Token 状态刷新
- 生命周期事件通知

---

如果你希望我帮你生成 `authx/` 的具体代码结构（例如 `token.go`、`store.go`）或编写单元测试、集成案例，请继续告诉我 ✅