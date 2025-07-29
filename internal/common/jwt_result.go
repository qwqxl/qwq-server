package common

import "gorm.io/gorm"

type JWTResult struct {
	gorm.Model
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Error        string `json:"error"`
	Platform     string `json:"platform"`
	DeviceID     string `json:"device_id"`
}
