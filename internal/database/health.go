package database

import (
	"context"
	"time"

	"go_sample_code/internal/ent"
	"go_sample_code/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// HealthStatus constants
const (
	StatusOK       = "ok"
	StatusDegraded = "degraded"
	StatusUp       = "up"
	StatusDown     = "down"
)

// RelationalStatus represents the status of the relational database
type RelationalStatus struct {
	Driver string `json:"driver"`
	Status string `json:"status"`
}

// RedisStatus represents the status of Redis
type RedisStatus struct {
	Status string `json:"status"`
}

// HealthData contains the health check data for all database dependencies
type HealthData struct {
	Status     string           `json:"status"`
	Relational RelationalStatus `json:"relational"`
	Redis      RedisStatus      `json:"redis"`
}

// HealthStatus represents the complete health check result
type HealthStatus struct {
	Data HealthData `json:"data"`
}

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	Check(ctx context.Context) *HealthStatus
}

// healthChecker implements HealthChecker interface
type healthChecker struct {
	entClient   *ent.Client
	redisClient *redis.Client
	driver      string
	log         logger.Logger
}

// NewHealthChecker creates a new health checker that aggregates health status
// of all database dependencies (relational DB and Redis)
func NewHealthChecker(
	entClient *ent.Client,
	redisClient *redis.Client,
	driver string,
	log logger.Logger,
) HealthChecker {
	return &healthChecker{
		entClient:   entClient,
		redisClient: redisClient,
		driver:      driver,
		log:         log,
	}
}

// Check performs health checks on all database dependencies and returns the aggregated status
func (h *healthChecker) Check(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Data: HealthData{
			Status: StatusOK,
			Relational: RelationalStatus{
				Driver: h.driver,
				Status: StatusUp,
			},
			Redis: RedisStatus{
				Status: StatusUp,
			},
		},
	}

	// Check relational database health
	if err := PingRelational(ctx, h.entClient); err != nil {
		h.log.Warn("relational database health check failed",
			zap.String("driver", h.driver),
			zap.Error(err),
		)
		status.Data.Status = StatusDegraded
		status.Data.Relational.Status = StatusDown
	}

	// Check Redis health
	if err := PingRedis(ctx, h.redisClient); err != nil {
		h.log.Warn("Redis health check failed",
			zap.Error(err),
		)
		status.Data.Status = StatusDegraded
		status.Data.Redis.Status = StatusDown
	}

	return status
}

// NewMockHealthChecker creates a mock health checker for testing purposes
func NewMockHealthChecker(healthStatus *HealthStatus) HealthChecker {
	return &mockHealthChecker{status: healthStatus}
}

// mockHealthChecker is a mock implementation for testing
type mockHealthChecker struct {
	status *HealthStatus
}

// Check returns the predefined health status
func (m *mockHealthChecker) Check(ctx context.Context) *HealthStatus {
	return m.status
}

// HealthCheckTimestamp returns the current timestamp for health check responses
func HealthCheckTimestamp() time.Time {
	return time.Now().UTC()
}
