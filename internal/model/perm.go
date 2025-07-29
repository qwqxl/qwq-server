package model

import (
	"time"
)

type UserID string
type RoleID string
type PermissionKey string

// 权限模型
type Permission struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         PermissionKey `gorm:"type:varchar(128);uniqueIndex;not null" json:"key"`
	Description string        `gorm:"type:varchar(1024)" json:"description"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// 角色模型
type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"type:varchar(128);uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:varchar(1024)" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

// 用户角色关联模型
type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
}

// 角色权限关联模型
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
}
