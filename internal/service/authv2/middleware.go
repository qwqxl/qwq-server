package authv2

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/pkg/util/network/client"
	"strings"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(HeaderAuthorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// 未提供或格式错误的Token
			c.Set(Identity, false)
			//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供或格式错误的Token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		ipaddr := client.GetClientIP(c.Request)
		claims, valid := parseToken(tokenString)
		if !valid {
			if handleTokenRefresh(c, tokenString, ipaddr) {
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token无效"})
			return
		}

		if !isDeviceValid(claims.UserID, claims.DeviceID, claims.IPAddress) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "已在其他设备登录"})
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyDeviceID, claims.DeviceID)
		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件（示例权限判断）
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}
		userID, _ := c.Get(ContextKeyUserID)
		if userID != "user_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			return
		}
		c.Next()
	}
}
