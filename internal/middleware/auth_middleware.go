package middleware

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"qwqserver/internal/app"
//	"qwqserver/internal/auth"
//	"qwqserver/internal/config"
//	"qwqserver/internal/global"
//	"strings"
//	"time"
//)
//
//type HandlerFunc func(c *gin.Context) (auth.CodeType, *global.HTTPResult)
//
//type ExcludeRouter struct {
//	IsValid     bool
//	HandlerFunc
//}
//
//var excludePaths = map[string]ExcludeRouter{
//	"/api/v1/user/register": {
//		IsValid: true,
//		HandlerFunc: func(c *gin.Context) (auth.CodeType, *global.HTTPResult) {
//			if isLoggedIn(c) {
//				c.AbortWithStatusJSON(http.StatusOK, global.HTTPResult{
//					Message: "已登录，请不要重复注册",
//				})
//				return auth.IdentitySkipped, nil
//			}
//			c.Next()
//			return auth.IdentitySkipped, nil
//		},
//	},
//	"/api/v1/user/login": {
//		IsValid: true,
//		HandlerFunc: func(c *gin.Context) (auth.CodeType, *global.HTTPResult) {
//			if isLoggedIn(c) {
//				c.AbortWithStatusJSON(http.StatusOK, global.HTTPResult{
//					Message: "已登录，请不要重复登录",
//				})
//				return auth.IdentitySkipped, nil
//			}
//			c.Next()
//			return auth.IdentitySkipped, nil
//		},
//	},
//}
//
//func Auth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		conf, err := config.New()
//		if err != nil {
//			// 忽略错误
//			return
//		}
//
//		// 测试环境
//		if conf.Server.Mode == "debug" {
//			c.Next()
//			return
//		}
//
//		fmt.Println("Server mode: ", conf.Server.Mode, "进来了")
//
//		// 1. 检查排除路径
//		if router, ok := excludePaths[c.Request.URL.Path]; ok && router.IsValid {
//			router.HandlerFunc(c)
//			return // 直接返回，不进入后续流程
//		}
//
//		// 2. 标准认证流程
//		authCode, res, err := JWTAuth(c)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, global.HTTPResult{
//				Code:    http.StatusUnauthorized,
//				Message: "认证失败: " + err.Error(),
//			})
//			return
//		}
//
//		// 3. 处理认证结果
//		switch authCode {
//		case auth.IdentityOK:
//			c.Next()
//		case auth.IdentitySkipped: // 通常不会发生
//			c.Next()
//		default:
//			if res == nil {
//				res = &global.HTTPResult{
//					Code:    http.StatusUnauthorized,
//					Message: "认证失败",
//				}
//			}
//			c.AbortWithStatusJSON(res.Code, res)
//		}
//	}
//}
//
//// 轻量级登录检查 (不刷新token/不设置上下文)
//func isLoggedIn(c *gin.Context) bool {
//	conf, err := config.New()
//	if err != nil {
//		// 忽略错误
//		return false
//	}
//	authConfig := conf.Auth
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
//	if err != nil {
//		return false
//	}
//
//	// 检查Redis中的token有效性
//	redisKey := authConfig.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
//
//	cachePool := app.New().CachePool
//	cacheClient, err := cachePool.GetClient()
//	if err != nil {
//		return false
//	}
//
//	storedToken, err := cacheClient.Get(context.Background(), redisKey)
//	return err == nil && storedToken == tokenString
//}
//
//func JWTAuth(c *gin.Context) (auth.CodeType, any, error) {
//	conf, err := config.New()
//	if err != nil {
//		// 忽略错误
//		return 500, nil, errors.New("文件配置错误")
//	}
//	authConfig := conf.Auth
//	// 获取Authorization header
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		return auth.IdentityErrNoToken, nil, errors.New("未提供认证令牌")
//	}
//
//	// 检查Bearer格式
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return auth.IdentityErrTokenFormat, nil, errors.New("令牌格式错误")
//	}
//
//	tokenString := parts[1]
//	claims, err := auth.ParseToken(tokenString)
//	if err != nil {
//		return auth.IdentityErrInvalidToken, nil, errors.New("无效令牌")
//	}
//
//	cachePool := app.New().CachePool
//	cacheClient, err := cachePool.GetClient()
//	if err != nil {
//		return auth.IdentityErrInvalidToken, nil, errors.New("获取缓存连接失败: " + err.Error())
//	}
//
//	// 检查Redis中token有效性
//	redisKey := authConfig.CacheTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
//	storedToken, err := cacheClient.Get(context.Background(), redisKey)
//	if err != nil || storedToken != tokenString {
//		return auth.IdentityErrTokenExpired, nil, errors.New("令牌已失效")
//	}
//
//	// 无感刷新
//	if time.Until(claims.ExpiresAt.Time) < authConfig.TokenRefreshInterval {
//		newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
//		if err == nil {
//			c.Header("New-Token", newToken)
//			err := cacheClient.Set(context.Background(), redisKey, newToken, authConfig.TokenExpireTime)
//			if err != nil {
//				return 0, nil, errors.New("令牌刷新失败")
//			}
//		}
//	}
//
//	// 设置上下文信息
//	c.Set(auth.IdentityStatusKey, auth.IdentityOK)
//	c.Set("user_id", claims.UserID)
//	c.Set("platform", claims.Platform)
//	c.Set("device_id", claims.DeviceID)
//
//	return auth.IdentityOK, claims, nil
//}
