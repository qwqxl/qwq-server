package user

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
	"time"
)

type DeleteRequest struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Platform string `json:"platform"`
	DeviceID string `json:"device_id"`
}

type DeleteResponse struct {
	ID       uint      `json:"id"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	DeleteAt time.Time `json:"delete_at,omitempty"`
	CreateAt time.Time `json:"create_at,omitempty"`
}

// Delete 删除用户（安全版本）
func (s *Service) Delete(ctx context.Context, req DeleteRequest) (interface{}, error) {
	if req.ID == 0 && req.Username == "" && req.Email == "" {
		return nil, fmt.Errorf("请选择要删除的用户（提供 ID、用户名或邮箱）")
	}

	repo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	var user *model.User

	// 按优先顺序选择查询方式
	if req.ID != 0 {
		user, err = repo.FindByID(ctx, req.ID)
	} else if req.Username != "" {
		user, err = repo.FindByUsername(ctx, req.Username)
	} else if req.Email != "" {
		user, err = repo.FindByEmail(ctx, req.Email)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("系统错误：用户对象为空")
	}

	if user.DeletedAt.Valid {
		return DeleteResponse{
			ID:       user.ID,
			Username: user.Username,
			DeleteAt: user.DeletedAt.Time, // user.DeletedAt.Time.Format(time.RFC3339) to string
		}, fmt.Errorf("用户已被删除")
	}

	// 删除用户
	if err := repo.Delete(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("删除用户失败: %w", err)
	}

	// 成功返回
	return DeleteResponse{ID: user.ID, Username: user.Username, CreateAt: user.CreatedAt}, nil
}
