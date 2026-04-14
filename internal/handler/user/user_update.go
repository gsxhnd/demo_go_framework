package user

import (
	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,max=64"`
	Avatar   *string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=32"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UserUpdate 更新用户
// PUT /api/users/:id
func (h *handler) UserUpdate(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.UserUpdate")
	defer span.End()

	var idParams UserIDParams
	if h.parseAndValidateParams(c, &idParams) {
		return nil
	}

	var req UpdateUserRequest
	if h.parseAndValidateBody(c, &req) {
		return nil
	}

	repoReq := &userrepo.UpdateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Phone:    req.Phone,
		IsActive: req.IsActive,
	}

	result, errNo := h.userService.UpdateUser(ctx, idParams.ID, repoReq)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to update user", zap.Int("id", idParams.ID), zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
