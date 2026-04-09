package health_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"go_sample_code/internal/database"
	healthhandler "go_sample_code/internal/handler/health"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_Healthy(t *testing.T) {
	log, err := logger.NewLogger(logger.DefaultConfig())
	require.NoError(t, err)

	// Create a mock health checker that returns healthy status
	healthyStatus := &database.HealthStatus{
		Data: database.HealthData{
			Status: database.StatusOK,
			Relational: database.RelationalStatus{
				Driver: database.DriverPostgres,
				Status: database.StatusUp,
			},
			Redis: database.RedisStatus{
				Status: database.StatusUp,
			},
		},
	}
	healthChecker := database.NewMockHealthChecker(healthyStatus)

	h := healthhandler.NewHandler(log, healthChecker)

	app := fiber.New()
	app.Get("/api/health", h.Check)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body map[string]any
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	assert.Equal(t, float64(0), body["code"])
	assert.Equal(t, "OK", body["message"])

	data, ok := body["data"].(map[string]any)
	require.True(t, ok, "data should be a map")
	assert.Equal(t, "ok", data["status"])

	relational, ok := data["relational"].(map[string]any)
	require.True(t, ok, "relational should be a map")
	assert.Equal(t, "postgres", relational["driver"])
	assert.Equal(t, "up", relational["status"])

	redis, ok := data["redis"].(map[string]any)
	require.True(t, ok, "redis should be a map")
	assert.Equal(t, "up", redis["status"])
}

func TestCheck_Degraded(t *testing.T) {
	log, err := logger.NewLogger(logger.DefaultConfig())
	require.NoError(t, err)

	// Create a mock health checker that returns degraded status
	degradedStatus := &database.HealthStatus{
		Data: database.HealthData{
			Status: database.StatusDegraded,
			Relational: database.RelationalStatus{
				Driver: database.DriverMySQL,
				Status: database.StatusDown,
			},
			Redis: database.RedisStatus{
				Status: database.StatusUp,
			},
		},
	}
	healthChecker := database.NewMockHealthChecker(degradedStatus)

	h := healthhandler.NewHandler(log, healthChecker)

	app := fiber.New()
	app.Get("/api/health", h.Check)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)

	var body map[string]any
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	assert.Equal(t, float64(0), body["code"])
	assert.Equal(t, "OK", body["message"])

	data, ok := body["data"].(map[string]any)
	require.True(t, ok, "data should be a map")
	assert.Equal(t, "degraded", data["status"])

	relational, ok := data["relational"].(map[string]any)
	require.True(t, ok, "relational should be a map")
	assert.Equal(t, "mysql", relational["driver"])
	assert.Equal(t, "down", relational["status"])
}

func TestCheck_RedisDegraded(t *testing.T) {
	log, err := logger.NewLogger(logger.DefaultConfig())
	require.NoError(t, err)

	// Create a mock health checker that returns degraded status (redis down)
	degradedStatus := &database.HealthStatus{
		Data: database.HealthData{
			Status: database.StatusDegraded,
			Relational: database.RelationalStatus{
				Driver: database.DriverPostgres,
				Status: database.StatusUp,
			},
			Redis: database.RedisStatus{
				Status: database.StatusDown,
			},
		},
	}
	healthChecker := database.NewMockHealthChecker(degradedStatus)

	h := healthhandler.NewHandler(log, healthChecker)

	app := fiber.New()
	app.Get("/api/health", h.Check)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)

	var body map[string]any
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	data, ok := body["data"].(map[string]any)
	require.True(t, ok)

	assert.Equal(t, "degraded", data["status"])

	redis, ok := data["redis"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "down", redis["status"])
}

// mockHealthChecker implements database.HealthChecker for testing
type mockHealthChecker struct {
	status *database.HealthStatus
}

func (m *mockHealthChecker) Check(ctx context.Context) *database.HealthStatus {
	return m.status
}
