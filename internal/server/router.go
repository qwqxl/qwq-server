package server

import (
	"github.com/gin-gonic/gin"
	"qwqserver/internal/handler"
	"qwqserver/internal/middleware"
	"qwqserver/internal/service/auth"
	"qwqserver/pkg/util/network/client"
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

	apiV1Group := r.Group("/api/v1")

	apiV1Group.GET("/info", func(c *gin.Context) {
		deviceInfo := client.JsonDeviceHandler(c.Writer, c.Request)
		c.JSON(200, deviceInfo)
	})

	// 认证路由

	userGroup := apiV1Group.Group("/user")
	{
		// 注册
		userGroup.POST(auth.RegisterPath, func(c *gin.Context) {
			res := handler.Register(c)
			c.JSON(res.Code, res)
		})
		// 登录
		userGroup.POST(auth.LoginPath, func(c *gin.Context) {
			res := handler.Login(c)
			c.JSON(res.Code, res)
		})
		// 登出
		userGroup.POST(auth.LogoutPath, func(c *gin.Context) {
			res := handler.Logout(c)
			c.JSON(res.Code, res)
		})
	}

}
