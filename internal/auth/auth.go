package auth

import "time"

const (
	LoginPath    = "/login"
	RegisterPath = "/register"
	LogoutPath   = "/logout"
	DelIDPath    = "/del"
	ProfilePath  = "/profile"
	Identity     = "/identity"
)

// Auth Middleware 鉴权状态码
// 未提供认证令牌
// 令牌格式错误
// 无效令牌
// 令牌已失效

const (
	IdentityStatusKey = "identity_status" // 身份认证状态Key
)

type CodeType uint64

const (
	IdentityErrAuthFailed   CodeType = 1 << iota
	IdentityErrNoToken               // ErrNoToken 未提供认证令牌
	IdentityErrTokenFormat           // ErrTokenFormat 令牌格式错误
	IdentityErrInvalidToken          // ErrInvalidToken 认证令牌无效
	IdentityErrTokenExpired          // ErrTokenExpired 令牌已失效
	IdentityOK                       // OK 认证成功
	IdentitySkipped                  // Skipped 跳过认证
)

const (
	PlatformSign = "qwq" // 平台标识
	JWTSecretKey = "your_secure_secret_key"

	TokenExpireTime      = 2 * time.Hour      // Token有效期
	RefreshTokenExpire   = 7 * 24 * time.Hour // RefreshToken有效期
	TokenRefreshInterval = 30 * time.Minute   // Token刷新间隔

	// Redis键名前缀
	RedisTokenPrefix      = "user_token:"
	RedisRefreshPrefix    = "refresh_token:"
	RedisUserDevicePrefix = "user_device:"
)
