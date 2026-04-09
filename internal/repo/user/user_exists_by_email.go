package user

import (
	"context"

	"go_sample_code/internal/ent/user"

	"go.uber.org/zap"
)

// ExistsByEmail 检查邮箱是否存在
func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.ExistsByEmail")
	defer span.End()

	count, err := r.client.User.Query().
		Where(user.Email(email)).
		Count(ctx)

	if err != nil {
		r.log.ErrorCtx(ctx, "failed to check email exists",
			zap.String("email", email),
			zap.Error(err),
		)
		return false, err
	}

	return count > 0, nil
}
