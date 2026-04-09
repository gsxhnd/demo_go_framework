package user

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname,omitempty" validate:"omitempty,max=64"`
	Avatar   string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,max=32"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,max=64"`
	Avatar   *string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=32"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// ListUsersRequest 分页查询请求
type ListUsersRequest struct {
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
	Keyword  string `query:"keyword" validate:"max=128"`
}

// UserIDParams 路径参数：用户 ID
type UserIDParams struct {
	ID int `params:"id" validate:"required,gt=0"`
}

// UsernameParams 路径参数：用户名
type UsernameParams struct {
	Username string `params:"username" validate:"required"`
}

// EmailParams 路径参数：邮箱
type EmailParams struct {
	Email string `params:"email" validate:"required,email"`
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
