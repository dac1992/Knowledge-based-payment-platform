// Package config 负责加载与校验运行时配置。
//
// 配置来源优先级（后者覆盖前者）：
//  1. 内置默认值
//  2. .env 文件
//  3. .env.local 文件（本地开发覆盖，不入 Git）
//  4. 真实环境变量
//
// 配置文件从当前工作目录向上最多查找 4 层，因此无论是
// `go run ./cmd/server`（cwd=services/api）还是运行编译后的二进制，
// 都能正确定位到 services/api/.env.local。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Version 服务版本号，随发版更新。
const Version = "0.1.0"

// 运行环境取值。
const (
	EnvLocal   = "local"
	EnvDev     = "dev"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

// Config 应用总配置。
type Config struct {
	App      App
	HTTP     HTTP
	Database Database
	JWT      JWT
	Redis    Redis
	Log      Log
}

// App 应用元信息。
type App struct {
	Name string
	Env  string
}

// IsProduction 是否生产环境。
func (a App) IsProduction() bool { return a.Env == EnvProd }

// IsDevelopment 是否本地或联调环境。
func (a App) IsDevelopment() bool { return a.Env == EnvLocal || a.Env == EnvDev }

// HTTP HTTP 服务配置。
type HTTP struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	CORSOrigins     []string
	RateLimitRPS    float64
	RateLimitBurst  int
}

// Addr 返回监听地址。
func (h HTTP) Addr() string { return fmt.Sprintf("%s:%d", h.Host, h.Port) }

// Database 数据库配置。
type Database struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	AutoMigrate     bool
	Seed            bool
}

// SafeDSN 返回脱敏后的连接串，用于日志输出。
func (d Database) SafeDSN() string {
	at := strings.LastIndex(d.DSN, "@")
	proto := strings.Index(d.DSN, "://")
	if at < 0 || proto < 0 || at < proto {
		return d.DSN
	}
	head := d.DSN[:proto+3]
	user := d.DSN[proto+3 : at]
	if i := strings.Index(user, ":"); i >= 0 {
		user = user[:i] + ":***"
	}
	return head + user + d.DSN[at:]
}

// JWT 认证配置。
type JWT struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// Redis 缓存与队列配置。
type Redis struct {
	Addr     string
	Password string
	DB       int
}

// Log 日志配置。
type Log struct {
	Level string
}

// Load 加载并校验配置。
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")

	if p, ok := findEnvFile(); ok {
		v.SetConfigFile(p)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("读取配置文件 %s 失败: %w", p, err)
		}
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setDefaults(v)

	cfg := &Config{
		App: App{
			Name: v.GetString("APP_NAME"),
			Env:  strings.ToLower(v.GetString("APP_ENV")),
		},
		HTTP: HTTP{
			Host:            v.GetString("HTTP_HOST"),
			Port:            v.GetInt("HTTP_PORT"),
			ReadTimeout:     v.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout:    v.GetDuration("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:     v.GetDuration("HTTP_IDLE_TIMEOUT"),
			ShutdownTimeout: v.GetDuration("HTTP_SHUTDOWN_TIMEOUT"),
			CORSOrigins:     splitAndTrim(v.GetString("HTTP_CORS_ORIGINS")),
			RateLimitRPS:    v.GetFloat64("HTTP_RATE_LIMIT_RPS"),
			RateLimitBurst:  v.GetInt("HTTP_RATE_LIMIT_BURST"),
		},
		Database: Database{
			DSN:             v.GetString("DATABASE_URL"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
			AutoMigrate:     v.GetBool("DB_AUTO_MIGRATE"),
			Seed:            v.GetBool("DB_SEED"),
		},
		JWT: JWT{
			Secret:     v.GetString("JWT_SECRET"),
			Issuer:     v.GetString("JWT_ISSUER"),
			AccessTTL:  v.GetDuration("JWT_ACCESS_TTL"),
			RefreshTTL: v.GetDuration("JWT_REFRESH_TTL"),
		},
		Redis: Redis{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		Log: Log{Level: v.GetString("LOG_LEVEL")},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	var problems []string

	switch c.App.Env {
	case EnvLocal, EnvDev, EnvStaging, EnvProd:
	default:
		problems = append(problems, fmt.Sprintf("APP_ENV 取值非法: %q（应为 local/dev/staging/prod）", c.App.Env))
	}

	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		problems = append(problems, fmt.Sprintf("HTTP_PORT 非法: %d", c.HTTP.Port))
	}

	if strings.TrimSpace(c.Database.DSN) == "" {
		problems = append(problems, "DATABASE_URL 未配置：请在 services/api/.env.local 中填入云 PostgreSQL 连接串")
	} else if !strings.HasPrefix(c.Database.DSN, "postgres://") && !strings.HasPrefix(c.Database.DSN, "postgresql://") {
		problems = append(problems, "DATABASE_URL 必须以 postgres:// 或 postgresql:// 开头")
	}

	if strings.TrimSpace(c.JWT.Secret) == "" {
		problems = append(problems, "JWT_SECRET 未配置")
	} else if !c.App.IsDevelopment() && len(c.JWT.Secret) < 32 {
		problems = append(problems, "非开发环境下 JWT_SECRET 长度需不少于 32 个字符")
	}

	if c.JWT.AccessTTL <= 0 {
		problems = append(problems, "JWT_ACCESS_TTL 必须大于 0")
	}
	if c.JWT.RefreshTTL <= c.JWT.AccessTTL {
		problems = append(problems, "JWT_REFRESH_TTL 必须大于 JWT_ACCESS_TTL")
	}

	if len(problems) > 0 {
		return fmt.Errorf("配置校验未通过:\n  - %s", strings.Join(problems, "\n  - "))
	}
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "yanmai-api")
	v.SetDefault("APP_ENV", EnvLocal)

	v.SetDefault("HTTP_HOST", "0.0.0.0")
	v.SetDefault("HTTP_PORT", 8080)
	v.SetDefault("HTTP_READ_TIMEOUT", "15s")
	v.SetDefault("HTTP_WRITE_TIMEOUT", "30s")
	v.SetDefault("HTTP_IDLE_TIMEOUT", "60s")
	v.SetDefault("HTTP_SHUTDOWN_TIMEOUT", "10s")
	v.SetDefault("HTTP_CORS_ORIGINS", "http://localhost:5173,http://localhost:5174,http://localhost:5175")
	v.SetDefault("HTTP_RATE_LIMIT_RPS", 50)
	v.SetDefault("HTTP_RATE_LIMIT_BURST", 100)

	v.SetDefault("DB_MAX_OPEN_CONNS", 20)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "30m")
	v.SetDefault("DB_AUTO_MIGRATE", true)
	v.SetDefault("DB_SEED", false)

	v.SetDefault("JWT_ISSUER", "yanmai-api")
	v.SetDefault("JWT_ACCESS_TTL", "2h")
	v.SetDefault("JWT_REFRESH_TTL", "168h")

	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("LOG_LEVEL", "info")
}

func splitAndTrim(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// findEnvFile 从当前目录向上查找配置文件，优先 .env.local。
func findEnvFile() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for i := 0; i < 4; i++ {
		for _, name := range []string{".env.local", ".env"} {
			p := filepath.Join(dir, name)
			if st, statErr := os.Stat(p); statErr == nil && !st.IsDir() {
				return p, true
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}
