package database

import (
	"context"
	"database/sql"
	"time"

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
	db                 *sql.DB
	redisClient        *redis.Client
	driver             string
	log                logger.Logger
	healthCheckTimeout time.Duration
}

// NewHealthChecker creates a new health checker that aggregates health status
// of all database dependencies (relational DB and Redis)
func NewHealthChecker(
	db *sql.DB,
	redisClient *redis.Client,
	driver string,
	log logger.Logger,
	healthCheckTimeout time.Duration,
) HealthChecker {
	if healthCheckTimeout <= 0 {
		healthCheckTimeout = DefaultHealthCheckTimeout
	}
	return &healthChecker{
		db:                 db,
		redisClient:        redisClient,
		driver:             driver,
		log:                log,
		healthCheckTimeout: healthCheckTimeout,
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

	// Create a context with timeout for health checks
	checkCtx, cancel := context.WithTimeout(ctx, h.healthCheckTimeout)
	defer cancel()

	// Check relational database health using sql.DB ping
	if err := PingRelational(checkCtx, h.db); err != nil {
		h.log.Warn("relational database health check failed",
			zap.String("driver", h.driver),
			zap.Error(err),
		)
		status.Data.Status = StatusDegraded
		status.Data.Relational.Status = StatusDown
	}

	// Check Redis health
	if err := PingRedis(checkCtx, h.redisClient); err != nil {
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
