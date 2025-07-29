package web

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"qwqserver/internal/config"
	"strings"
	"time"
)

type Server struct {
	config     *config.Web
	httpServer *http.Server
	listener   net.Listener
	transport  *http.Transport
}

func NewServer(cfg *config.Web) *Server {

	// 全局共享Transport，支持连接复用
	var transport = &http.Transport{
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Server{
		config:    cfg,
		transport: transport, // 保存共享transport
	}
}

func (s *Server) Start() error {
	// 创建主路由器
	mainHandler := http.NewServeMux()

	// 构建 hostname -> handler 映射表
	hostHandlers := make(map[string]http.Handler)

	for _, vhost := range s.config.VirtualHosts {
		var handler http.Handler

		// 1. 负载均衡处理
		if len(vhost.Backends) > 0 {
			lb := s.createLoadBalancer(vhost)
			handler = lb
			log.Printf("负载均衡: %s -> %d 后端 (%s)",
				vhost.Hostname, len(vhost.Backends), vhost.LBPolicy)

			// 2. 路由规则处理
		} else if len(vhost.Routes) > 0 {
			router := s.createRouter(vhost)
			handler = router
			log.Printf("路由配置: %s (%d 条规则)", vhost.Hostname, len(vhost.Routes))

			// 3. 旧版单个代理/静态文件处理
		} else if vhost.Proxy != "" {
			target, err := url.Parse(vhost.Proxy)
			if err != nil {
				return fmt.Errorf("invalid proxy target for %s: %w", vhost.Hostname, err)
			}
			// 在创建反向代理时复用Transport
			//handler = httputil.NewSingleHostReverseProxy(target)
			//handler.Transport = transport

			handler = httputil.NewSingleHostReverseProxy(target)
			log.Printf("反向代理: %s -> %s", vhost.Hostname, vhost.Proxy)
		} else if vhost.RootDir != "" {
			fs := http.FileServer(http.Dir(vhost.RootDir))
			handler = http.StripPrefix("/", fs)
			log.Printf("静态文件: %s -> %s", vhost.Hostname, vhost.RootDir)
		}

		hostHandlers[vhost.Hostname] = handler
	}

	// 临时 fallback 默认域名
	mainHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if strings.Contains(host, ":") {
			host, _, _ = net.SplitHostPort(host)
		}

		log.Printf("收到请求 Host: %s", host)

		// 临时 fallback
		for h, handler := range hostHandlers {
			if host == h {
				log.Printf("模糊匹配到 %s", h)
				handler.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Unknown host: "+host, http.StatusNotFound)
	})

	// ===== 新增中间件链 =====
	handlerChain := s.buildMiddlewareChain(mainHandler)

	// 创建HTTP服务器时使用中间件链
	s.httpServer = &http.Server{
		Addr:    s.config.ListenAddr,
		Handler: handlerChain, // 使用中间件包装的处理器

		// 安全超时设置
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// ===== TLS高级配置 =====
	// 构建增强TLS配置 监听
	var err error
	if s.config.TLS != nil && s.config.TLS.Enabled {
		if _, err = os.Stat(s.config.TLS.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS 证书文件不存在: %s", s.config.TLS.CertFile)
		}
		if _, err = os.Stat(s.config.TLS.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS 私钥文件不存在: %s", s.config.TLS.KeyFile)
		}
		// HTTPS
		cert, err := tls.LoadX509KeyPair(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS cert: %w", err)
		}

		s.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		s.listener, err = tls.Listen("tcp", s.config.ListenAddr, s.httpServer.TLSConfig)

		// 创建增强TLS配置
		s.httpServer.TLSConfig = s.buildTLSConfig(cert)
	} else {
		// HTTP
		s.listener, err = net.Listen("tcp", s.config.ListenAddr)
	}
	if err != nil {
		return fmt.Errorf("无法启动服务器: %w", err)
	}

	log.Printf("监听: %s (TLS: %v)", s.config.ListenAddr, s.config.TLS != nil)

	go func() {
		if err := s.httpServer.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() {
	if s.httpServer != nil {
		// 优雅关闭
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("优雅关闭失败: %v", err)
			s.httpServer.Close() // 强制关闭
		}

		// 关闭连接池
		if s.transport != nil {
			s.transport.CloseIdleConnections()
		}
	}
}
