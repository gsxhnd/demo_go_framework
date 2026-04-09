package user

import (
	"context"

	"go_sample_code/internal/ent"
	"go_sample_code/internal/ent/user"

	"go.uber.org/zap"
)

// GetByEmail 根据邮箱获取用户
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.GetByEmail")
	defer span.End()

	u, err := r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)

	if err != nil {
		r.log.WarnCtx(ctx, "user not found by email",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, err
	}

	return u, nil
}
