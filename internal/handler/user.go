package handler

import (
	"context"
	"qwqserver/internal/config"
	"qwqserver/internal/service"
	userService "qwqserver/internal/service/user"
	"qwqserver/internal/vo/vouser"
	"qwqserver/pkg/httpcore"
	"qwqserver/pkg/util/network/client"
)

// UserHandler 用户处理
type UserHandler struct {
	service.User
}

// Register 注册
func (handle *UserHandler) Register(c httpcore.Context) {
	req := &vouser.UserRegisterRequest{}

	// 绑定参数
	if err := c.ShouldBindJSON(req); err != nil || req.Username == "" || req.Password == "" {
		httpcore.Fail(c, "用户注册参数错误")
		return
	}
	// password4
	// Password123.
	// pwd4@iqwq.com

	// 调用服务层逻辑
	ctx := context.Background()
	data, err := handle.Service.Register(ctx, req)
	if err != nil {
		httpcore.Fail(c, err.Error())
		return
	}

	httpcore.Success(c, data)
}

// Login 用户登录
func (handle *UserHandler) Login(c httpcore.Context) {
	// 获取服务和配置
	serv := handle.Service
	authConfig, err := config.NewAuth()
	if err != nil {
		httpcore.Fail(c, "获取配置失败")
		return
	}

	// 获取设备 ID
	deviceInfo := client.JsonDeviceHandler(c.Response(), c.Request())
	deviceID := deviceInfo.DeviceType

	// 绑定请求体
	var req *vouser.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpcore.Fail(c, "请求体解析失败")
		return
	}

	// 设置平台签名和设备 ID（从服务端生成）
	req.Platform = authConfig.PlatformSign
	req.DeviceID = deviceID

	ctx := c.Request().Context()

	// 调用登录逻辑
	data, err := serv.Login(ctx, req)
	if err != nil {
		httpcore.Fail(c, err.Error())
		return
	}

	// 登录成功
	httpcore.Success(c, data)
}

// Logout 登出
func (handle *UserHandler) Logout(c httpcore.Context) {
	req := &vouser.UserLogoutRequest{}
	if err := c.ShouldBindJSON(req); err != nil || req.ID == 0 {
		httpcore.Fail(c, "用户登出参数错误")
		return
	}

	req.DeviceID = client.JsonDeviceHandler(c.Response(), c.Request()).DeviceType
	ca, err := config.NewAuth()
	if err != nil {
		httpcore.Fail(c, "获取配置失败")
		return
	}

	req.Platform = ca.PlatformSign

	ctx := c.Request().Context()

	data, err := handle.Service.Logout(ctx, req)
	if err != nil {
		httpcore.Fail(c, err.Error())
		return
	}

	httpcore.Success(c, data)
}

// Del 根据ID删除用户
func (handle *UserHandler) Del(c httpcore.Context) {
	req := userService.DeleteRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		httpcore.Fail(c, "用户删除参数错误")
		return
	}

	if req.ID == 0 && req.Username == "" && req.Email == "" {
		httpcore.Fail(c, "请提供用户ID、用户名或邮箱中的至少一个进行删除")
		return
	}

	// 自动补充 platform 和 device_id
	ca, err := config.NewAuth()
	if err != nil {
		httpcore.Fail(c, "获取配置失败")
		return
	}
	req.Platform = ca.PlatformSign
	req.DeviceID = client.JsonDeviceHandler(c.Response(), c.Request()).DeviceType

	ctx := c.Request().Context()
	data, err := handle.Service.Delete(ctx, req)
	if err != nil {
		httpcore.Fail(c, err.Error())
		return
	}

	httpcore.Success(c, data)
}

// Logout 登出
//func (handle *UserHandler) Logout(c *gin.Context) {
//	var req vouser.UserLogoutRequest
//	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
//		vo.Fail(c, "用户登出参数错误")
//		return
//	}
//
//	req.DeviceID = client.JsonDeviceHandler(c.Writer, c.Request).DeviceType
//	ca, err := config.NewAuth()
//	if err != nil {
//		vo.Fail(c, "获取配置失败")
//		return
//	}
//
//	req.Platform = ca.PlatformSign
//
//	ctx := context.Background()
//
//	data, err := handle.Service.Logout(ctx, &req)
//	if err != nil {
//		vo.Fail(c, err.Error())
//		return
//	}
//
//	vo.Success(c, data)
//}
//
//// Del 根据ID删除用户
//func (handle *UserHandler) Del(c *gin.Context) {
//	var req vouser.DelRequest
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		vo.Fail(c, "用户删除参数错误")
//		return
//	}
//
//	if req.ID == 0 && req.Username == "" && req.Email == "" {
//		vo.Fail(c, "请提供用户ID、用户名或邮箱中的至少一个进行删除")
//		return
//	}
//
//	// 自动补充 platform 和 device_id
//	ca, err := config.NewAuth()
//	if err != nil {
//		vo.Fail(c, "获取配置失败")
//		return
//	}
//	req.Platform = ca.PlatformSign
//	req.DeviceID = client.JsonDeviceHandler(c.Writer, c.Request).DeviceType
//
//	data, err := handle.Service.Del(c.Request.Context(), req)
//	if err != nil {
//		vo.Fail(c, err.Error())
//		return
//	}
//
//	vo.Success(c, data)
//}
