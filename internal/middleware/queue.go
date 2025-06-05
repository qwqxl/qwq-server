package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/config"
)

// 最大并发连接数
var maxConcurrent = config.New().Listen.MaxConcurrent

// 信号量控制器
var semaphore = make(chan struct{}, maxConcurrent)

// 请求中间件：队列限制
func QueueLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case semaphore <- struct{}{}:
			// 成功获取信号量，允许继续
			defer func() { <-semaphore }()
			c.Next()
		default:
			// 达到最大并发限制
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "too many concurrent requests",
			})
		}
	}
}
