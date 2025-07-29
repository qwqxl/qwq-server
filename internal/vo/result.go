// vo/result.go
package vo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Result 通用 API 响应结构体
type Result struct {
	Code      int         `json:"code"`           // 自定义业务码
	Message   string      `json:"message"`        // 提示信息
	Data      interface{} `json:"data,omitempty"` // 响应数据
	Timestamp int64       `json:"timestamp"`      // 时间戳
}

// Success 返回成功响应
func Success(c *gin.Context, data ...interface{}) {
	res := &Result{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	res.Timestamp = time.Now().UnixNano() / 1e6 // 时间戳
	c.JSON(http.StatusOK, res)
}

// Fail 返回失败响应（带自定义业务码和消息）
func Fail(c *gin.Context, msg string, codes ...int) {
	code := 1
	if len(codes) > 0 {
		code = codes[0]
	}
	res := &Result{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
	res.Timestamp = time.Now().UnixNano() / 1e6
	c.JSON(http.StatusOK, res)
}
