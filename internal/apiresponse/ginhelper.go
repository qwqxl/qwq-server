package apiresponse

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success 成功响应（HTTP 200）
func Success(c *gin.Context, code BaseCode, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: GetMessageByCode(code),
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的响应（HTTP 200）
func SuccessWithMessage(c *gin.Context, code BaseCode, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code.String(),
		Message: msg,
		Data:    data,
	})
}

// Fail 失败响应（HTTP 400）
func Fail(c *gin.Context, code BaseCode) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    code.String(),
		Message: GetMessageByCode(code),
	})
}

// Unauthorized 未认证（HTTP 401）
func Unauthorized(c *gin.Context, code BaseCode) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code.String(),
		Message: GetMessageByCode(code),
	})
}

// Forbidden 禁止访问（HTTP 403）
func Forbidden(c *gin.Context, code BaseCode) {
	c.JSON(http.StatusForbidden, Response{
		Code:    code.String(),
		Message: GetMessageByCode(code),
	})
}

// WithMessage 自定义消息响应（HTTP 状态码自定义）
func WithMessage(c *gin.Context, code BaseCode, msg string, status int) {
	c.JSON(status, Response{
		Code:    code.String(),
		Message: msg,
	})
}

// WithDataMessage 带自定义消息和数据响应
func WithDataMessage(c *gin.Context, code BaseCode, msg string, data interface{}, status int) {
	c.JSON(status, Response{
		Code:    code.String(),
		Message: msg,
		Data:    data,
	})
}
