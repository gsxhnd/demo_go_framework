// Package middleware HTTP 中间件实现
//
// 该文件包含 Metrics 中间件，用于自动收集 HTTP 请求指标。
//
// Metrics 中间件功能：
//   - 自动记录每个 HTTP 请求的指标
//   - 支持按 method、route、status_code 维度聚合
//   - 记录请求耗时分布（可用于计算 P50/P95/P99）
//   - 跟踪活跃请求数
//   - 自动识别错误请求（4xx, 5xx）
//
// 使用方式：
//
//	mp, err := metrics.NewMeterProvider(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	recorder, err := metrics.NewHTTPRecorder(mp)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	app.Use(middleware.Metrics(recorder))
//
// 中间件执行顺序建议：
//
//	Recovery -> RateLimit -> Metrics -> Logger -> [Auth] -> [RBAC] -> Handler
//
// 注意事项：
//   - 路由使用模板路径（如 /api/users/:id）而非原始路径（如 /api/users/123）
//   - 不要在指标属性中记录敏感信息（email, user_id, token 等）
//   - /api/health 等健康检查路径建议跳过 metrics 记录
package middleware

import (
	"strings"
	"time"

	"go_sample_code/pkg/metrics"

	"github.com/gofiber/fiber/v2"
)

// Metrics 创建 metrics 中间件
func Metrics(recorder *metrics.HTTPRecorder) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// 获取路由信息
		route := getRoute(c)

		// 增加活跃请求数
		recorder.ActiveRequestAdd(c.Context(), 1)

		// 执行后续处理
		err := c.Next()

		// defer 中减少活跃请求数
		defer recorder.ActiveRequestAdd(c.Context(), -1)

		// 获取响应信息
		statusCode := c.Response().StatusCode()
		durationMs := float64(time.Since(start).Milliseconds())
		protocol := getProtocol(c)
		hasError := err != nil || statusCode >= 500

		// 记录 metrics
		recorder.RecordRequest(c.Context(), metrics.HTTPRequestInfo{
			Method:     c.Method(),
			Route:      route,
			StatusCode: statusCode,
			DurationMs: durationMs,
			HasError:   hasError,
			Protocol:   protocol,
		})

		return err
	}
}

// getRoute 获取路由路径
// 优先使用 Route().Path 作为 route，避免高基数问题
func getRoute(c *fiber.Ctx) string {
	// 优先使用路由模板路径
	if c.Route().Path != "" {
		return c.Route().Path
	}
	// Fallback 到路径
	return c.Path()
}

// getProtocol 获取协议版本
func getProtocol(c *fiber.Ctx) string {
	proto := c.Protocol()
	if proto == "" {
		// 从请求中检测
		if strings.Contains(string(c.Request().Header.Protocol()), "1.1") {
			return "http/1.1"
		}
		if strings.Contains(string(c.Request().Header.Protocol()), "2") {
			return "http/2"
		}
		return "http/1.1"
	}
	return proto
}

// MetricsWithStatusCode 创建 metrics 中间件，支持自定义状态码提取
// 当业务逻辑在 c.Next() 后才设置状态码时使用
func MetricsWithStatusCode(recorder *metrics.HTTPRecorder, getStatusCode func(*fiber.Ctx) int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		route := getRoute(c)
		recorder.ActiveRequestAdd(c.Context(), 1)

		err := c.Next()
		defer recorder.ActiveRequestAdd(c.Context(), -1)

		statusCode := getStatusCode(c)
		durationMs := float64(time.Since(start).Milliseconds())
		protocol := getProtocol(c)
		hasError := err != nil || statusCode >= 500

		recorder.RecordRequest(c.Context(), metrics.HTTPRequestInfo{
			Method:     c.Method(),
			Route:      route,
			StatusCode: statusCode,
			DurationMs: durationMs,
			HasError:   hasError,
			Protocol:   protocol,
		})

		return err
	}
}

// GetStatusCodeFromResponse 从响应获取状态码
func GetStatusCodeFromResponse() func(*fiber.Ctx) int {
	return func(c *fiber.Ctx) int {
		return c.Response().StatusCode()
	}
}

// GetStatusCodeFromHeader 从响应头获取状态码
func GetStatusCodeFromHeader() func(*fiber.Ctx) int {
	return func(c *fiber.Ctx) int {
		return c.Response().StatusCode()
	}
}
