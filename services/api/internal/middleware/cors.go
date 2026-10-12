package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/pkg/httpx"
)

// CORS 跨域中间件。
//
// origins 为空表示允许任意来源（仅建议 local 环境使用）；
// 否则只回显白名单内的 Origin，并开启凭证传递。
func CORS(origins []string) gin.HandlerFunc {
	allowAll := len(origins) == 0
	whitelist := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		whitelist[strings.ToLower(strings.TrimSpace(o))] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		switch {
		case origin == "":
			// 非浏览器请求（curl / 服务间调用），无需跨域头
		case allowAll:
			c.Header("Access-Control-Allow-Origin", "*")
		default:
			if _, ok := whitelist[strings.ToLower(origin)]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Vary", "Origin")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With,"+httpx.TraceIDHeader)
		c.Header("Access-Control-Expose-Headers", httpx.TraceIDHeader)
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
