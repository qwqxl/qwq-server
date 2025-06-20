package authv2

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/util"
	"qwqserver/pkg/util/network/client"
	"qwqserver/pkg/util/passsec"
	"strings"
)

// 注册处理逻辑
func Register(c *gin.Context) *model.Result {
	res := &model.Result{}
	req := RegisterRequest{}

	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" || req.Nickname == "" || req.Email == "" {
		res.Code = http.StatusBadRequest
		res.Message = "注册用户参数错误"
		return res
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "获取数据库连接失败: " + err.Error()
		return res
	}

	// 检查用户名是否已存在
	if ok, _ := userRepo.ExistUsername(context.Background(), req.Username); ok {
		res.Code = http.StatusConflict
		res.Message = "您输入的用户名已经存在"
		return res
	}

	// 检查邮箱是否已存在
	if ok, _ := userRepo.ExistEmail(context.Background(), req.Email); ok {
		res.Code = http.StatusConflict
		res.Message = "您输入的邮箱已经存在"
		return res
	}

	// 密码哈希
	pwd, err := passsec.Hash(req.Password)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "密码哈希失败: " + err.Error()
		return res
	}

	// 创建新用户
	newUser := model.User{
		Username:     req.Username,
		Password:     pwd,
		PasswordSalt: pwd,
		PasswordHash: pwd,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Status:       1,
	}
	if err = userRepo.Create(context.Background(), &newUser); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "创建用户失败 Error: " + err.Error()
		return res
	}

	res.Code = http.StatusOK
	res.Data = gin.H{"user_id": newUser.ID}
	res.Message = "注册成功！"
	return res
}

// 登录处理函数
func Login(c *gin.Context) *model.Result {
	res := &model.Result{}
	req := LoginRequest{}
	userInfo := &model.User{}

	// 参数校验
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" || req.Password == "" {
		res.Code = http.StatusBadRequest
		res.Message = "invalid request"
		return res
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "获取数据库连接失败: " + err.Error()
		return res
	}

	// 验证输入的是用户名还是邮箱
	if util.IsEmail(req.Name) {
		if ok, _ := userRepo.ExistEmail(context.Background(), req.Name); !ok {
			res.Code = http.StatusUnauthorized
			res.Message = "Email不存在"
			return res
		}
		userInfo, err = userRepo.FindByEmail(context.Background(), req.Name)
		if err != nil {
			res.Code = http.StatusInternalServerError
			res.Message = "获取用户信息失败 Error: " + err.Error()
			return res
		}
	} else {
		if ok, _ := userRepo.ExistUsername(context.Background(), req.Name); !ok {
			res.Code = http.StatusUnauthorized
			res.Message = "用户名不存在"
			return res
		}
		userInfo, err = userRepo.FindByUsername(context.Background(), req.Name)
		if err != nil {
			res.Code = http.StatusInternalServerError
			res.Message = "获取用户信息失败 Error: " + err.Error()
			return res
		}
	}

	// 校验密码（应使用安全哈希，略）
	var ok bool
	if ok, err = passsec.Check(req.Password, userInfo.Password); err != nil {
		res.Code = http.StatusUnauthorized
		res.Message = "密码校验出错"
		return res
	}

	if !ok {
		res.Code = http.StatusUnauthorized
		res.Message = "密码错误"
		return res
	}

	// 禁用状态
	if userInfo.Status == 0 {
		res.Code = http.StatusForbidden
		res.Message = "用户已被封禁"
		return res
	}

	userID := fmt.Sprintf("%d", userInfo.ID)
	deviceID := client.JsonDeviceHandler(c.Writer, c.Request).DeviceType
	ipaddr := client.GetClientIP(c.Request)

	// 清除旧设备 session
	if err := TerminateOtherSessions(userID, deviceID); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "failed to terminate other sessions"
		return res
	}

	// 生成 token
	accessToken, refreshToken, err := GenerateTokenPair(userID, deviceID, ipaddr)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = "failed to generate tokens"
		return res
	}

	// 存储 session 信息到 Redis
	if err = saveTokenInfo(userID, deviceID, refreshToken, ipaddr, req.DeviceID); err != nil {
		cache.Delete(context.Background(), UserSessionCachePrefixToString(userID, deviceID, ipaddr))
		res.Code = http.StatusInternalServerError
		res.Message = "failed to save token info"
		return res
	}

	res.Code = http.StatusOK
	res.Message = "登录成功"
	// 成功响应
	res.Data = gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return res
}

// 登出处理函数
func Logout(c *gin.Context) *model.Result {
	userID, _ := c.Get(ContextKeyUserID)
	deviceID, _ := c.Get(ContextKeyDeviceID)
	ipaddr := client.GetClientIP(c.Request)
	//
	ctx := context.Background()
	key := UserSessionCachePrefixToString(userID.(string), deviceID.(string), ipaddr)
	fmt.Println(ipaddr, "ipaddr")
	if err := cache.Delete(ctx, key); err != nil {
		return &model.Result{
			Code:    http.StatusInternalServerError,
			Message: "登出失败",
		}
	}

	return &model.Result{
		Code:    http.StatusOK,
		Message: "成功登出",
	}
}

// 用户信息查询处理函数
func UserInfo(c *gin.Context) *model.Result {
	userID, _ := c.Get(ContextKeyUserID)
	deviceID, _ := c.Get(ContextKeyDeviceID)

	return &model.Result{
		Code:    http.StatusOK,
		Message: "ok",
		Data: gin.H{
			"user_id":   userID,
			"device_id": deviceID,
			"username":  strings.TrimPrefix(userID.(string), "user_"),
			"role":      "user",
		},
	}
}
