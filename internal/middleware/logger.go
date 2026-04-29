package middleware

import (
	"time"

	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Logger(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		// Use InfoCtx so the otelzap bridge picks up the active span's
		// TraceId/SpanId and writes them into the OTel Log record.
		log.InfoCtx(c.UserContext(), "",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("client_ip", c.IP()),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		)

		return err
	}
}
