package user

import (
	"context"
	"errors"
	"qwqserver/internal/app"
	"qwqserver/internal/auth"
	"qwqserver/internal/config"
	"qwqserver/pkg/cache"
)

// 刷新Token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	conf, _ := config.New()
	authConfig := conf.Auth
	claims, err := auth.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	cacheClient, err := app.Get[*cache.Client]()
	if err != nil {
		return "", "", errors.New("获取缓存客户端err：" + err.Error())
	}

	// 验证RefreshToken是否有效
	refreshKey := authConfig.CacheRefreshPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	storedRefresh, err := cacheClient.Get(ctx, refreshKey)
	if err != nil || storedRefresh != refreshToken {
		return "", "", errors.New("无效的刷新令牌")
	}

	// 生成新Token和RefreshToken
	newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := auth.GenerateRefreshToken(claims.UserID, claims.Platform, claims.DeviceID)
	if err != nil {
		return "", "", err
	}

	// 更新Redis中的Token
	tokenKey := authConfig.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	err = cacheClient.Set(ctx, tokenKey, newToken, authConfig.TokenExpireTime)
	if err != nil {
		return "", "", err
	}

	// 更新RefreshToken
	err = cacheClient.Set(ctx, refreshKey, newRefreshToken, authConfig.RefreshTokenExpire)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
}
