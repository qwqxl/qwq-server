package server

//import (
//	"context"
//	"errors"
//	"fmt"
//	"qwqserver/pkg/httpcore"
//
//	//v1 "qwqserver/internal/api/v1"
//	"qwqserver/internal/app"
//	"qwqserver/internal/config"
//	"qwqserver/pkg/httpadapter"
//	"qwqserver/pkg/logger"
//	"qwqserver/pkg/util"
//	"time"
//)

//type GinService struct {
//	name string
//}
//
//func NewGinService(name string) *GinService {
//	return &GinService{name: name}
//}
//
//func (s *GinService) Run(ctx context.Context) error {
//	// 获取日志
//
//	r := httpcore.Default()
//
//	l, err := app.Get[logger.Logger]()
//	if err != nil {
//		return errors.New("get logger error:" + err.Error())
//	}
//
//	conf, err := app.Get[*config.Config]()
//	if err != nil {
//		return errors.New("get config error:" + err.Error())
//	}
//
//	// 路由
//	//_, err := v1.Add("/api/v1")
//
//	srv := r.HTTPServer
//
//	// 监听停止信号
//	go func() {
//		<-ctx.Done()
//		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		if err := srv.Shutdown(shutdownCtx); err != nil {
//			l.Error("Gin服务器关闭失败: Error: %v", err)
//		}
//	}()
//
//	addr := conf.ListenAddress()
//
//	// 启动Gin服务器
//	l.Info("Gin服务启动中 %s", addr)
//
//	// Banner
//	banner := util.Banner{}
//	banner.Run()
//
//	err = r.Run(addr)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (s *GinService) Name() string {
//	return fmt.Sprintf("Gin服务(%s)", s.name)
//}
