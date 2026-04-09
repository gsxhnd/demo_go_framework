package health

import (
	"go_sample_code/internal/database"
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *handler) Check(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "health_heandler.check")
	defer span.End()
	c.SetUserContext(ctx)

	status := h.healthChecker.Check(ctx)

	if status.Data.Status != database.StatusDegraded {
		h.log.ErrorCtx(ctx, "health_check_failed", zap.Any("status", status))
	}

	// Determine HTTP status based on health status
	httpStatus := fiber.StatusOK
	if status.Data.Status != database.StatusOK {
		httpStatus = fiber.StatusServiceUnavailable
	}

	decoded := errno.Decode(status.Data, nil)
	return c.Status(httpStatus).JSON(decoded)
}
