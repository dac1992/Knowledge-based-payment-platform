package config

import (
	"strings"
	"testing"
	"time"
)

func baseConfig() *Config {
	return &Config{
		App: App{Name: "yanmai-api", Env: EnvLocal},
		HTTP: HTTP{
			Host: "0.0.0.0", Port: 8080,
			ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second,
			IdleTimeout: 60 * time.Second, ShutdownTimeout: 10 * time.Second,
			CORSOrigins:    []string{"http://localhost:5173"},
			RateLimitRPS:   50,
			RateLimitBurst: 100,
		},
		Database: Database{
			DSN: "postgres://u:p@127.0.0.1:5432/yanmai_dev?sslmode=disable",
		},
		JWT: JWT{
			Secret: "local-dev-secret-at-least-32-characters", Issuer: "yanmai-api",
			AccessTTL: 2 * time.Hour, RefreshTTL: 168 * time.Hour,
		},
		Log: Log{Level: "debug"},
	}
}

func TestValidateAcceptsValidConfig(t *testing.T) {
	if err := baseConfig().validate(); err != nil {
		t.Fatalf("合法配置不应报错: %v", err)
	}
}

func TestValidateRejectsMissingDatabaseURL(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.DSN = ""
	err := cfg.validate()
	if err == nil || !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Fatalf("缺少 DATABASE_URL 时应给出明确提示，实际: %v", err)
	}
}

func TestValidateRejectsNonPostgresDSN(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.DSN = "mysql://root@127.0.0.1:3306/yanmai"
	err := cfg.validate()
	if err == nil || !strings.Contains(err.Error(), "postgres://") {
		t.Fatalf("非 PostgreSQL 连接串应被拒绝，实际: %v", err)
	}
}

func TestValidateRejectsBadEnv(t *testing.T) {
	cfg := baseConfig()
	cfg.App.Env = "production" // 只接受 prod
	err := cfg.validate()
	if err == nil || !strings.Contains(err.Error(), "APP_ENV") {
		t.Fatalf("非法环境名应被拒绝，实际: %v", err)
	}
}

func TestValidateRejectsBadPort(t *testing.T) {
	cfg := baseConfig()
	cfg.HTTP.Port = 70000
	if err := cfg.validate(); err == nil {
		t.Fatal("非法端口应被拒绝")
	}
}

func TestValidateRequiresStrongSecretInProduction(t *testing.T) {
	cfg := baseConfig()
	cfg.App.Env = EnvProd
	cfg.JWT.Secret = "short"
	err := cfg.validate()
	if err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("生产环境弱密钥应被拒绝，实际: %v", err)
	}

	cfg.JWT.Secret = strings.Repeat("x", 32)
	if err := cfg.validate(); err != nil {
		t.Fatalf("32 位密钥在生产环境应被接受，实际: %v", err)
	}
}

func TestValidateRejectsRefreshTTLShorterThanAccessTTL(t *testing.T) {
	cfg := baseConfig()
	cfg.JWT.AccessTTL = 24 * time.Hour
	cfg.JWT.RefreshTTL = time.Hour
	err := cfg.validate()
	if err == nil || !strings.Contains(err.Error(), "JWT_REFRESH_TTL") {
		t.Fatalf("Refresh 有效期必须大于 Access，实际: %v", err)
	}
}

func TestValidateAggregatesAllProblems(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.DSN = ""
	cfg.JWT.Secret = ""
	cfg.HTTP.Port = 0

	err := cfg.validate()
	if err == nil {
		t.Fatal("多项配置错误应返回错误")
	}
	msg := err.Error()
	for _, want := range []string{"DATABASE_URL", "JWT_SECRET", "HTTP_PORT"} {
		if !strings.Contains(msg, want) {
			t.Errorf("错误信息应同时包含 %s，实际:\n%s", want, msg)
		}
	}
}

func TestSafeDSNMasksPassword(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.DSN = "postgres://yanmai:SuperSecret123@10.0.0.5:5432/yanmai_dev?sslmode=require"

	got := cfg.Database.SafeDSN()
	if strings.Contains(got, "SuperSecret123") {
		t.Fatalf("脱敏后不得包含明文密码: %s", got)
	}
	if !strings.Contains(got, "yanmai:***@10.0.0.5:5432") {
		t.Fatalf("脱敏结果格式不符: %s", got)
	}
}

func TestSafeDSNHandlesNoPassword(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.DSN = "postgres://postgres@127.0.0.1:5432/yanmai_dev?sslmode=disable"
	if got := cfg.Database.SafeDSN(); got != cfg.Database.DSN {
		t.Fatalf("无密码连接串应原样返回，实际: %s", got)
	}
}
