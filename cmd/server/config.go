package main

import (
	"fmt"
	"os"

	"go_sample_code/internal/database"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/metrics"
	"go_sample_code/pkg/trace"

	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

type ConfigPath string

// CommonConfig 通用配置
type CommonConfig struct {
	Listen string `yaml:"listen"`
}

// AppConfig 应用级配置结构
type AppConfig struct {
	Common   CommonConfig            `yaml:"common"`
	Database database.DatabaseConfig `yaml:"database"`
	Logger   logger.LoggerConfig     `yaml:"logger"`
	Trace    trace.TraceConfig       `yaml:"trace"`
	Metrics  metrics.MetricsConfig   `yaml:"metrics"`
}

// ApplyDefaults 应用默认值
func (c *AppConfig) ApplyDefaults() {
	// Common defaults
	if c.Common.Listen == "" {
		c.Common.Listen = ":8080"
	}

	// Database defaults
	c.Database.ApplyDefaults()

	// Logger defaults
	if c.Logger.Output == "" {
		c.Logger.Output = "console"
	}
	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}
	if c.Logger.OtelServiceName == "" {
		c.Logger.OtelServiceName = "demo-go-framework"
	}

	// Trace defaults
	if c.Trace.OtelServiceName == "" {
		c.Trace.OtelServiceName = "demo-go-framework"
	}
	if c.Trace.OtelServiceVersion == "" {
		c.Trace.OtelServiceVersion = "1.0.0"
	}
	if c.Trace.OtelEnvironment == "" {
		c.Trace.OtelEnvironment = "development"
	}
	if c.Trace.OtelSamplerType == "" {
		c.Trace.OtelSamplerType = "always_on"
	}

	// Metrics defaults
	c.Metrics.ApplyDefaults()
}

// Validate 验证配置
func (c *AppConfig) Validate() error {
	// Validate database
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config error: %w", err)
	}

	// Validate trace
	if err := c.Trace.Validate(); err != nil {
		return fmt.Errorf("trace config error: %w", err)
	}

	// Validate metrics
	if err := c.Metrics.Validate(); err != nil {
		return fmt.Errorf("metrics config error: %w", err)
	}

	return nil
}

// NewAppConfig 读取并解析应用配置
func NewAppConfig(cfgPath ConfigPath) (*AppConfig, error) {
	var cfg AppConfig

	data, err := os.ReadFile(string(cfgPath))
	if err != nil {
		// If config file doesn't exist, use defaults
		cfg = AppConfig{
			Common: CommonConfig{
				Listen: ":8080",
			},
			Database: database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverPostgres,
					Postgres: database.PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "postgres",
						DBName:   "demo",
						SSLMode:  "disable",
						Pool:     database.DefaultPoolConfig(),
					},
				},
				Redis: database.DefaultRedisConfig(),
			},
			Logger: logger.LoggerConfig{
				Output:          "console",
				Level:           "info",
				OtelServiceName: "demo-go-framework",
			},
			Trace: trace.TraceConfig{
				OtelEnable:         false,
				OtelServiceName:    "demo-go-framework",
				OtelServiceVersion: "1.0.0",
				OtelEnvironment:    "development",
			},
			Metrics: metrics.DefaultMetricsConfig(),
		}
		cfg.Database.Redis.Addr = "localhost:6379"
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Apply default values
	cfg.ApplyDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Config 兼容旧的 fx.Out 结构
type Config struct {
	fx.Out
	CommonConfig   *CommonConfig            `yaml:"common"`
	DatabaseConfig *database.DatabaseConfig `yaml:"database"`
}

// NewDatabaseConfig 从 AppConfig 提取数据库配置
func NewDatabaseConfig(cfg *AppConfig) *database.DatabaseConfig {
	return &cfg.Database
}

// NewLoggerConfig 从 AppConfig 提取日志配置
func NewLoggerConfig(cfg *AppConfig) *logger.LoggerConfig {
	return &cfg.Logger
}

// NewTraceConfig 从 AppConfig 提取追踪配置
func NewTraceConfig(cfg *AppConfig) *trace.TraceConfig {
	return &cfg.Trace
}

// NewMetricsConfig 从 AppConfig 提取 metrics 配置
func NewMetricsConfig(cfg *AppConfig) *metrics.MetricsConfig {
	return &cfg.Metrics
}
