package user

import (
	"context"

	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"go.uber.org/zap"
)

// ListUsers 分页获取用户列表
func (s *userService) ListUsers(ctx context.Context, req *userrepo.ListUsersRequest) (*ListUsersResponse, errno.Errno) {
	ctx, span := s.tracer.Start(ctx, "UserService.ListUsers")
	defer span.End()

	s.log.DebugCtx(ctx, "listing users",
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
	)

	users, total, err := s.userRepo.UserList(ctx, req)
	if err != nil {
		s.log.ErrorCtx(ctx, "failed to list users", zap.Error(err))
		return nil, errno.DatabaseError
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	totalPage := total / pageSize
	if total%pageSize > 0 {
		totalPage++
	}

	userResponses := make([]*UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, toUserResponse(u))
	}

	s.log.InfoCtx(ctx, "users listed successfully",
		zap.Int("total", total),
		zap.Int("count", len(users)),
	)

	return &ListUsersResponse{
		Users:     userResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, errno.OK
}
