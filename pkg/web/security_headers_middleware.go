package web

import (
	"fmt"
	"net/http"
	"qwqserver/internal/config"
)

// 安全头部中间件
func securityHeadersMiddlewareOld(conf *config.Web, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 基础安全头部
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// CSP可根据需要配置
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// HSTS (仅HTTPS)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security",
				fmt.Sprintf("max-age=%d; includeSubDomains", conf.TLS.HSTSMaxAge))
		}

		next.ServeHTTP(w, r)
	})
}

// securityHeadersMiddleware 安全头部中间件
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置安全相关HTTP头部
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// 内容安全策略 (根据应用需求调整)
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")

		// 仅当使用HTTPS时设置HSTS
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		// 防止缓存敏感信息
		if r.Method == http.MethodPost {
			w.Header().Set("Cache-Control", "no-store")
		}

		next.ServeHTTP(w, r)
	})
}
