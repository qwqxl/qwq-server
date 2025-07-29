package web

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

// logResponseWriter 包装http.ResponseWriter以捕获状态码
type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 捕获状态码
func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// AccessLogEntry 结构化日志条目
type AccessLogEntry struct {
	Timestamp  string  `json:"timestamp"`
	ClientIP   string  `json:"client_ip"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Query      string  `json:"query,omitempty"`
	StatusCode int     `json:"status_code"`
	Latency    float64 `json:"latency_ms"` // 毫秒
	UserAgent  string  `json:"user_agent"`
	Referer    string  `json:"referer,omitempty"`
}

// loggingMiddleware 日志记录中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建自定义ResponseWriter以捕获状态码
		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// 处理请求
		next.ServeHTTP(lrw, r)

		// 计算请求处理时间
		latency := time.Since(start).Seconds() * 1000

		// 获取客户端真实IP (考虑代理)
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			clientIP = realIP
		}

		// 解析主机和端口
		host, _, err := net.SplitHostPort(clientIP)
		if err == nil {
			clientIP = host
		}

		// 创建结构化日志条目
		logEntry := AccessLogEntry{
			Timestamp:  start.Format(time.RFC3339),
			ClientIP:   clientIP,
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			StatusCode: lrw.statusCode,
			Latency:    latency,
			UserAgent:  r.UserAgent(),
			Referer:    r.Referer(),
		}

		// 转换为JSON格式
		logData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("日志序列化错误: %v", err)
			return
		}

		// 输出结构化日志
		log.Println(string(logData))
	})
}
