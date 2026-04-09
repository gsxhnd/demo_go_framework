package health

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) Check(c *fiber.Ctx) error {
	status := h.healthChecker.Check(c.Context())

	// Determine HTTP status based on health status
	httpStatus := fiber.StatusOK
	if status.Data.Status != "ok" {
		httpStatus = fiber.StatusServiceUnavailable
	}

	decoded := errno.Decode(status.Data, nil)
	return c.Status(httpStatus).JSON(decoded)
}
