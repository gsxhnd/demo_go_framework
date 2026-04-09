package health

import (
	"go_sample_code/internal/database"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Check(c *fiber.Ctx) error
}

type handler struct {
	log           logger.Logger
	healthChecker database.HealthChecker
}

func NewHandler(log logger.Logger, healthChecker database.HealthChecker) Handler {
	return &handler{
		log:           log,
		healthChecker: healthChecker,
	}
}
