package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UsernameParams 路径参数：用户名
type UsernameParams struct {
	Username string `params:"username" validate:"required"`
}

// UserGetByUsername 根据用户名获取用户
// GET /api/users/username/:username
func (h *handler) UserGetByUsername(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.UserGetByUsername")
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
