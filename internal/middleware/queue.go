package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/app"
	"qwqserver/internal/config"
	"qwqserver/pkg/httpcore"
)

//var conf = app.Get[]()

// 最大并发连接数
var maxConcurrent int

// 信号量控制器
var semaphore chan struct{}

// QueueLimitMiddleware 请求中间件：队列限制
func QueueLimitMiddleware() httpcore.HandlerFunc {
	// 设置最大并发连接数
	if app.IsInitialized() {
		if conf, err := app.Get[*config.Config](); err != nil {
			maxConcurrent = 1
		} else {
			maxConcurrent = conf.Listen.QueueLimitMaxConcurrent
		}
		semaphore = make(chan struct{}, maxConcurrent)
	}
	return func(c httpcore.Context) {
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			c.Next()
		default:
			c.Abort()
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "too many concurrent requests",
			})
		}
	}
}

//func QueueLimitMiddleware() httpcore.HandlerFunc {
//	return func(c httpcore.Context) {
//		select {
//		case semaphore <- struct{}{}:
//			// 成功获取信号量，允许继续
//			defer func() { <-semaphore }()
//			c.Next()
//		default:
//			// 达到最大并发限制
//			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
//				"error": "too many concurrent requests",
//			})
//		}
//	}
//}
