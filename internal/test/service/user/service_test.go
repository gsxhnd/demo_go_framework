package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go_sample_code/internal/ent"
	"go_sample_code/internal/errno"
	userrepo "go_sample_code/internal/repo/user"
	userservice "go_sample_code/internal/service/user"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUserRepo 是用于测试的 UserRepo mock 实现
type mockUserRepo struct {
	createFn           func(ctx context.Context, req *userrepo.CreateUserRequest) (*ent.User, error)
	getByIDFn          func(ctx context.Context, id int) (*ent.User, error)
	getByUsernameFn    func(ctx context.Context, username string) (*ent.User, error)
	getByEmailFn       func(ctx context.Context, email string) (*ent.User, error)
	updateFn           func(ctx context.Context, id int, req *userrepo.UpdateUserRequest) (*ent.User, error)
	deleteFn           func(ctx context.Context, id int) error
	listFn             func(ctx context.Context, req *userrepo.ListUsersRequest) ([]*ent.User, int, error)
	existsByUsernameFn func(ctx context.Context, username string) (bool, error)
	existsByEmailFn    func(ctx context.Context, email string) (bool, error)
}

func (m *mockUserRepo) UserCreate(ctx context.Context, req *userrepo.CreateUserRequest) (*ent.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return &ent.User{ID: 1, Username: req.Username, Email: req.Email, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockUserRepo) UserGetByID(ctx context.Context, id int) (*ent.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return &ent.User{ID: id, Username: "testuser", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockUserRepo) UserGetByUsername(ctx context.Context, username string) (*ent.User, error) {
	if m.getByUsernameFn != nil {
		return m.getByUsernameFn(ctx, username)
	}
	return &ent.User{ID: 1, Username: username, Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockUserRepo) UserGetByEmail(ctx context.Context, email string) (*ent.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return &ent.User{ID: 1, Username: "testuser", Email: email, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockUserRepo) UserUpdate(ctx context.Context, id int, req *userrepo.UpdateUserRequest) (*ent.User, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return &ent.User{ID: id, Username: "updated", Email: "updated@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockUserRepo) UserDelete(ctx context.Context, id int) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockUserRepo) UserList(ctx context.Context, req *userrepo.ListUsersRequest) ([]*ent.User, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, req)
	}
	return []*ent.User{}, 0, nil
}

func (m *mockUserRepo) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.existsByUsernameFn != nil {
		return m.existsByUsernameFn(ctx, username)
	}
	return false, nil
}

func (m *mockUserRepo) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.existsByEmailFn != nil {
		return m.existsByEmailFn(ctx, email)
	}
	return false, nil
}

// newTestService 创建用于测试的 service 实例
func newTestService(t *testing.T, repo userrepo.UserRepo) userservice.UserService {
	logCfg := logger.DefaultConfig()
	log, err := logger.NewLogger(logCfg)
	require.NoError(t, err)

	tp, _ := trace.NewInMemoryProvider()
	tracer := trace.NewTracer(tp)

	return userservice.NewUserService(repo, log, tracer)
}

// =============================================================================
// CreateUser 测试
// =============================================================================

func TestCreateUser_Success(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.CreateUser(context.Background(), &userrepo.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, "newuser", result.Username)
	assert.Equal(t, "new@example.com", result.Email)
}

func TestCreateUser_UsernameExists(t *testing.T) {
	repo := &mockUserRepo{
		existsByUsernameFn: func(ctx context.Context, username string) (bool, error) {
			return true, nil
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.CreateUser(context.Background(), &userrepo.CreateUserRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.UserAlreadyExistsError.GetCode(), errNo.GetCode())
}

func TestCreateUser_EmailExists(t *testing.T) {
	repo := &mockUserRepo{
		existsByUsernameFn: func(ctx context.Context, username string) (bool, error) {
			return false, nil
		},
		existsByEmailFn: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.CreateUser(context.Background(), &userrepo.CreateUserRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.UserAlreadyExistsError.GetCode(), errNo.GetCode())
}

func TestCreateUser_CheckUsernameDBError(t *testing.T) {
	repo := &mockUserRepo{
		existsByUsernameFn: func(ctx context.Context, username string) (bool, error) {
			return false, errors.New("db error")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.CreateUser(context.Background(), &userrepo.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.DatabaseError.GetCode(), errNo.GetCode())
}

func TestCreateUser_CheckEmailDBError(t *testing.T) {
	repo := &mockUserRepo{
		existsByEmailFn: func(ctx context.Context, email string) (bool, error) {
			return false, errors.New("db error")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.CreateUser(context.Background(), &userrepo.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.DatabaseError.GetCode(), errNo.GetCode())
}

func TestCreateUser_CreateFailed(t *testing.T) {
	repo := &mockUserRepo{
		createFn: func(ctx context.Context, req *userrepo.CreateUserRequest) (*ent.User, error) {
			return nil, errors.New("insert failed")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.CreateUser(context.Background(), &userrepo.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.UserCreateFailedError.GetCode(), errNo.GetCode())
}

// =============================================================================
// GetUserByID 测试
// =============================================================================

func TestGetUserByID_Success(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByID(context.Background(), 1)
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, 1, result.ID)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByID(context.Background(), 0)
	assert.Nil(t, result)
	assert.Equal(t, errno.InvalidUserIDError.GetCode(), errNo.GetCode())
}

func TestGetUserByID_NegativeID(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByID(context.Background(), -1)
	assert.Nil(t, result)
	assert.Equal(t, errno.InvalidUserIDError.GetCode(), errNo.GetCode())
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id int) (*ent.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.GetUserByID(context.Background(), 999)
	assert.Nil(t, result)
	assert.Equal(t, errno.UserNotFoundError.GetCode(), errNo.GetCode())
}

// =============================================================================
// GetUserByEmail 测试
// =============================================================================

func TestGetUserByEmail_Success(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByEmail(context.Background(), "test@example.com")
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestGetUserByEmail_EmptyEmail(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByEmail(context.Background(), "")
	assert.Nil(t, result)
	assert.Equal(t, errno.InvalidEmailError.GetCode(), errNo.GetCode())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*ent.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.GetUserByEmail(context.Background(), "nonexistent@example.com")
	assert.Nil(t, result)
	assert.Equal(t, errno.UserNotFoundError.GetCode(), errNo.GetCode())
}

// =============================================================================
// GetUserByUsername 测试
// =============================================================================

func TestGetUserByUsername_Success(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByUsername(context.Background(), "testuser")
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, "testuser", result.Username)
}

func TestGetUserByUsername_EmptyUsername(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.GetUserByUsername(context.Background(), "")
	assert.Nil(t, result)
	assert.Equal(t, errno.InvalidUsernameError.GetCode(), errNo.GetCode())
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByUsernameFn: func(ctx context.Context, username string) (*ent.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.GetUserByUsername(context.Background(), "nonexistent")
	assert.Nil(t, result)
	assert.Equal(t, errno.UserNotFoundError.GetCode(), errNo.GetCode())
}

// =============================================================================
// UpdateUser 测试
// =============================================================================

func TestUpdateUser_Success(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	newNickname := "Updated Name"
	result, errNo := svc.UpdateUser(context.Background(), 1, &userrepo.UpdateUserRequest{
		Nickname: &newNickname,
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
}

func TestUpdateUser_InvalidID(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	newNickname := "Updated"
	result, errNo := svc.UpdateUser(context.Background(), 0, &userrepo.UpdateUserRequest{
		Nickname: &newNickname,
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.InvalidUserIDError.GetCode(), errNo.GetCode())
}

func TestUpdateUser_EmailConflict(t *testing.T) {
	repo := &mockUserRepo{
		existsByEmailFn: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
		getByIDFn: func(ctx context.Context, id int) (*ent.User, error) {
			return &ent.User{ID: 1, Email: "original@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	svc := newTestService(t, repo)

	newEmail := "conflict@example.com"
	result, errNo := svc.UpdateUser(context.Background(), 1, &userrepo.UpdateUserRequest{
		Email: &newEmail,
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.UserAlreadyExistsError.GetCode(), errNo.GetCode())
}

func TestUpdateUser_EmailSameAsOwn(t *testing.T) {
	repo := &mockUserRepo{
		existsByEmailFn: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
		getByIDFn: func(ctx context.Context, id int) (*ent.User, error) {
			return &ent.User{ID: 1, Email: "same@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	svc := newTestService(t, repo)

	newEmail := "same@example.com"
	result, errNo := svc.UpdateUser(context.Background(), 1, &userrepo.UpdateUserRequest{
		Email: &newEmail,
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
}

func TestUpdateUser_CheckEmailDBError(t *testing.T) {
	repo := &mockUserRepo{
		existsByEmailFn: func(ctx context.Context, email string) (bool, error) {
			return false, errors.New("db error")
		},
	}
	svc := newTestService(t, repo)

	newEmail := "test@example.com"
	result, errNo := svc.UpdateUser(context.Background(), 1, &userrepo.UpdateUserRequest{
		Email: &newEmail,
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.DatabaseError.GetCode(), errNo.GetCode())
}

func TestUpdateUser_UpdateFailed(t *testing.T) {
	repo := &mockUserRepo{
		updateFn: func(ctx context.Context, id int, req *userrepo.UpdateUserRequest) (*ent.User, error) {
			return nil, errors.New("update failed")
		},
	}
	svc := newTestService(t, repo)

	newNickname := "NewName"
	result, errNo := svc.UpdateUser(context.Background(), 1, &userrepo.UpdateUserRequest{
		Nickname: &newNickname,
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.UserUpdateFailedError.GetCode(), errNo.GetCode())
}

// =============================================================================
// DeleteUser 测试
// =============================================================================

func TestDeleteUser_Success(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	errNo := svc.DeleteUser(context.Background(), 1)
	assert.Equal(t, errno.OK, errNo)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	errNo := svc.DeleteUser(context.Background(), 0)
	assert.Equal(t, errno.InvalidUserIDError.GetCode(), errNo.GetCode())
}

func TestDeleteUser_NegativeID(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	errNo := svc.DeleteUser(context.Background(), -1)
	assert.Equal(t, errno.InvalidUserIDError.GetCode(), errNo.GetCode())
}

func TestDeleteUser_DeleteFailed(t *testing.T) {
	repo := &mockUserRepo{
		deleteFn: func(ctx context.Context, id int) error {
			return errors.New("delete failed")
		},
	}
	svc := newTestService(t, repo)

	errNo := svc.DeleteUser(context.Background(), 1)
	assert.Equal(t, errno.UserDeleteFailedError.GetCode(), errNo.GetCode())
}

// =============================================================================
// ListUsers 测试
// =============================================================================

func TestListUsers_Success(t *testing.T) {
	repo := &mockUserRepo{
		listFn: func(ctx context.Context, req *userrepo.ListUsersRequest) ([]*ent.User, int, error) {
			return []*ent.User{
				{ID: 1, Username: "user1", Email: "user1@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{ID: 2, Username: "user2", Email: "user2@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, 20, nil
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.ListUsers(context.Background(), &userrepo.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, 2, len(result.Users))
	assert.Equal(t, 20, result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 2, result.TotalPage)
}

func TestListUsers_EmptyList(t *testing.T) {
	svc := newTestService(t, &mockUserRepo{})

	result, errNo := svc.ListUsers(context.Background(), &userrepo.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, 0, len(result.Users))
	assert.Equal(t, 0, result.Total)
}

func TestListUsers_DBError(t *testing.T) {
	repo := &mockUserRepo{
		listFn: func(ctx context.Context, req *userrepo.ListUsersRequest) ([]*ent.User, int, error) {
			return nil, 0, errors.New("db error")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.ListUsers(context.Background(), &userrepo.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	assert.Nil(t, result)
	assert.Equal(t, errno.DatabaseError.GetCode(), errNo.GetCode())
}

func TestListUsers_DefaultPagination(t *testing.T) {
	// Test that page < 1 defaults to 1 and pageSize < 1 defaults to 10
	repo := &mockUserRepo{
		listFn: func(ctx context.Context, req *userrepo.ListUsersRequest) ([]*ent.User, int, error) {
			return []*ent.User{}, 0, nil
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.ListUsers(context.Background(), &userrepo.ListUsersRequest{
		Page:     0,
		PageSize: 0,
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
}

func TestListUsers_TotalPageCalculation(t *testing.T) {
	repo := &mockUserRepo{
		listFn: func(ctx context.Context, req *userrepo.ListUsersRequest) ([]*ent.User, int, error) {
			return []*ent.User{}, 25, nil
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.ListUsers(context.Background(), &userrepo.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NotNil(t, result)
	assert.Equal(t, errno.OK, errNo)
	assert.Equal(t, 3, result.TotalPage) // 25/10 = 2.5 → 3
}

// =============================================================================
// UserResponse / toUserResponse 测试
// =============================================================================

func TestToUserResponse_NilUser(t *testing.T) {
	// We test this indirectly - GetUserByID with not-found user
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id int) (*ent.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := newTestService(t, repo)

	result, errNo := svc.GetUserByID(context.Background(), 999)
	assert.Nil(t, result)
	assert.NotEqual(t, errno.OK, errNo)
}

func TestNewUserService(t *testing.T) {
	logCfg := logger.DefaultConfig()
	log, err := logger.NewLogger(logCfg)
	require.NoError(t, err)

	tp, _ := trace.NewInMemoryProvider()
	tracer := trace.NewTracer(tp)

	svc := userservice.NewUserService(&mockUserRepo{}, log, tracer)
	assert.NotNil(t, svc)
}
