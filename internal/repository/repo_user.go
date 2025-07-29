package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qwqserver/internal/model"
	"time"
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
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)

	// ----------- User 新增用户管理方法 ----------- //
	UpdateLoginInfo(ctx context.Context, userID uint, lastLogin time.Time, ip string) error
	IncrementFailedAttempt(ctx context.Context, userID uint, lastFailed time.Time) error
	ResetFailedAttempts(ctx context.Context, userID uint) error
	UpdatePassword(ctx context.Context, userID uint, passwordHash, passwordSalt string, iterations int) error
	UpdateStatus(ctx context.Context, userID uint, status uint8) error
	UpdatePermissions(ctx context.Context, userID uint, perms uint64) error
	ExistEmail(ctx context.Context, email string) (bool, error)
	ExistUsername(ctx context.Context, username string) (bool, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	*BaseRepository[model.User]
}

// NewUserRepository 创建新的用户仓库
func NewUserRepository() (UserRepository, error) {
	baseRepo, err := NewBaseRepository[model.User](
		WithCacheNS[model.User]("user"),
		WithCacheTTL[model.User](10*time.Minute),
	)
	if err != nil {
		return nil, err
	}
	return &userRepository{BaseRepository: baseRepo}, nil
}

// ========== 核心方法实现 ==========

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.BaseRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// 先获取旧用户数据用于清除缓存
	oldUser, err := r.FindByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 更新用户
	if err := r.BaseRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	// 清除所有相关缓存
	if err := r.clearUserCache(ctx, oldUser); err != nil {
		return fmt.Errorf("清除缓存失败: %w", err)
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	// 先获取用户数据用于清除缓存
	user, err := r.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 删除用户
	if err := r.BaseRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	// 清除所有相关缓存
	if err := r.clearUserCache(ctx, user); err != nil {
		return fmt.Errorf("清除缓存失败: %w", err)
	}

	return nil
}

// WithTransaction 在事务中执行用户操作
func (r *userRepository) WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error {
	return r.BaseRepository.WithTransaction(ctx, func(txRepo *BaseRepository[model.User]) error {
		txUserRepo := &userRepository{BaseRepository: txRepo}
		return fn(txUserRepo)
	})
}

// ========== 用户查询方法 ==========

// FindByEmail 根据邮箱查找用户（带缓存）
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}

	cacheKey := r.emailCacheKey(email)

	user, err := r.GetOrSet(ctx, cacheKey, func() (*model.User, time.Duration, error) {
		// 数据库查询
		u, err := r.BaseRepository.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
			return db.Where("email = ?", email)
		})

		if err != nil {
			return nil, 0, fmt.Errorf("数据库查询失败: %w", err)
		}
		if u == nil {
			return nil, 0, fmt.Errorf("邮箱 %s 对应的用户不存在", email)
		}

		return u, r.cacheTTL, nil
	})

	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// FindByUsername 根据用户名查找用户（带缓存）
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}

	cacheKey := r.usernameCacheKey(username)

	user, err := r.GetOrSet(ctx, cacheKey, func() (*model.User, time.Duration, error) {
		// 数据库查询
		u, err := r.BaseRepository.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
			return db.Where("username = ?", username)
		})

		if err != nil {
			return nil, 0, fmt.Errorf("数据库查询失败: %w", err)
		}
		if u == nil {
			return nil, 0, fmt.Errorf("用户名 %s 对应的用户不存在", username)
		}

		return u, r.cacheTTL, nil
	})

	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// List 获取用户列表（分页）
func (r *userRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 使用BaseRepository的分页方法
	users, total, err := r.BaseRepository.List(ctx, page, pageSize, "created_at DESC", nil)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 处理空列表
	if users == nil {
		users = []*model.User{}
	}

	return users, total, nil
}

// ========== 用户管理方法 ==========

// UpdateLoginInfo 更新用户登录信息
func (r *userRepository) UpdateLoginInfo(ctx context.Context, userID uint, lastLogin time.Time, ip string) error {
	if userID == 0 {
		return errors.New("用户ID不能为零")
	}

	updates := map[string]interface{}{
		"last_login_at":   lastLogin,
		"ip_address":      ip,
		"failed_attempts": 0, // 重置失败尝试次数
	}

	if err := r.PartialUpdate(ctx, userID, updates); err != nil {
		return fmt.Errorf("更新登录信息失败: %w", err)
	}

	// 清除用户缓存
	if err := r.DelCacheByID(ctx, userID); err != nil {
		return fmt.Errorf("清除缓存失败: %w", err)
	}

	return nil
}

// IncrementFailedAttempt 增加登录失败次数
func (r *userRepository) IncrementFailedAttempt(ctx context.Context, userID uint, lastFailed time.Time) error {
	if userID == 0 {
		return errors.New("用户ID不能为零")
	}

	updates := map[string]interface{}{
		"last_failed_attempt": lastFailed,
		"failed_attempts":     gorm.Expr("failed_attempts + ?", 1),
	}

	if err := r.PartialUpdate(ctx, userID, updates); err != nil {
		return fmt.Errorf("更新失败尝试次数失败: %w", err)
	}

	return nil
}

// ResetFailedAttempts 重置登录失败次数
func (r *userRepository) ResetFailedAttempts(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("用户ID不能为零")
	}

	updates := map[string]interface{}{
		"failed_attempts":     0,
		"last_failed_attempt": nil,
	}

	if err := r.PartialUpdate(ctx, userID, updates); err != nil {
		return fmt.Errorf("重置失败尝试次数失败: %w", err)
	}

	return nil
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, passwordHash, passwordSalt string, iterations int) error {
	if userID == 0 {
		return errors.New("用户ID不能为零")
	}
	if passwordHash == "" || passwordSalt == "" {
		return errors.New("密码哈希值和盐不能为空")
	}
	if iterations < 1000 {
		return errors.New("迭代次数必须至少为1000")
	}

	updates := map[string]interface{}{
		"password_hash": passwordHash,
		"password_salt": passwordSalt,
		"iterations":    iterations,
	}

	if err := r.PartialUpdate(ctx, userID, updates); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 清除用户缓存
	if err := r.DelCacheByID(ctx, userID); err != nil {
		return fmt.Errorf("清除缓存失败: %w", err)
	}

	return nil
}

// UpdateStatus 更新用户状态
func (r *userRepository) UpdateStatus(ctx context.Context, userID uint, status uint8) error {
	if userID == 0 {
		return errors.New("用户ID不能为零")
	}
	if status > 2 {
		return errors.New("无效的用户状态")
	}

	if err := r.PartialUpdate(ctx, userID, map[string]interface{}{"status": status}); err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	// 清除用户缓存
	if err := r.DelCacheByID(ctx, userID); err != nil {
		return fmt.Errorf("清除缓存失败: %w", err)
	}

	return nil
}

// UpdatePermissions 更新用户权限
func (r *userRepository) UpdatePermissions(ctx context.Context, userID uint, perms uint64) error {
	if userID == 0 {
		return errors.New("用户ID不能为零")
	}

	if err := r.PartialUpdate(ctx, userID, map[string]interface{}{"perms": perms}); err != nil {
		return fmt.Errorf("更新用户权限失败: %w", err)
	}

	// 清除用户缓存
	if err := r.DelCacheByID(ctx, userID); err != nil {
		return fmt.Errorf("清除缓存失败: %w", err)
	}

	return nil
}

// ========== 辅助方法 ==========

// clearUserCache 清除用户所有缓存
func (r *userRepository) clearUserCache(ctx context.Context, user *model.User) error {
	// 清除ID缓存
	if err := r.DelCacheByID(ctx, user.ID); err != nil {
		return err
	}

	// 清除邮箱缓存
	if user.Email != "" {
		if err := r.cacheClient.Del(ctx, r.emailCacheKey(user.Email)); err != nil {
			return err
		}
	}

	// 清除用户名缓存
	if user.Username != "" {
		if err := r.cacheClient.Del(ctx, r.usernameCacheKey(user.Username)); err != nil {
			return err
		}
	}

	return nil
}

// emailCacheKey 生成邮箱缓存键
func (r *userRepository) emailCacheKey(email string) string {
	return r.cacheKey("email", email)
}

// usernameCacheKey 生成用户名缓存键
func (r *userRepository) usernameCacheKey(username string) string {
	return r.cacheKey("username", username)
}

// exists

// ExistEmail 检查邮箱是否存在
func (r *userRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("邮箱不能为空")
	}
	return r.BaseRepository.Exists(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	})
}

// ExistUsername 检查用户名是否存在
func (r *userRepository) ExistUsername(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, errors.New("用户名不能为空")
	}
	return r.BaseRepository.Exists(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	})
}
