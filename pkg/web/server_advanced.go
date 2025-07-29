package web

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"qwqserver/internal/config"
	"strings"
)

// 创建负载均衡器
func (s *Server) createLoadBalancer(vhost *config.VirtualHost) http.Handler {
	targets := make([]*url.URL, 0, len(vhost.Backends))
	statuses := make([]*backendStatus, 0, len(vhost.Backends))

	for _, backend := range vhost.Backends {
		target, err := url.Parse(backend)
		if err != nil {
			log.Fatalf("无效的后端地址 %s: %v", backend, err)
		}

		status := &backendStatus{
			url:   target,
			alive: true,
		}

		targets = append(targets, target)
		statuses = append(statuses, status)

		// 启动健康检查
		//go s.healthCheck(status)
		go s.healthCheck(status, vhost.HealthPath, vhost.HealthCheck)
	}

	lb := &loadBalancer{
		policy:    vhost.LBPolicy,
		statuses:  statuses,
		transport: s.transport,
	}

	return lb
}

// 创建路由器
func (s *Server) createRouter(vhost *config.VirtualHost) http.Handler {
	router := http.NewServeMux()

	for _, route := range vhost.Routes {
		var handler http.Handler

		if route.Proxy != "" {
			target, _ := url.Parse(route.Proxy)
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.Transport = s.transport
			handler = http.StripPrefix(route.Path, proxy)
		} else if route.RootDir != "" {
			fs := http.FileServer(http.Dir(route.RootDir))
			handler = http.StripPrefix(route.Path, fs)
		}

		cleanPath := strings.TrimSuffix(route.Path, "/")
		router.Handle(cleanPath+"/", handler)
	}

	return router
}

func (s *Server) buildMiddlewareChain(handler http.Handler) http.Handler {
	// 优化后的顺序：
	// 1. 速率限制 (最先执行，拒绝恶意请求)
	// 2. 安全头部 (尽早设置安全相关头部)
	// 3. 日志记录 (记录处理后的请求)
	// 4. 响应压缩 (最后执行，避免压缩已处理的错误响应)

	handler = rateLimitMiddleware(handler)       // 速率限制
	handler = securityHeadersMiddleware(handler) // 安全头部
	handler = loggingMiddleware(handler)         // 日志记录
	handler = gzipMiddleware(handler)            // 响应压缩

	return handler
}

func tlsVersionFromString(ver string) uint16 {
	switch strings.ToUpper(ver) {
	case "TLS13":
		return tls.VersionTLS13
	default:
		return tls.VersionTLS12
	}
}

// 构建增强TLS配置
func (s *Server) buildTLSConfig(cert tls.Certificate) *tls.Config {
	webConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tlsVersionFromString(s.config.TLS.MinVersion),
		NextProtos:   []string{"h2", "http/1.1"},
	}

	// 严格加密套件
	if s.config.TLS.StrictCiphers {
		webConf.CipherSuites = []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		}
	}

	// HSTS设置
	if s.config.TLS.HSTSMaxAge > 0 {
		webConf.PreferServerCipherSuites = true
	}

	return webConf
}
