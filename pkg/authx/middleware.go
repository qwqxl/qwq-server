package authx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ( 
    // Gin 上下文中的 Claims 键
	ContextKeyClaims = "authx_claims"
)

// Middleware 创建一个 Gin 中间件，用于验证 JWT
func (ax *AuthX) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := ax.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(ContextKeyClaims, claims)
		c.Next()
	}
}

// GetClaims 从 Gin 上下文中获取 Claims
func GetClaims(c *gin.Context) *CustomClaims {
	if claims, exists := c.Get(ContextKeyClaims); exists {
		if customClaims, ok := claims.(*CustomClaims); ok {
			return customClaims
		}
	}
	return nil
}