package user

import (
	"context"

	"go_sample_code/internal/ent"
	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"
	"go_sample_code/pkg/logger"

	// "go_sample_code/pkg/trace"
	otel_trace "go.opentelemetry.io/otel/trace"
)

// UserService 用户服务接口
type UserService interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *userrepo.CreateUserRequest) (*UserResponse, errno.Errno)
	// GetUserByID 根据 ID 获取用户
	GetUserByID(ctx context.Context, id int) (*UserResponse, errno.Errno)
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*UserResponse, errno.Errno)
	// GetUserByEmail 根据邮箱获取用户
	GetUserByEmail(ctx context.Context, email string) (*UserResponse, errno.Errno)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, id int, req *userrepo.UpdateUserRequest) (*UserResponse, errno.Errno)
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id int) errno.Errno
	// ListUsers 分页获取用户列表
	ListUsers(ctx context.Context, req *userrepo.ListUsersRequest) (*ListUsersResponse, errno.Errno)
}

// userService 用户服务实现
type userService struct {
	userRepo userrepo.UserRepo
	log      logger.Logger
	tracer   otel_trace.Tracer
}

// NewUserService 创建用户服务实例（带追踪）
func NewUserService(userRepo userrepo.UserRepo, l logger.Logger, t otel_trace.Tracer) UserService {
	return &userService{
		userRepo: userRepo,
		log:      l,
		tracer:   t,
	}
}

// UserResponse 用户响应
type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Phone     string `json:"phone,omitempty"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Users     []*UserResponse `json:"users"`
	Total     int             `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}

// toUserResponse 将 ent.User 转换为 UserResponse
func toUserResponse(u *ent.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
