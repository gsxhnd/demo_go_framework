package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// @Summary      根据邮箱获取用户
// @Description  根据邮箱地址获取用户详情
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        email  path      string                true  "邮箱地址"
// @Success      200    {object}  SwaggerUserResponse   "获取成功"
// @Failure      400    {object}  SwaggerErrorResponse  "参数校验失败"
// @Failure      404    {object}  SwaggerErrorResponse  "用户不存在"
// @Router       /users/email/{email} [get]
func (h *handler) UserGetByEmail(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.UserGetByEmail")
	defer span.End()

	email := c.Params("email")
	if email == "" {
		h.log.WarnCtx(ctx, "email param is required")
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return nil
	}

	result, errNo := h.userService.GetUserByEmail(ctx, email)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by email", zap.String("email", email), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
