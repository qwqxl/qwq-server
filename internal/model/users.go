package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	//ID       uint   `gorm:"primaryKey;autoIncrement;comment:用户唯一标识符" json:"id"`
	Username string `gorm:"type:varchar(128);uniqueIndex;not null;comment:用户名" json:"username"`
	Nickname string `gorm:"type:varchar(1024);default:'新用户';comment:用户昵称" json:"nickname"`
	Email    string `gorm:"type:varchar(128);uniqueIndex;comment:邮箱地址" json:"email"`

	Password          string     `gorm:"type:varchar(1024);not null;comment:密码" json:"-"`
	PasswordHash      string     `gorm:"type:varchar(1024);not null;comment:密码哈希" json:"-"`
	PasswordSalt      string     `gorm:"type:varchar(1024);not null;comment:密码盐值" json:"-"`
	Iterations        int        `gorm:"default:10000;comment:哈希迭代次数" json:"-"`
	FailedAttempts    int        `gorm:"default:0;comment:登录失败次数" json:"failed_attempts"`
	LastLoginAt       *time.Time `gorm:"comment:上次登录时间;default:NULL" json:"last_login_at,omitempty"`
	LastFailedAttempt *time.Time `gorm:"comment:上次登录失败时间;default:NULL" json:"last_failed_attempt,omitempty"`
	Perms             uint64     `gorm:"type:BIGINT UNSIGNED;default:0;comment:权限位掩码" json:"perms"`
	Status            uint8      `gorm:"default:1;comment:状态 1=正常" json:"status"`
	CreatedAt         time.Time  `gorm:"autoCreateTime;comment:注册时间" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	IpAddress         string     `gorm:"type:varchar(1024);comment:用户IP" json:"ip_address"`
}

// RegisterAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:注册时间"`

// TableName table name
func (u *User) TableName() string {
	return "users"
}
