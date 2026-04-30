package middleware

import (
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"go_sample_code/internal/errno"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimitConfig 限流中间件配置
type RateLimitConfig struct {
	// 日志记录器
	Log logger.Logger
	// 是否启用限流
	Enabled bool
	// 每秒补充的令牌数
	RequestsPerSecond float64
	// 令牌桶容量（突发容量）
	Burst int
	// 清理过期 limiter 的间隔
	CleanupInterval time.Duration
	// limiter 条目过期时间
	EntryTTL time.Duration
	// 跳过限流的路径
	SkipPaths []string
	// 可信代理 CIDR 列表
	TrustedProxies []string
	// IP 来源 Header
	IPHeader string
	// IP 白名单
	IPWhitelist []string
	// 自定义 key 生成函数
	KeyFunc func(*fiber.Ctx) string
}

// limiterEntry 存储每个 key 对应的限流器
type limiterEntry struct {
	limiter    *rate.Limiter
	lastSeenAt time.Time
}

// rateLimiter 令牌桶限流器
type RateLimiter struct {
	mu              sync.RWMutex
	entries         map[string]*limiterEntry
	requestsPerSec  rate.Limit
	burst           int
	cleanupInterval time.Duration
	entryTTL        time.Duration
	log             logger.Logger
	stopCleanup     chan struct{}
}

// NewRateLimiter 创建限流器
func NewRateLimiter(cfg *RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		entries:         make(map[string]*limiterEntry),
		requestsPerSec:  rate.Limit(cfg.RequestsPerSecond),
		burst:           cfg.Burst,
		cleanupInterval: cfg.CleanupInterval,
		entryTTL:        cfg.EntryTTL,
		log:             cfg.Log,
		stopCleanup:     make(chan struct{}),
	}

	// 启动后台清理协程
	go rl.cleanupLoop()

	return rl
}

// cleanupLoop 定期清理过期的 limiter 条目
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanup 清理过期的条目
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, entry := range rl.entries {
		if now.Sub(entry.lastSeenAt) > rl.entryTTL {
			delete(rl.entries, key)
		}
	}
}

// Stop 停止限流器，清理后台协程
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.entries[key]
	if !exists {
		entry = &limiterEntry{
			limiter:    rate.NewLimiter(rl.requestsPerSec, rl.burst),
			lastSeenAt: time.Now(),
		}
		rl.entries[key] = entry
	}

	// 检查是否允许
	allowed := entry.limiter.Allow()

	// 更新 lastSeenAt
	entry.lastSeenAt = time.Now()

	return allowed
}

// EntryCount 返回当前条目数量（用于测试）
func (rl *RateLimiter) EntryCount() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return len(rl.entries)
}

// RateLimit 返回限流中间件处理函数
func RateLimit(cfg *RateLimitConfig) fiber.Handler {
	rl := NewRateLimiter(cfg)

	// 构建 skip paths 集合
	skipPaths := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}

	// 构建可信代理集合
	trustedProxies := make(map[string]bool)
	for _, proxy := range cfg.TrustedProxies {
		trustedProxies[proxy] = true
	}

	// 构建 IP 白名单集合
	ipWhitelist := make(map[string]bool)
	for _, ip := range cfg.IPWhitelist {
		ipWhitelist[ip] = true
	}

	// 默认 key 生成函数
	defaultKeyFunc := func(c *fiber.Ctx) string {
		clientIP := ClientIP(c, cfg)
		method := c.Method()
		route := c.Route().Path

		if route == "" {
			route = c.Path()
		}

		return clientIP + ":" + method + ":" + route
	}

	return func(c *fiber.Ctx) error {
		// 如果未启用限流，直接跳过
		if !cfg.Enabled {
			return c.Next()
		}

		// 检查是否跳过路径
		if skipPaths[c.Path()] {
			return c.Next()
		}

		// 获取客户端 IP
		clientIP := ClientIP(c, cfg)

		// 检查 IP 白名单
		if IsWhitelistedIP(clientIP, cfg, ipWhitelist) {
			return c.Next()
		}

		// 生成限流 key
		keyFunc := cfg.KeyFunc
		if keyFunc == nil {
			keyFunc = defaultKeyFunc
		}

		key := keyFunc(c)
		if key == "" {
			key = "unknown:" + c.Method() + ":" + c.Path()
			cfg.Log.WarnCtx(c.UserContext(), "failed to generate rate limit key, using fallback",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()))
		}

		// 检查是否允许
		if rl.Allow(key) {
			return c.Next()
		}

		// 限流被触发
		cfg.Log.WarnCtx(c.UserContext(), "rate limit exceeded",
			zap.String("key", key),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", clientIP))

		decoded := errno.Decode(nil, errno.RateLimitExceededError)
		return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
	}
}

// DefaultRateLimitConfig 返回默认限流配置
func DefaultRateLimitConfig(log logger.Logger) *RateLimitConfig {
	return &RateLimitConfig{
		Log:               log,
		Enabled:           true,
		RequestsPerSecond: 20,
		Burst:             50,
		CleanupInterval:   1 * time.Minute,
		EntryTTL:          10 * time.Minute,
		SkipPaths:         []string{"/api/health"},
		TrustedProxies:    []string{},
		IPHeader:          "X-Forwarded-For",
		IPWhitelist:       []string{"127.0.0.1"},
		KeyFunc:           nil,
	}
}

// ClientIP 提取客户端真实 IP
func ClientIP(c *fiber.Ctx, cfg *RateLimitConfig) string {
	remoteAddr := c.IP()

	// 如果配置了可信代理，尝试从 Header 解析
	if len(cfg.TrustedProxies) > 0 {
		if isTrustedProxy(remoteAddr, cfg) {
			ip := GetIPFromHeader(c, cfg.IPHeader)
			if ip != "" {
				return NormalizeIP(ip)
			}
		}
	}

	// 直接使用 remote addr
	return NormalizeIP(remoteAddr)
}

// GetIPFromHeader 从指定 Header 提取 IP
func GetIPFromHeader(c *fiber.Ctx, header string) string {
	// Fiber 提供了标准 header 的便捷方法
	header = textproto.CanonicalMIMEHeaderKey(header)

	var value string
	switch header {
	case "X-Forwarded-For":
		value = c.Get("X-Forwarded-For")
	case "X-Real-IP":
		value = c.Get("X-Real-IP")
	default:
		value = c.Get(header)
	}

	if value == "" {
		return ""
	}

	// X-Forwarded-For 可能包含多个 IP，取第一个
	parts := strings.Split(value, ",")
	for _, part := range parts {
		ip := strings.TrimSpace(part)
		parsed := net.ParseIP(ip)
		if parsed != nil {
			return ip
		}
	}

	return ""
}

// NormalizeIP 规范化 IP 地址
func NormalizeIP(raw string) string {
	ip := net.ParseIP(raw)
	if ip == nil {
		return raw
	}
	return ip.String()
}

// isTrustedProxy 检查 IP 是否来自可信代理
func isTrustedProxy(remoteIP string, cfg *RateLimitConfig) bool {
	for _, cidr := range cfg.TrustedProxies {
		if MatchCIDR(remoteIP, cidr) {
			return true
		}
	}
	return false
}

// IsWhitelistedIP 检查 IP 是否在白名单中
func IsWhitelistedIP(ip string, cfg *RateLimitConfig, whitelist map[string]bool) bool {
	// 检查直接匹配
	if whitelist[ip] {
		return true
	}

	// 检查 CIDR 匹配
	for _, cidr := range cfg.IPWhitelist {
		if MatchCIDR(ip, cidr) {
			return true
		}
	}

	return false
}

// MatchCIDR 检查 IP 是否匹配 CIDR
func MatchCIDR(ipStr, cidrStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	_, cidr, err := net.ParseCIDR(cidrStr)
	if err != nil {
		// 如果不是有效的 CIDR，尝试作为直接 IP 比较
		return ipStr == cidrStr
	}

	return cidr.Contains(ip)
}
