package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError 验证错误信息
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatErrors 将验证错误转换为友好的错误信息
// 返回字段名到错误消息的映射
func FormatErrors(err error) map[string]string {
	errors := make(map[string]string)

	if err == nil {
		return errors
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非验证错误，返回原始错误信息
		errors["_"] = err.Error()
		return errors
	}

	for _, e := range validationErrors {
		field := e.Field()
		msg := formatFieldError(e)
		errors[field] = msg
	}

	return errors
}

// FormatErrorsFlat 将验证错误转换为单行列表字符串
// 返回格式: "field1: message1; field2: message2; ..."
func FormatErrorsFlat(err error) string {
	if err == nil {
		return ""
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var parts []string
	for _, e := range validationErrors {
		field := e.Field()
		msg := formatFieldError(e)
		parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
	}

	return strings.Join(parts, "; ")
}

// formatFieldError 根据验证类型返回友好的错误消息
func formatFieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", e.Field())
	case "at_least_one_field":
		return "at least one field must be provided"
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}

// HasError 检查是否有错误（验证错误或普通错误）
func HasError(err error) bool {
	return err != nil
}

// ErrorCount 返回验证错误数量
func ErrorCount(err error) int {
	if err == nil {
		return 0
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return 1
	}
	return len(validationErrors)
}
