package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// @Summary      根据用户名获取用户
// @Description  根据用户名获取用户详情
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        username  path      string                true  "用户名"
// @Success      200       {object}  SwaggerUserResponse   "获取成功"
// @Failure      400       {object}  SwaggerErrorResponse  "参数校验失败"
// @Failure      404       {object}  SwaggerErrorResponse  "用户不存在"
// @Router       /users/username/{username} [get]
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
