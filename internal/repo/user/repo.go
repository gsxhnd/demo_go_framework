package user

import (
	"context"

	"go_sample_code/internal/ent"
	"go_sample_code/pkg/logger"

	otel_trace "go.opentelemetry.io/otel/trace"
)

// UserRepo 用户仓储接口
type UserRepo interface {
	// Create 创建用户
	Create(ctx context.Context, req *CreateUserRequest) (*ent.User, error)
	// GetByID 根据 ID 获取用户
	GetByID(ctx context.Context, id int) (*ent.User, error)
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*ent.User, error)
	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	// Update 更新用户
	Update(ctx context.Context, id int, req *UpdateUserRequest) (*ent.User, error)
	// Delete 删除用户
	Delete(ctx context.Context, id int) error
	// List 分页获取用户列表
	List(ctx context.Context, req *ListUsersRequest) ([]*ent.User, int, error)
	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname,omitempty" validate:"omitempty,max=100"`
	Avatar   string `json:"avatar,omitempty" validate:"omitempty,max=500"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,max=20"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,max=100"`
	Avatar   *string `json:"avatar,omitempty" validate:"omitempty,max=500"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// ListUsersRequest 分页查询请求
type ListUsersRequest struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword,omitempty"`
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

// withTrace 执行带追踪的操作
func (r *userRepo) withTrace(ctx context.Context, spanName string) (context.Context, func(err error)) {
	if r.tracer != nil {
		ctx, span := r.tracer.Start(ctx, spanName)
		return ctx, func(err error) {
			if err != nil {
				span.RecordError(err)
				span.End()
			} else {
				span.End()
			}
		}
	}
	return ctx, func(err error) {}
}
