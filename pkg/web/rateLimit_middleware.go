package web

import (
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
)

// RateLimiter 维护每个客户端的速率限制器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

// NewRateLimiter 创建速率限制器实例
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// GetLimiter 获取或创建客户端的速率限制器
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[ip]; exists {
		return limiter
	}

	// 创建新的限制器: 100 请求/秒，突发容量 200
	limiter := rate.NewLimiter(100, 200)
	rl.limiters[ip] = limiter
	return limiter
}

// 全局速率限制器实例
var globalRateLimiter = NewRateLimiter()

// rateLimitMiddleware 速率限制中间件
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取客户端IP
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		// 获取该IP的速率限制器
		limiter := globalRateLimiter.GetLimiter(ip)

		// 检查是否超过速率限制
		if !limiter.Allow() {
			http.Error(w, "请求过多，请稍后再试", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
