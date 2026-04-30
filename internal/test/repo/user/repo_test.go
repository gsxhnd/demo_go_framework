package user_test

import (
	"testing"

	userrepo "go_sample_code/internal/repo/user"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserRepo(t *testing.T) {
	logCfg := logger.DefaultConfig()
	log, err := logger.NewLogger(logCfg)
	require.NoError(t, err)

	tp, _ := trace.NewInMemoryProvider()
	tracer := trace.NewTracer(tp)

	repo := userrepo.NewUserRepo(nil, log, tracer)
	assert.NotNil(t, repo)
}

func TestCreateUserRequest(t *testing.T) {
	active := true
	req := userrepo.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Avatar:   "https://example.com/avatar.png",
		Phone:    "13800138000",
		IsActive: &active,
	}

	assert.Equal(t, "testuser", req.Username)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "Test User", req.Nickname)
	assert.Equal(t, "https://example.com/avatar.png", req.Avatar)
	assert.Equal(t, "13800138000", req.Phone)
	assert.True(t, *req.IsActive)
}

func TestCreateUserRequest_DefaultActive(t *testing.T) {
	// IsActive is nil by default - the repo should default to true
	req := userrepo.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Nil(t, req.IsActive)
}

func TestUpdateUserRequest(t *testing.T) {
	newEmail := "new@example.com"
	newNickname := "New Name"
	active := false

	req := userrepo.UpdateUserRequest{
		Email:    &newEmail,
		Nickname: &newNickname,
		IsActive: &active,
	}

	assert.Equal(t, "new@example.com", *req.Email)
	assert.Equal(t, "New Name", *req.Nickname)
	assert.False(t, *req.IsActive)
	assert.Nil(t, req.Password)
	assert.Nil(t, req.Avatar)
	assert.Nil(t, req.Phone)
}

func TestListUsersRequest(t *testing.T) {
	req := userrepo.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Keyword:  "test",
	}

	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.PageSize)
	assert.Equal(t, "test", req.Keyword)
}

func TestListUsersRequest_DefaultValues(t *testing.T) {
	req := userrepo.ListUsersRequest{}

	assert.Equal(t, 0, req.Page)
	assert.Equal(t, 0, req.PageSize)
	assert.Equal(t, "", req.Keyword)
}
