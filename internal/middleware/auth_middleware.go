package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/auth"
	"qwqserver/internal/common"
	"qwqserver/pkg/cache"
	"strings"
	"time"
)

type HandlerFunc func(c *gin.Context) (auth.CodeType, *common.HTTPResult)

var excludePaths = map[string]ExcludeRouter{
	"/api/v1/auth/register": {
		IsValid: true,
		HandlerFunc: func(c *gin.Context) (auth.CodeType, *common.HTTPResult) {
			if isLoggedIn(c) {
				c.AbortWithStatusJSON(http.StatusOK, common.HTTPResult{
					Msg: "已登录，请不要重复注册",
				})
				return auth.IdentitySkipped, nil
			}
			c.Next()
			return auth.IdentitySkipped, nil
		},
	},
	"/api/v1/auth/login": {
		IsValid: true,
		HandlerFunc: func(c *gin.Context) (auth.CodeType, *common.HTTPResult) {
			if isLoggedIn(c) {
				c.AbortWithStatusJSON(http.StatusOK, common.HTTPResult{
					Msg: "已登录，请不要重复登录",
				})
				return auth.IdentitySkipped, nil
			}
			c.Next()
			return auth.IdentitySkipped, nil
		},
	},
}

type ExcludeRouter struct {
	IsValid     bool
	HandlerFunc HandlerFunc
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查排除路径
		if router, ok := excludePaths[c.Request.URL.Path]; ok && router.IsValid {
			router.HandlerFunc(c)
			return // 直接返回，不进入后续流程
		}

		// 2. 标准认证流程
		authCode, res := JWTAuth(c)

		// 3. 处理认证结果
		switch authCode {
		case auth.IdentityOK:
			c.Next()
		case auth.IdentitySkipped: // 通常不会发生
			c.Next()
		default:
			if res == nil {
				res = &common.HTTPResult{
					Code: http.StatusUnauthorized,
					Msg:  "认证失败",
				}
			}
			c.AbortWithStatusJSON(res.Code, res)
		}
	}
}

// 轻量级登录检查 (不刷新token/不设置上下文)
func isLoggedIn(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}

	tokenString := parts[1]
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		return false
	}

	// 检查Redis中的token有效性
	redisKey := common.RedisTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	storedToken, err := cache.Get(redisKey)
	return err == nil && storedToken == tokenString
}

func JWTAuth(c *gin.Context) (auth.CodeType, *common.HTTPResult) {
	// 获取Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return auth.IdentityErrNoToken, &common.HTTPResult{
			Code: http.StatusUnauthorized,
			Msg:  "未提供认证令牌",
		}
	}

	// 检查Bearer格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return auth.IdentityErrTokenFormat, &common.HTTPResult{
			Code: http.StatusUnauthorized,
			Msg:  "令牌格式错误",
		}
	}

	tokenString := parts[1]
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		return auth.IdentityErrInvalidToken, &common.HTTPResult{
			Code: http.StatusUnauthorized,
			Msg:  "无效令牌",
		}
	}

	// 检查Redis中token有效性
	redisKey := common.RedisTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	storedToken, err := cache.Get(redisKey)
	if err != nil || storedToken != tokenString {
		return auth.IdentityErrTokenExpired, &common.HTTPResult{
			Code: http.StatusUnauthorized,
			Msg:  "令牌已失效",
		}
	}

	// 无感刷新
	if time.Until(claims.ExpiresAt.Time) < common.TokenRefreshInterval {
		newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
		if err == nil {
			c.Header("New-Token", newToken)
			cache.Set(redisKey, newToken, common.TokenExpireTime)
		}
	}

	// 设置上下文信息
	c.Set(auth.IdentityStatusKey, auth.IdentityOK)
	c.Set("user_id", claims.UserID)
	c.Set("platform", claims.Platform)
	c.Set("device_id", claims.DeviceID)

	return auth.IdentityOK, &common.HTTPResult{
		Code: http.StatusOK,
		Msg:  "认证成功",
		Data: claims,
	}
}
