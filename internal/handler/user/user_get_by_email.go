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

	email := c.Params("email")
	if email == "" {
		h.log.WarnCtx(ctx, "email is empty")
		return c.Status(errno.InvalidEmailError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.InvalidEmailError))
	}

	result, errNo := h.userService.GetUserByEmail(ctx, email)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by email", zap.String("email", email), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
