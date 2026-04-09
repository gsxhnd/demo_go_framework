package rbac

import "fmt"

// RBAC 权限配置
type RBACConfig struct {
	// 模型配置文件路径
	ModelPath string `yaml:"model_path"`
	// 策略文件路径
	PolicyPath string `yaml:"policy_path"`
	// 是否启用 RBAC
	Enabled bool `yaml:"enabled"`
	// 默认角色
	DefaultRole string `yaml:"default_role"`
}

// ABACConfig ABAC 配置
type ABACConfig struct {
	// 是否启用 ABAC
	Enabled bool `yaml:"enabled"`
	// ABAC 模型配置文件路径
	ModelPath string `yaml:"model_path"`
}

// Config 完整的权限配置
type Config struct {
	RBAC RBACConfig `yaml:"rbac"`
	ABAC ABACConfig `yaml:"abac"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		RBAC: RBACConfig{
			Enabled:     true,
			ModelPath:   "pkg/rbac/model.conf",
			PolicyPath:  "pkg/rbac/policy.csv",
			DefaultRole: "guest",
		},
		ABAC: ABACConfig{
			Enabled:   true,
			ModelPath: "pkg/rbac/abac_model.conf",
		},
	}
}

// Validate 验证配置有效性
func (c *Config) Validate() error {
	if c.RBAC.Enabled {
		if c.RBAC.ModelPath == "" {
			return fmt.Errorf("rbac model_path is required when RBAC is enabled")
		}
		if c.RBAC.PolicyPath == "" {
			return fmt.Errorf("rbac policy_path is required when RBAC is enabled")
		}
	}
	if c.ABAC.Enabled && c.ABAC.ModelPath == "" {
		return fmt.Errorf("abac model_path is required when ABAC is enabled")
	}
	return nil
}
