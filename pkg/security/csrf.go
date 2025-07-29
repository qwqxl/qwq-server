package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/config"
)

// CSRFMiddleware 提供CSRF保护
func CSRFMiddleware(cfg config.CSRF) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 验证配置
	if len(cfg.Secret) < 32 {
		panic("CSRF secret must be at least 32 bytes long")
	}
	if cfg.TokenLength < 16 {
		cfg.TokenLength = 32
	}

	return func(c *gin.Context) {
		// 安全方法白名单（不需要CSRF验证）
		safeMethods := []string{"GET", "HEAD", "OPTIONS"}
		currentMethod := c.Request.Method
		isSafeMethod := false
		for _, m := range safeMethods {
			if m == currentMethod {
				isSafeMethod = true
				break
			}
		}

		// 获取或创建CSRF令牌
		var csrfToken string
		cookie, err := c.Cookie(cfg.CookieName)
		if err == nil && cookie != "" {
			// 使用现有令牌
			csrfToken = cookie
		} else {
			// 生成新令牌
			csrfToken = generateCSRFToken(cfg.TokenLength)
			c.SetCookie(
				cfg.CookieName,
				csrfToken,
				cfg.CookieMaxAge,
				"/",
				"",
				cfg.CookieSecure,
				cfg.CookieHTTPOnly,
			)
		}

		// 将令牌添加到上下文
		c.Set("csrf_token", csrfToken)
		c.Header("X-CSRF-Token", csrfToken)

		// 对于安全方法，直接通过
		if isSafeMethod {
			c.Next()
			return
		}

		// 验证CSRF令牌
		valid := false

		// 1. 检查请求头
		headerToken := c.GetHeader("X-CSRF-Token")
		if headerToken != "" && headerToken == csrfToken {
			valid = true
		}

		// 2. 检查表单字段
		if !valid {
			formToken := c.PostForm("csrf_token")
			if formToken != "" && formToken == csrfToken {
				valid = true
			}
		}

		// 3. 检查查询参数
		if !valid {
			queryToken := c.Query("csrf_token")
			if queryToken != "" && queryToken == csrfToken {
				valid = true
			}
		}

		if !valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "CSRF token validation failed",
				"message": "Invalid or missing CSRF token",
			})
			return
		}

		c.Next()
	}
}

// generateCSRFToken 生成安全的CSRF令牌
func generateCSRFToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate CSRF token: %v", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}

// VerifyCSRFToken 验证CSRF令牌（用于中间件外的手动验证）
func VerifyCSRFToken(c *gin.Context, token string) bool {
	cfg, exists := c.Get("csrf_config")
	if !exists {
		return false
	}

	csrfConfig := cfg.(config.CSRF)
	cookie, err := c.Cookie(csrfConfig.CookieName)
	if err != nil {
		return false
	}

	return token == cookie
}
