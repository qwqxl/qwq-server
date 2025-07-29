package authx

import "errors"

var (
	// ErrInvalidToken 表示令牌无效
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired 表示令牌已过期
	ErrTokenExpired = errors.New("token expired")
	// ErrSessionNotFound 表示会话未找到
	ErrSessionNotFound = errors.New("session not found")
	// ErrInvalidConfig 表示配置无效
	ErrInvalidConfig = errors.New("invalid config")
)