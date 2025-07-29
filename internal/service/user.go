package service

import (
	"context"
	"fmt"
	"qwqserver/internal/app"
	"qwqserver/internal/auth"
	"qwqserver/internal/config"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	userService "qwqserver/internal/service/user"
	"qwqserver/internal/vo/vouser"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/util/passsec"
	"strconv"
)

type User struct {
	userService.Service
}

func (s *User) LoginOld(req *vouser.UserLoginRequest) (data any, msg error, err error) {
	token := vouser.UserTokenData{}
	cfg, _ := config.New()
	authConfig := cfg.Auth
	mUser := &model.User{}

	if req.Username == "" && req.Email == "" {
		// 用户登录用户名或邮箱未提供（必须填写其一）
		err = fmt.Errorf("用户登录用户名或邮箱未提供（必须填写其一）")
		return
	}

	if req.Password == "" {
		// 密码未提供
		err = fmt.Errorf("用户登录密码未提供（必须填写密码）")
		return
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		// 数据库连接失败
		err = fmt.Errorf("数据库连接失败: %v", err)
		return
	}

	if req.Username != "" {
		if mUser, err = userRepo.FindByUsername(context.Background(), req.Username); mUser == nil {
			// 用户不存在
			err = fmt.Errorf("用户不存在")
			return
		}
	} else {
		if mUser, _ = userRepo.FindByEmail(context.Background(), req.Email); mUser == nil {
			// 邮箱未注册
			err = fmt.Errorf("邮箱未注册")
			return
		}
	}

	// 校验密码（应使用安全哈希，略）
	var ok bool
	if ok, err = passsec.Check(req.Password, mUser.Password); err != nil {
		// 密码校验错误（哈希异常等）
		err = fmt.Errorf("服务器错误，密码校验错误")
		return
	}

	if !ok {
		// 密码错误（用户输入不对）
		err = fmt.Errorf("密码错误，请重新输入")
		return
	}
	// 清除密码
	req.Password = ""

	// 禁用状态
	if mUser.Status == 0 {
		// 用户账号被禁用
		err = fmt.Errorf("用户账号被禁用")
		return
	}

	uid := strconv.Itoa(int(mUser.ID))
	token.ID = mUser.ID

	// 生成Token和RefreshToken
	newToken, err := auth.GenerateToken(uid, req.Platform, req.DeviceID)
	if err != nil {
		// Token 生成失败（JWT 过程异常）
		err = fmt.Errorf("服务器错误，Token 创建失败")
		return
	}
	// token
	token.AccessToken = newToken

	refreshToken, err := auth.GenerateRefreshToken(uid, req.Platform, req.DeviceID)
	if err != nil {
		// Refresh Token 生成失败
		err = fmt.Errorf("服务器错误，Refresh Token 创建失败")
		return
	}
	// refresh token
	token.RefreshToken = refreshToken

	cacheClient, err := app.Get[*cache.Client]()
	if err != nil {
		// 获取Redis连接失败
		err = fmt.Errorf("服务器错误，获取Redis连接失败：" + err.Error())
		return
	}

	ctx := context.Background()

	// 存储Token到Redis
	tokenKey := authConfig.CacheTokenPrefix + uid + ":" + req.Platform + ":" + req.DeviceID
	err = cacheClient.SetJSON(ctx, tokenKey, token, authConfig.TokenExpireTime)
	if err != nil {
		// Access Token 存储失败（例如 Redis 写入失败）
		err = fmt.Errorf("服务器错误，Access Token 缓存失败")
		return
	}

	// 存储RefreshToken
	refreshKey := authConfig.CacheRefreshPrefix + uid + ":" + req.Platform + ":" + req.DeviceID
	err = cacheClient.SetJSON(ctx, refreshKey, refreshToken, authConfig.RefreshTokenExpire)
	if err != nil {
		// Refresh Token 存储失败
		err = fmt.Errorf("服务器错误，Refresh Token 缓存失败")
		return
	}

	// 记录用户设备关系
	userDeviceKey := authConfig.CacheUserDevicePrefix + uid
	err = cacheClient.HSet(ctx, userDeviceKey, req.Platform, req.DeviceID)
	if err != nil {
		// 记录用户设备关系失败
		err = fmt.Errorf("服务器错误，缓存用户设备关系失败")
		return
	}

	// 登录成功
	data = token
	err = nil
	return
}

// Logout 登出
//func (s *UserService) Logout(req *model.UserLogoutRequest) (any, error) {
//	if req.ID == 0 || req.Platform == "" || req.DeviceID == "" {
//		return nil, fmt.Errorf("登出请求参数无效")
//	}
//
//	conf := config.New()
//	authConfig := conf.Auth
//
//	userRepo, err := repository.NewUserRepository()
//	if err != nil {
//		return nil, fmt.Errorf("数据库连接失败: %w", err)
//	}
//
//	// 判断用户是否存在
//	mUser, err := userRepo.FindByID(context.Background(), req.ID)
//	if err != nil || mUser == nil {
//		return nil, fmt.Errorf("用户不存在")
//	}
//
//	cachePool := app.New().CachePool
//	cacheClient, err := cachePool.GetClient()
//	if err != nil {
//		return nil, fmt.Errorf("获取缓存连接失败: %w", err)
//	}
//
//	ctx := context.Background()
//
//	// 删除 Token
//	tokenKey := fmt.Sprintf("%s%d:%s:%s", authConfig.CacheTokenPrefix, req.ID, req.Platform, req.DeviceID)
//	if err := cacheClient.Del(ctx, tokenKey); err != nil {
//		return nil, fmt.Errorf("删除 AccessToken 出错: %w", err)
//	}
//
//	// 删除 RefreshToken
//	refreshKey := fmt.Sprintf("%s%d:%s:%s", authConfig.CacheRefreshPrefix, req.ID, req.Platform, req.DeviceID)
//	if err := cacheClient.Del(ctx, refreshKey); err != nil {
//		return nil, fmt.Errorf("删除 RefreshToken 出错: %w", err)
//	}
//
//	// 删除设备记录
//	userDeviceKey := fmt.Sprintf("%s%d", authConfig.CacheUserDevicePrefix, req.ID)
//	if err := cacheClient.HDel(ctx, userDeviceKey, req.Platform); err != nil {
//		return nil, fmt.Errorf("删除设备记录出错: %w", err)
//	}
//
//	// 登出成功
//	return model.UserLogoutResponse{ID: req.ID}, nil
//}

// 修改密码
//func (s *AuthService) ChangePassword(userID, oldPassword, newPassword string) error {
//	user, err := model.GetUserByID(userID)
//	if err != nil {
//		return err
//	}
//
//	// 验证旧密码
//	if user.Password != oldPassword {
//		return errors.New("旧密码错误")
//	}
//
//	// 更新密码
//	user.Password = newPassword
//	return model.UpdateUser(user)
//}
