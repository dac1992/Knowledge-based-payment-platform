package modules

import (
	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
)

// Placeholder 描述一个尚未实现的模块。
type Placeholder struct {
	name  string
	phase int
	title string
}

// NewPlaceholder 把一个占位模块实现为 Module：注册一个只读探针路由，
// 访问后返回 501 并说明计划交付阶段，便于联调期确认路由挂载是否正确。
func NewPlaceholder(name string, phase int, title string) Module {
	return &Placeholder{name: name, phase: phase, title: title}
}

// Name 模块名。
func (p *Placeholder) Name() string { return p.name }

// RegisterRoutes 注册占位路由。
func (p *Placeholder) RegisterRoutes(rg *gin.RouterGroup) {
	rg.Group("/"+p.name).GET("", func(c *gin.Context) {
		httpx.Fail(c, errs.New(errs.CodeNotImplemented).
			WithMessage("模块 %s（%s）计划在 Phase %d 交付，当前为骨架占位", p.name, p.title, p.phase))
	})
}
