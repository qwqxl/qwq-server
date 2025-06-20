package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/auth"
	"qwqserver/internal/common"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/util/passsec"
	"strconv"
)

type AuthService struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// Del 删除用户
func (s *AuthService) Del(uid uint, platform, deviceID string) (res *common.HTTPResult) {
	// 初始化 返回结果
	res = &common.HTTPResult{}
	user := &model.User{}

	if uid == 0 {
		res.Code = 400
		res.Msg = "用户ID不能为: 0"
		return
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = 500
		res.Msg = "获取数据库连接失败: " + err.Error()
		return
	}

	// 检查用户是否存在
	if user, err = userRepo.FindByID(context.Background(), uid); err != nil || user == nil {
		res.Code = 404
		res.Msg = "用户不存在"
		return
	}

	if err = userRepo.Delete(context.Background(), user.ID); err != nil {
		res.Code = 500
		res.Msg = "删除用户失败: " + err.Error()
		return
	}

	res.Code = 200
	res.Msg = "删除用户成功"
	res.Data = map[string]any{
		"uid":      user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	// 登出
	return s.Logout(uid, platform, deviceID)
}

// 注册
func (s *AuthService) Register() *common.HTTPResult {
	res := &common.HTTPResult{}
	req := s

	if req.Username == "" || req.Password == "" || req.Nickname == "" || req.Email == "" {
		res.Code = http.StatusBadRequest
		res.Msg = "注册用户参数错误"
		return res
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "获取数据库连接失败: " + err.Error()
		return res
	}

	// 检查用户名是否已存在
	if ok, _ := userRepo.ExistUsername(context.Background(), req.Username); ok {
		res.Code = http.StatusConflict
		res.Msg = "您输入的用户名已经存在"
		return res
	}

	// 检查邮箱是否已存在
	if ok, _ := userRepo.ExistEmail(context.Background(), req.Email); ok {
		res.Code = http.StatusConflict
		res.Msg = "您输入的邮箱已经存在"
		return res
	}

	// 密码哈希
	pwd, err := passsec.Hash(req.Password)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "密码哈希失败: " + err.Error()
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
		res.Msg = "创建用户失败 Error: " + err.Error()
		return res
	}

	res.Code = http.StatusOK
	res.Data = gin.H{"user_id": newUser.ID}
	res.Msg = "注册成功！"
	return res
}

// 登录
func (s *AuthService) Login(platform, deviceID string) (res *common.HTTPResult) {
	res = &common.HTTPResult{}
	token := &common.JWTResult{
		DeviceID: deviceID,
		Platform: platform,
	}
	mUser := &model.User{}
	req := s

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "获取数据库连接失败: " + err.Error()
		return
	}

	if req.Username != "" {
		if mUser, _ = userRepo.FindByUsername(context.Background(), req.Username); mUser == nil {
			res.Code = http.StatusNotFound
			res.Msg = "用户名不存在"
			return
		}
	} else {
		if mUser, _ = userRepo.FindByEmail(context.Background(), req.Email); mUser == nil {
			res.Code = http.StatusNotFound
			res.Msg = "邮箱不存在"
			return
		}
	}

	// 校验密码（应使用安全哈希，略）
	var ok bool
	if ok, err = passsec.Check(req.Password, mUser.Password); err != nil {
		res.Code = http.StatusUnauthorized
		res.Msg = "密码校验出错"
		return
	}

	if !ok {
		res.Code = http.StatusUnauthorized
		res.Msg = "密码错误"
		return
	}

	// 禁用状态
	if mUser.Status == 0 {
		res.Code = http.StatusForbidden
		res.Msg = "用户已被封禁"
		return
	}

	uid := strconv.Itoa(int(mUser.ID))
	token.ID = mUser.ID

	// 生成Token和RefreshToken
	newToken, err := auth.GenerateToken(uid, platform, deviceID)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "生成Token出错"
		token.AccessToken = ""
		token.RefreshToken = ""
		return
	}
	// token
	token.AccessToken = newToken

	refreshToken, err := auth.GenerateRefreshToken(uid, platform, deviceID)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "生成RefreshToken出错"
		token.AccessToken = ""
		token.RefreshToken = ""
		return
	}
	// refresh token
	token.RefreshToken = refreshToken

	// 存储Token到Redis
	tokenKey := common.RedisTokenPrefix + uid + ":" + platform + ":" + deviceID
	cache.Set(tokenKey, token, common.TokenExpireTime)

	// 存储RefreshToken
	refreshKey := common.RedisRefreshPrefix + uid + ":" + platform + ":" + deviceID
	cache.Set(refreshKey, refreshToken, common.RefreshTokenExpire)

	// 记录用户设备关系
	userDeviceKey := common.RedisUserDevicePrefix + uid
	cache.HSet(userDeviceKey, platform, deviceID)

	res.Code = http.StatusOK
	res.Msg = "登录成功"
	res.Data = token

	return
}

// 登出
func (s *AuthService) Logout(uid uint, platform, deviceID string) (res *common.HTTPResult) {
	//
	res = &common.HTTPResult{}
	mUser := &model.User{}
	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "获取数据库连接失败: " + err.Error()
		return
	}

	// 判断用户是否存在
	if mUser, err = userRepo.FindByID(context.Background(), uid); err != nil || mUser == nil {
		res.Code = http.StatusNotFound
		res.Msg = "用户不存在"
		return
	}

	// 删除Token
	tokenKey := fmt.Sprintf("%v:%v:%v:%v", common.RedisTokenPrefix, uid, platform, deviceID)
	if err = cache.Del(tokenKey); err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "删除Token出错" + err.Error()
		return
	}

	// 删除RefreshToken
	//refreshKey := common.RedisRefreshPrefix + userID + ":" + platform + ":" + deviceID
	refreshKey := fmt.Sprintf("%v:%v:%v:%v", common.RedisRefreshPrefix, uid, platform, deviceID)
	if err = cache.Del(refreshKey); err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "删除RefreshToken出错" + err.Error()
		return
	}

	// 删除设备记录
	userDeviceKey := fmt.Sprintf("%v%v", common.RedisUserDevicePrefix, uid)
	cache.HDel(userDeviceKey, platform)
	res.Msg = "登出成功"
	res.Code = http.StatusOK
	res.Data = map[string]any{
		"uid":      uid,
		"username": mUser.Username,
		"email":    mUser.Email,
		"nickname": mUser.Nickname,
	}
	return
}

// 刷新Token
func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := auth.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 验证RefreshToken是否有效
	refreshKey := common.RedisRefreshPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	storedRefresh, err := cache.Get(refreshKey)
	if err != nil || storedRefresh != refreshToken {
		return "", "", errors.New("无效的刷新令牌")
	}

	// 生成新Token和RefreshToken
	newToken, err := auth.GenerateToken(claims.UserID, claims.Platform, claims.DeviceID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := auth.GenerateRefreshToken(claims.UserID, claims.Platform, claims.DeviceID)
	if err != nil {
		return "", "", err
	}

	// 更新Redis中的Token
	tokenKey := common.RedisTokenPrefix + claims.UserID + ":" + claims.Platform + ":" + claims.DeviceID
	cache.Set(tokenKey, newToken, common.TokenExpireTime)

	// 更新RefreshToken
	cache.Set(refreshKey, newRefreshToken, common.RefreshTokenExpire)

	return newToken, newRefreshToken, nil
}

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
