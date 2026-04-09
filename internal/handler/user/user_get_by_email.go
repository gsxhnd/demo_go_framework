package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetByEmail 根据邮箱获取用户
// GET /api/users/email/:email
func (h *handler) GetByEmail(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.GetByEmail")
	defer span.End()

	var params EmailParams
	if h.parseAndValidateParams(c, &params) {
		return nil
	}

	result, errNo := h.userService.GetUserByEmail(ctx, params.Email)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by email", zap.String("email", params.Email), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
