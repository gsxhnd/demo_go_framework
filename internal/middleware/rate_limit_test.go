package middleware

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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
func setupTestApp(cfg *RateLimitConfig) *fiber.App {
	app := fiber.New()

	if cfg == nil {
		cfg = DefaultRateLimitConfig(&mockLogger{})
	}

	app.Use(RateLimit(cfg))

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

// TestRateLimit_Disabled 禁用时全部放行
func TestRateLimit_Disabled(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.Enabled = false
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	// 发送多个请求，全部应该成功
	for i := 0; i < 10; i++ {
		req := mustNewRequest("GET", "/api/users")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "请求 %d 应该成功", i+1)
	}
}

// TestRateLimit_SkipPaths 跳过路径不受限流
func TestRateLimit_SkipPaths(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	// 高频访问跳过路径，应该全部成功
	for i := 0; i < 20; i++ {
		req := mustNewRequest("GET", "/api/health")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "健康检查请求 %d 应该成功", i+1)
	}
}

// TestRateLimit_BurstAllowed 在突发容量范围内允许通过
func TestRateLimit_BurstAllowed(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 1
	cfg.Burst = 5

	app := setupTestApp(cfg)

	// 突发容量范围内请求应该成功
	for i := 0; i < 5; i++ {
		req := mustNewRequest("GET", "/api/users")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "请求 %d 应该成功", i+1)
	}
}

// TestRateLimit_Exceeded 返回 429
func TestRateLimit_Exceeded(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 1
	cfg.Burst = 2

	app := setupTestApp(cfg)

	// 先消耗突发容量
	for i := 0; i < 2; i++ {
		req := mustNewRequest("GET", "/api/users")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}

	// 等待令牌补充
	time.Sleep(100 * time.Millisecond)

	// 超出额度的请求应该返回 429
	req := mustNewRequest("GET", "/api/users")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
}

// TestRateLimit_DifferentKeys 不同 key 互不影响
func TestRateLimit_DifferentKeys(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	// GET /api/users 消耗额度
	req1 := mustNewRequest("GET", "/api/users")
	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// POST /api/users 应该不受影响
	req2 := mustNewRequest("POST", "/api/users")
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode, "不同 METHOD 应该有不同的限流桶")
}

// TestRateLimit_IPWhitelist 白名单 IP 不受限流
func TestRateLimit_IPWhitelist(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.IPWhitelist = []string{"127.0.0.1", "192.168.1.0/24"}
	whitelist := make(map[string]bool)
	for _, ip := range cfg.IPWhitelist {
		whitelist[ip] = true
	}

	// 测试白名单 IP 直接匹配
	assert.True(t, isWhitelistedIP("127.0.0.1", cfg, whitelist), "127.0.0.1 应该在白名单中")

	// 测试白名单 CIDR 匹配
	assert.True(t, isWhitelistedIP("192.168.1.100", cfg, whitelist), "192.168.1.100 应该匹配 192.168.1.0/24")
	assert.True(t, isWhitelistedIP("192.168.1.200", cfg, whitelist), "192.168.1.200 应该匹配 192.168.1.0/24")

	// 测试非白名单 IP
	assert.False(t, isWhitelistedIP("10.0.0.1", cfg, whitelist), "10.0.0.1 不应该在白名单中")
	assert.False(t, isWhitelistedIP("8.8.8.8", cfg, whitelist), "8.8.8.8 不应该在白名单中")
}

// TestRateLimit_Cleanup 清理过期 key
func TestRateLimit_Cleanup(t *testing.T) {
	cfg := &RateLimitConfig{
		Log:               &mockLogger{},
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

	rl := newRateLimiter(cfg)
	defer rl.Stop()

	// 创建多个 limiter
	rl.allow("key1")
	rl.allow("key2")
	rl.allow("key3")

	rl.mu.RLock()
	initialCount := len(rl.entries)
	rl.mu.RUnlock()
	assert.Equal(t, 3, initialCount, "应该有 3 个 limiter")

	// 等待 TTL 过期
	time.Sleep(200 * time.Millisecond)

	// 触发清理
	rl.cleanup()

	rl.mu.RLock()
	finalCount := len(rl.entries)
	rl.mu.RUnlock()
	assert.Equal(t, 0, finalCount, "清理后应该没有 limiter")
}

// TestClientIP_DirectConnection 直接连接场景
func TestClientIP_DirectConnection(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.TrustedProxies = []string{}

	app := fiber.New()
	app.Use(RateLimit(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(ClientIP(c, cfg))
	})

	req := mustNewRequest("GET", "/")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// TestClientIP_TrustedProxy 可信代理场景
func TestClientIP_TrustedProxy(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.TrustedProxies = []string{"127.0.0.0/8"}
	cfg.IPHeader = "X-Forwarded-For"

	app := fiber.New()
	app.Use(RateLimit(cfg))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(ClientIP(c, cfg))
	})

	req := mustNewRequest("GET", "/")
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.1, 127.0.0.1")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// TestMatchCIDR CIDR 匹配测试
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
		{"192.168.1.1", "192.168.1.1", true}, // 直接 IP 匹配
		{"invalid", "192.168.0.0/16", false},
	}

	for _, tt := range tests {
		t.Run(tt.ip+"_"+tt.cidr, func(t *testing.T) {
			result := matchCIDR(tt.ip, tt.cidr)
			assert.Equal(t, tt.expect, result)
		})
	}
}

// TestNormalizeIP IP 规范化测试
func TestNormalizeIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"127.0.0.1", "127.0.0.1"},
		{"192.168.1.100", "192.168.1.100"},
		{"::1", "::1"},
		{"2001:db8::1", "2001:db8::1"},
		{"invalid", "invalid"}, // 无效 IP 返回原值
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeIP(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRateLimit_429Response 验证 429 响应格式
func TestRateLimit_429Response(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1

	app := setupTestApp(cfg)

	// 消耗突发容量
	req := mustNewRequest("GET", "/api/users")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 超出额度
	req = mustNewRequest("GET", "/api/users")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

// TestRateLimit_DefaultKeyFunc 默认 key 生成函数
func TestRateLimit_DefaultKeyFunc(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})

	app := fiber.New(fiber.Config{
		AppName: "Test App",
	})
	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		return c.SendString("user")
	})

	handler := RateLimit(cfg)
	app.Use(handler)

	req := mustNewRequest("GET", "/api/users/123")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// TestRateLimit_CustomKeyFunc 自定义 key 生成函数
func TestRateLimit_CustomKeyFunc(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 1
	cfg.Burst = 1
	cfg.KeyFunc = func(c *fiber.Ctx) string {
		// 只按路径限流，忽略 IP
		return "global:" + c.Path()
	}

	app := setupTestApp(cfg)

	// 第一个请求成功
	req := mustNewRequest("GET", "/api/users")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 第二个请求应该被限流（因为使用相同的 key）
	req = mustNewRequest("GET", "/api/users")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
}

// TestRateLimit_ConcurrentAccess 并发访问测试
func TestRateLimit_ConcurrentAccess(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.RequestsPerSecond = 10
	cfg.Burst = 10

	app := setupTestApp(cfg)

	// 并发发送请求
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

	// 等待所有请求完成
	for i := 0; i < 20; i++ {
		<-done
	}

	// 在突发容量范围内应该有成功的请求
	successCount := 0
	rateLimitCount := 0
	for _, status := range results {
		if status == fiber.StatusOK {
			successCount++
		} else if status == fiber.StatusTooManyRequests {
			rateLimitCount++
		}
	}

	// 应该有成功的请求
	assert.Greater(t, successCount, 0, "应该有成功的请求")
	assert.LessOrEqual(t, successCount, cfg.Burst, "成功请求数不应超过突发容量")
}

// TestDefaultRateLimitConfig 默认配置验证
func TestDefaultRateLimitConfig(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})

	assert.True(t, cfg.Enabled)
	assert.Equal(t, float64(20), cfg.RequestsPerSecond)
	assert.Equal(t, 50, cfg.Burst)
	assert.Equal(t, 1*time.Minute, cfg.CleanupInterval)
	assert.Equal(t, 10*time.Minute, cfg.EntryTTL)
	assert.Contains(t, cfg.SkipPaths, "/api/health")
	assert.Contains(t, cfg.IPWhitelist, "127.0.0.1")
}

// TestRateLimit_StopCleanup 停止清理协程
func TestRateLimit_StopCleanup(t *testing.T) {
	cfg := DefaultRateLimitConfig(&mockLogger{})
	cfg.CleanupInterval = 10 * time.Millisecond

	rl := newRateLimiter(cfg)

	// 添加一些数据
	rl.allow("key1")
	rl.allow("key2")

	// 停止
	rl.Stop()

	// 验证不会 panic
	time.Sleep(50 * time.Millisecond)
}

// mustNewRequest 创建一个测试请求
func mustNewRequest(method, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	return req
}
