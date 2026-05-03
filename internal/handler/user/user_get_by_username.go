package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserGetByUsername 根据用户名获取用户
// GET /api/users/username/:username
func (h *handler) UserGetByUsername(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.UserGetByUsername")
	defer span.End()

	username := c.Params("username")
	if username == "" {
		h.log.WarnCtx(ctx, "username param is required")
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return nil
	}

	result, errNo := h.userService.GetUserByUsername(ctx, username)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by username", zap.String("username", username), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
