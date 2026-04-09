package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthStatus_Constants(t *testing.T) {
	assert.Equal(t, "ok", StatusOK)
	assert.Equal(t, "degraded", StatusDegraded)
	assert.Equal(t, "up", StatusUp)
	assert.Equal(t, "down", StatusDown)
}

func TestNewMockHealthChecker(t *testing.T) {
	ctx := context.Background()

	t.Run("healthy status", func(t *testing.T) {
		healthyStatus := &HealthStatus{
			Data: HealthData{
				Status: StatusOK,
				Relational: RelationalStatus{
					Driver: DriverPostgres,
					Status: StatusUp,
				},
				Redis: RedisStatus{
					Status: StatusUp,
				},
			},
		}

		checker := NewMockHealthChecker(healthyStatus)
		result := checker.Check(ctx)

		assert.Equal(t, StatusOK, result.Data.Status)
		assert.Equal(t, DriverPostgres, result.Data.Relational.Driver)
		assert.Equal(t, StatusUp, result.Data.Relational.Status)
		assert.Equal(t, StatusUp, result.Data.Redis.Status)
	})

	t.Run("degraded status - relational down", func(t *testing.T) {
		degradedStatus := &HealthStatus{
			Data: HealthData{
				Status: StatusDegraded,
				Relational: RelationalStatus{
					Driver: DriverMySQL,
					Status: StatusDown,
				},
				Redis: RedisStatus{
					Status: StatusUp,
				},
			},
		}

		checker := NewMockHealthChecker(degradedStatus)
		result := checker.Check(ctx)

		assert.Equal(t, StatusDegraded, result.Data.Status)
		assert.Equal(t, DriverMySQL, result.Data.Relational.Driver)
		assert.Equal(t, StatusDown, result.Data.Relational.Status)
		assert.Equal(t, StatusUp, result.Data.Redis.Status)
	})

	t.Run("degraded status - redis down", func(t *testing.T) {
		degradedStatus := &HealthStatus{
			Data: HealthData{
				Status: StatusDegraded,
				Relational: RelationalStatus{
					Driver: DriverPostgres,
					Status: StatusUp,
				},
				Redis: RedisStatus{
					Status: StatusDown,
				},
			},
		}

		checker := NewMockHealthChecker(degradedStatus)
		result := checker.Check(ctx)

		assert.Equal(t, StatusDegraded, result.Data.Status)
		assert.Equal(t, StatusUp, result.Data.Relational.Status)
		assert.Equal(t, StatusDown, result.Data.Redis.Status)
	})
}

func TestHealthData_Structure(t *testing.T) {
	data := HealthData{
		Status: StatusOK,
		Relational: RelationalStatus{
			Driver: DriverPostgres,
			Status: StatusUp,
		},
		Redis: RedisStatus{
			Status: StatusUp,
		},
	}

	assert.Equal(t, StatusOK, data.Status)
	assert.Equal(t, DriverPostgres, data.Relational.Driver)
	assert.Equal(t, StatusUp, data.Relational.Status)
	assert.Equal(t, StatusUp, data.Redis.Status)
}

func TestHealthStatus_JSONStructure(t *testing.T) {
	status := &HealthStatus{
		Data: HealthData{
			Status: StatusDegraded,
			Relational: RelationalStatus{
				Driver: DriverMySQL,
				Status: StatusDown,
			},
			Redis: RedisStatus{
				Status: StatusUp,
			},
		},
	}

	// Verify the structure is correct
	assert.NotNil(t, status)
	assert.Equal(t, StatusDegraded, status.Data.Status)
	assert.Equal(t, DriverMySQL, status.Data.Relational.Driver)
	assert.Equal(t, StatusDown, status.Data.Relational.Status)
	assert.Equal(t, StatusUp, status.Data.Redis.Status)
}
