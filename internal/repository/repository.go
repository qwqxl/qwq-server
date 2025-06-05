package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"qwqserver/pkg/database"

	"gorm.io/gorm"
)

// BaseRepository 提供基础的 CRUD 操作
type BaseRepository[T any] struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewBaseRepository 创建新的基础 Repository
func NewBaseRepository[T any]() (*BaseRepository[T], error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}
	return &BaseRepository[T]{db: db}, nil
}

// WithTransaction 在事务中执行操作
func (r *BaseRepository[T]) WithTransaction(ctx context.Context, fn func(txRepo *BaseRepository[T]) error) error {
	return database.WithTransaction(ctx, func(tx *gorm.DB) error {
		txRepo := &BaseRepository[T]{db: tx}
		return fn(txRepo)
	})
}

// Create 创建新记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("创建记录失败: %w", err)
	}
	return nil
}

// Update 更新记录
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("更新记录失败: %w", err)
	}
	return nil
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		return fmt.Errorf("删除记录失败: %w", err)
	}
	return nil
}

// FindByID 根据ID查找记录
func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查找记录失败: %w", err)
	}
	return &entity, nil
}

// FindAll 查找所有记录
func (r *BaseRepository[T]) FindAll(ctx context.Context) ([]*T, error) {
	var entities []*T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查找所有记录失败: %w", err)
	}
	return entities, nil
}

// Query 自定义查询
func (r *BaseRepository[T]) Query(ctx context.Context, query func(db *gorm.DB) *gorm.DB) ([]*T, error) {
	var entities []*T
	db := r.db.WithContext(ctx)
	if err := query(db).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	return entities, nil
}

// First 查找第一条匹配记录
func (r *BaseRepository[T]) First(ctx context.Context, query func(db *gorm.DB) *gorm.DB) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)
	if err := query(db).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查找记录失败: %w", err)
	}
	return &entity, nil
}

// Count 统计记录数量
func (r *BaseRepository[T]) Count(ctx context.Context, query func(db *gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx)
	if err := query(db).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计记录失败: %w", err)
	}
	return count, nil
}
