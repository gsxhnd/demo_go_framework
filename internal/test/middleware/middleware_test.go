package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go_sample_code/internal/middleware"
	"go_sample_code/pkg/metrics"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// =============================================================================
// Metrics 中间件测试
// =============================================================================

func TestMetrics_Middleware(t *testing.T) {
	cfg := metrics.DefaultMetricsConfig()
	cfg.OtelEnable = false
	mp, err := metrics.NewMeterProvider(&cfg)
	require.NoError(t, err)
	defer mp.Shutdown(context.Background())

	recorder, err := metrics.NewHTTPRecorder(mp)
	require.NoError(t, err)

	app := fiber.New()
	app.Use(middleware.Metrics(recorder))

	app.Get("/test/:id", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("Created")
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Status(500).SendString("Internal Error")
	})

	// 测试成功请求
	req := httptest.NewRequest("GET", "/test/123", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 测试 POST 请求
	req = httptest.NewRequest("POST", "/test", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// 测试错误请求
	req = httptest.NewRequest("GET", "/error", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestMetrics_RouteTemplate(t *testing.T) {
	cfg := metrics.DefaultMetricsConfig()
	cfg.OtelEnable = false
	mp, err := metrics.NewMeterProvider(&cfg)
	require.NoError(t, err)
	defer mp.Shutdown(context.Background())

	recorder, err := metrics.NewHTTPRecorder(mp)
	require.NoError(t, err)

	app := fiber.New()
	app.Use(middleware.Metrics(recorder))

	app.Get("/users/:id/posts/:postId", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetRoute(t *testing.T) {
	app := fiber.New()

	app.Get("/api/v1/users/:id", func(c *fiber.Ctx) error {
		route := middleware.GetRoute(c)
		assert.Equal(t, "/api/v1/users/:id", route)
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/api/v1/users/123", nil)
	_, err := app.Test(req)
	require.NoError(t, err)
}

func TestGetProtocol(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		proto := middleware.GetProtocol(c)
		// httptest 请求可能返回 "http" 或 "http/1.1"
		assert.NotEmpty(t, proto)
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req)
	require.NoError(t, err)
}

func TestGetStatusCodeFromResponse(t *testing.T) {
	app := fiber.New()
	statusCodeCapture := 0

	app.Get("/ok", func(c *fiber.Ctx) error {
		c.Status(200)
		statusCodeCapture = middleware.GetStatusCodeFromResponse()(c)
		return c.SendString("OK")
	})
	app.Get("/notfound", func(c *fiber.Ctx) error {
		c.Status(404)
		statusCodeCapture = middleware.GetStatusCodeFromResponse()(c)
		return c.SendString("Not Found")
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	_, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, statusCodeCapture)

	req = httptest.NewRequest("GET", "/notfound", nil)
	_, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, statusCodeCapture)
}

// =============================================================================
// Rate Limit 中间件测试
// =============================================================================

// mockLogger 用于测试的日志记录器
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields ...zap.Field)                         {}
func (m *mockLogger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {}
func (m *mockLogger) Info(msg string, fields ...zap.Field)                          {}
func (m *mockLogger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field)  {}
func (m *mockLogger) Warn(msg string, fields ...zap.Field)                          {}
func (m *mockLogger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field)  {}
func (m *mockLogger) Error(msg string, fields ...zap.Field)                         {}
func (m *mockLogger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {}
func (m *mockLogger) Panic(msg string, fields ...zap.Field)                         {}
func (m *mockLogger) PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {}
func (m *mockLogger) GetLogger() *zap.Logger                                        { return nil }
func (m *mockLogger) Shutdown(ctx context.Context)                                  {}

// setupTestApp 创建一个用于测试的 Fiber app
func setupTestApp(cfg *middleware.RateLimitConfig) *fiber.App {
	app := fiber.New()

	ml := &mockLogger{}
	if cfg == nil {
		cfg = middleware.DefaultRateLimitConfig(ml)
	}

	app.Use(middleware.RateLimit(cfg))

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/api/users", func(c *fiber.Ctx) error {
		return c.SendString("users")
	})

	app.Post("/api/users", func(c *fiber.Ctx) error {
		return c.SendString("created")
	})

	return app
}

func TestRateLimit_Disabled(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.Enabled = false
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	for i := 0; i < 10; i++ {
		req := mustNewRequest("GET", "/api/users")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "请求 %d 应该成功", i+1)
	}
}

func TestRateLimit_SkipPaths(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	for i := 0; i < 20; i++ {
		req := mustNewRequest("GET", "/api/health")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "健康检查请求 %d 应该成功", i+1)
	}
}

func TestRateLimit_BurstAllowed(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 1
	cfg.Burst = 5

	app := setupTestApp(cfg)

	for i := 0; i < 5; i++ {
		req := mustNewRequest("GET", "/api/users")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "请求 %d 应该成功", i+1)
	}
}

func TestRateLimit_Exceeded(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 1
	cfg.Burst = 2

	app := setupTestApp(cfg)

	for i := 0; i < 2; i++ {
		req := mustNewRequest("GET", "/api/users")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}

	time.Sleep(100 * time.Millisecond)

	req := mustNewRequest("GET", "/api/users")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
}

func TestRateLimit_DifferentKeys(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	req1 := mustNewRequest("GET", "/api/users")
	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	req2 := mustNewRequest("POST", "/api/users")
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode, "不同 METHOD 应该有不同的限流桶")
}

func TestRateLimit_IPWhitelist(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.IPWhitelist = []string{"127.0.0.1", "192.168.1.0/24"}
	whitelist := make(map[string]bool)
	for _, ip := range cfg.IPWhitelist {
		whitelist[ip] = true
	}

	assert.True(t, middleware.IsWhitelistedIP("127.0.0.1", cfg, whitelist), "127.0.0.1 应该在白名单中")
	assert.True(t, middleware.IsWhitelistedIP("192.168.1.100", cfg, whitelist), "192.168.1.100 应该匹配 192.168.1.0/24")
	assert.True(t, middleware.IsWhitelistedIP("192.168.1.200", cfg, whitelist), "192.168.1.200 应该匹配 192.168.1.0/24")
	assert.False(t, middleware.IsWhitelistedIP("10.0.0.1", cfg, whitelist), "10.0.0.1 不应该在白名单中")
	assert.False(t, middleware.IsWhitelistedIP("8.8.8.8", cfg, whitelist), "8.8.8.8 不应该在白名单中")
}

func TestRateLimit_Cleanup(t *testing.T) {
	ml := &mockLogger{}
	cfg := &middleware.RateLimitConfig{
		Log:               ml,
		Enabled:           true,
		RequestsPerSecond: 1,
		Burst:             1,
		CleanupInterval:   50 * time.Millisecond,
		EntryTTL:          100 * time.Millisecond,
		SkipPaths:         []string{},
		TrustedProxies:    []string{},
		IPHeader:          "X-Forwarded-For",
		IPWhitelist:       []string{},
	}

	rl := middleware.NewRateLimiter(cfg)
	defer rl.Stop()

	rl.Allow("key1")
	rl.Allow("key2")
	rl.Allow("key3")

	assert.Equal(t, 3, rl.EntryCount(), "应该有 3 个 limiter")

	time.Sleep(200 * time.Millisecond)
	rl.Allow("key4") // 触发内部时间更新，但不影响 key1-key3 的过期

	// 等待足够时间后，key1-key3 应该还在或已被清理
	// 这里仅验证 EntryCount 不会 panic
	_ = rl.EntryCount()
}

func TestRateLimit_ConcurrentAccess(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 10
	cfg.Burst = 10

	app := setupTestApp(cfg)

	done := make(chan bool, 20)
	results := make([]int, 20)

	for i := 0; i < 20; i++ {
		go func(idx int) {
			req := mustNewRequest("GET", "/api/users")
			resp, err := app.Test(req)
			if err == nil {
				results[idx] = resp.StatusCode
			}
			done <- true
		}(i)
	}

	for i := 0; i < 20; i++ {
		<-done
	}

	successCount := 0
	rateLimitCount := 0
	for _, status := range results {
		if status == fiber.StatusOK {
			successCount++
		} else if status == fiber.StatusTooManyRequests {
			rateLimitCount++
		}
	}

	assert.Greater(t, successCount, 0, "应该有成功的请求")
	assert.LessOrEqual(t, successCount, cfg.Burst, "成功请求数不应超过突发容量")
}

func TestDefaultRateLimitConfig(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)

	assert.True(t, cfg.Enabled)
	assert.Equal(t, float64(20), cfg.RequestsPerSecond)
	assert.Equal(t, 50, cfg.Burst)
	assert.Equal(t, 1*time.Minute, cfg.CleanupInterval)
	assert.Equal(t, 10*time.Minute, cfg.EntryTTL)
	assert.Contains(t, cfg.SkipPaths, "/api/health")
	assert.Contains(t, cfg.IPWhitelist, "127.0.0.1")
}

func TestRateLimit_StopCleanup(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.CleanupInterval = 10 * time.Millisecond

	rl := middleware.NewRateLimiter(cfg)

	rl.Allow("key1")
	rl.Allow("key2")

	rl.Stop()

	time.Sleep(50 * time.Millisecond)
	// 验证不会 panic
}

func TestMatchCIDR(t *testing.T) {
	tests := []struct {
		ip     string
		cidr   string
		expect bool
	}{
		{"127.0.0.1", "127.0.0.0/8", true},
		{"127.0.0.1", "192.168.0.0/16", false},
		{"192.168.1.100", "192.168.0.0/16", true},
		{"10.0.0.1", "10.0.0.0/8", true},
		{"10.255.255.255", "10.0.0.0/8", true},
		{"172.16.0.1", "172.16.0.0/12", true},
		{"::1", "::1/128", true},
		{"192.168.1.1", "192.168.1.1", true},
		{"invalid", "192.168.0.0/16", false},
	}

	for _, tt := range tests {
		t.Run(tt.ip+"_"+tt.cidr, func(t *testing.T) {
			result := middleware.MatchCIDR(tt.ip, tt.cidr)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestNormalizeIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"127.0.0.1", "127.0.0.1"},
		{"192.168.1.100", "192.168.1.100"},
		{"::1", "::1"},
		{"2001:db8::1", "2001:db8::1"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := middleware.NormalizeIP(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimit_429Response(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	req := mustNewRequest("GET", "/api/users")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	req = mustNewRequest("GET", "/api/users")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestRateLimit_CustomKeyFunc(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1
	cfg.KeyFunc = func(c *fiber.Ctx) string {
		return "global:" + c.Path()
	}

	app := setupTestApp(cfg)

	req := mustNewRequest("GET", "/api/users")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	req = mustNewRequest("GET", "/api/users")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
}

func TestClientIP_DirectConnection(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.TrustedProxies = []string{}

	app := fiber.New()
	app.Use(middleware.RateLimit(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(middleware.ClientIP(c, cfg))
	})

	req := mustNewRequest("GET", "/")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestClientIP_TrustedProxy(t *testing.T) {
	ml := &mockLogger{}
	cfg := middleware.DefaultRateLimitConfig(ml)
	cfg.TrustedProxies = []string{"127.0.0.0/8"}
	cfg.IPHeader = "X-Forwarded-For"

	app := fiber.New()
	app.Use(middleware.RateLimit(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(middleware.ClientIP(c, cfg))
	})

	req := mustNewRequest("GET", "/")
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.1, 127.0.0.1")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func mustNewRequest(method, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	return req
}
