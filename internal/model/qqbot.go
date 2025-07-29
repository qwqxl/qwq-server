package model

import "encoding/json"

// 消息内容项
type BotMessageItem struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"` // 或根据 CQCode 类型更精细定义
}

// 发送者信息
type BotSender struct {
	UserID   json.Number `json:"user_id"`
	Nickname string      `json:"nickname"`
	// 可扩展字段：Card, Role, Level 等
}

// 主请求体
type BotMessageRequest struct {
	PostType   string           `json:"post_type"`
	MessageID  json.Number      `json:"message_id"`
	GroupID    json.Number      `json:"group_id"`
	UserID     json.Number      `json:"user_id"`
	RawMessage string           `json:"raw_message"`
	Message    []BotMessageItem `json:"message"`
	Sender     BotSender        `json:"sender"`
}
