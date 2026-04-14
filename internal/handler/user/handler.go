package user

import (
	"go_sample_code/internal/service/user"
	"go_sample_code/pkg/logger"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	otel_trace "go.opentelemetry.io/otel/trace"
)

// Handler 用户处理器接口
type Handler interface {
	// UserCreate 创建用户
	UserCreate(c *fiber.Ctx) error
	// UserGetByID 根据 ID 获取用户
	UserGetByID(c *fiber.Ctx) error
	// UserGetByUsername 根据用户名获取用户
	UserGetByUsername(c *fiber.Ctx) error
	// UserGetByEmail 根据邮箱获取用户
	UserGetByEmail(c *fiber.Ctx) error
	// UserUpdate 更新用户
	UserUpdate(c *fiber.Ctx) error
	// UserDelete 删除用户
	UserDelete(c *fiber.Ctx) error
	// UserList 分页获取用户列表
	UserList(c *fiber.Ctx) error
}

// handler 用户处理器实现
type handler struct {
	userService user.UserService
	log         logger.Logger
	tracer      otel_trace.Tracer
	validate    *validator.Validate
}

// NewHandler 创建用户处理器实例
func NewHandler(userService user.UserService, log logger.Logger, tracer otel_trace.Tracer, validate *validator.Validate) Handler {
	return &handler{
		userService: userService,
		log:         log,
		tracer:      tracer,
		validate:    validate,
	}
}
