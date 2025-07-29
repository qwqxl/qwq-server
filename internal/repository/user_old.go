package repository

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"gorm.io/gorm"
//	"qwqserver/internal/model"
//	"time"
//)
//
//// UserRepository 用户领域仓库接口
//type UserRepository interface {
//	FindByID(ctx context.Context, id uint) (*model.User, error)
//	Create(ctx context.Context, user *model.User) error
//	Update(ctx context.Context, user *model.User) error
//	Delete(ctx context.Context, id uint) error
//	WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error
//
//	// ---------- User 相关操作 ----------- //
//
//	FindByEmail(ctx context.Context, email string) (*model.User, error)
//	FindByUsername(ctx context.Context, username string) (*model.User, error)
//	//ExistEmail(ctx context.Context, email string) (bool, error)
//	//ExistUsername(ctx context.Context, username string) (bool, error)
//	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
//
//	// ----------- User 新增用户管理方法 ----------- //
//
//	UpdateLoginInfo(ctx context.Context, userID uint, lastLogin time.Time, ip string) error
//	IncrementFailedAttempt(ctx context.Context, userID uint, lastFailed time.Time) error
//	ResetFailedAttempts(ctx context.Context, userID uint) error
//	UpdatePassword(ctx context.Context, userID uint, passwordHash, passwordSalt string, iterations int) error
//	UpdateStatus(ctx context.Context, userID uint, status uint8) error
//	UpdatePermissions(ctx context.Context, userID uint, perms uint64) error
//}
//
//// userRepository 用户仓库实现
//type userRepository struct {
//	*BaseRepository[model.User]
//}
//
//// NewUserRepository 创建新的用户仓库
//func NewUserRepository() (UserRepository, error) {
//	baseRepo, err := NewBaseRepository[model.User]()
//	if err != nil {
//		return nil, err
//	}
//	return &userRepository{BaseRepository: baseRepo}, nil
//}
//
//// List 获取用户列表
//func (r userRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
//	// 参数验证
//	if page < 1 {
//		page = 1
//	}
//	if pageSize < 1 || pageSize > 100 {
//		pageSize = 20 // 默认分页大小
//	}
//	offset := (page - 1) * pageSize
//
//	// 获取总数
//	total, err := r.Count(ctx, nil)
//	if err != nil {
//		return nil, 0, fmt.Errorf("failed to count users: %w", err)
//	}
//
//	// 获取分页数据
//	users, err := r.Query(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Order("created_at DESC").Offset(offset).Limit(pageSize)
//	})
//
//	if err != nil {
//		return nil, 0, fmt.Errorf("failed to list users: %w", err)
//	}
//
//	// 处理空列表
//	if users == nil {
//		users = []*model.User{}
//	}
//
//	return users, total, nil
//}
//
//// FindByEmailOld 根据邮箱查找用户
//func (r userRepository) FindByEmailOld(ctx context.Context, email string) (*model.User, error) {
//	u, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where("email = ?", email)
//	})
//	if err != nil {
//		return nil, err
//	}
//	return u, nil
//}
//
//// FindByEmail 根据邮箱查找用户
//func (r userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
//	if email == "" {
//		return nil, errors.New("email cannot be empty")
//	}
//
//	u, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where("email = ?", email)
//	})
//
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, fmt.Errorf("user with email %s not found", email)
//		}
//		return nil, fmt.Errorf("failed to find user by email: %w", err)
//	}
//	return u, nil
//}
//
//// FindByUsernameOld 根据用户名查找用户
//func (r userRepository) FindByUsernameOld(ctx context.Context, username string) (*model.User, error) {
//	u, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where("username = ?", username)
//	})
//	if err != nil {
//		return nil, err
//	}
//	return u, nil
//}
//
//// FindByUsername 根据用户名查找用户
//func (r userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
//	if username == "" {
//		return nil, errors.New("username cannot be empty")
//	}
//
//	u, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where("username = ?", username)
//	})
//
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, fmt.Errorf("user %s not found", username)
//		}
//		return nil, fmt.Errorf("failed to find user by username: %w", err)
//	}
//	return u, nil
//}
//
//// ExistEmailOld 判断用户是否存在
//func (r userRepository) ExistEmailOld(ctx context.Context, email string) (bool, error) {
//	user, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where("email = ?", email)
//	})
//	if user == nil {
//		return false, err
//	}
//	return true, nil
//}
//
//// ExistEmail 判断邮箱是否存在
//
//// ExistUsernameOld 检测用户名是否存在 存在则为true
//func (r userRepository) ExistUsernameOld(ctx context.Context, username string) (bool, error) {
//	user, err := r.First(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where("username = ?", username)
//	})
//	if user == nil {
//		return false, err
//	}
//	return true, nil
//}
//
//// ExistEmail 判断邮箱是否存在
////func (r userRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
////	return r.ExistsByFields(ctx, "email", email)
////}
////
////// ExistUsername 判断用户名是否存在
////func (r userRepository) ExistUsername(ctx context.Context, username string) (bool, error) {
////	return r.FieldExists(ctx, "username", username)
////}
//
//// ListOld 获取用户列表
//func (r userRepository) ListOld(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
//	offset := (page - 1) * pageSize
//
//	// 获取总数
//	total, err := r.Count(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Model(&model.User{})
//	})
//	if err != nil {
//		return nil, 0, err
//	}
//
//	// 获取分页数据
//	users, err := r.Query(ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Offset(offset).Limit(pageSize)
//	})
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return users, total, nil
//}
//
//// WithTransaction 在事务中执行用户操作
//func (r userRepository) WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error {
//	return r.BaseRepository.WithTransaction(ctx, func(txRepo *BaseRepository[model.User]) error {
//		txUserRepo := &userRepository{BaseRepository: txRepo}
//		return fn(txUserRepo)
//	})
//}
//
//// --------- user 新增用户管理方法 --------- //
//
//// UpdateLoginInfo 更新用户登录信息
//func (r userRepository) UpdateLoginInfo(ctx context.Context, userID uint, lastLogin time.Time, ip string) error {
//	if userID == 0 {
//		return errors.New("user ID cannot be zero")
//	}
//
//	return r.db.WithContext(ctx).Model(&model.User{}).
//		Where("id = ?", userID).
//		Updates(map[string]interface{}{
//			"last_login_at":   lastLogin,
//			"ip_address":      ip,
//			"failed_attempts": 0, // 重置失败尝试次数
//		}).Error
//}
//
//// IncrementFailedAttempt 增加登录失败次数
//func (r userRepository) IncrementFailedAttempt(ctx context.Context, userID uint, lastFailed time.Time) error {
//	if userID == 0 {
//		return errors.New("user ID cannot be zero")
//	}
//
//	return r.db.WithContext(ctx).Model(&model.User{}).
//		Where("id = ?", userID).
//		Updates(map[string]interface{}{
//			"last_failed_attempt": lastFailed,
//			"failed_attempts":     gorm.Expr("failed_attempts + ?", 1),
//		}).Error
//}
//
//// ResetFailedAttempts 重置登录失败次数
//func (r userRepository) ResetFailedAttempts(ctx context.Context, userID uint) error {
//	if userID == 0 {
//		return errors.New("user ID cannot be zero")
//	}
//
//	return r.db.WithContext(ctx).Model(&model.User{}).
//		Where("id = ?", userID).
//		Updates(map[string]interface{}{
//			"failed_attempts":     0,
//			"last_failed_attempt": nil,
//		}).Error
//}
//
//// UpdatePassword 更新用户密码
//func (r userRepository) UpdatePassword(ctx context.Context, userID uint, passwordHash, passwordSalt string, iterations int) error {
//	if userID == 0 {
//		return errors.New("user ID cannot be zero")
//	}
//	if passwordHash == "" || passwordSalt == "" {
//		return errors.New("password hash and salt cannot be empty")
//	}
//	if iterations < 1000 {
//		return errors.New("iterations must be at least 1000")
//	}
//
//	return r.db.WithContext(ctx).Model(&model.User{}).
//		Where("id = ?", userID).
//		Updates(map[string]interface{}{
//			"password_hash": passwordHash,
//			"password_salt": passwordSalt,
//			"iterations":    iterations,
//		}).Error
//}
//
//// UpdateStatus 更新用户状态
//func (r userRepository) UpdateStatus(ctx context.Context, userID uint, status uint8) error {
//	if userID == 0 {
//		return errors.New("user ID cannot be zero")
//	}
//	if status > 2 { // 假设状态值范围是0-2
//		return errors.New("invalid user status")
//	}
//
//	return r.db.WithContext(ctx).Model(&model.User{}).
//		Where("id = ?", userID).
//		Update("status", status).Error
//}
//
//// UpdatePermissions 更新用户权限
//func (r userRepository) UpdatePermissions(ctx context.Context, userID uint, perms uint64) error {
//	if userID == 0 {
//		return errors.New("user ID cannot be zero")
//	}
//
//	return r.db.WithContext(ctx).Model(&model.User{}).
//		Where("id = ?", userID).
//		Update("perms", perms).Error
//}
