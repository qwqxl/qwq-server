package vouser

type UserTokenData struct {
	ID           uint   `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Platform     string `json:"platform"`
	DeviceID     string `json:"device_id"`
}

type UserLoginRequest struct {
	ID       uint   `json:"id"`                 // 用户ID
	Username string `json:"username"`           // 用户名
	Email    string `json:"email"`              // 邮箱
	Password string `json:"password,omitempty"` // 密码

	Platform string `json:"platform"`  // 平台
	DeviceID string `json:"device_id"` // 设备ID
}
