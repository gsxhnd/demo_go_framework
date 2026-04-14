package validator

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStruct 用于测试的简单结构体
type TestStruct struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=0,lte=150"`
}

// TestStructWithOptional 用于测试可选字段的结构体
type TestStructWithOptional struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=32"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

// TestStructAtLeastOneField 用于测试至少一个字段的结构体
type TestStructAtLeastOneField struct {
	Field1 *string `json:"field1,omitempty"`
	Field2 *string `json:"field2,omitempty"`
	Field3 *int    `json:"field3,omitempty"`
}

func TestNew(t *testing.T) {
	v := New()
	require.NotNil(t, v)
	assert.IsType(t, &validator.Validate{}, v)
}

func TestDefault(t *testing.T) {
	v := New()
	require.NotNil(t, v)
	assert.IsType(t, &validator.Validate{}, v)
}

func TestNew_WithStructValidation(t *testing.T) {
	v := New(
		WithStructValidation(TestStructAtLeastOneField{}, func(sl validator.StructLevel) {
			req := sl.Current().Interface().(TestStructAtLeastOneField)
			if req.Field1 == nil && req.Field2 == nil && req.Field3 == nil {
				sl.ReportError(reflect.ValueOf(req), "TestStructAtLeastOneField", "", "at_least_one_field", "")
			}
		}),
	)

	require.NotNil(t, v)

	// 测试空结构体应该失败
	err := v.Struct(TestStructAtLeastOneField{})
	assert.Error(t, err)
	assert.True(t, HasError(err))

	// 测试至少有一个字段应该成功
	strVal := "test"
	err = v.Struct(TestStructAtLeastOneField{Field1: &strVal})
	assert.NoError(t, err)
}

func TestValidator_BasicValidation(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		data    TestStruct
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid struct",
			data:    TestStruct{Name: "John", Email: "john@example.com", Age: 25},
			wantErr: false,
		},
		{
			name:    "missing name",
			data:    TestStruct{Name: "", Email: "john@example.com", Age: 25},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name:    "invalid email",
			data:    TestStruct{Name: "John", Email: "invalid-email", Age: 25},
			wantErr: true,
			errMsg:  "email must be a valid email address",
		},
		{
			name:    "age too high",
			data:    TestStruct{Name: "John", Email: "john@example.com", Age: 200},
			wantErr: true,
		},
		{
			name:    "negative age",
			data:    TestStruct{Name: "John", Email: "john@example.com", Age: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					errMsg := FormatErrorsFlat(err)
					assert.Contains(t, errMsg, tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_OptionalFields(t *testing.T) {
	v := New()

	// 空结构体应该成功
	err := v.Struct(TestStructWithOptional{})
	assert.NoError(t, err)

	// 有效密码应该成功
	pwd := "password123"
	err = v.Struct(TestStructWithOptional{Password: &pwd})
	assert.NoError(t, err)

	// 密码太短应该失败
	shortPwd := "123"
	err = v.Struct(TestStructWithOptional{Password: &shortPwd})
	assert.Error(t, err)
}

func TestFormatErrors(t *testing.T) {
	v := New()

	// 测试空错误
	errors := FormatErrors(nil)
	assert.Empty(t, errors)

	// 测试非验证错误
	errors = FormatErrors(assert.AnError)
	assert.Contains(t, errors, "_")
	assert.Equal(t, assert.AnError.Error(), errors["_"])

	// 测试验证错误
	invalidStruct := TestStruct{Name: "", Email: "invalid"}
	err := v.Struct(invalidStruct)
	require.Error(t, err)

	errors = FormatErrors(err)
	assert.NotEmpty(t, errors)
	assert.Contains(t, errors, "name")
	assert.Contains(t, errors, "email")
}

func TestFormatErrorsFlat(t *testing.T) {
	v := New()

	// 测试空错误
	result := FormatErrorsFlat(nil)
	assert.Empty(t, result)

	// 测试非验证错误
	result = FormatErrorsFlat(assert.AnError)
	assert.Equal(t, assert.AnError.Error(), result)

	// 测试验证错误
	invalidStruct := TestStruct{Name: "", Email: "invalid"}
	err := v.Struct(invalidStruct)
	require.Error(t, err)

	result = FormatErrorsFlat(err)
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "email")
}

func TestHasError(t *testing.T) {
	v := New()

	// 无错误
	assert.False(t, HasError(nil))

	// 非验证错误
	assert.True(t, HasError(assert.AnError))

	// 验证错误
	invalidStruct := TestStruct{Name: "", Email: "invalid"}
	err := v.Struct(invalidStruct)
	assert.True(t, HasError(err))
}

func TestErrorCount(t *testing.T) {
	v := New()

	// 无错误
	assert.Equal(t, 0, ErrorCount(nil))

	// 非验证错误
	assert.Equal(t, 1, ErrorCount(assert.AnError))

	// 单个验证错误
	singleErr := TestStruct{Name: "", Email: "valid@example.com"}
	err := v.Struct(singleErr)
	assert.Equal(t, 1, ErrorCount(err))

	// 多个验证错误
	multiErr := TestStruct{Name: "", Email: "invalid"}
	err = v.Struct(multiErr)
	assert.GreaterOrEqual(t, ErrorCount(err), 2)
}
