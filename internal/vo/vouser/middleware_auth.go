package vouser

type MiddlewareAuthRequest struct {
	DeviceID string `json:"device_id"`
	Email    string `json:"email"`
	ID       uint   `json:"id"`
	Platform string `json:"platform"`
	Username string `json:"username"`
}

type MiddlewareAuthJWTResponse struct {
	ID uint `json:"id"`
}
