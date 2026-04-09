package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetByUsername 根据用户名获取用户
// GET /api/users/username/:username
func (h *handler) GetByUsername(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.GetByUsername")
	defer span.End()

	var params UsernameParams
	if h.parseAndValidateParams(c, &params) {
		return nil
	}

	result, errNo := h.userService.GetUserByUsername(ctx, params.Username)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by username", zap.String("username", params.Username), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
