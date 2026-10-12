package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/internal/middleware"
	"github.com/yanmai/api/pkg/httpx"
)

// Handler 认证模块 HTTP 层。
type Handler struct{ svc Service }

// NewHandler 构造 Handler。
func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

// Register 用户注册
//
// @Summary      注册
// @Description  创建学员账号并直接返回登录态（Access + Refresh 双 Token）
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "注册参数"
// @Success      201   {object}  httpx.Envelope{data=TokenResponse}
// @Failure      400   {object}  httpx.Envelope
// @Failure      409   {object}  httpx.Envelope
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if !httpx.Bind(c, &req) {
		return
	}

	result, err := h.svc.Register(c.Request.Context(), RegisterInput{
		Username: req.Username,
		Password: req.Password,
		Phone:    req.Phone,
		Email:    req.Email,
		Nickname: req.Nickname,
	})
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, toTokenResponse(result))
}

// Login 用户登录
//
// @Summary      登录
// @Description  支持使用登录名或手机号登录，返回 Access + Refresh 双 Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "登录参数"
// @Success      200   {object}  httpx.Envelope{data=TokenResponse}
// @Failure      400   {object}  httpx.Envelope
// @Failure      401   {object}  httpx.Envelope  "账号或密码不正确"
// @Failure      403   {object}  httpx.Envelope  "账号被禁用"
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !httpx.Bind(c, &req) {
		return
	}

	result, err := h.svc.Login(c.Request.Context(), LoginInput{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, toTokenResponse(result))
}

// Refresh 刷新令牌
//
// @Summary      刷新令牌
// @Description  用 Refresh Token 换取新的令牌对；会回源校验账号仍存在且可用
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshRequest  true  "刷新参数"
// @Success      200   {object}  httpx.Envelope{data=TokenResponse}
// @Failure      401   {object}  httpx.Envelope
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if !httpx.Bind(c, &req) {
		return
	}

	result, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, toTokenResponse(result))
}

// Me 当前登录账号
//
// @Summary      当前登录账号
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  httpx.Envelope
// @Failure      401  {object}  httpx.Envelope
// @Router       /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	p := middleware.MustCurrentUser(c)
	profile, err := h.svc.Profile(c.Request.Context(), p.UserID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, profile)
}

// Logout 登出
//
// @Summary      登出
// @Description  Phase 1 为无状态登出，由客户端丢弃令牌；服务端黑名单在 Phase 2 引入 Redis 后实现
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  httpx.Envelope
// @Failure      401  {object}  httpx.Envelope
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	httpx.OK(c, gin.H{
		"revoked": false,
		"hint":    "服务端暂无令牌黑名单，请客户端清除本地令牌",
	})
}

func toTokenResponse(r *Result) TokenResponse {
	return TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  r.AccessToken,
		ExpiresIn:    r.ExpiresIn,
		RefreshToken: r.RefreshToken,
		User:         r.User,
	}
}
