package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qwqserver/internal/model"
	"time"
)

// PostRepository Post领域仓库接口
type PostRepository interface {
	FindByID(ctx context.Context, id uint) (*model.Post, error)
	Create(ctx context.Context, post *model.Post) error
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uint) error
	WithTransaction(ctx context.Context, fn func(repo PostRepository) error) error

	// --------- Post 相关操作 --------- //

	// ListByUserID 获取用户的所有帖子
	ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*model.Post, int64, error)

	// ListByCategory 获取分类下的帖子
	ListByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]*model.Post, int64, error)

	// Search 搜索帖子
	Search(ctx context.Context, query string, page, pageSize int) ([]*model.Post, int64, error)

	// IncrementViewCount 增加帖子浏览量
	IncrementViewCount(ctx context.Context, id uint) error

	// IncrementLikeCount 增加帖子点赞数
	IncrementLikeCount(ctx context.Context, id uint) error

	// DecrementLikeCount 减少帖子点赞数
	DecrementLikeCount(ctx context.Context, id uint) error

	// IncrementCommentCount 增加帖子评论数
	IncrementCommentCount(ctx context.Context, id uint) error

	// DecrementCommentCount 减少帖子评论数
	DecrementCommentCount(ctx context.Context, id uint) error

	// PinPost 置顶帖子
	PinPost(ctx context.Context, id uint) error

	// UnpinPost 取消置顶帖子
	UnpinPost(ctx context.Context, id uint) error

	// ListPopular 获取热门帖子
	ListPopular(ctx context.Context, days int, limit int) ([]*model.Post, error)

	// ListLatest 获取最新帖子
	ListLatest(ctx context.Context, limit int) ([]*model.Post, error)

	// ListRecommended 获取推荐帖子
	ListRecommended(ctx context.Context, userID uint, limit int) ([]*model.Post, error)

	// Exists 检查帖子是否存在
	Exists(ctx context.Context, id uint) (bool, error)
}

// postRepository Post仓库实现
type postRepository struct {
	*BaseRepository[model.Post]
}

// NewPostRepository 创建新的Post仓库
func NewPostRepository() (PostRepository, error) {
	baseRepo, err := NewBaseRepository[model.Post]()
	if err != nil {
		return nil, err
	}
	return &postRepository{BaseRepository: baseRepo}, nil
}

// WithTransaction 在事务中执行Post操作
func (r *postRepository) WithTransaction(ctx context.Context, fn func(repo PostRepository) error) error {
	return r.BaseRepository.WithTransaction(ctx, func(txRepo *BaseRepository[model.Post]) error {
		txPostRepo := &postRepository{BaseRepository: txRepo}
		return fn(txPostRepo)
	})
}

// ListByUserID 获取用户的所有帖子
func (r *postRepository) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*model.Post, int64, error) {
	offset := (page - 1) * pageSize

	// 获取总数
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计用户帖子数量失败: %w", err)
	}

	// 获取分页数据
	var posts []*model.Post
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户帖子失败: %w", err)
	}

	return posts, total, nil
}

// ListByCategory 获取分类下的帖子
func (r *postRepository) ListByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]*model.Post, int64, error) {
	offset := (page - 1) * pageSize

	// 获取总数
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计分类帖子数量失败: %w", err)
	}

	// 获取分页数据
	var posts []*model.Post
	if err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Order("is_pinned DESC, created_at DESC"). // 置顶帖子优先
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("查询分类帖子失败: %w", err)
	}

	return posts, total, nil
}

// Search 搜索帖子
func (r *postRepository) Search(ctx context.Context, query string, page, pageSize int) ([]*model.Post, int64, error) {
	offset := (page - 1) * pageSize
	searchQuery := "%" + query + "%"

	// 获取总数
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("title LIKE ? OR content LIKE ?", searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计搜索结果失败: %w", err)
	}

	// 获取分页数据
	var posts []*model.Post
	if err := r.db.WithContext(ctx).
		Where("title LIKE ? OR content LIKE ?", searchQuery, searchQuery).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("搜索帖子失败: %w", err)
	}

	return posts, total, nil
}

// IncrementViewCount 增加帖子浏览量 （未实现）
func (r *postRepository) IncrementViewCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("增加浏览量失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在")
	}

	return nil
}

// IncrementLikeCount 增加帖子点赞数 （未实现）
func (r *postRepository) IncrementLikeCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		Update("like_count", gorm.Expr("like_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("增加点赞数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在")
	}

	return nil
}

// DecrementLikeCount 减少帖子点赞数 （未实现）
func (r *postRepository) DecrementLikeCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ? AND like_count > 0", id).
		Update("like_count", gorm.Expr("like_count - ?", 1))

	if result.Error != nil {
		return fmt.Errorf("减少点赞数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在或点赞数已为0")
	}

	return nil
}

// IncrementCommentCount 增加帖子评论数 （未实现）
func (r *postRepository) IncrementCommentCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		Update("comment_count", gorm.Expr("comment_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("增加评论数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在")
	}

	return nil
}

// DecrementCommentCount 减少帖子评论数（未实现）
func (r *postRepository) DecrementCommentCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ? AND comment_count > 0", id).
		Update("comment_count", gorm.Expr("comment_count - ?", 1))

	if result.Error != nil {
		return fmt.Errorf("减少评论数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在或评论数已为0")
	}

	return nil
}

// PinPost 置顶帖子（半未实现）
func (r *postRepository) PinPost(ctx context.Context, id uint) error {
	// 先取消当前置顶的帖子
	if err := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("is_pinned = ?", true).
		Update("is_pinned", false).Error; err != nil {
		return fmt.Errorf("取消当前置顶失败: %w", err)
	}

	// 置顶新帖子
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		Update("is_pinned", true)

	if result.Error != nil {
		return fmt.Errorf("置顶帖子失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在")
	}

	return nil
}

// UnpinPost 取消置顶帖子（半未实现）
func (r *postRepository) UnpinPost(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		Update("is_pinned", false)

	if result.Error != nil {
		return fmt.Errorf("取消置顶失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("帖子不存在")
	}

	return nil
}

// ListPopular 获取热门帖子（未实现）
func (r *postRepository) ListPopular(ctx context.Context, days int, limit int) ([]*model.Post, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var posts []*model.Post
	err := r.db.WithContext(ctx).
		Where("created_at >= ?", startDate).
		Order("like_count DESC, comment_count DESC, view_count DESC").
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, fmt.Errorf("获取热门帖子失败: %w", err)
	}

	return posts, nil
}

// ListLatest 获取最新帖子
func (r *postRepository) ListLatest(ctx context.Context, limit int) ([]*model.Post, error) {
	var posts []*model.Post
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, fmt.Errorf("获取最新帖子失败: %w", err)
	}

	return posts, nil
}

// ListRecommended 获取推荐帖子（未实现）
func (r *postRepository) ListRecommended(ctx context.Context, userID uint, limit int) ([]*model.Post, error) {
	// 这是一个简化的推荐算法实现
	// 实际项目中可以根据用户兴趣、历史行为等实现更复杂的推荐逻辑

	var posts []*model.Post
	err := r.db.WithContext(ctx).
		Joins("JOIN user_interests ON user_interests.category_id = posts.category_id").
		Where("user_interests.user_id = ?", userID).
		Order("posts.like_count DESC, posts.comment_count DESC").
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, fmt.Errorf("获取推荐帖子失败: %w", err)
	}

	return posts, nil
}

// Exists 检查帖子是否存在
func (r *postRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("检查帖子存在失败: %w", err)
	}

	return count > 0, nil
}
