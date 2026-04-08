package user

import (
	"context"

	"go_sample_code/internal/ent/user"

	"go.uber.org/zap"
)

// ExistsByEmail 检查邮箱是否存在
func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, end := r.withTrace(ctx, "UserRepo.ExistsByEmail")
	defer end(nil)

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
