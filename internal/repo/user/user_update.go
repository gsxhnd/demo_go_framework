package user

import (
	"context"

	"go_sample_code/internal/ent"

	"go.uber.org/zap"
)

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,max=100"`
	Avatar   *string `json:"avatar,omitempty" validate:"omitempty,max=500"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UserUpdate 更新用户
func (r *userRepo) UserUpdate(ctx context.Context, id int, req *UpdateUserRequest) (*ent.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.UserUpdate")
	defer span.End()

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
