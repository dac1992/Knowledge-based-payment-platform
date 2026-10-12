package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
)

// RateLimit 基于内存令牌桶的按 IP 限流。
//
// 说明：单实例内存限流，适合 Phase 1；多实例部署时应替换为 Redis 分布式限流。
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	if rps <= 0 || burst <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	var (
		mu       sync.Mutex
		visitors = make(map[string]*rate.Limiter)
		lastGC   = time.Now()
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, ok := visitors[ip]
		if !ok {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			visitors[ip] = limiter
		}
		if time.Since(lastGC) > 5*time.Minute {
			for k, v := range visitors {
				if v.Tokens() >= float64(burst) {
					delete(visitors, k)
				}
			}
			lastGC = time.Now()
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.Header("Retry-After", "1")
			httpx.Fail(c, errs.New(errs.CodeTooManyRequest))
			return
		}
		c.Next()
	}
}
