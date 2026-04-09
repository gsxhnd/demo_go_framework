package user

import (
	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Delete 删除用户
// DELETE /api/users/:id
func (h *handler) Delete(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.Delete")
	defer span.End()

	var params UserIDParams
	if h.parseAndValidateParams(c, &params) {
		return nil
	}

	errNo := h.userService.DeleteUser(ctx, params.ID)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to delete user", zap.Int("id", params.ID), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(fiber.Map{"message": "user deleted successfully"}, nil))
}
