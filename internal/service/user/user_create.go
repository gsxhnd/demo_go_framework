package user

import (
	"context"

	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"

	"go.uber.org/zap"
)

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *userrepo.CreateUserRequest) (*UserResponse, errno.Errno) {
	ctx, end := s.withTrace(ctx, "UserService.CreateUser")
	defer end(nil)

	s.log.InfoCtx(ctx, "creating user",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.log.ErrorCtx(ctx, "failed to check username exists", zap.Error(err))
		return nil, errno.DatabaseError
	}
	if exists {
		s.log.WarnCtx(ctx, "username already exists", zap.String("username", req.Username))
		return nil, errno.UserAlreadyExistsError.WithData("username")
	}

	// 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.log.ErrorCtx(ctx, "failed to check email exists", zap.Error(err))
		return nil, errno.DatabaseError
	}
	if exists {
		s.log.WarnCtx(ctx, "email already exists", zap.String("email", req.Email))
		return nil, errno.UserAlreadyExistsError.WithData("email")
	}

	// 创建用户
	u, err := s.userRepo.Create(ctx, req)
	if err != nil {
		s.log.ErrorCtx(ctx, "failed to create user", zap.Error(err))
		return nil, errno.UserCreateFailedError
	}

	s.log.InfoCtx(ctx, "user created successfully",
		zap.Int("user_id", u.ID),
	)

	return toUserResponse(u), errno.OK
}
