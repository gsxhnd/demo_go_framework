package health

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) Check(c *fiber.Ctx) error {
	decoded := errno.Decode("ok", nil)
	return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
}
