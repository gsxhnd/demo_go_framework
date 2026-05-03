package health

import (
	"go_sample_code/internal/database"
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// @Summary      健康检查
// @Description  检查数据库和 Redis 连接状态
// @Tags         健康检查
// @Accept       json
// @Produce      json
// @Success      200  {object}  SwaggerHealthResponse  "服务正常"
// @Failure      503  {object}  SwaggerHealthResponse  "服务不可用"
// @Router       /health [get]
func (h *handler) Check(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "health_heandler.check")
	defer span.End()
	c.SetUserContext(ctx)

	status := h.healthChecker.Check(ctx)

	if status.Data.Status != database.StatusOK {
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
