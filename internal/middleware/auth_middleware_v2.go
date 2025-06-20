package middleware

//
//import (
//	"net/http"
//	"qwqserver/internal/auth"
//	"qwqserver/internal/common"
//	"qwqserver/pkg/cache"
//
//	"strings"
//	"time"
//
//	"github.com/gin-gonic/gin"
//)
//
//type HandlerFunc func(c *gin.Context) (auth.CodeType, *common.HTTPResult)
//
//// 定义需要排除认证的路径
//var excludePaths = map[string]ExcludeRouter{
//	"/api/v1/auth/register": {
//		IsValid: true,
//		HandlerFunc: func(c *gin.Context) (auth.CodeType, *common.HTTPResult) {
//			authIdentity, res := JWTAuth(c)
//			// 认证成功
//			if authIdentity != auth.IdentityOK {
//				// 注册路由 不存在认证 允许放行
//				c.Next()
//			} else {
//				// 注册路由 已存在认证 拒绝访问
//				c.Abort()
//				res.Code = http.StatusOK
//				res.Msg = "已登录，请不要重复注册"
//			}
//			return auth.IdentitySkipped, res
//		},
//	},
//	"/api/v1/auth/login": {
//		IsValid: true,
//		HandlerFunc: func(c *gin.Context) (auth.CodeType, *common.HTTPResult) {
//			authIdentity, res := JWTAuth(c)
//			if authIdentity != auth.IdentityOK {
//				// 登录路由 不存在认证 允许放行
//				c.Next()
//			} else {
//				// 登录路由 已存在认证 拒绝访问
//				c.Abort()
//				res.Code = http.StatusOK
//				res.Msg = "已登录，请不要重复登录"
//			}
//			return auth.IdentitySkipped, res
//		},
//	},
//}
//
//type ExcludeRouter struct {
//	IsValid     bool
//	HandlerFunc HandlerFunc
//}
//
//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var authCode auth.CodeType
//		var res *common.HTTPResult
//
//		// 1. 检查当前路径是否需要特殊处理
//		if router, ok := excludePaths[c.Request.URL.Path]; ok && router.IsValid {
//			authCode, res = router.HandlerFunc(c)
//		} else {
//			// 2. 标准认证流程
//			authCode, res = JWTAuth(c)
//		}
//
//		// 3. 统一处理认证结果
//		switch authCode {
//		case auth.IdentityOK:
//			// 认证成功或跳过认证，继续后续处理
//			c.Next()
//		case auth.IdentitySkipped:
//			// 跳过认证，继续后续处理
//		default:
//			if res == nil {
//				res = &common.HTTPResult{
//					Code: http.StatusUnauthorized,
//					Msg:  "认证失败",
//				}
//			}
//			c.AbortWithStatusJSON(res.Code, res)
//		}
//	}
//}
//
//func JWTAuth(c *gin.Context) (auth.CodeType, *common.HTTPResult) {
//	// 获取Authorization header
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		return auth.IdentityErrNoToken, &common.HTTPResult{
//			Code: http.StatusUnauthorized,
//			Msg:  "未提供认证令牌",
//		}
//	}
//
//	// 检查Bearer格式
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return auth.IdentityErrTokenFormat, &common.HTTPResult{
//			Code: http.StatusUnauthorized,
//			Msg:  "令牌格式错误",
//		}
//	}
//
//	tokenString := parts[1]
//	claims, err := auth.ParseToken(tokenString)
//	if err != nil {
//		return auth.IdentityErrInvalidToken, &common.HTTPResult{
//			Code: http.StatusUnauthorized,
//			Msg:  "无效令牌",
//		}
//	}
//
//	// 检查Redis中是否存在该token
//	redisKey := common.RedisTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
//	storedToken, err := cache.Get(redisKey)
//	if err != nil || storedToken != tokenString {
//		return auth.IdentityErrTokenExpired, &common.HTTPResult{
//			Code: http.StatusUnauthorized,
//			Msg:  "令牌已失效",
//		}
//	}
//
//	// 无感刷新：检查是否需要刷新token
//	if time.Until(claims.ExpiresAt.Time) < common.TokenRefreshInterval {
//		newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
//		if err == nil {
//			c.Header("New-Token", newToken)
//			cache.Set(redisKey, newToken, common.TokenExpireTime)
//		}
//	}
//
//	// 设置上下文信息
//	c.Set(auth.IdentityStatusKey, auth.IdentityOK)
//	c.Set("user_id", claims.UserID)
//	c.Set("platform", claims.Platform)
//	c.Set("device_id", claims.DeviceID)
//
//	return auth.IdentityOK, &common.HTTPResult{
//		Code: http.StatusOK,
//		Msg:  "认证成功",
//		Data: claims,
//	}
//}
