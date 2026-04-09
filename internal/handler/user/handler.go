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
	// Create 创建用户
	Create(c *fiber.Ctx) error
	// GetByID 根据 ID 获取用户
	GetByID(c *fiber.Ctx) error
	// GetByUsername 根据用户名获取用户
	GetByUsername(c *fiber.Ctx) error
	// GetByEmail 根据邮箱获取用户
	GetByEmail(c *fiber.Ctx) error
	// Update 更新用户
	Update(c *fiber.Ctx) error
	// Delete 删除用户
	Delete(c *fiber.Ctx) error
	// List 分页获取用户列表
	List(c *fiber.Ctx) error
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
