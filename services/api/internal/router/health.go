package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yanmai/api/internal/config"
	"github.com/yanmai/api/internal/database"
	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
	"github.com/yanmai/api/pkg/response"
)

// registerHealth 注册健康检查端点。
//
//	GET /health        —— 综合健康（含数据库探测），永远 200，用 status 字段表达降级
//	GET /health/live   —— 存活探针，进程在即 200
//	GET /health/ready  —— 就绪探针，数据库不可用时返回 503（供负载均衡摘流）
func registerHealth(r *gin.Engine, cfg *config.Config, db *gorm.DB) {
	r.GET("/health", func(c *gin.Context) {
		httpx.OK(c, buildHealth(cfg, db))
	})

	r.GET("/health/live", func(c *gin.Context) {
		httpx.OK(c, gin.H{"status": "alive", "time": time.Now().Format(time.RFC3339)})
	})

	r.GET("/health/ready", func(c *gin.Context) {
		body := buildHealth(cfg, db)
		if body["status"] != "ok" {
			c.JSON(http.StatusServiceUnavailable, response.New(
				int(errs.CodeInternal), "服务暂不可用", body, httpx.TraceID(c),
			))
			return
		}
		httpx.OK(c, body)
	})
}

func buildHealth(cfg *config.Config, db *gorm.DB) gin.H {
	dbStatus := "up"
	dbError := ""
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := database.Ping(ctx, db); err != nil {
		dbStatus = "down"
		dbError = err.Error()
	}

	overall := "ok"
	if dbStatus != "up" {
		overall = "degraded"
	}

	// 先组装 database 子对象，避免对其做类型断言后再写入
	dbInfo := gin.H{
		"status":  dbStatus,
		"latency": time.Since(start).String(),
	}
	if dbError != "" {
		dbInfo["error"] = dbError
	}

	return gin.H{
		"status":   overall,
		"service":  cfg.App.Name,
		"env":      cfg.App.Env,
		"version":  config.Version,
		"time":     time.Now().Format(time.RFC3339),
		"uptime":   time.Since(start).String(),
		"database": dbInfo,
	}
}
