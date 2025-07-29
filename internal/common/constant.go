package common

import "time"

const (
	PlatformSignString = "qwq" // 平台标识
	JWTSecretKey       = "your_secure_secret_key"

	TokenExpireTime      = 2 * time.Hour      // Token有效期
	RefreshTokenExpire   = 7 * 24 * time.Hour // RefreshToken有效期
	TokenRefreshInterval = 30 * time.Minute   // Token刷新间隔

	// Redis键名前缀
	RedisTokenPrefix      = "user_token:"
	RedisRefreshPrefix    = "refresh_token:"
	RedisUserDevicePrefix = "user_device:"
)

func PlatformSign() string {
	return PlatformSignString
}
