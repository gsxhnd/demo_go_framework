package user

import (
	"context"

	"go_sample_code/internal/ent/user"

	"go.uber.org/zap"
)

// ExistsByUsername 检查用户名是否存在
func (r *userRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	ctx, end := r.withTrace(ctx, "UserRepo.ExistsByUsername")
	defer end(nil)

	count, err := r.client.User.Query().
		Where(user.Username(username)).
		Count(ctx)

	if err != nil {
		r.log.ErrorCtx(ctx, "failed to check username exists",
			zap.String("username", username),
			zap.Error(err),
		)
		return false, err
	}

	return count > 0, nil
}
