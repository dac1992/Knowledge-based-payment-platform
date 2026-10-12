package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/internal/config"
	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/jwt"
	"github.com/yanmai/api/pkg/logger"
)

const testSecret = "integration-test-secret-at-least-32-chars"

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	TraceID string          `json:"traceId"`
}

func testConfig() *config.Config {
	return &config.Config{
		App: config.App{Name: "yanmai-api", Env: config.EnvLocal},
		HTTP: config.HTTP{
			Host: "127.0.0.1", Port: 8080,
			ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second,
			IdleTimeout: 5 * time.Second, ShutdownTimeout: 5 * time.Second,
			CORSOrigins:    []string{"http://localhost:5173"},
			RateLimitRPS:   10000,
			RateLimitBurst: 10000,
		},
		JWT: config.JWT{
			Secret: testSecret, Issuer: "yanmai-api",
			AccessTTL: time.Hour, RefreshTTL: 24 * time.Hour,
		},
		Log: config.Log{Level: "error"},
	}
}

func newTestEngine(t *testing.T, cfg *config.Config) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	// DB 传 nil：本套集成测试覆盖与数据库无关的路由、鉴权、错误处理与中间件行为
	return New(Deps{Config: cfg, DB: nil, Logger: logger.NewNop()})
}

func call(t *testing.T, e *gin.Engine, method, path, body string, headers map[string]string) (*httptest.ResponseRecorder, envelope) {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	var env envelope
	_ = json.Unmarshal(w.Body.Bytes(), &env)
	return w, env
}

func TestHealthEndpointReportsDatabaseStatus(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/health", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("状态码应为 200，实际 %d", w.Code)
	}
	if env.Code != int(errs.CodeOK) {
		t.Fatalf("业务码应为 0，实际 %d", env.Code)
	}
	if env.TraceID == "" {
		t.Fatal("响应应携带 traceId")
	}
	if got := w.Header().Get("X-Trace-Id"); got == "" {
		t.Fatal("响应头应携带 X-Trace-Id")
	}

	var data struct {
		Status   string `json:"status"`
		Service  string `json:"service"`
		Version  string `json:"version"`
		Database struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		} `json:"database"`
	}
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("解析 data 失败: %v", err)
	}
	if data.Service != "yanmai-api" {
		t.Errorf("service 字段不符: %s", data.Service)
	}
	if data.Status != "degraded" || data.Database.Status != "down" {
		t.Errorf("数据库不可用时状态应为 degraded/down，实际 %s/%s", data.Status, data.Database.Status)
	}
}

func TestHealthReadyReturns503WhenDatabaseDown(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, _ := call(t, e, http.MethodGet, "/health/ready", "", nil)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("数据库不可用时 /health/ready 应返回 503，实际 %d", w.Code)
	}
}

func TestHealthLiveAlwaysOK(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/health/live", "", nil)
	if w.Code != http.StatusOK || env.Code != int(errs.CodeOK) {
		t.Fatalf("存活探针应始终 200，实际 %d / %d", w.Code, env.Code)
	}
}

func TestMetaListsAllTwelveModules(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/api/v1/meta/modules", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("状态码应为 200，实际 %d", w.Code)
	}

	var data struct {
		Modules []string `json:"modules"`
		Queue   bool     `json:"queue"`
	}
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("解析 data 失败: %v", err)
	}
	if len(data.Modules) != 12 {
		t.Fatalf("应挂载 12 个模块，实际 %d：%v", len(data.Modules), data.Modules)
	}

	want := []string{"auth", "user", "course", "order", "payment", "learning", "marketing", "member", "live", "circle", "finance", "stats"}
	for i, name := range want {
		if data.Modules[i] != name {
			t.Errorf("第 %d 个模块应为 %s，实际 %s", i+1, name, data.Modules[i])
		}
	}
	if data.Queue {
		t.Error("未配置 Redis 时 queue 应为 false")
	}
}

func TestUnknownRouteReturnsUnified404(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/api/v1/does-not-exist", "", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("状态码应为 404，实际 %d", w.Code)
	}
	if env.Code != int(errs.CodeNotFound) {
		t.Fatalf("业务码应为 %d，实际 %d", errs.CodeNotFound, env.Code)
	}
}

func TestMethodNotAllowedReturns405(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodDelete, "/api/v1/auth/login", "", nil)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("状态码应为 405，实际 %d", w.Code)
	}
	if env.Code != int(errs.CodeMethodNotAllow) {
		t.Fatalf("业务码应为 %d，实际 %d", errs.CodeMethodNotAllow, env.Code)
	}
}

func TestProtectedRouteRequiresToken(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/api/v1/users/me", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("状态码应为 401，实际 %d", w.Code)
	}
	if env.Code != int(errs.CodeTokenMissing) {
		t.Fatalf("业务码应为 %d，实际 %d", errs.CodeTokenMissing, env.Code)
	}
}

func TestProtectedRouteRejectsMalformedToken(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/api/v1/users/me", "", map[string]string{
		"Authorization": "Bearer not-a-real-token",
	})
	if w.Code != http.StatusUnauthorized || env.Code != int(errs.CodeTokenInvalid) {
		t.Fatalf("非法令牌应返回 401/CodeTokenInvalid，实际 %d / %d", w.Code, env.Code)
	}
}

func TestProtectedRouteRejectsNonBearerScheme(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodGet, "/api/v1/users/me", "", map[string]string{
		"Authorization": "Basic dXNlcjpwYXNz",
	})
	if w.Code != http.StatusUnauthorized || env.Code != int(errs.CodeTokenMissing) {
		t.Fatalf("非 Bearer 方案应返回 401/CodeTokenMissing，实际 %d / %d", w.Code, env.Code)
	}
}

func TestAccessProtectedRouteRejectsRefreshToken(t *testing.T) {
	e := newTestEngine(t, testConfig())

	mgr := jwt.NewManager(testSecret, "yanmai-api", time.Hour, 24*time.Hour)
	refresh, err := mgr.GenerateRefresh(1, "admin", "admin")
	if err != nil {
		t.Fatalf("签发 Refresh Token 失败: %v", err)
	}

	w, env := call(t, e, http.MethodGet, "/api/v1/users/me", "", map[string]string{
		"Authorization": "Bearer " + refresh.Value,
	})
	if w.Code != http.StatusUnauthorized || env.Code != int(errs.CodeTokenTypeMismatch) {
		t.Fatalf("Refresh Token 不得访问业务接口，实际 %d / %d", w.Code, env.Code)
	}
}

func TestProtectedRouteAcceptsValidAccessToken(t *testing.T) {
	e := newTestEngine(t, testConfig())

	mgr := jwt.NewManager(testSecret, "yanmai-api", time.Hour, 24*time.Hour)
	access, err := mgr.GenerateAccess(1, "admin", "admin")
	if err != nil {
		t.Fatalf("签发 Access Token 失败: %v", err)
	}

	// DB 为 nil，鉴权通过后会因查询失败返回 500 —— 这里只验证「已越过鉴权层」
	w, env := call(t, e, http.MethodGet, "/api/v1/users/me", "", map[string]string{
		"Authorization": "Bearer " + access.Value,
	})
	if w.Code == http.StatusUnauthorized {
		t.Fatalf("合法令牌不应被鉴权层拒绝（业务码 %d）", env.Code)
	}
	if env.Code == int(errs.CodeTokenMissing) || env.Code == int(errs.CodeTokenInvalid) || env.Code == int(errs.CodeTokenTypeMismatch) {
		t.Fatalf("合法令牌不应产生令牌类错误，实际 %d", env.Code)
	}
}

func TestLoginValidationFailureReturns400(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, env := call(t, e, http.MethodPost, "/api/v1/auth/login", `{"login":"ab"}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("参数校验失败应返回 400，实际 %d", w.Code)
	}
	if env.Code != int(errs.CodeInvalidParam) {
		t.Fatalf("业务码应为 %d，实际 %d", errs.CodeInvalidParam, env.Code)
	}
	if !strings.Contains(env.Message, "login") {
		t.Errorf("校验提示应指出具体字段，实际: %s", env.Message)
	}
}

func TestRegisterValidationFailureReturns400(t *testing.T) {
	e := newTestEngine(t, testConfig())

	body := `{"username":"validuser","password":"123","email":"not-an-email"}`
	w, env := call(t, e, http.MethodPost, "/api/v1/auth/register", body, nil)
	if w.Code != http.StatusBadRequest || env.Code != int(errs.CodeInvalidParam) {
		t.Fatalf("应返回 400/CodeInvalidParam，实际 %d / %d", w.Code, env.Code)
	}
}

func TestPlaceholderModuleReturns501(t *testing.T) {
	e := newTestEngine(t, testConfig())

	for _, name := range []string{"course", "order", "marketing", "stats"} {
		w, env := call(t, e, http.MethodGet, "/api/v1/"+name, "", nil)
		if w.Code != http.StatusNotImplemented {
			t.Errorf("模块 %s 应返回 501，实际 %d", name, w.Code)
		}
		if env.Code != int(errs.CodeNotImplemented) {
			t.Errorf("模块 %s 业务码应为 %d，实际 %d", name, errs.CodeNotImplemented, env.Code)
		}
	}
}

func TestCORSPreflightForWhitelistedOrigin(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, _ := call(t, e, http.MethodOptions, "/api/v1/auth/login", "", map[string]string{
		"Origin":                        "http://localhost:5173",
		"Access-Control-Request-Method": "POST",
	})
	if w.Code != http.StatusNoContent {
		t.Fatalf("预检请求应返回 204，实际 %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("白名单来源应被回显，实际 %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("应允许携带凭证，实际 %q", got)
	}
}

func TestCORSDoesNotEchoUnlistedOrigin(t *testing.T) {
	e := newTestEngine(t, testConfig())

	w, _ := call(t, e, http.MethodOptions, "/api/v1/auth/login", "", map[string]string{
		"Origin": "http://evil.example.com",
	})
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("非白名单来源不应被回显，实际 %q", got)
	}
}

func TestRateLimitReturns429(t *testing.T) {
	cfg := testConfig()
	cfg.HTTP.RateLimitRPS = 1
	cfg.HTTP.RateLimitBurst = 1
	e := newTestEngine(t, cfg)

	_, _ = call(t, e, http.MethodGet, "/health/live", "", nil)
	w, env := call(t, e, http.MethodGet, "/health/live", "", nil)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("超过限流阈值应返回 429，实际 %d", w.Code)
	}
	if env.Code != int(errs.CodeTooManyRequest) {
		t.Fatalf("业务码应为 %d，实际 %d", errs.CodeTooManyRequest, env.Code)
	}
}

func TestTraceIDIsPropagatedFromHeader(t *testing.T) {
	e := newTestEngine(t, testConfig())

	const tid = "trace-from-client-0001"
	w, env := call(t, e, http.MethodGet, "/health/live", "", map[string]string{
		"X-Trace-Id": tid,
	})
	if env.TraceID != tid || w.Header().Get("X-Trace-Id") != tid {
		t.Fatalf("应透传客户端 traceId，实际 body=%q header=%q", env.TraceID, w.Header().Get("X-Trace-Id"))
	}
}
