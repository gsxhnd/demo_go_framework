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

		log.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
		)

		return err
	}
}
