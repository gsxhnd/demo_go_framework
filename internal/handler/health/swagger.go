package health

import "go_sample_code/internal/database"

// SwaggerHealthResponse wraps the health check response for swagger documentation
type SwaggerHealthResponse struct {
	Code    int                 `json:"code" example:"0"`
	Message string              `json:"message" example:"OK"`
	Data    database.HealthData `json:"data"`
}
