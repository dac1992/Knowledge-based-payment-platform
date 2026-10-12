package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/internal/middleware"
	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
)

// Handler 用户模块 HTTP 层。
type Handler struct{ svc Service }

// NewHandler 构造 Handler。
func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

type UpdateProfileRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=32"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=512"`
	Email    *string `json:"email" binding:"omitempty,email,max=128"`
	Phone    *string `json:"phone" binding:"omitempty,max=20"`
}

type ListUsersRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=64"`
	Role     string `form:"role" binding:"omitempty,oneof=admin lecturer student"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1"`
}

// Me 当前登录用户资料
//
// @Summary      当前用户资料
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  httpx.Envelope{data=Profile}
// @Failure      401  {object}  httpx.Envelope
// @Router       /users/me [get]
func (h *Handler) Me(c *gin.Context) {
	p := middleware.MustCurrentUser(c)
	profile, err := h.svc.Profile(c.Request.Context(), p.UserID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, profile)
}

// UpdateMe 修改个人资料
//
// @Summary      修改个人资料
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpdateProfileRequest  true  "待修改字段，未传的字段保持不变"
// @Success      200   {object}  httpx.Envelope{data=Profile}
// @Failure      400   {object}  httpx.Envelope
// @Failure      409   {object}  httpx.Envelope  "手机号已被占用"
// @Router       /users/me [patch]
func (h *Handler) UpdateMe(c *gin.Context) {
	var req UpdateProfileRequest
	if !httpx.Bind(c, &req) {
		return
	}

	p := middleware.MustCurrentUser(c)
	// UpdateProfileRequest 与 UpdateInput 字段一致，直接转换
	profile, err := h.svc.UpdateProfile(c.Request.Context(), p.UserID, UpdateInput(req))
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, profile)
}

// List 用户列表
//
// @Summary      用户列表（分页）
// @Description  仅管理员可访问，支持关键字、角色、状态筛选
// @Tags         admin-user
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int     false  "页码，从 1 开始"        default(1)
// @Param        pageSize  query     int     false  "每页条数，1-100"        default(20)
// @Param        keyword   query     string  false  "登录名 / 昵称 / 手机号 模糊匹配"
// @Param        role      query     string  false  "角色"                   Enums(admin, lecturer, student)
// @Param        status    query     int     false  "状态"                   Enums(0, 1)
// @Success      200       {object}  httpx.Envelope{data=httpx.PageEnvelope}
// @Failure      401       {object}  httpx.Envelope
// @Failure      403       {object}  httpx.Envelope
// @Router       /admin/users [get]
func (h *Handler) List(c *gin.Context) {
	var req ListUsersRequest
	if !httpx.BindQuery(c, &req) {
		return
	}

	list, total, err := h.svc.List(c.Request.Context(), ListQuery(req))
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Page(c, list, total, req.Page, req.PageSize)
}

// GetByID 用户详情
//
// @Summary      用户详情
// @Tags         admin-user
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "用户 ID"
// @Success      200  {object}  httpx.Envelope{data=Profile}
// @Failure      401  {object}  httpx.Envelope
// @Failure      403  {object}  httpx.Envelope
// @Failure      404  {object}  httpx.Envelope
// @Router       /admin/users/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		httpx.Fail(c, errs.New(errs.CodeInvalidParam).WithMessage("用户 ID 非法"))
		return
	}

	profile, err := h.svc.Profile(c.Request.Context(), id)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, profile)
}
