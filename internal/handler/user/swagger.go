package user

// SwaggerErrorResponse represents an error response
type SwaggerErrorResponse struct {
	Code    int    `json:"code" example:"1003"`
	Message string `json:"message" example:"Request Validate Error"`
	Data    any    `json:"data"`
}

// SwaggerUserResponse wraps a single user response
type SwaggerUserResponse struct {
	Code    int          `json:"code" example:"0"`
	Message string       `json:"message" example:"OK"`
	Data    UserResponse `json:"data"`
}

// SwaggerUserListResponse wraps a user list response
type SwaggerUserListResponse struct {
	Code    int               `json:"code" example:"0"`
	Message string            `json:"message" example:"OK"`
	Data    ListUsersResponse `json:"data"`
}

// SwaggerDeleteResponse wraps a delete response
type SwaggerDeleteResponse struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message" example:"OK"`
	Data    DeleteMessage `json:"data"`
}

// DeleteMessage represents the delete success message
type DeleteMessage struct {
	Message string `json:"message" example:"user deleted successfully"`
}
