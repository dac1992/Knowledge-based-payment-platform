package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/httpx"
	"github.com/yanmai/api/pkg/jwt"
)

// CtxUserKey 当前登录主体在 gin.Context 中的键名。
const CtxUserKey = "auth.principal"

// Principal 当前请求的登录主体。
type Principal struct {
	UserID   uint64
	Username string
	Role     string
}

// RequireAuth 校验 Authorization: Bearer <access_token>，通过后把主体写入上下文。
func RequireAuth(m *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := bearerToken(c)
		if raw == "" {
			httpx.Fail(c, errs.New(errs.CodeTokenMissing))
			return
		}

		claims, err := m.ParseAccess(raw)
		if err != nil {
			httpx.Fail(c, err)
			return
		}

		c.Set(CtxUserKey, &Principal{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		})
		c.Next()
	}
}

// RequireRole 在 RequireAuth 之后使用，限制只有指定角色可访问。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		p, ok := CurrentUser(c)
		if !ok {
			httpx.Fail(c, errs.New(errs.CodeUnauthorized))
			return
		}
		if _, ok := allowed[p.Role]; !ok {
			httpx.Fail(c, errs.New(errs.CodePermissionDenied).
				WithMessage("角色 %s 无权访问该接口", p.Role))
			return
		}
		c.Next()
	}
}

// CurrentUser 读取当前登录主体。
func CurrentUser(c *gin.Context) (*Principal, bool) {
	v, exists := c.Get(CtxUserKey)
	if !exists {
		return nil, false
	}
	p, ok := v.(*Principal)
	return p, ok
}

// MustCurrentUser 读取当前登录主体，缺失时返回零值（调用方需自行保证已在 RequireAuth 之后）。
func MustCurrentUser(c *gin.Context) *Principal {
	if p, ok := CurrentUser(c); ok {
		return p
	}
	return &Principal{}
}

func bearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
