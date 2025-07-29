package middleware

import (
	"context"
	"errors"
	"net/http"
	"qwqserver/internal/app"
	"qwqserver/internal/auth"
	"qwqserver/internal/config"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/httpcore"
	"strings"
	"time"
)

// 定义认证状态常量
const (
	IdentityOK = iota
	IdentitySkipped
	IdentityErrNoToken
	IdentityErrTokenFormat
	IdentityErrInvalidToken
	IdentityErrTokenExpired
	IdentityErrAlreadyLoggedIn
	IdentityErrInternal
)

type HandlerFunc func(c httpcore.Context) (int, error)

type ExcludeRouter struct {
	IsValid     bool
	HandlerFunc HandlerFunc
}

var excludePaths = map[string]ExcludeRouter{
	"/api/v1/user/register": {
		IsValid: true,
		HandlerFunc: func(c httpcore.Context) (int, error) {
			if isLoggedIn(c) {
				return IdentityErrAlreadyLoggedIn, errors.New("已登录，请不要重复注册")
			}
			return IdentitySkipped, nil
		},
	},
	"/api/v1/user/login": {
		IsValid: true,
		HandlerFunc: func(c httpcore.Context) (int, error) {
			if isLoggedIn(c) {
				return IdentityErrAlreadyLoggedIn, errors.New("已登录，请不要重复登录")
			}
			return IdentitySkipped, nil
		},
	},
}

func Auth() httpcore.HandlerFunc {
	return func(c httpcore.Context) {
		conf, err := config.New()
		if err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, map[string]any{
				"code":    http.StatusInternalServerError,
				"message": "服务器配置错误",
			})
			return
		}

		// 本地开发模式跳过认证
		if conf.Server.Mode == "debug" {
			c.Next()
			return
		}

		// 排除路径
		if router, ok := excludePaths[c.Request().URL.Path]; ok && router.IsValid {
			code, err := router.HandlerFunc(c)
			if err != nil {
				c.Abort()
				c.JSON(http.StatusBadRequest, map[string]any{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}
			if code == IdentitySkipped {
				c.Next()
				return
			}
		}

		// 执行认证
		code, claims, err := JWTAuth(c)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, map[string]any{
				"code":    http.StatusUnauthorized,
				"message": "认证失败: " + err.Error(),
			})
			return
		}

		if code == IdentityOK {
			c.Set("user_id", claims.UserID)
			c.Set("platform", claims.Platform)
			c.Set("device_id", claims.DeviceID)
			c.Next()
		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, map[string]any{
				"code":    http.StatusUnauthorized,
				"message": "认证失败",
			})
		}
	}
}

func isLoggedIn(c httpcore.Context) bool {
	conf, err := config.New()
	if err != nil {
		return false
	}

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}

	token := parts[1]
	claims, err := auth.ParseToken(token)
	if err != nil || claims == nil {
		return false
	}

	redisKey := conf.Auth.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	cacheClient, err := app.Get[*cache.Client]()
	if err != nil {
		return false
	}

	storedToken, err := cacheClient.Get(context.Background(), redisKey)
	return err == nil && storedToken == token
}

func JWTAuth(c httpcore.Context) (int, *auth.CustomClaims, error) {
	conf, err := config.New()
	if err != nil {
		return IdentityErrInternal, nil, errors.New("服务器配置错误")
	}

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return IdentityErrNoToken, nil, errors.New("未提供认证令牌")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return IdentityErrTokenFormat, nil, errors.New("令牌格式错误")
	}

	token := parts[1]
	claims, err := auth.ParseToken(token)
	if err != nil || claims == nil {
		return IdentityErrInvalidToken, nil, errors.New("无效令牌")
	}

	cacheClient, err := app.Get[*cache.Client]()
	if err != nil {
		return IdentityErrInternal, nil, errors.New("缓存客户端获取失败：" + err.Error())
	}

	redisKey := conf.Auth.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	storedToken, err := cacheClient.Get(context.Background(), redisKey)
	if err != nil || storedToken != token {
		return IdentityErrTokenExpired, nil, errors.New("令牌已失效")
	}

	// 静默刷新令牌
	if time.Until(claims.ExpiresAt.Time) < 5*time.Minute {
		newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
		if err == nil {
			c.Header("New-Token", newToken)
			_ = cacheClient.Set(context.Background(), redisKey, newToken, conf.Auth.TokenExpireTime)
		}
	}

	return IdentityOK, claims, nil
}

//
//type HandlerFunc func(c *gin.Context) (int, error)
//
//type ExcludeRouter struct {
//	IsValid     bool
//	HandlerFunc HandlerFunc
//}
//
//var excludePaths = map[string]ExcludeRouter{
//	"/api/v1/user/register": {
//		IsValid: true,
//		HandlerFunc: func(c *gin.Context) (int, error) {
//			if isLoggedIn(c) {
//				return IdentityErrAlreadyLoggedIn, errors.New("已登录，请不要重复注册")
//			}
//			return IdentitySkipped, nil
//		},
//	},
//	"/api/v1/user/login": {
//		IsValid: true,
//		HandlerFunc: func(c *gin.Context) (int, error) {
//			if isLoggedIn(c) {
//				return IdentityErrAlreadyLoggedIn, errors.New("已登录，请不要重复登录")
//			}
//			return IdentitySkipped, nil
//		},
//	},
//}
//
//func Auth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		conf, err := config.New()
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//				"code":    http.StatusInternalServerError,
//				"message": "服务器配置错误",
//			})
//			return
//		}
//
//		// 测试环境跳过认证
//		if conf.Server.Mode == "debug" {
//			c.Next()
//			return
//		}
//
//		// 1. 检查排除路径
//		if router, ok := excludePaths[c.Request.URL.Path]; ok && router.IsValid {
//			code, err := router.HandlerFunc(c)
//			if err != nil {
//				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
//					"code":    http.StatusBadRequest,
//					"message": err.Error(),
//				})
//				return
//			}
//			if code == IdentitySkipped {
//				c.Next()
//				return
//			}
//		}
//
//		// 2. 标准认证流程
//		authCode, claims, err := JWTAuth(c)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//				"code":    http.StatusUnauthorized,
//				"message": "认证失败: " + err.Error(),
//			})
//			return
//		}
//
//		// 3. 处理认证结果
//		switch authCode {
//		case IdentityOK:
//			// 设置用户上下文信息
//			c.Set("user_id", claims.UserID)
//			c.Set("platform", claims.Platform)
//			c.Set("device_id", claims.DeviceID)
//			c.Next()
//		default:
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//				"code":    http.StatusUnauthorized,
//				"message": "认证失败",
//			})
//		}
//	}
//}
//
//// 轻量级登录检查
//func isLoggedIn(c *gin.Context) bool {
//	conf, err := config.New()
//	if err != nil {
//		return false
//	}
//
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		return false
//	}
//
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return false
//	}
//
//	tokenString := parts[1]
//	claims, err := auth.ParseToken(tokenString)
//	if err != nil || claims == nil {
//		return false
//	}
//
//	// 检查Redis中的token有效性
//	redisKey := conf.Auth.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
//
//	cacheClient, err := app.Get[*cache.Client]()
//	if err != nil {
//		return false
//	}
//
//	storedToken, err := cacheClient.Get(context.Background(), redisKey)
//	return err == nil && storedToken == tokenString
//}
//
//func JWTAuth(c *gin.Context) (int, *auth.CustomClaims, error) {
//	conf, err := config.New()
//	if err != nil {
//		return IdentityErrInternal, nil, errors.New("服务器配置错误")
//	}
//	authConfig := conf.Auth
//
//	// 获取Authorization header
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		return IdentityErrNoToken, nil, errors.New("未提供认证令牌")
//	}
//
//	// 检查Bearer格式
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return IdentityErrTokenFormat, nil, errors.New("令牌格式错误")
//	}
//
//	tokenString := parts[1]
//	claims, err := auth.ParseToken(tokenString)
//	if err != nil || claims == nil {
//		return IdentityErrInvalidToken, nil, errors.New("无效令牌")
//	}
//
//	cacheClient, err := app.Get[*cache.Client]()
//	if err != nil {
//		return IdentityErrInternal, nil, errors.New("获取缓存客户端err：" + err.Error())
//	}
//
//	// 检查Redis中token有效性
//	redisKey := authConfig.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
//	storedToken, err := cacheClient.Get(context.Background(), redisKey)
//	if err != nil {
//		return IdentityErrTokenExpired, nil, errors.New("令牌已失效")
//	}
//	if storedToken != tokenString {
//		return IdentityErrTokenExpired, nil, errors.New("令牌已失效")
//	}
//
//	// 无感刷新 - 在令牌过期前刷新
//	expirationWindow := 5 * time.Minute // 提前5分钟刷新
//	if time.Until(claims.ExpiresAt.Time) < expirationWindow {
//		newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
//		if err == nil {
//			// 设置新令牌到响应头
//			c.Header("New-Token", newToken)
//
//			// 更新Redis中的令牌
//			err := cacheClient.Set(context.Background(), redisKey, newToken, authConfig.TokenExpireTime)
//			if err != nil {
//				fmt.Println("令牌刷新失败:", err)
//			}
//		}
//	}
//
//	return IdentityOK, claims, nil
//}
