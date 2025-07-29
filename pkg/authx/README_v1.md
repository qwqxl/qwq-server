# 🔐 AuthX — 企业级 Go JWT 鉴权框架

> 作者：xiaolin
>
> 描述：`AuthX` 是一个基于 Go 构建的模块化、可扩展、状态驱动的企业级 JWT 鉴权工具包，支持单点登录（SSO）、单端登录（SLO）、无感刷新、多端隔离、生命周期钩子、Redis 会话存储，适用于微服务与中大型系统的用户登录场景。

---

## 🚀 核心特性

*   ✅ **AccessToken + RefreshToken**：采用双层令牌机制，兼顾安全与性能。
*   🧠 **Redis 中心化会话管理**：所有会话状态存储于 Redis，实现强制下线、状态同步等功能。
*   🧍‍♂️ **SSO/SLO 模式切换**：支持单点登录（SSO）和单端登录（SLO），可根据业务需求灵活配置。
*   🔄 **无感刷新**：AccessToken 过期后可使用 RefreshToken 自动刷新，用户体验无中断。
*   📱 **多端设备隔离**：支持 PC、Web、Mobile、App 等多端同时登录，也可配置为互踢模式。
*   🔧 **生命周期钩子**：提供登录、退出、刷新、踢出等关键节点的钩子函数，便于扩展业务逻辑（如记录日志、发送通知）。
*   🧩 **框架无关设计**：核心逻辑与 Web 框架解耦，提供 Gin 中间件实现，也可轻松扩展至 Echo、Fiber 等其他框架。
*   🔐 **自定义 Claims**：支持在 JWT 的 Claims 中扩展自定义字段，满足不同业务场景需求。
*   📦 **高可扩展性**：预留了黑名单、OpenTelemetry 链路追踪、i18n 等扩展接口。

---

## 🏗️ 架构设计

`AuthX` 的核心设计思想是将 **无状态的 JWT** 与 **有状态的 Redis 会话** 相结合，实现了灵活、可控的鉴权机制。

1.  **登录 (Login)**
    *   用户提供凭证，验证通过后，`AuthX` 生成 `AccessToken` 和 `RefreshToken`。
    *   同时，在 Redis 中创建一个会话记录，key 与 JWT 中的唯一标识对应，并设置与 `RefreshToken` 相同的过期时间。
2.  **鉴权 (Authentication)**
    *   客户端携带 `AccessToken` 访问受保护资源。
    *   服务端中间件验证 `AccessToken` 的签名、过期时间等。
    *   **（重要）** 中间件还会查询 Redis 中是否存在对应的会话，如果会话不存在（例如已被主动踢出），即使 `AccessToken` 本身有效，验证也会失败。
3.  **刷新 (Refresh)**
    *   当 `AccessToken` 过期时，客户端使用 `RefreshToken` 请求新的 `AccessToken`。
    *   服务端验证 `RefreshToken`，并检查 Redis 中对应的会话是否存在。
    *   验证通过后，生成新的 `AccessToken` 和 `RefreshToken`，并刷新 Redis 中会话的过期时间。
4.  **登出/踢出 (Logout/Kick-out)**
    *   主动删除 Redis 中的会话记录，即可实现对应用户/设备的强制下线。

---

## 📦 模块结构

```txt
authx/
├── authx.go          // 框架主入口与核心调度逻辑
├── config.go         // 所有配置项的结构体定义
├── token.go          // JWT 的签发、验证与管理
├── claims.go         // 自定义 Claims 结构，用于扩展 JWT 荷载
├── store.go          // 基于 Redis 的会话存储层
├── middleware.go     // Gin 框架的 HTTP 中间件
├── lifecycle.go      // 生命周期钩子（Hooks）的定义
├── errors.go         // 自定义错误类型封装
└── model.go          // API 的输入输出数据结构
```

---

## 🛠️ 使用指南

### 1. 安装

```bash
go get github.com/your-repo/qwqserver/pkg/authx
```

### 2. 初始化 AuthX

```go
package main

import (
	"time"
	"github.com/redis/go-redis/v9"
	"qwqserver/pkg/authx"
)

func main() {
	// 初始化 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 初始化 AuthX
	ax, err := authx.New(&authx.Config{
		JWTSecret:   "your-super-secret-key",
		AccessTTL:   15 * time.Minute,
		RefreshTTL:  7 * 24 * time.Hour,
		RedisPrefix: "authx:session",
		EnableSSO:   true, // 设置为 true 开启单点登录
		RedisClient: redisClient,
		Hooks: authx.LifecycleHooks{
			OnLogin: func(userID, platform, device string) error {
				// 在此处理登录成功后的逻辑，例如记录日志
				return nil
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// ... 后续业务逻辑
}
```

### 3. 用户登录

```go
// 登录请求
loginInput := &authx.LoginInput{
	UserID:     "user-123",
	Platform:   "web",
	DeviceSign: "chrome-108",
}

// 执行登录
tokenPair, err := ax.Login(context.Background(), loginInput)
if err != nil {
	// 处理错误
}

// 将 tokenPair 返回给客户端
// tokenPair.AccessToken
// tokenPair.RefreshToken
```

### 4. 在 Gin 中使用中间件

```go
import "github.com/gin-gonic/gin"

func setupRouter(ax *authx.AuthX) *gin.Engine {
	r := gin.Default()

	// 应用 AuthX 中间件
	authRequired := r.Group("/api")
	authRequired.Use(ax.Middleware())
	{
		authRequired.GET("/me", func(c *gin.Context) {
			// 从 Gin 上下文中安全地获取 Claims
			claims := authx.GetClaims(c)
			if claims == nil {
				c.JSON(401, gin.H{"error": "invalid token"})
				return
			}
			c.JSON(200, gin.H{
				"user_id":  claims.UserID,
				"platform": claims.PlatformSign,
				"device":   claims.DeviceSign,
			})
		})
	}

	return r
}
```

### 5. 刷新 Token

```go
// 从客户端获取 refreshToken
refreshToken := "..."

// 刷新令牌
newTokenPair, err := ax.Refresh(context.Background(), refreshToken)
if err != nil {
	// 如果刷新失败，通常意味着会话已过期或被吊销，需要用户重新登录
}

// 将新的 tokenPair 返回给客户端
```

### 6. 主动踢出用户

你可以踢出某个用户在特定平台或特定设备上的会话。

```go
// 踢出 user-123 在 web 平台上的 chrome-108 这个设备的会话
err := ax.KickOut(context.Background(), "user-123", "web", "chrome-108")
if err != nil {
	// 处理错误
}
```

---

## ⚙️ API 参考

### `authx.New(*Config) (*AuthX, error)`

初始化一个新的 `AuthX` 实例。

### `(*AuthX).Login(context.Context, *LoginInput) (*TokenPair, error)`

处理用户登录，返回 `AccessToken` 和 `RefreshToken`。

### `(*AuthX).Logout(context.Context, *CustomClaims) error`

用户主动登出，会删除 Redis 中的会话。

### `(*AuthX).Refresh(context.Context, string) (*TokenPair, error)`

使用 `RefreshToken` 字符串刷新会话，获取新的令牌对。

### `(*AuthX).ValidateToken(context.Context, string) (*CustomClaims, error)`

验证 `AccessToken` 字符串，返回解析后的 `CustomClaims`。

### `(*AuthX).KickOut(context.Context, userID, platform, deviceSign string) error`

强制踢出指定用户会话。

### `(*AuthX).Middleware() gin.HandlerFunc`

返回一个 Gin 中间件。

### `authx.GetClaims(*gin.Context) *CustomClaims`

从 `gin.Context` 中安全地获取 `CustomClaims`。

---

## 🗃 Redis Key 设计

`AuthX` 使用结构化的 Redis Key 来管理会话，格式如下：

```
{RedisPrefix}:{UserID}:{Platform}:{DeviceSign}
```

**示例**：

```
authx:session:user-123:web:chrome-108
```

这种设计使得可以精确地控制每个用户的每个会话，例如：

*   **查询**：轻松查询某个用户的所有在线设备。
*   **踢出**：通过删除特定的 key 来踢出单个设备，或使用 `SCAN` 匹配模式 `authx:session:user-123:*` 来踢出该用户的所有设备。

---

## 💡 扩展建议

*   **黑名单机制**：虽然 `AuthX` 通过 Redis 会话实现了主动失效，但你也可以额外实现一个黑名单，用于存放已明确注销的 JWT ID，防止在 Redis 恢复或延迟期间被重放。
*   **加密算法**：默认使用 `HS256`，对于安全性要求更高的场景，可以修改 `token.go` 以支持 `RS256` 等非对称加密算法。
*   **多租户支持**：在 `RedisPrefix` 或 `CustomClaims` 中增加 `TenantID` 字段，即可轻松实现多租户隔离。