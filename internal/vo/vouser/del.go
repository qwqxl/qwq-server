package vouser

type DelRequest struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Platform string `json:"platform"`
	DeviceID string `json:"device_id"`
}

type DelResponse struct {
	ID uint `json:"id"`
}
