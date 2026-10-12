// Package modules 定义业务模块的统一契约。
//
// 架构约定（模块化单体）：
//   - 每个业务域一个子包（internal/modules/<domain>），内部按 handler / service /
//     repository / dto / model / router 六件套分层
//   - 模块只允许操作自己的数据表；跨模块协作必须通过对方的 Service 接口
//   - 本包只放契约，禁止 import 任何子包，避免循环依赖
package modules

import "github.com/gin-gonic/gin"

// Module 一个可挂载的业务模块。
type Module interface {
	// Name 模块标识，需与路由前缀一致，例如 "auth" / "user"。
	Name() string
	// RegisterRoutes 把模块自身的路由挂载到全局 API 路由组下。
	RegisterRoutes(rg *gin.RouterGroup)
}

// Registry 模块注册表，维护挂载顺序并对外暴露模块清单。
type Registry struct {
	modules []Module
}

// NewRegistry 构造注册表。
func NewRegistry(ms ...Module) *Registry {
	r := &Registry{modules: make([]Module, 0, len(ms))}
	for _, m := range ms {
		r.Add(m)
	}
	return r
}

// Add 注册模块（nil 会被忽略，便于按配置裁剪）。
func (r *Registry) Add(m Module) {
	if m == nil {
		return
	}
	r.modules = append(r.modules, m)
}

// All 返回全部模块。
func (r *Registry) All() []Module { return r.modules }

// Names 返回模块名列表。
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.modules))
	for _, m := range r.modules {
		names = append(names, m.Name())
	}
	return names
}

// Mount 按注册顺序把所有模块挂载到路由组。
func (r *Registry) Mount(rg *gin.RouterGroup) {
	for _, m := range r.modules {
		m.RegisterRoutes(rg)
	}
}
