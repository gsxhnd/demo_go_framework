package user

import (
	"context"

	"go_sample_code/internal/errno"

	"go.uber.org/zap"
)

// GetUserByID 根据 ID 获取用户
func (s *userService) GetUserByID(ctx context.Context, id int) (*UserResponse, errno.Errno) {
	ctx, end := s.withTrace(ctx, "UserService.GetUserByID")
	defer end(nil)

	if id <= 0 {
		s.log.WarnCtx(ctx, "invalid user id", zap.Int("id", id))
		return nil, errno.InvalidUserIDError
	}

	s.log.DebugCtx(ctx, "getting user by id", zap.Int("id", id))

	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.WarnCtx(ctx, "user not found", zap.Int("id", id))
		return nil, errno.UserNotFoundError.WithData(id)
	}

	return toUserResponse(u), errno.OK
}
