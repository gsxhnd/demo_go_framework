package user

import (
	"strconv"

	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetByID 根据 ID 获取用户
// GET /api/users/:id
func (h *handler) GetByID(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.GetByID")
	defer span.End()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.WarnCtx(ctx, "invalid user id", zap.String("id", idStr))
		return c.Status(errno.InvalidUserIDError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.InvalidUserIDError))
	}

	result, errNo := h.userService.GetUserByID(ctx, id)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to get user by id", zap.Int("id", id), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
