// Package jwt 封装 Access / Refresh 双 Token 的签发与校验。
//
// 约定：
//   - Token 中显式携带 typ 声明（access / refresh），校验时强制比对类型，
//     避免用 refresh token 直接访问业务接口。
//   - 解析失败统一翻译为 errs 业务错误码，调用方无需关心 jwt 库细节。
package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/yanmai/api/pkg/errs"
)

// Token 类型。
const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

// Claims 自定义声明。
type Claims struct {
	UserID    uint64 `json:"uid"`
	Username  string `json:"uname"`
	Role      string `json:"role"`
	TokenType string `json:"typ"`
	jwtlib.RegisteredClaims
}

// Token 签发结果。
type Token struct {
	Value     string
	ExpiresAt time.Time
	JTI       string
}

// Manager Token 管理器。
type Manager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewManager 构造 Token 管理器。
func NewManager(secret, issuer string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// AccessTTL 返回 Access Token 有效期。
func (m *Manager) AccessTTL() time.Duration { return m.accessTTL }

// RefreshTTL 返回 Refresh Token 有效期。
func (m *Manager) RefreshTTL() time.Duration { return m.refreshTTL }

// GenerateAccess 签发 Access Token。
func (m *Manager) GenerateAccess(userID uint64, username, role string) (*Token, error) {
	return m.generate(userID, username, role, TypeAccess, m.accessTTL)
}

// GenerateRefresh 签发 Refresh Token。
func (m *Manager) GenerateRefresh(userID uint64, username, role string) (*Token, error) {
	return m.generate(userID, username, role, TypeRefresh, m.refreshTTL)
}

func (m *Manager) generate(userID uint64, username, role, typ string, ttl time.Duration) (*Token, error) {
	now := time.Now()
	exp := now.Add(ttl)
	jti := uuid.NewString()

	claims := Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: typ,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ID:        jti,
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(exp),
		},
	}

	signed, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternal, err).WithMessage("签发令牌失败")
	}
	return &Token{Value: signed, ExpiresAt: exp, JTI: jti}, nil
}

// ParseAccess 解析并校验 Access Token。
func (m *Manager) ParseAccess(raw string) (*Claims, error) { return m.parse(raw, TypeAccess) }

// ParseRefresh 解析并校验 Refresh Token。
func (m *Manager) ParseRefresh(raw string) (*Claims, error) { return m.parse(raw, TypeRefresh) }

func (m *Manager) parse(raw, want string) (*Claims, error) {
	if raw == "" {
		return nil, errs.New(errs.CodeTokenMissing)
	}

	claims := &Claims{}
	_, err := jwtlib.ParseWithClaims(raw, claims, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名算法: %v", t.Header["alg"])
		}
		return m.secret, nil
	}, jwtlib.WithIssuer(m.issuer), jwtlib.WithExpirationRequired())

	if err != nil {
		switch {
		case errors.Is(err, jwtlib.ErrTokenExpired):
			return nil, errs.Wrap(errs.CodeTokenExpired, err)
		case errors.Is(err, jwtlib.ErrTokenMalformed),
			errors.Is(err, jwtlib.ErrTokenSignatureInvalid),
			errors.Is(err, jwtlib.ErrTokenNotValidYet),
			errors.Is(err, jwtlib.ErrTokenInvalidIssuer):
			return nil, errs.Wrap(errs.CodeTokenInvalid, err)
		default:
			return nil, errs.Wrap(errs.CodeTokenInvalid, err)
		}
	}

	if claims.TokenType != want {
		return nil, errs.New(errs.CodeTokenTypeMismatch).
			WithMessage("需要 %s 令牌，实际为 %s", want, claims.TokenType)
	}
	return claims, nil
}
