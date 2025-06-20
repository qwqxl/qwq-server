package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"qwqserver/internal/auth"
	"qwqserver/internal/handler"
	"qwqserver/internal/middleware"
	"qwqserver/pkg/util"
)

func RouterApiV1() {
	r := New()

	// 限流
	r.Use(middleware.QueueLimitMiddleware())

	// 健康检测
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	mdDir := util.WorkDir("resources/docs")
	mdTemplateDir := util.WorkDir("resources/templates/markdown")

	fmt.Println("mdDir: ", mdDir, "mdTemplateDir: ", mdTemplateDir)

	// 添加静态文件服务 - 关键部分
	r.Static("/docs/img", mdDir)

	// 加载模板
	r.LoadHTMLGlob(filepath.Join(mdTemplateDir, "*.html"))

	// 获取所有 Markdown 文件
	r.GET("/docs", func(c *gin.Context) {
		util.ListMarkdownFiles(mdDir, c)
	})

	// 渲染 Markdown 文件
	r.GET("/docs/:filename", func(c *gin.Context) {
		util.RenderMarkdown(mdDir, c)
	})

	apiV1Group := r.Group("/api/v1")

	// 认证中间件
	apiV1Group.Use(middleware.AuthMiddleware())

	// 认证路由
	authGroup := apiV1Group.Group("/auth")
	{
		// 注册
		authGroup.POST(auth.RegisterPath, func(c *gin.Context) {
			// get user service handler
			handle := handler.NewUserHandler()
			res := handle.Register(c)
			c.JSON(res.Code, res)
		})
		// 登录
		authGroup.POST(auth.LoginPath, func(c *gin.Context) {
			// get user service handler
			handle := handler.NewUserHandler()
			res := handle.Login(c)
			c.JSON(res.Code, res)
		})
		// 登出
		authGroup.POST(auth.LogoutPath, func(c *gin.Context) {
			// get user service handler
			handle := handler.NewUserHandler()
			res := handle.Logout(c)
			c.JSON(res.Code, res)
		})
		// --------- 用户操作 --------- //
		// 删除用户
		authGroup.DELETE(auth.DelIDPath, func(c *gin.Context) {
			// get user service handler
			handle := handler.NewUserHandler()
			res := handle.DelID(c)
			c.JSON(res.Code, res)
		})

	}

	// 文章路由
	postGroup := apiV1Group.Group("/post")
	{
		// 创建文章
		postGroup.POST("/create", func(c *gin.Context) {
			handle := handler.NewPost()
			res := handle.Create(c)
			c.JSON(res.Code, res)
		})
	}

}
