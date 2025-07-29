package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/app"
	"qwqserver/internal/auth"
	"qwqserver/internal/config"
	"qwqserver/internal/handler"
	"qwqserver/internal/middleware"
	"qwqserver/pkg/httpcore"
)

type ApiModule struct{}

func (m *ApiModule) RegisterRoutes(e *httpcore.Engine) {
	conf, err := app.Get[*config.Config]()
	if err != nil {
		panic("register routes err: " + err.Error())
		return
	}

	// 设置全局自定义 404 处理
	e.SetNotFoundHandler(func(c httpcore.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom 404 - The requested resource was not found",
			"path":  c.Request().URL.Path,
			"suggestions": []string{
				"/api/v1/users",
				"/api/v1/products",
				"/docs",
			},
			"documentation": "https://api.iqwq.com/docs",
		})
	})

	// static
	e.AddStatic("/", "resources/public")

	r := e.Group("/api/v1")

	// group
	//pluginGroup := r.Group("/plugin") // plugin group
	postGroup := r.Group("/post")   // post group
	userGroup := r.Group("/user")   // auth group
	debugGroup := r.Group("/debug") // utils group

	// middleware
	authMiddlewareHandle := middleware.Auth()                                   // 认证中间件
	r.Use(middleware.QueueLimitMiddleware(), middleware.Cors(conf.Server.Cors)) // 限流, cors跨域配置
	//apiV1Group.Group("/security").Use(security.CSRFMiddleware(conf.Server.CSRF)) // csrf
	userGroup.Use(authMiddlewareHandle) // 认证中间件
	//swaggerGroup.Use(authMiddlewareHandle) // 认证中间件

	// 用户/Auth路由
	{
		userHandler := handler.UserHandler{}
		userGroup.POST(auth.RegisterPath, userHandler.Register) // 用户注册
		userGroup.POST(auth.LoginPath, userHandler.Login)       // 用户登录
		userGroup.POST(auth.LogoutPath, userHandler.Logout)     // 用户登出
		// --------- other operate --------- //
		userGroup.DELETE(auth.DelIDPath, userHandler.Del) // 删除用户
	}

	// 文章路由
	{
		postHandler := handler.PostHandler{}
		postGroup.POST("/create", postHandler.Create) // 创建文章
	}

	// debug
	{
		debugGroup.GET("/connpool", func(c httpcore.Context) {
			//engConnPool, err := e.GetConnPool()
			//if err != nil {
			//	httpcore.Fail(c, "")
			//}
			//stats := engConnPool.Stats()
			//httpcore.Success(gin.H{
			//	"idle":   stats.IdleCount,
			//	"active": stats.ActiveCount,
			//	"total":  stats.TotalCount,
			//})
		})
	}

}

//func Add(relativePaths ...string) (*server.HTTPServer, error) {
//	r := server.NewHTTPServer()
//	conf, err := app.Get[*config.Config]()
//	if err != nil {
//		return nil, errors.New("api v1 get config error: " + err.Error())
//	}
//
//	relativePath := strings.Join(relativePaths, "")
//
//	// 健康检测
//	//r.GET("/ping", Ping)
//
//	// group
//	apiV1Group := r.Group(relativePath)          // api group
//	pluginGroup := apiV1Group.Group("/plugin")   // plugin group
//	swaggerGroup := apiV1Group.Group("/swagger") // swagger api doc group
//	postGroup := apiV1Group.Group("/post")       // post group
//	authGroup := apiV1Group.Group("/user")       // auth group
//	captchaGroup := apiV1Group.Group("/captcha")
//
//	// middleware
//	authMiddlewareHandle := middleware.Auth()                                    // 认证中间件
//	r.Use(middleware.QueueLimitMiddleware(), middleware.Cors(conf.Server.Cors))  // 限流, cors跨域配置
//	apiV1Group.Group("/security").Use(security.CSRFMiddleware(conf.Server.CSRF)) // csrf
//	authGroup.Use(authMiddlewareHandle)                                          // 认证中间件
//	swaggerGroup.Use(authMiddlewareHandle)                                       // 认证中间件
//
//	{
//		pluginGroup.POST("/register", handler.RegisterPlugin)     // 注册插件
//		pluginGroup.POST("/:name/execute", handler.ExecutePlugin) // 执行插件
//		pluginGroup.POST("/leave", handler.LeavePlugin)           // 离开插件
//	}
//
//	// other routers
//	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger 路由
//	//handler.Docs("/docs", s, authMiddlewareHandle)                          // Docs 路由
//
//	// 用户/Auth路由
//	{
//		userHandler := handler.UserHandler{}
//		authGroup.POST(auth.RegisterPath, userHandler.Register) // 用户注册
//		authGroup.POST(auth.LoginPath, userHandler.Login)       // 用户登录
//		authGroup.POST(auth.LogoutPath, userHandler.Logout)     // 用户登出
//		// --------- other operate --------- //
//		authGroup.DELETE(auth.DelIDPath, userHandler.Del) // 删除用户
//	}
//
//	// captcha 验证码路由
//	{
//		//captchaGroup.GET("/image", service.CaptchaImage)
//
//		captchaGroup.GET("/view", func(c *gin.Context) {
//			c.Header("Content-Type", "text/html")
//			c.String(http.StatusOK, `
//		<html>
//			<body>
//				<h3>图形验证码预览</h3>
//				<img src="image">
//			</body>
//		</html>
//	`)
//		})
//
//	}
//
//	// 文章路由
//	{
//		// 创建文章
//		postGroup.POST("/create", func(c *gin.Context) {
//			handle := handler.NewPost()
//			res := handle.Create(c)
//			c.JSON(res.Code, res)
//		})
//	}
//
//	return r, nil
//}
