package user

import (
	"context"

	"go_sample_code/internal/ent"

	"go.uber.org/zap"
)

// GetByID 根据 ID 获取用户
func (r *userRepo) GetByID(ctx context.Context, id int) (*ent.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.GetByID")
	defer span.End()

	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		r.log.WarnCtx(ctx, "user not found by id",
			zap.Int("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return u, nil
}
