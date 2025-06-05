package handler

import (
	"github.com/gin-gonic/gin"
	"qwqserver/internal/model"
	"qwqserver/internal/service/auth"
)

// Register 注册
func Register(c *gin.Context) *model.Result {
	return auth.Register(c)
}

// Login 登录
func Login(c *gin.Context) *model.Result {
	return auth.Login(c)
}

// Logout 登出
func Logout(c *gin.Context) *model.Result {
	return auth.Logout(c)
}
