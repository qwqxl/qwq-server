package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"qwqserver/internal/model"
	"qwqserver/pkg/database"
	"time"
)

// AuthRepository 用户认证仓库接口
type AuthRepository interface {
	// 用户认证
	Authenticate(ctx context.Context, email, password string) (*model.User, error)

	// 登录记录
	RecordLogin(ctx context.Context, userID uint, ipAddress string, deviceInfo string) error

	// 检查登录尝试次数
	CheckLoginAttempts(ctx context.Context, email string) (int, error)

	// 记录登录失败尝试
	RecordLoginFailure(ctx context.Context, email, ipAddress string) error

	// 重置登录失败计数
	ResetLoginFailures(ctx context.Context, email string) error

	// 生成访问令牌
	GenerateAccessToken(ctx context.Context, userID uint, expiresIn time.Duration) (string, error)

	// 验证访问令牌
	VerifyAccessToken(ctx context.Context, token string) (*model.User, error)

	// 吊销访问令牌
	RevokeAccessToken(ctx context.Context, token string) error
}

// authRepository 用户认证仓库实现
type authRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func (r *authRepository) CheckLoginAttempts(ctx context.Context, email string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *authRepository) RecordLoginFailure(ctx context.Context, email, ipAddress string) error {
	//TODO implement me
	panic("implement me")
}

func (r *authRepository) ResetLoginFailures(ctx context.Context, email string) error {
	//TODO implement me
	panic("implement me")
}

func (r *authRepository) GenerateAccessToken(ctx context.Context, userID uint, expiresIn time.Duration) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *authRepository) VerifyAccessToken(ctx context.Context, token string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *authRepository) RevokeAccessToken(ctx context.Context, token string) error {
	//TODO implement me
	panic("implement me")
}

// NewAuthRepository 创建新的认证仓库
func NewAuthRepository() (AuthRepository, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}
	return &authRepository{db: db}, nil
}

// Authenticate 用户认证
func (r *authRepository) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	userRepo, err := NewUserRepository()
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 检查账户状态
	if user.Status == 1 {
		return nil, errors.New("账户已被锁定")
	}
	if user.Status == 2 {
		return nil, errors.New("账户未激活")
	}

	// 验证密码
	if password != user.Password {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// RecordLogin 记录登录
func (r *authRepository) RecordLogin(ctx context.Context, userID uint, ipAddress string, deviceInfo string) error {
	//loginRecord := model.User{
	//	LoginAt:    time.Now(),
	//	IPAddress:  ipAddress,
	//	DeviceInfo: deviceInfo,
	//}
	//
	//if err := r.db.WithContext(ctx).Create(&loginRecord).Error; err != nil {
	//	return fmt.Errorf("记录登录失败: %w", err)
	//}
	//
	//// 更新用户最后登录时间
	//if err := r.db.WithContext(ctx).Model(&model.User{}).
	//	Where("id = ?", userID).
	//	Update("last_login_at", time.Now()).Error; err != nil {
	//	return fmt.Errorf("更新最后登录时间失败: %w", err)
	//}

	return nil
}
