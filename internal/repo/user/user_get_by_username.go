package user

import (
	"context"

	"go_sample_code/internal/ent"
	"go_sample_code/internal/ent/user"

	"go.uber.org/zap"
)

// GetByUsername 根据用户名获取用户
func (r *userRepo) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	ctx, end := r.withTrace(ctx, "UserRepo.GetByUsername")
	defer end(nil)

	u, err := r.client.User.Query().
		Where(user.Username(username)).
		Only(ctx)

	if err != nil {
		r.log.WarnCtx(ctx, "user not found by username",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, err
	}

	return u, nil
}
