package user

import (
	"context"
	"errors"
	"fmt"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	"qwqserver/internal/vo/vouser"
	"qwqserver/pkg/util"
	"qwqserver/pkg/util/passsec"
	"time"
)

// 注册服务实现
func (s *Service) Register(ctx context.Context, req *vouser.UserRegisterRequest) (*vouser.UserRegisterResponse, error) {
	// 参数验证
	//if err := validateRegisterRequest(req); err != nil {
	//	return nil, err
	//}

	// 创建用户仓库
	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("创建用户仓库失败: %w", err)
	}

	// 检查用户名是否已存在
	exist, err := userRepo.ExistUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exist {
		return nil, errors.New("用户名已被使用")
	}

	// 检查邮箱是否已存在
	exist, err = userRepo.ExistEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exist {
		return nil, errors.New("邮箱已被注册")
	}

	// 模拟密码加密
	//hashedPwd := req.Password

	// 密码加密
	hashedPwd, err := passsec.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户模型
	newUser := &model.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Password:     hashedPwd,
		PasswordHash: hashedPwd,
		PasswordSalt: hashedPwd,
		Iterations:   0,
		Status:       uint8(model.UserStatusActive),
	}

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	// 创建用户
	if err := userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 构造响应数据（脱敏返回）
	res := &vouser.UserRegisterResponse{}

	res.ID = newUser.ID
	res.Username = newUser.Username
	res.Nickname = newUser.Nickname
	res.Email = newUser.Email
	res.CreatedAt = newUser.CreatedAt

	return res, nil
}

// validateRegisterRequest 验证注册请求参数
func validateRegisterRequest(req *vouser.UserRegisterRequest) error {
	if req == nil {
		return errors.New("请求不能为空")
	}

	if req.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(req.Username) < 4 || len(req.Username) > 20 {
		return errors.New("用户名长度必须在4-20个字符之间")
	}

	if req.Password == "" {
		return errors.New("密码不能为空")
	}
	if len(req.Password) < 8 {
		return errors.New("密码长度至少8位")
	}

	if err := util.ValidatePasswordComplexity(req.Password); err != nil {
		return err
	}

	if req.Nickname == "" {
		return errors.New("昵称不能为空")
	}
	if len(req.Nickname) > 30 {
		return errors.New("昵称长度不能超过30个字符")
	}

	if req.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if !util.IsValidEmailFormat(req.Email) {
		return errors.New("邮箱格式不正确")
	}

	allowedSuffixes := []string{
		"qq.com",
		"163.com",
		"gmail.com",
		"outlook.com",
		"iqwq.com",
	}

	if !util.IsAllowedEmailSuffix(req.Email, allowedSuffixes) {
		return fmt.Errorf("仅支持以下邮箱后缀注册: %v", allowedSuffixes)
	}

	return nil
}
