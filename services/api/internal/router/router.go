// Package router 负责 HTTP 引擎的组装：全局中间件、健康检查、业务模块挂载。
//
// Phase 1 采用「手动构造注入」显式装配依赖，链路清晰、易于排查；
// 模块数量增长后（Phase 3+）再引入 google/wire 生成注入代码。
package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/yanmai/api/internal/config"
	"github.com/yanmai/api/internal/middleware"
	"github.com/yanmai/api/internal/modules"
	"github.com/yanmai/api/internal/modules/auth"
	"github.com/yanmai/api/internal/modules/circle"
	"github.com/yanmai/api/internal/modules/course"
	"github.com/yanmai/api/internal/modules/finance"
	"github.com/yanmai/api/internal/modules/learning"
	"github.com/yanmai/api/internal/modules/live"
	"github.com/yanmai/api/internal/modules/marketing"
	"github.com/yanmai/api/internal/modules/member"
	"github.com/yanmai/api/internal/modules/order"
	"github.com/yanmai/api/internal/modules/payment"
	"github.com/yanmai/api/internal/modules/stats"
	"github.com/yanmai/api/internal/modules/user"
	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
	"github.com/yanmai/api/pkg/jwt"
	"github.com/yanmai/api/pkg/queue"
)

// Deps 路由装配所需的依赖集合。
type Deps struct {
	Config *config.Config
	DB     *gorm.DB
	Logger *zap.Logger
	Queue  *queue.Client
}

// New 组装并返回 Gin 引擎。
func New(deps Deps) *gin.Engine {
	cfg, log := deps.Config, deps.Logger

	if !cfg.App.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(
		middleware.TraceID(),
		middleware.Logger(log),
		middleware.Recovery(log),
		middleware.CORS(cfg.HTTP.CORSOrigins),
		middleware.RateLimit(cfg.HTTP.RateLimitRPS, cfg.HTTP.RateLimitBurst),
	)

	r.NoRoute(func(c *gin.Context) {
		httpx.Fail(c, errs.New(errs.CodeNotFound).
			WithMessage("接口不存在: %s %s", c.Request.Method, c.Request.URL.Path))
	})
	r.NoMethod(func(c *gin.Context) {
		httpx.Fail(c, errs.New(errs.CodeMethodNotAllow).
			WithMessage("请求方法不被支持: %s %s", c.Request.Method, c.Request.URL.Path))
	})

	registerHealth(r, cfg, deps.DB)

	tokenMgr := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	requireAuth := middleware.RequireAuth(tokenMgr)

	// ---- 依赖装配 ----
	userRepo := user.NewRepository(deps.DB)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	authSvc := auth.NewService(userSvc, tokenMgr)
	authHandler := auth.NewHandler(authSvc)

	// ---- 模块注册（顺序即挂载顺序）----
	registry := modules.NewRegistry(
		auth.NewModule(authHandler, requireAuth),
		user.NewModule(userHandler, requireAuth),
		course.New(),    // Phase 3
		order.New(),     // Phase 5
		payment.New(),   // Phase 5
		learning.New(),  // Phase 6
		marketing.New(), // Phase 7
		member.New(),    // Phase 8
		live.New(),      // Phase 9
		circle.New(),    // Phase 11
		finance.New(),   // Phase 12
		stats.New(),     // Phase 12
	)

	api := r.Group("/api/v1")
	api.GET("/meta/modules", func(c *gin.Context) {
		httpx.OK(c, gin.H{
			"service": cfg.App.Name,
			"version": config.Version,
			"env":     cfg.App.Env,
			"modules": registry.Names(),
			"queue":   deps.Queue != nil && deps.Queue.Enabled(),
			"time":    time.Now().Format(time.RFC3339),
		})
	})
	registry.Mount(api)

	log.Info("路由装配完成", zap.Int("modules", len(registry.Names())), zap.Strings("moduleNames", registry.Names()))
	return r
}
