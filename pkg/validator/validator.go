package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Validate go-playground/validator/v10 的类型别名，用于使用我们封装的验证器
type Validate = validator.Validate

// StructLevel go-playground/validator/v10 的类型别名
type StructLevel = validator.StructLevel

// StructValidationFunc 结构级校验函数类型
type StructValidationFunc = validator.StructLevelFunc

// option 配置选项
type option struct {
	structValidations map[reflect.Type]StructValidationFunc
}

// Option 配置验证器的选项
type Option func(*option)

// WithStructValidation 注册结构级校验
func WithStructValidation(v interface{}, fn StructValidationFunc) Option {
	return func(o *option) {
		if o.structValidations == nil {
			o.structValidations = make(map[reflect.Type]StructValidationFunc)
		}
		o.structValidations[reflect.TypeOf(v)] = fn
	}
}

// New 创建验证器实例
func New(opts ...Option) *Validate {
	v := validator.New(validator.WithRequiredStructEnabled(), validator.WithPrivateFieldValidation())

	// 应用选项
	cfg := &option{}
	for _, opt := range opts {
		opt(cfg)
	}

	// 注册结构级校验
	for typ, fn := range cfg.structValidations {
		v.RegisterStructValidation(fn, reflect.Zero(typ).Interface())
	}

	return v
}
