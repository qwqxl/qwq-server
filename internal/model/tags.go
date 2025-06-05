package model

// Tag 标签模型
type Tag struct {
	TagID uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"size:50;unique;not null"`
	Posts []Post `gorm:"many2many:post_tags;"`
}

// tag 标签模型表名
func (t *Tag) TableName() string {
	return "tags"
}
