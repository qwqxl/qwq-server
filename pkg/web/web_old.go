package web

//
//import (
//	"crypto/tls"
//	"errors"
//	"fmt"
//	"log"
//	"net"
//	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"os"
//	"os/signal"
//	"qwqserver/internal/config"
//	"strings"
//	"sync"
//	"syscall"
//)
//
//func NewServerOld(cfg *config.Web) *Server {
//
//	return &Server{config: cfg}
//}
//
//func (s *Server) StartOld() error {
//	// 创建主路由器
//	mainHandler := http.NewServeMux()
//
//	// 构建 hostname -> handler 映射表
//	hostHandlers := make(map[string]http.Handler)
//
//	for _, vhost := range s.config.VirtualHosts {
//		var handler http.Handler
//
//		switch {
//		case vhost.Proxy != "":
//			target, err := url.Parse(vhost.Proxy)
//			if err != nil {
//				return fmt.Errorf("invalid proxy target for %s: %w", vhost.Hostname, err)
//			}
//			// 在创建反向代理时复用Transport
//			//handler = httputil.NewSingleHostReverseProxy(target)
//			//handler.Transport = transport
//
//			handler = httputil.NewSingleHostReverseProxy(target)
//			log.Printf("反向代理: %s -> %s", vhost.Hostname, vhost.Proxy)
//
//		case vhost.RootDir != "":
//			fs := http.FileServer(http.Dir(vhost.RootDir))
//			handler = http.StripPrefix("/", fs)
//			log.Printf("静态文件: %s -> %s", vhost.Hostname, vhost.RootDir)
//
//		default:
//			return fmt.Errorf("virtual host %s has no backend configuration", vhost.Hostname)
//		}
//
//		hostHandlers[vhost.Hostname] = handler
//	}
//
//	// 临时 fallback 默认域名
//	mainHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		host := r.Host
//		if strings.Contains(host, ":") {
//			host, _, _ = net.SplitHostPort(host)
//		}
//
//		log.Printf("收到请求 Host: %s", host)
//
//		// 临时 fallback
//		for h, handler := range hostHandlers {
//			if strings.HasPrefix(host, h) {
//				log.Printf("模糊匹配到 %s", h)
//				handler.ServeHTTP(w, r)
//				return
//			}
//		}
//
//		http.Error(w, "Unknown host: "+host, http.StatusNotFound)
//	})
//
//	// 创建 HTTP 服务器
//	s.httpServer = &http.Server{
//		Addr:    s.config.ListenAddr,
//		Handler: mainHandler,
//	}
//
//	// 启动监听
//	var err error
//	if s.config.TLS != nil && s.config.TLS.Enabled {
//		if _, err := os.Stat(s.config.TLS.CertFile); os.IsNotExist(err) {
//			return fmt.Errorf("TLS 证书文件不存在: %s", s.config.TLS.CertFile)
//		}
//		if _, err := os.Stat(s.config.TLS.KeyFile); os.IsNotExist(err) {
//			return fmt.Errorf("TLS 私钥文件不存在: %s", s.config.TLS.KeyFile)
//		}
//		// HTTPS
//		cert, err := tls.LoadX509KeyPair(s.config.TLS.CertFile, s.config.TLS.KeyFile)
//		if err != nil {
//			return fmt.Errorf("failed to load TLS cert: %w", err)
//		}
//
//		s.httpServer.TLSConfig = &tls.Config{
//			Certificates: []tls.Certificate{cert},
//		}
//
//		s.listener, err = tls.Listen("tcp", s.config.ListenAddr, s.httpServer.TLSConfig)
//	} else {
//		// HTTP
//		s.listener, err = net.Listen("tcp", s.config.ListenAddr)
//	}
//
//	if err != nil {
//		return err
//	}
//
//	log.Printf("监听: %s (TLS: %v)", s.config.ListenAddr, s.config.TLS != nil)
//
//	go func() {
//		if err := s.httpServer.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
//			log.Fatalf("服务器错误: %v", err)
//		}
//	}()
//
//	return nil
//}
//
//func (s *Server) StopOld() {
//	if s.httpServer != nil {
//		_ = s.httpServer.Close()
//	}
//}
//
//
//package web
//
//import (
//"log"
//"os"
//"os/signal"
//"qwqserver/internal/config"
//"strings"
//"sync"
//"syscall"
//)
//
//func (s *Server) matchHost(requestHost, configHost string) bool {
//	// 去除端口号
//	reqHost := strings.Split(requestHost, ":")[0]
//	return reqHost == configHost
//}
//
//func Run(webConfigs []*config.Web) {
//
//	// 启动所有服务器
//	var servers []*Server
//	var wg sync.WaitGroup
//
//	for _, cfg := range webConfigs {
//		server := NewServer(cfg)
//		if err := server.Start(); err != nil {
//			log.Fatalf("启动服务器失败: %v", err)
//		}
//		servers = append(servers, server)
//		wg.Add(1)
//	}
//
//	log.Println("所有服务器已启动")
//
//	// 等待中断信号
//	waitForInterrupt()
//
//	// 优雅关闭
//	for _, s := range servers {
//		s.Stop()
//	}
//	wg.Wait()
//	log.Println("所有服务器已关闭")
//}
//
//// 等待中断信号 (Ctrl+C)
//func waitForInterrupt() {
//	sigCh := make(chan os.Signal, 1)
//	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
//	<-sigCh
//}
