package user

import (
	"context"

	"go.uber.org/zap"
)

// Delete 删除用户
func (r *userRepo) Delete(ctx context.Context, id int) error {
	ctx, end := r.withTrace(ctx, "UserRepo.Delete")
	defer end(nil)

	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		r.log.ErrorCtx(ctx, "failed to delete user",
			zap.Int("user_id", id),
			zap.Error(err),
		)
		return err
	}

	r.log.InfoCtx(ctx, "user deleted successfully",
		zap.Int("user_id", id),
	)

	return nil
}
