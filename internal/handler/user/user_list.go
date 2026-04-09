package user

import (
	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// List 分页获取用户列表
// GET /api/users?page=1&page_size=10&keyword=xxx
func (h *handler) List(c *fiber.Ctx) error {
	ctx, span := h.tracer.Start(c.UserContext(), "UserHandler.List")
	defer span.End()

	var req ListUsersRequest
	if err := c.QueryParser(&req); err != nil {
		h.log.ErrorCtx(ctx, "failed to parse query parameters", zap.Error(err))
		_ = c.Status(errno.RequestParserError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestParserError))
		return nil
	}

	// 应用默认值（只有未传参数时才使用默认值）
	if c.Query("page") == "" {
		req.Page = 1
	}
	if c.Query("page_size") == "" {
		req.PageSize = 10
	}

	// 校验（此时显式非法值会触发校验失败）
	if err := h.validate.Struct(req); err != nil {
		h.log.WarnCtx(ctx, "query validation failed", zap.Error(err))
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return nil
	}

	repoReq := &userrepo.ListUsersRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}

	result, errNo := h.userService.ListUsers(ctx, repoReq)
	if errNo.GetCode() != errno.OK.Code {
		h.log.ErrorCtx(ctx, "failed to list users", zap.Int("code", errNo.GetCode()))
		return c.Status(errNo.GetHTTPStatus()).JSON(errno.Decode(nil, errNo))
	}

	return c.Status(fiber.StatusOK).JSON(errno.Decode(result, nil))
}
