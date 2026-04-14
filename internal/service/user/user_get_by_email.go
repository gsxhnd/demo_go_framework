package user

import (
	"context"

	"go_sample_code/internal/errno"

	"go.uber.org/zap"
)

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*UserResponse, errno.Errno) {
	ctx, span := s.tracer.Start(ctx, "UserService.GetUserByEmail")
	defer span.End()

	if email == "" {
		s.log.WarnCtx(ctx, "empty email")
		return nil, errno.InvalidEmailError
	}

	s.log.DebugCtx(ctx, "getting user by email", zap.String("email", email))

	u, err := s.userRepo.UserGetByEmail(ctx, email)
	if err != nil {
		s.log.WarnCtx(ctx, "user not found", zap.String("email", email))
		return nil, errno.UserNotFoundError
	}

	return toUserResponse(u), errno.OK
}
