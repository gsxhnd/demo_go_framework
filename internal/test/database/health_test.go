package database_test

import (
	"context"
	"testing"

	"go_sample_code/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestHealthStatus_Constants(t *testing.T) {
	assert.Equal(t, "ok", database.StatusOK)
	assert.Equal(t, "degraded", database.StatusDegraded)
	assert.Equal(t, "up", database.StatusUp)
	assert.Equal(t, "down", database.StatusDown)
}

func TestNewMockHealthChecker(t *testing.T) {
	ctx := context.Background()

	t.Run("healthy status", func(t *testing.T) {
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

		checker := database.NewMockHealthChecker(healthyStatus)
		result := checker.Check(ctx)

		assert.Equal(t, database.StatusOK, result.Data.Status)
		assert.Equal(t, database.DriverPostgres, result.Data.Relational.Driver)
		assert.Equal(t, database.StatusUp, result.Data.Relational.Status)
		assert.Equal(t, database.StatusUp, result.Data.Redis.Status)
	})

	t.Run("degraded status - relational down", func(t *testing.T) {
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

		checker := database.NewMockHealthChecker(degradedStatus)
		result := checker.Check(ctx)

		assert.Equal(t, database.StatusDegraded, result.Data.Status)
		assert.Equal(t, database.DriverMySQL, result.Data.Relational.Driver)
		assert.Equal(t, database.StatusDown, result.Data.Relational.Status)
		assert.Equal(t, database.StatusUp, result.Data.Redis.Status)
	})

	t.Run("degraded status - redis down", func(t *testing.T) {
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

		checker := database.NewMockHealthChecker(degradedStatus)
		result := checker.Check(ctx)

		assert.Equal(t, database.StatusDegraded, result.Data.Status)
		assert.Equal(t, database.StatusUp, result.Data.Relational.Status)
		assert.Equal(t, database.StatusDown, result.Data.Redis.Status)
	})
}

func TestHealthData_Structure(t *testing.T) {
	data := database.HealthData{
		Status: database.StatusOK,
		Relational: database.RelationalStatus{
			Driver: database.DriverPostgres,
			Status: database.StatusUp,
		},
		Redis: database.RedisStatus{
			Status: database.StatusUp,
		},
	}

	assert.Equal(t, database.StatusOK, data.Status)
	assert.Equal(t, database.DriverPostgres, data.Relational.Driver)
	assert.Equal(t, database.StatusUp, data.Relational.Status)
	assert.Equal(t, database.StatusUp, data.Redis.Status)
}

func TestHealthStatus_JSONStructure(t *testing.T) {
	status := &database.HealthStatus{
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

	assert.NotNil(t, status)
	assert.Equal(t, database.StatusDegraded, status.Data.Status)
	assert.Equal(t, database.DriverMySQL, status.Data.Relational.Driver)
	assert.Equal(t, database.StatusDown, status.Data.Relational.Status)
	assert.Equal(t, database.StatusUp, status.Data.Redis.Status)
}
