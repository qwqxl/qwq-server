package repository

import (
	"context"
	"gorm.io/gorm"
	"qwqserver/internal/model"
)

// UserRepository 用户领域仓库接口
type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error

	// ---------- User 相关操作 ----------- //

	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	ExistEmail(ctx context.Context, email string) (bool, error)
	ExistUsername(ctx context.Context, username string) (bool, error)
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	*BaseRepository[model.User]
}

// FindByEmail 根据邮箱查找用户
func (r userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindByUsername 根据用户名查找用户
func (r userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	u, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ExistEmail 判断用户是否存在
func (r userRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	user, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	})
	if user == nil {
		return false, err
	}
	return true, nil
}

// ExistUsername 检测用户名是否存在 存在则为true
func (r userRepository) ExistUsername(ctx context.Context, username string) (bool, error) {
	user, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	})
	if user == nil {
		return false, err
	}
	return true, nil
}

// List 获取用户列表
func (r userRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize

	// 获取总数
	total, err := r.Count(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.User{})
	})
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	users, err := r.Query(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// WithTransaction 在事务中执行用户操作
func (r userRepository) WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error {
	return r.BaseRepository.WithTransaction(ctx, func(txRepo *BaseRepository[model.User]) error {
		txUserRepo := &userRepository{BaseRepository: txRepo}
		return fn(txUserRepo)
	})
}

// NewUserRepository 创建新的用户仓库
func NewUserRepository() (UserRepository, error) {
	baseRepo, err := NewBaseRepository[model.User]()
	if err != nil {
		return nil, err
	}
	return &userRepository{BaseRepository: baseRepo}, nil
}
