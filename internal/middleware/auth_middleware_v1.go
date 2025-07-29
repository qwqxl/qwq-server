package middleware

//
//import (
//	"errors"
//	"github.com/golang-jwt/jwt/v5"
//	"net/http"
//	"qwqserver/internal/auth"
//	"strings"
//
//	"github.com/gin-gonic/gin"
//)
//
//func Auth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 1. 从Header获取Token
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
//			return
//		}
//
//		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
//		claims, err := auth.ParseToken(tokenString)
//		if err != nil {
//			// 2. Token过期时尝试刷新
//			if errors.Is(err, jwt.ErrTokenExpired) {
//				handleTokenRefresh(c, claims)
//				return
//			}
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//			return
//		}
//
//		// 3. 验证是否为最新Token
//		if isValid, _ := validateTokenFreshness(claims.UserID, claims.DeviceID, tokenString); !isValid {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
//			return
//		}
//
//		// 4. 设置用户信息到上下文
//		c.Set("userID", claims.UserID)
//		c.Set("deviceID", claims.DeviceID)
//		c.Next()
//	}
//}
//
//// 刷新Token
//func handleTokenRefresh(c *gin.Context, claims *auth.Claims) {
//	refreshToken := c.GetHeader("X-Refresh-Token")
//	if refreshToken == "" {
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
//		return
//	}
//
//	// 验证Refresh Token
//	refreshClaims, err := auth.ParseToken(refreshToken)
//	if err != nil || refreshClaims.UserID != claims.UserID || refreshClaims.DeviceID != claims.DeviceID {
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
//		return
//	}
//
//	// 检查Redis中的Token是否匹配
//	// StoreAccess, storedRefresh := auth.GetTokens(claims.UserID, claims.DeviceID)
//	_, storedRefresh, err := auth.GetTokens(claims.UserID, claims.DeviceID)
//	if err != nil || storedRefresh != refreshToken {
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token mismatch"})
//		return
//	}
//
//	// 生成新Token对
//	newAccess, newRefresh, err := auth.GenerateTokens(claims.UserID, claims.DeviceID)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
//		return
//	}
//
//	// 更新Redis
//	if err := auth.StoreTokens(claims.UserID, claims.DeviceID, newAccess, newRefresh); err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token storage failed"})
//		return
//	}
//
//	// 返回新Token给客户端
//	c.Header("New-Access-Token", newAccess)
//	c.Header("New-Refresh-Token", newRefresh)
//	c.Set("userID", claims.UserID)
//	c.Set("deviceID", claims.DeviceID)
//	c.Next()
//}
//
//func validateTokenFreshness(userID, deviceID, accessToken string) (bool, error) {
//	storedAccess, _, err := auth.GetTokens(userID, deviceID)
//	if err != nil {
//		return false, err
//	}
//	return storedAccess == accessToken, nil
//}
