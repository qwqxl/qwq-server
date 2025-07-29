package authx

// LoginInput 登录时所需的输入参数
type LoginInput struct {
	UserID     string // 用户 ID
	Platform   string // 平台标识
	DeviceSign string // 设备标识
}

// TokenPair 包含 AccessToken 和 RefreshToken
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}