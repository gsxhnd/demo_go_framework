package user

import (
	"go_sample_code/internal/errno"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// @Summary      删除用户
// @Description  根据用户 ID 删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int                    true  "用户 ID"
// @Success      200  {object}  SwaggerDeleteResponse  "删除成功"
// @Failure      400  {object}  SwaggerErrorResponse   "参数校验失败"
// @Failure      404  {object}  SwaggerErrorResponse   "用户不存在"
// @Router       /users/{id} [delete]
func (h *handler) UserDelete(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.UserDelete")
	defer span.End()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.log.WarnCtx(ctx, "invalid id param", zap.String("id", idStr))
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return nil
	}

	errNo := h.userService.DeleteUser(ctx, id)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to delete user", zap.Int("id", id), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(fiber.Map{"message": "user deleted successfully"}, nil))
}
