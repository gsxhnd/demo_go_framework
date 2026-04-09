package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetByID 根据 ID 获取用户
// GET /api/users/:id
func (h *handler) GetByID(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.GetByID")
	defer span.End()

	var params UserIDParams
	if h.parseAndValidateParams(c, &params) {
		return nil
	}

	result, errNo := h.userService.GetUserByID(ctx, params.ID)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by id", zap.Int("id", params.ID), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
