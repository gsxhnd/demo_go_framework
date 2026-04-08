package user

import (
	"context"

	"go_sample_code/internal/ent"

	"go.uber.org/zap"
)

// Create 创建用户
func (r *userRepo) Create(ctx context.Context, req *CreateUserRequest) (*ent.User, error) {
	ctx, end := r.withTrace(ctx, "UserRepo.Create")
	defer end(nil)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	created, err := r.client.User.Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(req.Password).
		SetNickname(req.Nickname).
		SetAvatar(req.Avatar).
		SetPhone(req.Phone).
		SetIsActive(isActive).
		Save(ctx)

	if err != nil {
		r.log.ErrorCtx(ctx, "failed to create user",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}

	r.log.InfoCtx(ctx, "user created successfully",
		zap.Int("user_id", created.ID),
		zap.String("username", created.Username),
	)

	return created, nil
}
