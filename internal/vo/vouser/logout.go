package vouser

type UserLogoutRequest struct {
	ID       uint   `json:"id"`
	DeviceID string `json:"device_id"`
	Platform string `json:"platform"`
}

type UserLogoutResponse struct {
	ID uint `json:"id"`
}
