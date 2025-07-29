package user

import (
	"context"
	"fmt"
	"qwqserver/internal/app"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	"qwqserver/internal/vo/vouser"
	"qwqserver/pkg/authx"
	"qwqserver/pkg/util/passsec"
	"strconv"
)

// Login 登录
func (s *Service) Login(ctx context.Context, req *vouser.UserLoginRequest) (any, error) {
	if req.Username == "" && req.Email == "" {
		return nil, fmt.Errorf("必须提供用户名或邮箱")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("必须提供密码")
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	var mUser *model.User
	if req.Username != "" {
		mUser, err = userRepo.FindByUsername(ctx, req.Username)
	} else {
		mUser, err = userRepo.FindByEmail(ctx, req.Email)
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if mUser == nil {
		return nil, fmt.Errorf("用户不存在或邮箱未注册")
	}

	// 模拟校验密码
	//if ok := req.Password == mUser.Password; !ok {
	//	return nil, fmt.Errorf("密码错误，请重新输入")
	//}

	// 校验密码
	ok, err := passsec.Check(req.Password, mUser.Password)
	if err != nil {
		return nil, fmt.Errorf("服务器错误，密码校验失败: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("密码错误，请重新输入")
	}

	if mUser.Status == 0 {
		return nil, fmt.Errorf("用户账号被禁用")
	}

	uidStr := strconv.Itoa(int(mUser.ID))
	// 生成 Token
	//accessToken, err := auth.GenerateToken(uidStr, req.Platform, req.DeviceID) // old
	//if err != nil {
	//	return nil, fmt.Errorf("生成 Access Token 失败: %w", err)
	//}
	//refreshToken, err := auth.GenerateRefreshToken(uidStr, req.Platform, req.DeviceID)
	//if err != nil {
	//	return nil, fmt.Errorf("生成 Refresh Token 失败: %w", err)
	//}
	//
	////构造返回结构体
	//tokenData := &vouser.UserTokenData{
	//	ID:           mUser.ID,
	//	Platform:     req.Platform,
	//	DeviceID:     req.DeviceID,
	//	AccessToken:  accessToken,
	//	RefreshToken: refreshToken,
	//}
	//
	//cacheClient, err := app.Get[*cache.Client]()
	//if err != nil {
	//	return nil, errors.New("获取缓存客户端err：" + err.Error())
	//}
	//
	//authConfig, _ := config.NewAuth()
	////ctx := context.Background()
	//
	//// 缓存 AccessToken
	//tokenKey := auth.BuildTokenKey(authConfig.CacheTokenPrefix, uidStr, req.Platform, req.DeviceID)
	//if err := cacheClient.SetJSON(ctx, tokenKey, tokenData, authConfig.TokenExpireTime); err != nil {
	//	return nil, fmt.Errorf("缓存 AccessToken 失败: %w", err)
	//}
	//
	//// 缓存 RefreshToken
	//refreshKey := auth.BuildTokenKey(authConfig.CacheRefreshPrefix, uidStr, req.Platform, req.DeviceID)
	//if err := cacheClient.SetJSON(ctx, refreshKey, refreshToken, authConfig.RefreshTokenExpire); err != nil {
	//	return nil, fmt.Errorf("缓存 RefreshToken 失败: %w", err)
	//}
	//
	//// 记录设备
	//deviceKey := authConfig.CacheUserDevicePrefix + uidStr
	//if err := cacheClient.HSet(ctx, deviceKey, req.Platform, req.DeviceID); err != nil {
	//	return nil, fmt.Errorf("缓存用户设备失败: %w", err)
	//}

	loginInput := &authx.LoginInput{
		UserID:     "user-" + uidStr,
		Platform:   "qwq",
		DeviceSign: "web/pc",
	}

	ax, _ := app.Get[*authx.AuthX]()

	// 执行登录
	tokenPair, err := ax.Login(ctx, loginInput)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}
