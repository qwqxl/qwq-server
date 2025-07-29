package server

import (
	"context"
	"fmt"
	v1 "qwqserver/internal/api/v1"
	"qwqserver/pkg/httpcore"
	"qwqserver/pkg/logger"
	"qwqserver/pkg/util"
	"time"
)

type GinServer struct {
	name   string
	Engine *httpcore.Engine
	Logger logger.Logger
}

func NewGinServer(name string, engine *httpcore.Engine, l logger.Logger) *GinServer {
	engine.RegisterModules(&v1.ApiModule{})

	return &GinServer{
		name:   name,
		Engine: engine,
		Logger: l,
	}
}

func (s *GinServer) Run(ctx context.Context) error {
	// 获取日志

	r := s.Engine

	// 启动Gin服务器
	s.Logger.Info("Gin服务启动中")

	// Banner
	banner := util.Banner{}
	banner.Run()

	// 启动服务（优雅关停）
	if err := r.RunWithGracefulShutdown(":5000", 5*time.Second); err != nil {
		s.Logger.Error("服务关闭失败: %v", err)
		return err
	}
	return nil
}

func (s *GinServer) Name() string {
	return fmt.Sprintf("Gin服务(%s)", s.name)
}
