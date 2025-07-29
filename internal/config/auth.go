package config

import "time"

type Auth struct {
	Type string `yaml:"type" qwq-default:"jwt"` // 认证方式

	PlatformSign string `yaml:"platform_sign" qwq-default:"qwq"`                 // 平台标识
	Issuer       string `yaml:"issuer" qwq-default:"qwq"`                        // 签发者
	SecretKey    string `yaml:"secret_key" qwq-default:"your_secure_secret_key"` // JWT密钥

	TokenExpireTime      time.Duration `yaml:"token_expire_time" qwq-default:"2h"`       // 令牌有效期
	RefreshTokenExpire   time.Duration `yaml:"refresh_token_expire" qwq-default:"168h"`  // 刷新令牌有效期
	TokenRefreshInterval time.Duration `yaml:"token_refresh_interval" qwq-default:"24h"` // 令牌刷新间隔

	CacheTokenPrefix      string `yaml:"cache_token_prefix" qwq-default:"user_token:"`        // 缓存令牌前缀
	CacheRefreshPrefix    string `yaml:"cache_refresh_prefix" qwq-default:"refresh_token:"`   // 缓存刷新令牌前缀
	CacheUserDevicePrefix string `yaml:"cache_user_device_prefix" qwq-default:"user_device:"` // 缓存用户设备前缀
}
