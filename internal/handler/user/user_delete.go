package user

import (
	"go_sample_code/internal/errno"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserDelete 删除用户
// DELETE /api/users/:id
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
