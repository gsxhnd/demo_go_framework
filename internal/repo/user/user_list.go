package user

import (
	"context"

	"go_sample_code/internal/ent"
	"go_sample_code/internal/ent/user"

	"go.uber.org/zap"
)

// List 分页获取用户列表
func (r *userRepo) List(ctx context.Context, req *ListUsersRequest) ([]*ent.User, int, error) {
	ctx, end := r.withTrace(ctx, "UserRepo.List")
	defer end(nil)

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	query := r.client.User.Query()

	// 关键字搜索
	if req.Keyword != "" {
		query = query.Where(
			user.Or(
				user.UsernameContains(req.Keyword),
				user.EmailContains(req.Keyword),
				user.NicknameContains(req.Keyword),
			),
		)
	}

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil {
		r.log.ErrorCtx(ctx, "failed to count users",
			zap.Error(err),
		)
		return nil, 0, err
	}

	// 获取分页数据
	users, err := query.
		Order(ent.Desc(user.FieldCreatedAt)).
		Offset(offset).
		Limit(pageSize).
		All(ctx)

	if err != nil {
		r.log.ErrorCtx(ctx, "failed to list users",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
		return nil, 0, err
	}

	r.log.InfoCtx(ctx, "users listed successfully",
		zap.Int("total", total),
		zap.Int("count", len(users)),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
	)

	return users, total, nil
}
