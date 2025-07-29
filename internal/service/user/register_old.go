package user

//
//import (
//	"context"
//	"fmt"
//	"qwqserver/internal/model"
//	"qwqserver/internal/repository"
//	"qwqserver/internal/vo/vouser"
//	"qwqserver/pkg/util/passsec"
//)
//
//// 注册
//func (s *Service) Register(ctx context.Context, req vouser.UserRegisterRequest) (any, error) {
//	if req.Username == "" || req.Password == "" || req.Nickname == "" || req.Email == "" {
//		return nil, fmt.Errorf("注册参数缺失，用户名、密码、昵称、邮箱不能为空")
//	}
//
//	userRepo, err := repository.NewUserRepository()
//	if err != nil {
//		return nil, fmt.Errorf("数据库连接失败: %w", err)
//	}
//
//	// 用户名检查
//	if ok, _ := userRepo.ExistUsername(ctx, req.Username); ok {
//		return nil, fmt.Errorf("用户名已存在")
//	}
//
//	// 邮箱检查
//	if ok, _ := userRepo.ExistEmail(ctx, req.Email); ok {
//		return nil, fmt.Errorf("邮箱已被注册")
//	}
//
//	// 密码加密
//	hashedPwd, err := passsec.Hash(req.Password)
//	if err != nil {
//		return nil, fmt.Errorf("密码加密失败: %w", err)
//	}
//
//	newUser := &model.User{
//		Username:     req.Username,
//		Nickname:     req.Nickname,
//		Email:        req.Email,
//		Password:     hashedPwd,
//		PasswordHash: hashedPwd,
//		PasswordSalt: hashedPwd, // ⚠️ 建议将 salt 独立存储
//		Status:       1,
//	}
//
//	if err := userRepo.Create(ctx, newUser); err != nil {
//		return nil, fmt.Errorf("用户创建失败: %w", err)
//	}
//
//	// 构造响应数据（脱敏返回）
//	//res := &model.UserLoginResponse{}
//	return newUser, nil
//}
