// Package middleware 提供横切关注点的 Gin 中间件。
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/yanmai/api/pkg/httpx"
)

// TraceID 为每个请求生成（或透传）链路 ID，写入上下文与响应头。
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := strings.TrimSpace(c.GetHeader(httpx.TraceIDHeader))
		if tid == "" {
			tid = strings.ReplaceAll(uuid.NewString(), "-", "")
		}
		c.Set(httpx.TraceIDKey, tid)
		c.Header(httpx.TraceIDHeader, tid)
		c.Next()
	}
}
