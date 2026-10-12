package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/internal/modules"
)

// Module 认证模块的路由装配。
type Module struct {
	handler *Handler
	auth    gin.HandlerFunc
}

// NewModule 构造认证模块。
func NewModule(handler *Handler, requireAuth gin.HandlerFunc) modules.Module {
	return &Module{handler: handler, auth: requireAuth}
}

// Name 模块名。
func (m *Module) Name() string { return "auth" }

// RegisterRoutes 挂载认证相关路由。
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/register", m.handler.Register)
	g.POST("/login", m.handler.Login)
	g.POST("/refresh", m.handler.Refresh)
	g.GET("/me", m.auth, m.handler.Me)
	g.POST("/logout", m.auth, m.handler.Logout)
}
