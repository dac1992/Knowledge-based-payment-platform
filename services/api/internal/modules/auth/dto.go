package auth

import "github.com/yanmai/api/internal/modules/user"

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Email    string `json:"email" binding:"omitempty,email,max=128"`
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
}

// LoginRequest 登录请求；Login 支持登录名或手机号。
type LoginRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=1,max=64"`
}

// RefreshRequest 刷新令牌请求。
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// TokenResponse 登录/注册/刷新统一返回。
type TokenResponse struct {
	TokenType    string       `json:"tokenType"`
	AccessToken  string       `json:"accessToken"`
	ExpiresIn    int64        `json:"expiresIn"`
	RefreshToken string       `json:"refreshToken"`
	User         user.Profile `json:"user"`
}
