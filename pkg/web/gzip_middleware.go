package web

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipResponseWriter 包装gzip.Writer
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write 实现gzip压缩写入
func (g gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

// gzipMiddleware 响应压缩中间件
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查客户端是否支持gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// 创建gzip writer
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer gz.Close()

		// 设置响应头部
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		// 使用gzip writer包装响应
		//gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		//next.ServeHTTP(gzw, r)

		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)

		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/") && !strings.Contains(contentType, "json") {
			w.Header().Del("Content-Encoding") // 删除压缩标志
		}
	})
}
