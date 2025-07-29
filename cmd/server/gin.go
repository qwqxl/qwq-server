package main

//type GinService struct {
//	name string
//}
//
//func (s *GinService) Run(ctx context.Context) error {
//	// 获取日志
//
//	r := server.NewHTTPServer()
//
//	// 路由
//	_, err := v1.Add("/api/v1")
//
//	srv := r.Server
//
//	// 监听停止信号
//	go func() {
//		<-ctx.Done()
//		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		if err := srv.Shutdown(shutdownCtx); err != nil {
//			r.Logger.Error("Gin服务器关闭失败: Error: %v", err)
//		}
//	}()
//
//	// 启动Gin服务器
//	//logger.Info("Gin服务启动中 %s", address)
//
//	// Banner
//	banner := util.Banner{}
//	banner.Run()
//
//	err = r.Run()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (s *GinService) Name() string {
//	return fmt.Sprintf("Gin服务(%s)", s.name)
//}
