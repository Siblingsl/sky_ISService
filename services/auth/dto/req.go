package dto

// AdminLoginRequest 登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// VerifyTokenRequest 用于验证 Token 请求
type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
