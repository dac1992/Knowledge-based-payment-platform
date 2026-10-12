package user

import (
	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/internal/middleware"
	"github.com/yanmai/api/internal/modules"
)

// Module 用户模块的路由装配。
type Module struct {
	handler *Handler
	auth    gin.HandlerFunc
}

// NewModule 构造用户模块。
func NewModule(handler *Handler, requireAuth gin.HandlerFunc) modules.Module {
	return &Module{handler: handler, auth: requireAuth}
}

// Name 模块名。
func (m *Module) Name() string { return "user" }

// RegisterRoutes 挂载用户相关路由。
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	personal := rg.Group("/users")
	personal.GET("/me", m.auth, m.handler.Me)
	personal.PATCH("/me", m.auth, m.handler.UpdateMe)

	admin := rg.Group("/admin/users")
	admin.Use(m.auth, middleware.RequireRole(RoleAdmin))
	admin.GET("", m.handler.List)
	admin.GET("/:id", m.handler.GetByID)
}
