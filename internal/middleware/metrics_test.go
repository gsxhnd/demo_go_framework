package middleware

import (
	"net/http/httptest"
	"testing"

	"go_sample_code/pkg/metrics"

	"github.com/gofiber/fiber/v2"
)

func TestMetrics_Middleware(t *testing.T) {
	// 创建 MeterProvider
	cfg := metrics.DefaultMetricsConfig()
	cfg.OtelEnable = false
	mp, err := metrics.NewMeterProvider(&cfg)
	if err != nil {
		t.Fatalf("failed to create meter provider: %v", err)
	}
	defer mp.Shutdown(nil)

	// 创建 HTTPRecorder
	recorder, err := metrics.NewHTTPRecorder(mp)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}

	// 创建 Fiber app
	app := fiber.New()

	// 注册中间件
	app.Use(Metrics(recorder))

	// 注册测试路由
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
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// 测试 POST 请求
	req = httptest.NewRequest("POST", "/test", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// 测试错误请求
	req = httptest.NewRequest("GET", "/error", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestMetrics_RouteTemplate(t *testing.T) {
	cfg := metrics.DefaultMetricsConfig()
	cfg.OtelEnable = false
	mp, err := metrics.NewMeterProvider(&cfg)
	if err != nil {
		t.Fatalf("failed to create meter provider: %v", err)
	}
	defer mp.Shutdown(nil)

	recorder, err := metrics.NewHTTPRecorder(mp)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}

	app := fiber.New()
	app.Use(Metrics(recorder))

	// 测试路由模板
	app.Get("/users/:id/posts/:postId", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetRoute(t *testing.T) {
	app := fiber.New()

	app.Get("/api/v1/users/:id", func(c *fiber.Ctx) error {
		route := getRoute(c)
		if route != "/api/v1/users/:id" {
			t.Errorf("expected route '/api/v1/users/:id', got '%s'", route)
		}
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/api/v1/users/123", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
}

func TestGetProtocol(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		proto := getProtocol(c)
		// 默认为 http/1.1
		if proto != "http/1.1" {
			t.Logf("got protocol: %s", proto)
		}
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
}

func TestGetStatusCodeFromResponse(t *testing.T) {
	app := fiber.New()
	statusCodeCapture := 0

	app.Get("/ok", func(c *fiber.Ctx) error {
		// 在处理函数内部设置状态码后再获取
		c.Status(200)
		statusCodeCapture = GetStatusCodeFromResponse()(c)
		return c.SendString("OK")
	})
	app.Get("/notfound", func(c *fiber.Ctx) error {
		c.Status(404)
		statusCodeCapture = GetStatusCodeFromResponse()(c)
		return c.SendString("Not Found")
	})

	// 测试正常状态码
	req := httptest.NewRequest("GET", "/ok", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
	if statusCodeCapture != 200 {
		t.Errorf("expected captured status code 200, got %d", statusCodeCapture)
	}

	// 测试 404
	req = httptest.NewRequest("GET", "/notfound", nil)
	_, err = app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}
	if statusCodeCapture != 404 {
		t.Errorf("expected captured status code 404, got %d", statusCodeCapture)
	}
}
