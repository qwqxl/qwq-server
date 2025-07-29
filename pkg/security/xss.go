package security

import (
	"fmt"
	"qwqserver/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// XSSMiddleware 提供XSS和点击劫持保护
func XSSMiddleware(cfg config.XSS) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// 设置X-XSS-Protection头
		if cfg.XSSProtection != "" {
			c.Header("X-XSS-Protection", cfg.XSSProtection)
		} else {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		// 设置X-Content-Type-Options头
		if cfg.ContentTypeNosniff != "" {
			c.Header("X-Content-Type-Options", cfg.ContentTypeNosniff)
		} else {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// 设置X-Frame-Options头
		if cfg.XFrameOptions != "" {
			c.Header("X-Frame-Options", cfg.XFrameOptions)
		} else {
			c.Header("X-Frame-Options", "DENY")
		}

		// 设置Content-Security-Policy头
		if cfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}

		// 设置HSTS头（仅HTTPS）
		if c.Request.TLS != nil && cfg.HSTSMaxAge > 0 {
			hstsValue := fmt.Sprintf("max-age=%d", cfg.HSTSMaxAge)
			if cfg.HSTSIncludeSubdomains {
				hstsValue += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// 清理用户输入防止XSS
		c.Request.ParseForm()
		for key, values := range c.Request.Form {
			for i, value := range values {
				c.Request.Form[key][i] = sanitizeInput(value)
			}
		}

		// 清理JSON请求体
		if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
			// 在后续处理中清理JSON数据
			c.Set("xss_sanitize", true)
		}

		c.Next()
	}
}

// sanitizeInput 清理用户输入防止XSS攻击
func sanitizeInput(input string) string {
	// 替换危险字符
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")
	input = strings.ReplaceAll(input, "&", "&amp;")

	// 移除危险属性
	dangerousPatterns := []string{
		"javascript:", "vbscript:", "expression(", "onload", "onerror",
		"onclick", "onmouseover", "onfocus", "onblur", "eval(",
	}

	for _, pattern := range dangerousPatterns {
		input = strings.ReplaceAll(input, pattern, "")
	}

	return input
}

// SanitizeJSON 清理JSON数据防止XSS
func SanitizeJSON(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			sanitized[key] = sanitizeInput(v)
		case map[string]interface{}:
			sanitized[key] = SanitizeJSON(v)
		case []interface{}:
			sanitized[key] = sanitizeArray(v)
		default:
			sanitized[key] = v
		}
	}

	return sanitized
}

func sanitizeArray(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, item := range arr {
		switch v := item.(type) {
		case string:
			result[i] = sanitizeInput(v)
		case map[string]interface{}:
			result[i] = SanitizeJSON(v)
		case []interface{}:
			result[i] = sanitizeArray(v)
		default:
			result[i] = v
		}
	}
	return result
}
