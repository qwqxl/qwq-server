package post

import (
	"context"
	"errors"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
)

type CreateRequest struct {
	model.Post
	//AuthorID      uint64     `gorm:"not null;column:author_id;comment:作者ID" json:"author_id"`
	//Author        model.User `gorm:"foreignKey:AuthorID;references:ID" json:"author"`
	//Mod           string     `gorm:"size:1024;not null;comment:内容模型" json:"mod"`
	//Title         string     `gorm:"size:1024;comment:标题" json:"title"`
	//Content       string     `gorm:"type:longtext;comment:内容" json:"content"`
	//Status        string     `gorm:"type:enum('draft','published','pending','trash');default:'draft';comment:状态" json:"status"`
	//IsSticky      bool       `gorm:"default:false;comment:是否置顶" json:"is_sticky"`
	//CommentStatus string     `gorm:"type:enum('open','closed');default:'open';comment:评论状态" json:"comment_status"`
	//CreatedAt     time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	//UpdatedAt     time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	//PublishedAt   time.Time  `gorm:"comment:发布时间" json:"published_at"`
	//IsDeleted     bool       `gorm:"default:false;comment:删除标记" json:"is_deleted"`
}

type CreateResponse struct {
	model.Post
}

// Create 创建文章
func (s *Service) Create(req CreateRequest) (interface{}, error) {
	res := CreateResponse{}
	//
	postRepo, err := repository.NewPostRepository()
	if err != nil {
		return nil, errors.New("获取数据库连接失败: " + err.Error())
	}

	//mPost := &model.Post{}
	if err = postRepo.Create(context.Background(), &req.Post); err != nil {
		return nil, errors.New("创建文章失败: " + err.Error())
	}

	res.Post = req.Post
	return res, nil
}
