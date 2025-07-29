package server

import (
	"context"
	"fmt"
	"qwqserver/internal/config"
	"qwqserver/pkg/web"
)

type WebServer struct {
	name string
	conf *config.Config
}

func NewWebServer(name string, conf *config.Config) *WebServer {
	return &WebServer{
		name: name,
		conf: conf,
	}
}

func (s *WebServer) Run(ctx context.Context) error {
	if !config.IsInitialized() {
		return fmt.Errorf("请先初始化配置")
	}
	web.Run(s.conf.Web)
	return nil
}

func (s *WebServer) Name() string {
	return fmt.Sprintf("Gin服务(%s)", s.name)
}
