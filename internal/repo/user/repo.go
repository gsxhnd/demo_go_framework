package user

import (
	"context"

	"go_sample_code/internal/ent"
	"go_sample_code/pkg/logger"

	otel_trace "go.opentelemetry.io/otel/trace"
)

// UserRepo 用户仓储接口
type UserRepo interface {
	// UserCreate 创建用户
	UserCreate(ctx context.Context, req *CreateUserRequest) (*ent.User, error)
	// UserGetByID 根据 ID 获取用户
	UserGetByID(ctx context.Context, id int) (*ent.User, error)
	// UserGetByUsername 根据用户名获取用户
	UserGetByUsername(ctx context.Context, username string) (*ent.User, error)
	// UserGetByEmail 根据邮箱获取用户
	UserGetByEmail(ctx context.Context, email string) (*ent.User, error)
	// UserUpdate 更新用户
	UserUpdate(ctx context.Context, id int, req *UpdateUserRequest) (*ent.User, error)
	// UserDelete 删除用户
	UserDelete(ctx context.Context, id int) error
	// UserList 分页获取用户列表
	UserList(ctx context.Context, req *ListUsersRequest) ([]*ent.User, int, error)
	// UserExistsByUsername 检查用户名是否存在
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
	// UserExistsByEmail 检查邮箱是否存在
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}

// userRepo Ent 用户仓储实现
type userRepo struct {
	client *ent.Client
	log    logger.Logger
	tracer otel_trace.Tracer
}

// NewUserRepo 创建用户仓储实例
func NewUserRepo(client *ent.Client, l logger.Logger, tracer otel_trace.Tracer) UserRepo {
	return &userRepo{
		client: client,
		log:    l,
		tracer: tracer,
	}
}
