package middleware

import (
	"fmt"
	"net/http"
	"qwqserver/internal/config"
	"qwqserver/pkg/httpcore"
	"regexp"
	"strings"
)

func Cors(cfg config.Cors) httpcore.HandlerFunc {
	// 编译正则
	var originRegexps []*regexp.Regexp
	for _, pattern := range cfg.AllowOriginPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			originRegexps = append(originRegexps, re)
		}
	}

	return func(c httpcore.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}

		req := c.Request()
		//writer := c.Response()
		origin := req.Header.Get("Origin")
		requestMethod := req.Method

		// 判断是否允许 Origin
		allowed := false
		if origin != "" {
			for _, o := range cfg.AllowOrigins {
				if o == "*" && !cfg.AllowCredentials {
					allowed = true
					break
				}
				if o == origin {
					allowed = true
					break
				}
			}
			if !allowed {
				for _, re := range originRegexps {
					if re.MatchString(origin) {
						allowed = true
						break
					}
				}
			}
			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
				if cfg.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			}
		}

		// OPTIONS 预检请求处理
		if requestMethod == http.MethodOptions {
			if origin == "" {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			methods := strings.Join(cfg.AllowMethods, ", ")
			if methods == "" {
				methods = "GET, POST, PUT, DELETE, PATCH, OPTIONS"
			}
			c.Header("Access-Control-Allow-Methods", methods)

			allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
			if allowHeaders == "" {
				allowHeaders = "Origin, Content-Type, Authorization, Accept"
			}
			c.Header("Access-Control-Allow-Headers", allowHeaders)

			if len(cfg.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
			}
			if cfg.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
			}

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// gin old
//func Cors(cfg config.Cors) gin.HandlerFunc {
//	// 预编译正则表达式
//	var originRegexps []*regexp.Regexp
//	for _, pattern := range cfg.AllowOriginPatterns {
//		re, err := regexp.Compile(pattern)
//		if err == nil {
//			originRegexps = append(originRegexps, re)
//		}
//	}
//
//	return func(c *gin.Context) {
//		if !cfg.Enabled {
//			c.Next()
//			return
//		}
//
//		origin := c.Request.Header.Get("Origin")
//		requestMethod := c.Request.Method
//
//		// 处理来源域
//		if origin != "" {
//			allowed := false
//
//			// 检查精确匹配
//			for _, o := range cfg.AllowOrigins {
//				if o == "*" {
//					if cfg.AllowCredentials {
//						// 凭证模式下不能使用通配符
//						continue
//					}
//					allowed = true
//					break
//				}
//				if o == origin {
//					allowed = true
//					break
//				}
//			}
//
//			// 检查正则匹配
//			if !allowed {
//				for _, re := range originRegexps {
//					if re.MatchString(origin) {
//						allowed = true
//						break
//					}
//				}
//			}
//
//			if allowed {
//				c.Header("Access-Control-Allow-Origin", origin)
//				if cfg.AllowCredentials {
//					c.Header("Access-Control-Allow-Credentials", "true")
//				}
//			}
//		}
//
//		// 处理预检请求
//		if requestMethod == "OPTIONS" {
//			if origin == "" {
//				c.AbortWithStatus(http.StatusForbidden)
//				return
//			}
//
//			// 设置允许的方法
//			methods := strings.Join(cfg.AllowMethods, ", ")
//			if methods == "" {
//				methods = "GET, POST, PUT, DELETE, PATCH, OPTIONS"
//			}
//			c.Header("Access-Control-Allow-Methods", methods)
//
//			// 设置允许的请求头
//			allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
//			if allowHeaders == "" {
//				allowHeaders = "Origin, Content-Type, Authorization, Accept"
//			}
//			c.Header("Access-Control-Allow-Headers", allowHeaders)
//
//			// 设置暴露的响应头
//			if len(cfg.ExposeHeaders) > 0 {
//				c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
//			}
//
//			// 设置预检请求缓存时间
//			if cfg.MaxAge > 0 {
//				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
//			}
//
//			c.AbortWithStatus(http.StatusNoContent)
//			return
//		}
//
//		c.Next()
//	}
//}
