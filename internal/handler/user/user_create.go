package user

import (
	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname,omitempty" validate:"omitempty,max=64"`
	Avatar   string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,max=32"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// @Summary      创建用户
// @Description  创建一个新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request  body      CreateUserRequest     true  "用户信息"
// @Success      201      {object}  SwaggerUserResponse   "创建成功"
// @Failure      400      {object}  SwaggerErrorResponse  "参数校验失败"
// @Failure      500      {object}  SwaggerErrorResponse  "服务器内部错误"
// @Router       /users [post]
func (h *handler) UserCreate(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.UserCreate")
	defer span.End()

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.ErrorCtx(ctx, "failed to parse request body", zap.Error(err))
		_ = c.Status(errno.RequestParserError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestParserError))
		return nil
	}
	if err := h.validate.Struct(req); err != nil {
		h.log.WarnCtx(ctx, "request validation failed", zap.Error(err))
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
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
