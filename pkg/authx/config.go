package authx

import (
	"time"

	"qwqserver/pkg/cache"
)

// Config AuthX 的配置结构
type Config struct {
	JWTSecret    string         // JWT 密钥
	AccessTTL    time.Duration  // AccessToken 有效期
	RefreshTTL   time.Duration  // RefreshToken 有效期
	EnableSSO    bool           // 是否启用 SSO（单点登录）
	RedisPrefix  string         // Redis 键前缀
	CacheClient  *cache.Client  // Cache 客户端实例
	Hooks        LifecycleHooks // 生命周期钩子
	UseBlacklist bool           // 是否启用黑名单模式
}
