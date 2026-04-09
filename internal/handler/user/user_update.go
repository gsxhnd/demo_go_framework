package user

import (
	"strconv"

	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Update 更新用户
// PUT /api/users/:id
func (h *handler) Update(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.Update")
	defer span.End()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.WarnCtx(ctx, "invalid user id", zap.String("id", idStr))
		return c.Status(errno.InvalidUserIDError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.InvalidUserIDError))
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.ErrorCtx(ctx, "failed to parse update user request", zap.Error(err))
		return c.Status(errno.RequestParserError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestParserError))
	}

	repoReq := &userrepo.UpdateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Phone:    req.Phone,
		IsActive: req.IsActive,
	}

	result, errNo := h.userService.UpdateUser(ctx, id, repoReq)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to update user", zap.Int("id", id), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
