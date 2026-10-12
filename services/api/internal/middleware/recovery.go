package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
)

// Recovery 捕获 panic，记录堆栈并返回统一的 500 响应，避免进程被拖垮或细节外泄。
func Recovery(log *zap.Logger) gin.HandlerFunc {
	l := log.Named("recovery")
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				l.Error("请求处理发生 panic",
					zap.Any("panic", r),
					zap.String("traceId", httpx.TraceID(c)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Stack("stack"),
				)
				httpx.Fail(c, errs.New(errs.CodeInternal).WithCause(fmt.Errorf("panic: %v", r)))
			}
		}()

		c.Next()
	}
}
