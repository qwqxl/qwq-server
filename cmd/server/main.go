package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"qwqserver/internal/app"
	"qwqserver/internal/config"
	"qwqserver/internal/model"
	"qwqserver/internal/server"
	"qwqserver/pkg/httpcore"
	"qwqserver/pkg/httpcore/httpconnpool"
	"qwqserver/pkg/qqbot"
	"qwqserver/pkg/qwqlog"
	"qwqserver/pkg/svcmgr"
	"syscall"
	"time"
)

func main() {
	// 捕获全局panic
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println("main global panic: ", rec)
		}
	}()

	// init config
	cfg, err := config.New()
	if err != nil {
		panic("Config 配置文件初始化错误")
	}

	// 配置全局日志
	qwqlog.ConfigureGlobal(
		qwqlog.WithBaseDir("./logs"),
		qwqlog.WithLevel(cfg.Level),
		qwqlog.WithDefaultGroup("main"),
	)
	// 获取默认日志实例
	logger := qwqlog.Default()

	// 添加日志输出器
	logger.AddConsoleWriter()
	logger.AddFileWriter("app/app.log", 100, 10, 30)
	logger.AddJSONWriter("app/app.json")

	logger.Debug("Config is singleton initialized: %v", config.IsInitialized())

	// 初始化App
	logger.Debug("app init")
	app.Initialized(&app.Application{
		Config: cfg,
		Logger: logger,
	})

	defer app.Initialized().Close()

	// 创建带自动重启的服务管理器
	mgr := svcmgr.New(svcmgr.Config{
		StopTimeout:    15 * time.Second,
		RestartDelay:   2 * time.Second,
		MaxRestarts:    5,
		RestartOnError: true,
	})

	{
		// 启动web服务器
		wbs := server.NewWebServer("WEB", cfg)
		mgr.Add(wbs)

		// 启动http服务器
		heConnPool := httpconnpool.NewServerConnPool(100, 1000, 90*time.Second) // 创建连接池实例
		eng := httpcore.DefaultWithLogger(logger)
		eng.SetConnPool(heConnPool)
		rs := server.NewGinServer("GIN", eng, logger)
		mgr.Add(rs)
	}

	// 启动服务
	if err = mgr.Start(); err != nil {
		panic(err)
	}

	// 错误处理
	go func() {
		for mErr := range mgr.ErrChan() {
			logger.Error("Service error: %v", mErr)
		}
	}()

	// 优雅停止
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Debug("Shutting down services...")
	mgr.Stop()

	logger.Info("All services stopped")
}

type HelloModule struct{}

func (m *HelloModule) RegisterRoutes(e *httpcore.Engine) {

	//e.POST("/bot", func(c httpcore.Context) {
	//	req := c.Request()
	//
	//	// 1. 读取 Body 内容
	//	body, err := io.ReadAll(req.Body)
	//	if err != nil {
	//		c.JSON(400, "读取请求失败")
	//		return
	//	}
	//	defer req.Body.Close()
	//
	//	// 2. 解析 JSON
	//	var jsonReq map[string]interface{}
	//	if err := json.Unmarshal(body, &jsonReq); err != nil {
	//		c.JSON(400, "JSON 解析失败")
	//		return
	//	}
	//
	//	fmt.Printf("请求体内容: %+v\n", jsonReq)
	//
	//	c.JSON(200, httpcore.Success("Hello, World!"))
	//})

	e.POST("/bot", func(c httpcore.Context) {
		var msg model.BotMessageRequest
		decoder := json.NewDecoder(c.Request().Body)
		decoder.UseNumber() // 保证 ID 不被转成 float64
		if err := decoder.Decode(&msg); err != nil {
			c.JSON(400, "解析失败: "+err.Error())
			return
		}

		groupID, _ := msg.GroupID.Int64()
		userID := msg.UserID.String()
		text := msg.RawMessage

		fmt.Printf("Group: %d, User: %s, Text: %s\n", groupID, userID, text)

		// ✅ 关键逻辑：收到“测试”则回复“测试成功”
		if text == "测试" {
			go qqbot.SendMsg("http://localhost:3001", groupID, "测试成功")
		}

		httpcore.Success(c, "ok.")
	})

}
