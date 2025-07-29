package auth

type JWTResult struct {
	ID           uint   `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Platform     string `json:"platform"`
	DeviceID     string `json:"device_id"`
}
