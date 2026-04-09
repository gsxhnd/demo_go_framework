package rbac

import "net/http"

// 权限相关错误定义
var (
	// Permission denied errors
	ErrPermissionDenied = &PermissionError{
		HTTPStatus: http.StatusForbidden,
		Code:       4001,
		Message:    "Permission denied",
	}

	ErrRoleNotFound = &PermissionError{
		HTTPStatus: http.StatusNotFound,
		Code:       4002,
		Message:    "Role not found",
	}

	ErrRoleAlreadyExists = &PermissionError{
		HTTPStatus: http.StatusConflict,
		Code:       4003,
		Message:    "Role already exists",
	}

	ErrRoleAssignFailed = &PermissionError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       4004,
		Message:    "Failed to assign role",
	}

	ErrPolicyInvalid = &PermissionError{
		HTTPStatus: http.StatusBadRequest,
		Code:       4005,
		Message:    "Invalid policy",
	}

	ErrSubjectInvalid = &PermissionError{
		HTTPStatus: http.StatusBadRequest,
		Code:       4006,
		Message:    "Invalid subject",
	}

	ErrObjectInvalid = &PermissionError{
		HTTPStatus: http.StatusBadRequest,
		Code:       4007,
		Message:    "Invalid object",
	}

	ErrActionInvalid = &PermissionError{
		HTTPStatus: http.StatusBadRequest,
		Code:       4008,
		Message:    "Invalid action",
	}

	ErrRBACNotEnabled = &PermissionError{
		HTTPStatus: http.StatusServiceUnavailable,
		Code:       4009,
		Message:    "RBAC is not enabled",
	}

	ErrABACNotEnabled = &PermissionError{
		HTTPStatus: http.StatusServiceUnavailable,
		Code:       4010,
		Message:    "ABAC is not enabled",
	}
)

// PermissionError 权限错误
type PermissionError struct {
	HTTPStatus int
	Code       int
	Message    string
}

// Error 实现 error 接口
func (e *PermissionError) Error() string {
	return e.Message
}

// GetHTTPStatus 获取 HTTP 状态码
func (e *PermissionError) GetHTTPStatus() int {
	return e.HTTPStatus
}

// GetCode 获取错误码
func (e *PermissionError) GetCode() int {
	return e.Code
}

// GetMessage 获取错误消息
func (e *PermissionError) GetMessage() string {
	return e.Message
}
