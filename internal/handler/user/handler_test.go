package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"go_sample_code/internal/errno"
	userhandler "go_sample_code/internal/handler/user"
	userrepo "go_sample_code/internal/repo/user"
	userservice "go_sample_code/internal/service/user"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/trace"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUserService 是一个用于测试的 mock 实现
type mockUserService struct{}

func (m *mockUserService) CreateUser(ctx context.Context, req *userrepo.CreateUserRequest) (*userservice.UserResponse, errno.Errno) {
	return &userservice.UserResponse{ID: 1, Username: req.Username, Email: req.Email}, errno.OK
}

func (m *mockUserService) GetUserByID(ctx context.Context, id int) (*userservice.UserResponse, errno.Errno) {
	if id <= 0 {
		return nil, errno.UserNotFoundError
	}
	return &userservice.UserResponse{ID: id, Username: "testuser"}, errno.OK
}

func (m *mockUserService) GetUserByUsername(ctx context.Context, username string) (*userservice.UserResponse, errno.Errno) {
	if username == "" {
		return nil, errno.InvalidUsernameError
	}
	return &userservice.UserResponse{ID: 1, Username: username}, errno.OK
}

func (m *mockUserService) GetUserByEmail(ctx context.Context, email string) (*userservice.UserResponse, errno.Errno) {
	if email == "" {
		return nil, errno.InvalidEmailError
	}
	return &userservice.UserResponse{ID: 1, Email: email}, errno.OK
}

func (m *mockUserService) UpdateUser(ctx context.Context, id int, req *userrepo.UpdateUserRequest) (*userservice.UserResponse, errno.Errno) {
	if id <= 0 {
		return nil, errno.InvalidUserIDError
	}
	return &userservice.UserResponse{ID: id, Username: "updated"}, errno.OK
}

func (m *mockUserService) DeleteUser(ctx context.Context, id int) errno.Errno {
	if id <= 0 {
		return errno.InvalidUserIDError
	}
	return errno.OK
}

func (m *mockUserService) ListUsers(ctx context.Context, req *userrepo.ListUsersRequest) (*userservice.ListUsersResponse, errno.Errno) {
	return &userservice.ListUsersResponse{
		Users:     []*userservice.UserResponse{},
		Total:     0,
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: 0,
	}, errno.OK
}

// newTestValidator 创建用于测试的 validator 实例
func newTestValidator() *validator.Validate {
	v := validator.New()
	// 设置 TagNameFunc，优先使用 json/query/params 标签名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Tag.Get("query")
		}
		if name == "" {
			name = fld.Tag.Get("params")
		}
		if name == "" {
			name = fld.Name
		}
		// 移除 omitempty 等额外标签
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		return name
	})

	// 注册 UpdateUserRequest 结构级校验：至少传一个可更新字段
	v.RegisterStructValidation(func(sl validator.StructLevel) {
		req := sl.Current().Interface().(userhandler.UpdateUserRequest)
		if req.Email == nil && req.Password == nil && req.Nickname == nil &&
			req.Avatar == nil && req.Phone == nil && req.IsActive == nil {
			sl.ReportError(reflect.ValueOf(req), "UpdateUserRequest", "", "at_least_one_field", "")
		}
	}, userhandler.UpdateUserRequest{})

	return v
}

// 创建测试 handler 的辅助函数
func newTestHandler(t *testing.T) userhandler.Handler {
	logCfg := logger.DefaultConfig()
	log, err := logger.NewLogger(logCfg)
	require.NoError(t, err)

	// 使用内存 tracer provider 用于测试
	tp, _ := trace.NewInMemoryProvider()
	tracer := trace.NewTracer(tp)

	v := newTestValidator()

	return userhandler.NewHandler(&mockUserService{}, log, tracer, v)
}

func TestCreate_MissingUsername(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Post("/api/users", h.Create)

	// 缺少 username
	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, float64(errno.RequestValidateError.GetCode()), result["code"])
}

func TestCreate_InvalidEmail(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Post("/api/users", h.Create)

	// 无效的 email 格式
	body := `{"username":"testuser","email":"invalid-email","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreate_PasswordTooShort(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Post("/api/users", h.Create)

	// password 太短
	body := `{"username":"testuser","email":"test@example.com","password":"123"}`
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdate_InvalidID(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Put("/api/users/:id", h.Update)

	// 无效的 id（非数字）
	body := `{"nickname":"newname"}`
	req := httptest.NewRequest("PUT", "/api/users/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdate_IDLessThanOrEqualZero(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Put("/api/users/:id", h.Update)

	// id <= 0
	body := `{"nickname":"newname"}`
	req := httptest.NewRequest("PUT", "/api/users/0", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdate_EmptyBody(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Put("/api/users/:id", h.Update)

	// 空 body（没有传任何可更新字段）
	body := `{}`
	req := httptest.NewRequest("PUT", "/api/users/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdate_InvalidEmail(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Put("/api/users/:id", h.Update)

	// 无效的 email 格式
	body := `{"email":"invalid-email"}`
	req := httptest.NewRequest("PUT", "/api/users/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdate_PasswordTooShort(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Put("/api/users/:id", h.Update)

	// password 太短
	body := `{"password":"123"}`
	req := httptest.NewRequest("PUT", "/api/users/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetByID_InvalidID(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users/:id", h.GetByID)

	// 无效的 id（非数字）
	req := httptest.NewRequest("GET", "/api/users/abc", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetByID_IDLessThanOrEqualZero(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users/:id", h.GetByID)

	// id <= 0
	req := httptest.NewRequest("GET", "/api/users/0", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDelete_InvalidID(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Delete("/api/users/:id", h.Delete)

	// 无效的 id（非数字）
	req := httptest.NewRequest("DELETE", "/api/users/abc", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDelete_IDLessThanOrEqualZero(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Delete("/api/users/:id", h.Delete)

	// id <= 0
	req := httptest.NewRequest("DELETE", "/api/users/0", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetByEmail_InvalidEmail(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users/email/:email", h.GetByEmail)

	// 无效的 email 格式
	req := httptest.NewRequest("GET", "/api/users/email/invalid-email", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestList_PageZero(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users", h.List)

	// page = 0 应该失败
	req := httptest.NewRequest("GET", "/api/users?page=0", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestList_PageSizeExceedsLimit(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users", h.List)

	// page_size = 101 应该失败
	req := httptest.NewRequest("GET", "/api/users?page_size=101", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestList_DefaultValues(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users", h.List)

	// 不传分页参数应该使用默认值
	req := httptest.NewRequest("GET", "/api/users", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// 应该成功返回 200
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestList_ValidParams(t *testing.T) {
	h := newTestHandler(t)

	app := fiber.New()
	app.Get("/api/users", h.List)

	// 有效参数
	req := httptest.NewRequest("GET", "/api/users?page=1&page_size=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
