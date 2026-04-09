package user

import (
	"context"

	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"go.uber.org/zap"
)

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, id int, req *userrepo.UpdateUserRequest) (*UserResponse, errno.Errno) {
	ctx, end := s.withTrace(ctx, "UserService.UpdateUser")
	defer end(nil)

	if id <= 0 {
		s.log.WarnCtx(ctx, "invalid user id", zap.Int("id", id))
		return nil, errno.InvalidUserIDError
	}

	s.log.InfoCtx(ctx, "updating user", zap.Int("id", id))

	// 如果更新邮箱，检查邮箱是否已被使用
	if req.Email != nil {
		exists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			s.log.ErrorCtx(ctx, "failed to check email exists", zap.Error(err))
			return nil, errno.DatabaseError
		}
		if exists {
			// 检查是否是自己的邮箱
			u, err := s.userRepo.GetByID(ctx, id)
			if err == nil && u.Email != *req.Email {
				s.log.WarnCtx(ctx, "email already in use", zap.String("email", *req.Email))
				return nil, errno.UserAlreadyExistsError
			}
		}
	}

	u, err := s.userRepo.Update(ctx, id, req)
	if err != nil {
		s.log.ErrorCtx(ctx, "failed to update user", zap.Error(err))
		return nil, errno.UserUpdateFailedError
	}

	s.log.InfoCtx(ctx, "user updated successfully", zap.Int("id", id))

	return toUserResponse(u), errno.OK
}
