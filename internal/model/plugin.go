package model

import (
	"gorm.io/gorm"
)

// 插件状态常量
const (
	PluginDisabled = iota
	PluginActive
	PluginError
)

// Plugin 插件
type Plugin struct {
	gorm.Model
	UUID      string `gorm:"uniqueIndex;size:36"` // 插件唯一标识
	Name      string `gorm:"size:100;not null"`   // 插件名称
	Type      string `gorm:"size:20;not null"`    // builtin/grpc/tcp/http
	Status    int    `gorm:"default:0"`           // 状态
	Config    string `gorm:"type:text"`           // JSON配置
	Endpoint  string `gorm:"size:255"`            // gRPC地址
	Version   string `gorm:"size:20"`             // 版本号
	Signature string `gorm:"size:64"`             // 安全签名
}

/*
CREATE TABLE plugins (
    id          VARCHAR(36) PRIMARY KEY,      -- UUID
    name        VARCHAR(50) UNIQUE NOT NULL,  -- 插件名
    rpc_address VARCHAR(100) NOT NULL,        -- gRPC地址(ip:port)
    status      ENUM('active','disabled') DEFAULT 'active',
    config      JSON,                         -- 插件配置
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
*/
