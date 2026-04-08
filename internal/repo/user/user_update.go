package user

import (
	"context"

	"go_sample_code/internal/ent"

	"go.uber.org/zap"
)

// Update 更新用户
func (r *userRepo) Update(ctx context.Context, id int, req *UpdateUserRequest) (*ent.User, error) {
	ctx, end := r.withTrace(ctx, "UserRepo.Update")
	defer end(nil)

	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		r.log.WarnCtx(ctx, "user not found for update",
			zap.Int("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	update := u.Update()
	if req.Email != nil {
		update.SetEmail(*req.Email)
	}
	if req.Password != nil {
		update.SetPassword(*req.Password)
	}
	if req.Nickname != nil {
		update.SetNickname(*req.Nickname)
	}
	if req.Avatar != nil {
		update.SetAvatar(*req.Avatar)
	}
	if req.Phone != nil {
		update.SetPhone(*req.Phone)
	}
	if req.IsActive != nil {
		update.SetIsActive(*req.IsActive)
	}

	updated, err := update.Save(ctx)
	if err != nil {
		r.log.ErrorCtx(ctx, "failed to update user",
			zap.Int("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	r.log.InfoCtx(ctx, "user updated successfully",
		zap.Int("user_id", id),
	)

	return updated, nil
}
