package handler

import (
	"github.com/gin-gonic/gin"
	"qwqserver/internal/model"
	"qwqserver/internal/service"
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

// UserDeleteID 删除用户
func UserDeleteID(c *gin.Context) *model.Result {
	req := struct {
		UID uint `json:"uid"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		return &model.Result{
			Code:    400,
			Message: "用户删除参数错误",
		}
	}

	return service.UserDelete(req.UID)
}
