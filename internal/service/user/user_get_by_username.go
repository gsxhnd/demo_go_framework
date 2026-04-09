package user

import (
	"context"

	"go_sample_code/internal/errno"

	"go.uber.org/zap"
)

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*UserResponse, errno.Errno) {
	ctx, end := s.withTrace(ctx, "UserService.GetUserByUsername")
	defer end(nil)

	if username == "" {
		s.log.WarnCtx(ctx, "empty username")
		return nil, errno.InvalidUsernameError
	}

	s.log.DebugCtx(ctx, "getting user by username", zap.String("username", username))

	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.log.WarnCtx(ctx, "user not found", zap.String("username", username))
		return nil, errno.UserNotFoundError
	}

	return toUserResponse(u), errno.OK
}
