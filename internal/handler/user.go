package handler

import (
	"github.com/gin-gonic/gin"
	"qwqserver/internal/common"
	"qwqserver/internal/service"
	"qwqserver/pkg/util/network/client"
)

// UserHandlerInterface 用户处理接口
type UserHandlerInterface interface {
	Register(c *gin.Context) *common.HTTPResult
	Login(c *gin.Context) *common.HTTPResult
	Logout(c *gin.Context) *common.HTTPResult
	DelID(c *gin.Context) *common.HTTPResult
}

// UserHandler 用户处理
type UserHandler struct {
	Service *service.AuthService
	HandleBaseImpl
}

// NewUserHandler 创建用户处理
func NewUserHandler() UserHandlerInterface {
	return &UserHandler{
		Service: &service.AuthService{},
	}
}

// Register 注册
func (handle *UserHandler) Register(c *gin.Context) *common.HTTPResult {
	serv := handle.Service
	if err := c.ShouldBindJSON(serv); err != nil || serv.Username == "" || serv.Password == "" {
		return &common.HTTPResult{
			Code: 400,
			Msg:  "用户注册参数错误",
		}
	}
	return serv.Register()
}

// Login 登录
func (handle *UserHandler) Login(c *gin.Context) *common.HTTPResult {
	serv := handle.Service
	if err := c.ShouldBindJSON(serv); err != nil || serv.Password == "" {
		return &common.HTTPResult{
			Code: 400,
			Msg:  "用户登录参数错误",
		}
	}
	if serv.Username == "" && serv.Email == "" {
		return &common.HTTPResult{
			Code: 400,
			Msg:  "请使用邮箱或者用户名进行登录",
		}
	}
	deviceID := client.JsonDeviceHandler(c.Writer, c.Request).DeviceType
	return serv.Login("qwq", deviceID)
}

// Logout 登出
func (handle *UserHandler) Logout(c *gin.Context) *common.HTTPResult {
	serv := handle.Service
	if err := c.ShouldBindJSON(serv); err != nil || serv.ID == 0 {
		return &common.HTTPResult{
			Code: 400,
			Msg:  "用户登出参数错误",
		}
	}

	deviceID := client.JsonDeviceHandler(c.Writer, c.Request).DeviceType
	return serv.Logout(serv.ID, common.PlatformSign(), deviceID)
}

// DelID 根据ID删除用户
func (handle *UserHandler) DelID(c *gin.Context) *common.HTTPResult {
	serv := handle.Service
	if err := c.ShouldBindJSON(serv); err != nil {
		return &common.HTTPResult{
			Code: 400,
			Msg:  "用户删除参数错误",
		}
	}
	deviceID := client.JsonDeviceHandler(c.Writer, c.Request).DeviceType
	return serv.Del(serv.ID, common.PlatformSign(), deviceID)
}
