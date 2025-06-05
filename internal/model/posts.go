package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 帖子模型
type Post struct {
	gorm.Model
	// ID            uint64     `gorm:"primaryKey;autoIncrement;comment:帖子ID"`
	AuthorUID     uint64     `gorm:"not null;column:author_id;comment:作者ID"`
	Author        User       `gorm:"foreignKey:AuthorUID;references:ID"`
	Mod           string     `gorm:"size:1024;not null;comment:内容模型"`
	Title         string     `gorm:"size:1024;comment:标题"`
	Content       string     `gorm:"type:longtext;comment:内容"`
	Status        string     `gorm:"type:enum('draft','published','pending','trash');default:'draft';comment:状态"`
	IsSticky      bool       `gorm:"default:false;comment:是否置顶"`
	CommentStatus string     `gorm:"type:enum('open','closed');default:'open';comment:评论状态"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;comment:更新时间"`
	PublishedAt   *time.Time `gorm:"comment:发布时间"`
	IsDeleted     bool       `gorm:"default:false;comment:删除标记"`
}

// TableName 帖子模型
func (p *Post) TableName() string {
	return "posts"
}
