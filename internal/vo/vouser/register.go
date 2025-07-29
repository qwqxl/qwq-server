package vouser

import "qwqserver/internal/model"

// 注册请求
type UserRegisterRequest struct {
	model.User
	EmailVerificationCode string `json:"email_verification_code"`
	ImgVerificationCode   string `json:"img_verification_code"`
	Password              string `json:"password"`
}

type UserRegisterResponse struct {
	model.User
	Password string `json:"-"`
}
