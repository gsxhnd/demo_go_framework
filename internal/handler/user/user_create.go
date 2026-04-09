package user

import (
	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Create 创建用户
// POST /api/users
func (h *handler) Create(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.Create")
	defer span.End()

	var req CreateUserRequest
	if h.parseAndValidateBody(c, &req) {
		return nil
	}

	repoReq := &userrepo.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Phone:    req.Phone,
		IsActive: req.IsActive,
	}

	result, errNo := h.userService.CreateUser(ctx, repoReq)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to create user", zap.Int("code", errNo.GetCode()), zap.String("message", errNo.GetMessage()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusCreated).JSON(errno.Decode(result, nil))
}
