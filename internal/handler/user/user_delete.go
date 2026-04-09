package user

import (
	"strconv"

	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Delete 删除用户
// DELETE /api/users/:id
func (h *handler) Delete(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.Delete")
	defer span.End()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.WarnCtx(ctx, "invalid user id", zap.String("id", idStr))
		return c.Status(errno.InvalidUserIDError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.InvalidUserIDError))
	}

	errNo := h.userService.DeleteUser(ctx, id)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to delete user", zap.Int("id", id), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(fiber.Map{"message": "user deleted successfully"}, nil))
}
