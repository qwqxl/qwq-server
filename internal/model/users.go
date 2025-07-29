package model

import (
	"gorm.io/gorm"
	"time"
)

type UserStatus uint8

const (
	// 正常状态
	UserStatusActive UserStatus = iota + 1 // 1: 账号正常启用中

	// 限制状态
	UserStatusDisabled            // 2: 账号被禁用（管理员操作）
	UserStatusLocked              // 3: 账号已锁定（如失败次数过多）
	UserStatusPendingVerification // 4: 等待验证（如注册邮箱未验证）
	UserStatusSuspended           // 5: 账号被暂时冻结（如封禁 X 天）
	UserStatusArchived            // 6: 账号已归档（非活跃用户）

	// 注销状态
	UserStatusDeleted // 7: 用户主动注销
	UserStatusBanned  // 8: 永久封禁（严重违规）
)

var UserStatusText = map[UserStatus]string{
	UserStatusActive:              "正常",
	UserStatusDisabled:            "已禁用",
	UserStatusLocked:              "已锁定",
	UserStatusPendingVerification: "待验证",
	UserStatusSuspended:           "已冻结",
	UserStatusArchived:            "已归档",
	UserStatusDeleted:             "已注销",
	UserStatusBanned:              "永久封禁",
}

type User struct {
	gorm.Model
	//ID       uint   `gorm:"primaryKey;autoIncrement;comment:用户唯一标识符" json:"id"`
	Username string `gorm:"type:varchar(128);uniqueIndex;not null;comment:用户名" json:"username"`
	Nickname string `gorm:"type:varchar(1024);default:'新用户';comment:用户昵称" json:"nickname"`
	Email    string `gorm:"type:varchar(128);uniqueIndex;comment:邮箱地址" json:"email"`

	Password          string     `gorm:"type:varchar(1024);not null;comment:密码" json:"password"`
	PasswordHash      string     `gorm:"type:varchar(1024);not null;comment:密码哈希" json:"-"`
	PasswordSalt      string     `gorm:"type:varchar(1024);not null;comment:密码盐值" json:"-"`
	Iterations        int        `gorm:"default:10000;comment:哈希迭代次数" json:"-"`
	FailedAttempts    int        `gorm:"default:0;comment:登录失败次数" json:"failed_attempts"`
	LastLoginAt       *time.Time `gorm:"comment:上次登录时间;default:NULL" json:"last_login_at,omitempty"`
	LastFailedAttempt *time.Time `gorm:"comment:上次登录失败时间;default:NULL" json:"last_failed_attempt,omitempty"`
	Perms             uint64     `gorm:"type:BIGINT UNSIGNED;default:0;comment:权限位掩码" json:"perms"`
	IpAddress         string     `gorm:"type:varchar(1024);comment:用户IP" json:"ip_address"`

	Status        uint8      `gorm:"default:1;comment:状态 1=正常" json:"status"`
	BanReason     string     `gorm:"type:varchar(1024);comment:封禁原因" json:"ban_reason,omitempty"`
	BanExpiresAt  *time.Time `gorm:"comment:封禁过期时间;default:NULL" json:"ban_expires_at,omitempty"`
	VerifyToken   string     `gorm:"type:varchar(128);comment:验证用Token" json:"-"`
	VerifyExpires *time.Time `gorm:"comment:验证Token过期时间" json:"-"`

	//CreatedAt         time.Time
	//UpdatedAt         time.Time
	//DeletedAt         *time.Time

	/* ----------- perm 扩展用户模型 ----------- */

	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

// RegisterAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:注册时间"`

// TableName table name
func (u *User) TableName() string {
	return "users"
}
