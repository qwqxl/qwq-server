package user

import (
	"context"
	"fmt"
	"qwqserver/internal/app"
	"qwqserver/internal/repository"
	"qwqserver/internal/vo/vouser"
	"qwqserver/pkg/authx"
	"strconv"
)

// Logout 登出
func (s *Service) Logout(ctx context.Context, req *vouser.UserLogoutRequest) (any, error) {
	if req.ID == 0 || req.Platform == "" || req.DeviceID == "" {
		return nil, fmt.Errorf("登出请求参数无效")
	}

	//conf, _ := config.New()
	//authConfig := conf.Auth

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 判断用户是否存在
	mUser, err := userRepo.FindByID(ctx, req.ID)
	if err != nil || mUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	//cacheClient, err := app.Get[*cache.Client]()
	//if err != nil {
	//	return nil, fmt.Errorf("获取缓存客户端连接失败: %w", err)
	//}

	//ctx := context.Background()

	// 删除 Token
	//tokenKey := fmt.Sprintf("%s%d:%s:%s", authConfig.CacheTokenPrefix, req.ID, req.Platform, req.DeviceID)
	//if err := cacheClient.Del(ctx, tokenKey); err != nil {
	//	return nil, fmt.Errorf("删除 AccessToken 出错: %w", err)
	//}
	//
	//// 删除 RefreshToken
	//refreshKey := fmt.Sprintf("%s%d:%s:%s", authConfig.CacheRefreshPrefix, req.ID, req.Platform, req.DeviceID)
	//if err := cacheClient.Del(ctx, refreshKey); err != nil {
	//	return nil, fmt.Errorf("删除 RefreshToken 出错: %w", err)
	//}
	//
	//// 删除设备记录
	//userDeviceKey := fmt.Sprintf("%s%d", authConfig.CacheUserDevicePrefix, req.ID)
	//if err := cacheClient.HDel(ctx, userDeviceKey, req.Platform); err != nil {
	//	return nil, fmt.Errorf("删除设备记录出错: %w", err)
	//}

	uidStr := strconv.Itoa(int(mUser.ID))

	logoutInput := &authx.CustomClaims{
		UserID:       "user-" + uidStr,
		PlatformSign: "qwq",
		DeviceSign:   "web/pc",
	}

	ax, _ := app.Get[*authx.AuthX]()

	// 执行退出
	err = ax.Logout(ctx, logoutInput)

	// 退出成功
	return map[string]interface{}{
		"user_id": mUser.ID,
	}, err
}
