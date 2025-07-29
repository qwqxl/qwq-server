package httpcore

import (
	"net/http"
	"time"
)

// Result 通用 API 响应结构体
type Result struct {
	Code      int         `json:"code"`           // 自定义业务码
	Message   string      `json:"message"`        // 提示信息
	Data      interface{} `json:"data,omitempty"` // 响应数据
	Timestamp int64       `json:"timestamp"`      // 时间戳
}

func (h *Result) NewTimestamp() int64 {
	return time.Now().UnixNano() / 1e6 // 时间戳
}

func NewResult(code int, msg string, data interface{}) *Result {
	res := &Result{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	res.Timestamp = res.NewTimestamp()
	return res
}

// Success 返回成功响应
func Success(c Context, data ...interface{}) {
	res := NewResult(0, "success", data)
	c.JSON(http.StatusOK, res)
}

// Fail 返回失败响应（带自定义业务码和消息）
func Fail(c Context, msg string, codes ...int) {
	code := 1
	if len(codes) > 0 {
		code = codes[0]
	}
	res := NewResult(code, msg, nil)
	c.JSON(http.StatusOK, res)
}
