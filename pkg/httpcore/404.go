package httpcore

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SetNotFoundHandler 设置全局 404 处理函数
func (e *Engine) SetNotFoundHandler(handler HandlerFunc) {
	e.notFoundHandler = handler
	e.engine.NoRoute(func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: e.logger})
	})
	e.logger.Debug("Global 404 handler set")
}

// SetGroupNotFoundHandler 为路由组设置自定义 404 处理函数
func (g *routerGroupWrapper) SetNotFoundHandler(handler HandlerFunc) {
	// 为路由组创建子引擎
	groupEngine := gin.New()

	// 应用父路由组的中间件
	for _, h := range g.group.Handlers {
		groupEngine.Use(h)
	}

	// 设置 404 处理
	groupEngine.NoRoute(func(c *gin.Context) {
		handler(&contextWrapper{ctx: c, logIns: g.log})
	})

	// 替换路由组的处理引擎
	g.group.Handlers = nil
	g.group.Handlers = append(g.group.Handlers, groupEngine.Handlers...)
	g.log.Debug("Group 404 handler set for path: %s", g.group.BasePath())
}

// defaultNotFoundHandler 默认 404 处理函数
func defaultNotFoundHandler(c Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":      http.StatusNotFound,
		"message":   "Resource not found",
		"path":      c.Request().URL.Path,
		"method":    c.Request().Method,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// notFoundMiddleware 自定义 404 处理中间件
func (e *Engine) notFoundMiddleware() HandlerFunc {
	return func(c Context) {
		// 继续处理请求链
		c.Next()

		// 检查请求是否未被处理且状态为 404
		if !c.IsHandled() && c.Response().(gin.ResponseWriter).Status() == http.StatusNotFound {
			// 调用自定义 404 处理函数
			if e.notFoundHandler != nil {
				e.notFoundHandler(c)
			} else {
				defaultNotFoundHandler(c)
			}
		}
	}
}
