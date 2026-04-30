package health

import (
	"go_sample_code/internal/database"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	otel_trace "go.opentelemetry.io/otel/trace"
)

type Handler interface {
	Check(c *fiber.Ctx) error
}

type handler struct {
	log           logger.Logger
	healthChecker database.HealthChecker
	tracer        otel_trace.Tracer
}

func NewHandler(log logger.Logger, healthChecker database.HealthChecker, tracer otel_trace.Tracer) Handler {
	return &handler{
		log:           log,
		healthChecker: healthChecker,
		tracer:        tracer,
	}
}
