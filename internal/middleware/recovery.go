package middleware

import (
	"go_sample_code/internal/errno"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Path()),
				)

				err := errno.InternalServerError
				decoded := errno.Decode(nil, err)
				c.Status(decoded.GetHTTPStatus()).JSON(decoded)
			}
		}()

		return c.Next()
	}
}
