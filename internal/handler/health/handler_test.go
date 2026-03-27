package health_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	healthhandler "go_sample_code/internal/handler/health"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	log, err := logger.NewLogger(logger.DefaultConfig())
	require.NoError(t, err)

	h := healthhandler.NewHandler(log)

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

	assert.Equal(t, "ok", body["data"])
}
