package authx

import (
	"context"
)

// AuthX 是鉴权框架的核心结构
type AuthX struct {
	config         *Config
	tokenManager   *TokenManager
	sessionStore   *SessionStore
	lifecycleHooks *LifecycleHooks
}

// New 创建一个新的 AuthX 实例
func New(config *Config) (*AuthX, error) {
	if config.CacheClient == nil || config.JWTSecret == "" {
		return nil, ErrInvalidConfig
	}

	return &AuthX{
		config:         config,
		tokenManager:   NewTokenManager(config.JWTSecret, config.AccessTTL, config.RefreshTTL),
		sessionStore:   NewSessionStore(config.CacheClient, config.RedisPrefix),
		lifecycleHooks: &config.Hooks,
	}, nil
}

// Login 处理用户登录，生成令牌并创建会话
func (ax *AuthX) Login(ctx context.Context, input *LoginInput) (*TokenPair, error) {
	claims := &CustomClaims{
		UserID:       input.UserID,
		DeviceSign:   input.DeviceSign,
		PlatformSign: input.Platform,
	}

	// 如果不是 SSO 模式，则先删除旧会话
	if !ax.config.EnableSSO {
		ax.sessionStore.DeleteSession(ctx, input.UserID, input.Platform, input.DeviceSign)
	}

	// 创建新会话
	err := ax.sessionStore.CreateSession(ctx, input.UserID, input.Platform, input.DeviceSign, ax.config.RefreshTTL)
	if err != nil {
		return nil, err
	}

	// 生成令牌
	tokens, err := ax.tokenManager.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// 触发登录钩子
	if ax.lifecycleHooks.OnLogin != nil {
		ax.lifecycleHooks.OnLogin(input.UserID, input.Platform, input.DeviceSign)
	}

	return tokens, nil
}

// Logout 处理用户登出，删除会话
func (ax *AuthX) Logout(ctx context.Context, claims *CustomClaims) error {
	err := ax.sessionStore.DeleteSession(ctx, claims.UserID, claims.PlatformSign, claims.DeviceSign)
	if err != nil {
		return err
	}

	// 触发登出钩子
	if ax.lifecycleHooks.OnLogout != nil {
		ax.lifecycleHooks.OnLogout(claims.UserID, claims.PlatformSign, claims.DeviceSign)
	}

	return nil
}

// Refresh 刷新 AccessToken
func (ax *AuthX) Refresh(ctx context.Context, refreshTokenString string) (*TokenPair, error) {
	claims, err := ax.tokenManager.VerifyToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// 检查会话是否存在
	_, err = ax.sessionStore.GetSession(ctx, claims.UserID, claims.PlatformSign, claims.DeviceSign)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// 更新会话过期时间
	err = ax.sessionStore.UpdateSessionExpiration(ctx, claims.UserID, claims.PlatformSign, claims.DeviceSign, ax.config.RefreshTTL)
	if err != nil {
		return nil, err
	}

	// 生成新令牌
	newTokens, err := ax.tokenManager.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// 触发刷新钩子
	if ax.lifecycleHooks.OnRefresh != nil {
		ax.lifecycleHooks.OnRefresh(claims.UserID, claims.PlatformSign, claims.DeviceSign)
	}

	return newTokens, nil
}

// ValidateToken 验证 AccessToken
func (ax *AuthX) ValidateToken(ctx context.Context, accessTokenString string) (*CustomClaims, error) {
	return ax.tokenManager.VerifyToken(accessTokenString)
}

// KickOut 踢出用户
func (ax *AuthX) KickOut(ctx context.Context, userID, platform, deviceSign string) error {
	err := ax.sessionStore.DeleteSession(ctx, userID, platform, deviceSign)
	if err != nil {
		return err
	}

	// 触发踢出钩子
	if ax.lifecycleHooks.OnKick != nil {
		ax.lifecycleHooks.OnKick(userID, platform, deviceSign)
	}

	return nil
}
