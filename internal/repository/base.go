package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qwqserver/internal/app"
	"reflect"
	"sync"
	"time"

	"qwqserver/pkg/cache"
	"qwqserver/pkg/database"
)

// BaseRepository 基础仓库结构
type BaseRepository[T any] struct {
	db          *gorm.DB
	cacheClient *cache.Client
	cacheTTL    time.Duration
	cacheNS     string
	useCache    bool
	modelType   reflect.Type
	lock        sync.RWMutex
}

// RepoOption 仓库配置函数
type RepoOption[T any] func(*BaseRepository[T])

// NewBaseRepository 创建基础仓库
func NewBaseRepository[T any](opts ...RepoOption[T]) (*BaseRepository[T], error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.New("获取数据库连接失败：" + err.Error())
	}

	cacheClient, err := app.Get[*cache.Client]()
	if err != nil {
		return nil, errors.New("获取缓存客户端失败：" + err.Error())
	}

	// 获取泛型类型信息
	var model T
	modelType := reflect.TypeOf(model)
	if modelType == nil {
		return nil, errors.New("无法获取模型类型信息")
	}

	repo := &BaseRepository[T]{
		db:          db,
		cacheClient: cacheClient,
		cacheTTL:    5 * time.Minute,
		cacheNS:     "repo",
		useCache:    true,
		modelType:   modelType,
	}

	for _, opt := range opts {
		opt(repo)
	}

	return repo, nil
}

// ========== CRUD 操作 ==========

// Create 创建记录
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
	return r.delCache(ctx, entity)
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		return fmt.Errorf("删除记录失败: %w", err)
	}
	return r.DelCacheByID(ctx, id)
}

// FindByID 根据ID查找
func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	cacheKey := r.cacheKey("id", fmt.Sprint(id))

	// 尝试从缓存获取
	if r.useCache {
		var cached T
		if err := r.cacheClient.GetJSON(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		} else if !errors.Is(err, cache.ErrKeyNotFound) {
			return nil, fmt.Errorf("缓存读取失败: %w", err)
		}
	}

	// 从数据库查询
	var result T
	if err := r.db.WithContext(ctx).First(&result, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("记录不存在")
		}
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	// 设置缓存
	if r.useCache {
		if err := r.cacheClient.SetJSON(ctx, cacheKey, result, r.cacheTTL); err != nil {
			return nil, fmt.Errorf("缓存设置失败: %w", err)
		}
	}

	return &result, nil
}

// FindOne 条件查询单条记录
func (r *BaseRepository[T]) FindOne(ctx context.Context, where func(db *gorm.DB) *gorm.DB) (*T, error) {
	var result T
	if err := where(r.db.WithContext(ctx)).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	return &result, nil
}

// FindAll 查询所有记录
func (r *BaseRepository[T]) FindAll(ctx context.Context, where func(db *gorm.DB) *gorm.DB) ([]*T, error) {
	var results []*T
	if err := where(r.db.WithContext(ctx)).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	return results, nil
}

// List 分页查询
func (r *BaseRepository[T]) List(ctx context.Context, page, pageSize int, order string, where func(db *gorm.DB) *gorm.DB) ([]*T, int64, error) {
	// 计算总数
	var total int64
	if err := where(r.db.Model(new(T))).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取总数失败: %w", err)
	}

	// 分页查询
	var results []*T
	offset := (page - 1) * pageSize
	db := where(r.db.WithContext(ctx)).Offset(offset).Limit(pageSize)
	if order != "" {
		db = db.Order(order)
	}

	if err := db.Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询失败: %w", err)
	}

	return results, total, nil
}

// Exists 判断记录是否存在
func (r *BaseRepository[T]) Exists(ctx context.Context, where func(db *gorm.DB) *gorm.DB) (bool, error) {
	var count int64
	if err := where(r.db.WithContext(ctx).Model(new(T))).Count(&count).Error; err != nil {
		return false, fmt.Errorf("查询失败: %w", err)
	}
	return count > 0, nil
}

// Count 统计记录数量
func (r *BaseRepository[T]) Count(ctx context.Context, where func(db *gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	if err := where(r.db.WithContext(ctx).Model(new(T))).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计失败: %w", err)
	}
	return count, nil
}

// ========== 缓存操作 ==========

// GetOrSet 获取或设置缓存
func (r *BaseRepository[T]) GetOrSet(ctx context.Context, key string, fn func() (*T, time.Duration, error)) (*T, error) {
	cacheKey := r.cacheKey("custom", key)

	// 尝试从缓存获取
	if r.useCache {
		var cached T
		if err := r.cacheClient.GetJSON(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		} else if !errors.Is(err, cache.ErrKeyNotFound) {
			return nil, fmt.Errorf("缓存读取失败: %w", err)
		}
	}

	// 执行加载函数
	result, ttl, err := fn()
	if err != nil {
		return nil, fmt.Errorf("加载函数执行失败: %w", err)
	}

	// 设置缓存
	if r.useCache {
		actualTTL := r.cacheTTL
		if ttl > 0 {
			actualTTL = ttl
		}

		if err := r.cacheClient.SetJSON(ctx, cacheKey, result, actualTTL); err != nil {
			return nil, fmt.Errorf("缓存设置失败: %w", err)
		}
	}

	return result, nil
}

// DelCacheByID 根据ID删除缓存
func (r *BaseRepository[T]) DelCacheByID(ctx context.Context, id uint) error {
	return r.delCacheByKey(ctx, "id", fmt.Sprint(id))
}

// DelCacheByKey 删除自定义键缓存
func (r *BaseRepository[T]) DelCacheByKey(ctx context.Context, key string) error {
	return r.delCacheByKey(ctx, "custom", key)
}

// DelCacheByTag 根据标签删除缓存
func (r *BaseRepository[T]) DelCacheByTag(ctx context.Context, tag string) error {
	pattern := r.cacheKey("tag", tag) + ":*"
	return r.deleteKeysByPattern(ctx, pattern)
}

// ========== 事务支持 ==========

// WithTransaction 执行事务
func (r *BaseRepository[T]) WithTransaction(ctx context.Context, fn func(txRepo *BaseRepository[T]) error) error {
	return database.WithTransaction(ctx, func(tx *gorm.DB) error {
		txRepo := &BaseRepository[T]{
			db:          tx,
			cacheClient: r.cacheClient,
			useCache:    r.useCache,
			cacheTTL:    r.cacheTTL,
			cacheNS:     r.cacheNS,
			modelType:   r.modelType,
		}
		return fn(txRepo)
	})
}

// ========== 高级功能 ==========

// PartialUpdate 部分更新
func (r *BaseRepository[T]) PartialUpdate(ctx context.Context, id uint, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("部分更新失败: %w", err)
	}

	// 延迟双删策略
	go func() {
		time.Sleep(1 * time.Second)
		_ = r.DelCacheByID(context.Background(), id)
	}()

	return nil
}

// OptimisticUpdate 乐观锁更新
func (r *BaseRepository[T]) OptimisticUpdate(ctx context.Context, entity *T, versionField string) error {
	currentVersion := reflect.ValueOf(entity).Elem().FieldByName(versionField).Uint()
	newVersion := currentVersion + 1

	updateData := map[string]interface{}{
		versionField: newVersion,
	}

	result := r.db.WithContext(ctx).Model(entity).Where("id = ? AND "+versionField+" = ?",
		getIDValue(entity), currentVersion).Updates(updateData)

	if result.Error != nil {
		return fmt.Errorf("更新失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("乐观锁冲突: 记录已被修改")
	}

	// 更新本地版本号
	reflect.ValueOf(entity).Elem().FieldByName(versionField).SetUint(newVersion)
	return r.delCache(ctx, entity)
}

// WithLock 分布式锁执行
func (r *BaseRepository[T]) WithLock(ctx context.Context, lockKey string, ttl time.Duration, fn func(args ...any) error) error {
	return r.cacheClient.WithSafeLock(ctx, lockKey, ttl, fn)
}

// ========== 辅助方法 ==========

// cacheKey 生成缓存键
func (r *BaseRepository[T]) cacheKey(keyType, value string) string {
	typeName := r.modelType.Name()
	if typeName == "" {
		typeName = "unknown"
	}
	return fmt.Sprintf("%s:%s:%s:%s", r.cacheNS, typeName, keyType, value)
}

// delCacheByKey 删除缓存
func (r *BaseRepository[T]) delCacheByKey(ctx context.Context, keyType, key string) error {
	if !r.useCache {
		return nil
	}
	cacheKey := r.cacheKey(keyType, key)
	return r.cacheClient.Del(ctx, cacheKey)
}

// deleteKeysByPattern 按模式删除缓存
func (r *BaseRepository[T]) deleteKeysByPattern(ctx context.Context, pattern string) error {
	if !r.useCache {
		return nil
	}
	keys, err := r.cacheClient.Redis.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("获取缓存键失败: %w", err)
	}

	if len(keys) > 0 {
		if err := r.cacheClient.Del(ctx, keys...); err != nil {
			return fmt.Errorf("删除缓存失败: %w", err)
		}
	}
	return nil
}

// delCache 删除实体缓存
func (r *BaseRepository[T]) delCache(ctx context.Context, entity *T) error {
	id := getIDValue(entity)
	return r.DelCacheByID(ctx, id)
}

// getIDValue 获取实体ID值
func getIDValue[T any](entity *T) uint {
	val := reflect.ValueOf(entity).Elem()
	idField := val.FieldByName("ID")
	if !idField.IsValid() || idField.Kind() != reflect.Uint {
		return 0
	}
	return uint(idField.Uint())
}

// ========== 配置选项 ==========

// WithCacheEnabled 启用/禁用缓存
func WithCacheEnabled[T any](enabled bool) RepoOption[T] {
	return func(repo *BaseRepository[T]) {
		repo.useCache = enabled
	}
}

// WithCacheTTL 设置缓存TTL
func WithCacheTTL[T any](ttl time.Duration) RepoOption[T] {
	return func(repo *BaseRepository[T]) {
		repo.cacheTTL = ttl
	}
}

// WithCacheNS 设置缓存命名空间
func WithCacheNS[T any](ns string) RepoOption[T] {
	return func(repo *BaseRepository[T]) {
		repo.cacheNS = ns
	}
}

// WithDB 使用自定义数据库连接
func WithDB[T any](db *gorm.DB) RepoOption[T] {
	return func(repo *BaseRepository[T]) {
		repo.db = db
	}
}

// WithRedisClient 使用自定义Redis客户端
func WithRedisClient[T any](client *cache.Client) RepoOption[T] {
	return func(repo *BaseRepository[T]) {
		repo.cacheClient = client
	}
}

// WithReadReplica 使用只读副本
func WithReadReplica[T any](replicaName string) RepoOption[T] {
	return func(repo *BaseRepository[T]) {
		if db, err := database.GetDBByName(replicaName); err == nil {
			repo.db = db
		}
	}
}
