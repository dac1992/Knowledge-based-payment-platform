package jwt

import (
	"testing"
	"time"

	"github.com/yanmai/api/pkg/errs"
)

func newTestManager(accessTTL, refreshTTL time.Duration) *Manager {
	return NewManager("unit-test-secret-key-at-least-32-chars", "yanmai-api", accessTTL, refreshTTL)
}

func TestGenerateAndParseAccess(t *testing.T) {
	m := newTestManager(time.Hour, 24*time.Hour)

	token, err := m.GenerateAccess(42, "tester", "admin")
	if err != nil {
		t.Fatalf("签发 Access Token 失败: %v", err)
	}
	if token.Value == "" || token.JTI == "" {
		t.Fatal("签发结果不完整")
	}

	claims, err := m.ParseAccess(token.Value)
	if err != nil {
		t.Fatalf("解析 Access Token 失败: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "tester" || claims.Role != "admin" {
		t.Fatalf("声明内容不符: %+v", claims)
	}
	if claims.TokenType != TypeAccess {
		t.Fatalf("令牌类型应为 %s，实际 %s", TypeAccess, claims.TokenType)
	}
}

func TestParseRejectsRefreshTokenAsAccess(t *testing.T) {
	m := newTestManager(time.Hour, 24*time.Hour)

	refresh, err := m.GenerateRefresh(7, "tester", "student")
	if err != nil {
		t.Fatalf("签发 Refresh Token 失败: %v", err)
	}

	if _, err := m.ParseAccess(refresh.Value); !errs.Is(err, errs.CodeTokenTypeMismatch) {
		t.Fatalf("用 Refresh Token 访问应报类型不匹配，实际: %v", err)
	}
}

func TestParseExpiredToken(t *testing.T) {
	// Access TTL 为负数 → 签发即过期
	m := newTestManager(-time.Minute, time.Hour)

	token, err := m.GenerateAccess(1, "tester", "student")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if _, err := m.ParseAccess(token.Value); !errs.Is(err, errs.CodeTokenExpired) {
		t.Fatalf("应报令牌过期，实际: %v", err)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	issuer := newTestManager(time.Hour, 24*time.Hour)
	verifier := NewManager("another-secret-key-also-32-characters", "yanmai-api", time.Hour, 24*time.Hour)

	token, err := issuer.GenerateAccess(1, "tester", "student")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if _, err := verifier.ParseAccess(token.Value); !errs.Is(err, errs.CodeTokenInvalid) {
		t.Fatalf("密钥不匹配应报令牌无效，实际: %v", err)
	}
}

func TestParseRejectsWrongIssuer(t *testing.T) {
	issuer := newTestManager(time.Hour, 24*time.Hour)
	verifier := NewManager("unit-test-secret-key-at-least-32-chars", "other-issuer", time.Hour, 24*time.Hour)

	token, err := issuer.GenerateAccess(1, "tester", "student")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if _, err := verifier.ParseAccess(token.Value); !errs.Is(err, errs.CodeTokenInvalid) {
		t.Fatalf("签发者不匹配应报令牌无效，实际: %v", err)
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	m := newTestManager(time.Hour, 24*time.Hour)

	if _, err := m.ParseAccess("not-a-jwt"); !errs.Is(err, errs.CodeTokenInvalid) {
		t.Fatalf("非法令牌应报令牌无效，实际: %v", err)
	}
	if _, err := m.ParseAccess(""); !errs.Is(err, errs.CodeTokenMissing) {
		t.Fatalf("空令牌应报缺少凭证，实际: %v", err)
	}
}
