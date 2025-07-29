package web

import (
	"log"
	"os"
	"os/signal"
	"qwqserver/internal/config"
	"strings"
	"sync"
	"syscall"
)

func (s *Server) matchHost(requestHost, configHost string) bool {
	// 去除端口号
	reqHost := strings.Split(requestHost, ":")[0]
	return reqHost == configHost
}

func Run(webConfigs []*config.Web) {

	// 启动所有服务器
	var servers []*Server
	var wg sync.WaitGroup

	for _, cfg := range webConfigs {
		server := NewServer(cfg)
		if err := server.Start(); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
		servers = append(servers, server)
		wg.Add(1)
	}

	log.Println("所有服务器已启动")

	// 等待中断信号
	waitForInterrupt()

	// 优雅关闭
	for _, s := range servers {
		s.Stop()
	}
	wg.Wait()
	log.Println("所有服务器已关闭")
}

// 等待中断信号 (Ctrl+C)
func waitForInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
