package user

import (
	"context"

	"go_sample_code/internal/errno"

	"go.uber.org/zap"
)

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id int) errno.Errno {
	ctx, span := s.tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()

	if id <= 0 {
		s.log.WarnCtx(ctx, "invalid user id", zap.Int("id", id))
		return errno.InvalidUserIDError
	}

	s.log.InfoCtx(ctx, "deleting user", zap.Int("id", id))

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		s.log.ErrorCtx(ctx, "failed to delete user", zap.Error(err))
		return errno.UserDeleteFailedError
	}

	s.log.InfoCtx(ctx, "user deleted successfully", zap.Int("id", id))

	return errno.OK
}
