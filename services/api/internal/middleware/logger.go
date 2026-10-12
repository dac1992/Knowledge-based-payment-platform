package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yanmai/api/pkg/httpx"
)

// Logger 输出结构化访问日志。
func Logger(log *zap.Logger) gin.HandlerFunc {
	l := log.Named("http")
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		fields := []zap.Field{
			zap.String("traceId", httpx.TraceID(c)),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
			zap.Int("size", c.Writer.Size()),
		}
		if query != "" {
			fields = append(fields, zap.String("query", query))
		}
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		switch {
		case c.Writer.Status() >= 500:
			l.Error("请求处理失败", fields...)
		case c.Writer.Status() >= 400:
			l.Warn("请求被拒绝", fields...)
		default:
			l.Info("请求完成", fields...)
		}
	}
}
