package database

import (
	"errors"
	"fmt"
	"time"
)

// Driver type constants
const (
	DriverMySQL    = "mysql"
	DriverPostgres = "postgres"
)

// ValidDrivers contains all valid relational database drivers
var ValidDrivers = []string{DriverMySQL, DriverPostgres}

// PoolConfig holds connection pool settings for relational databases
type PoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

// DefaultPoolConfig returns the default connection pool configuration
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// MySQLConfig holds MySQL-specific configuration
type MySQLConfig struct {
	Host     string     `yaml:"host"`
	Port     int        `yaml:"port"`
	User     string     `yaml:"user"`
	Password string     `yaml:"password"`
	DBName   string     `yaml:"dbname"`
	Params   string     `yaml:"params"`
	Pool     PoolConfig `yaml:"pool"`
}

// PostgresConfig holds PostgreSQL-specific configuration
type PostgresConfig struct {
	Host     string     `yaml:"host"`
	Port     int        `yaml:"port"`
	User     string     `yaml:"user"`
	Password string     `yaml:"password"`
	DBName   string     `yaml:"dbname"`
	SSLMode  string     `yaml:"sslmode"`
	Pool     PoolConfig `yaml:"pool"`
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	MaxRetries   int           `yaml:"max_retries"`
}

// DefaultRedisConfig returns the default Redis configuration
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// RelationalConfig holds relational database configuration
type RelationalConfig struct {
	Driver   string         `yaml:"driver"`
	MySQL    MySQLConfig    `yaml:"mysql"`
	Postgres PostgresConfig `yaml:"postgres"`
}

// DatabaseConfig holds all database configuration
type DatabaseConfig struct {
	Relational         RelationalConfig `yaml:"relational"`
	Redis              RedisConfig      `yaml:"redis"`
	HealthCheckTimeout time.Duration    `yaml:"health_check_timeout"`
}

// DefaultHealthCheckTimeout returns the default health check timeout
const DefaultHealthCheckTimeout = 5 * time.Second

// ApplyDefaults applies default values to the configuration without modifying
// the original struct passed in. Returns a new config with defaults applied.
func (c *DatabaseConfig) ApplyDefaults() {
	// Apply MySQL pool defaults
	if c.Relational.Driver == DriverMySQL && c.Relational.MySQL.Pool.MaxOpenConns == 0 {
		c.Relational.MySQL.Pool = DefaultPoolConfig()
	}

	// Apply PostgreSQL pool and SSL defaults
	if c.Relational.Driver == DriverPostgres {
		if c.Relational.Postgres.Pool.MaxOpenConns == 0 {
			c.Relational.Postgres.Pool = DefaultPoolConfig()
		}
		if c.Relational.Postgres.SSLMode == "" {
			c.Relational.Postgres.SSLMode = "disable"
		}
	}

	// Apply Redis defaults
	if c.Redis.PoolSize == 0 {
		defaultCfg := DefaultRedisConfig()
		c.Redis.PoolSize = defaultCfg.PoolSize
		c.Redis.MinIdleConns = defaultCfg.MinIdleConns
		c.Redis.DialTimeout = defaultCfg.DialTimeout
		c.Redis.ReadTimeout = defaultCfg.ReadTimeout
		c.Redis.WriteTimeout = defaultCfg.WriteTimeout
	}

	// Apply health check timeout default
	if c.HealthCheckTimeout == 0 {
		c.HealthCheckTimeout = DefaultHealthCheckTimeout
	}
}

// Validate checks if the database configuration is valid without modifying it
func (c *DatabaseConfig) Validate() error {
	// Validate relational database driver
	if !isValidDriver(c.Relational.Driver) {
		return fmt.Errorf("invalid relational driver: %s, must be one of: %v",
			c.Relational.Driver, ValidDrivers)
	}

	// Validate driver-specific configuration
	switch c.Relational.Driver {
	case DriverMySQL:
		if err := c.validateMySQL(); err != nil {
			return err
		}
	case DriverPostgres:
		if err := c.validatePostgres(); err != nil {
			return err
		}
	}

	// Validate Redis configuration
	if err := c.validateRedis(); err != nil {
		return err
	}

	return nil
}

// GetHealthCheckTimeout returns the configured health check timeout or default
func (c *DatabaseConfig) GetHealthCheckTimeout() time.Duration {
	if c.HealthCheckTimeout > 0 {
		return c.HealthCheckTimeout
	}
	return DefaultHealthCheckTimeout
}

// getPoolConfig returns the connection pool configuration for the selected driver
func (c *DatabaseConfig) getPoolConfig() PoolConfig {
	switch c.Relational.Driver {
	case DriverMySQL:
		return c.Relational.MySQL.Pool
	case DriverPostgres:
		return c.Relational.Postgres.Pool
	default:
		return PoolConfig{}
	}
}

// isValidDriver checks if the driver name is valid
func isValidDriver(driver string) bool {
	for _, d := range ValidDrivers {
		if d == driver {
			return true
		}
	}
	return false
}

// validateMySQL validates MySQL-specific configuration (pure validation, no side effects)
func (c *DatabaseConfig) validateMySQL() error {
	cfg := c.Relational.MySQL

	if cfg.Host == "" {
		return errors.New("mysql host is required")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("mysql port must be between 1 and 65535")
	}
	if cfg.User == "" {
		return errors.New("mysql user is required")
	}
	if cfg.Password == "" {
		return errors.New("mysql password is required")
	}
	if cfg.DBName == "" {
		return errors.New("mysql dbname is required")
	}

	return nil
}

// validatePostgres validates PostgreSQL-specific configuration (pure validation, no side effects)
func (c *DatabaseConfig) validatePostgres() error {
	cfg := c.Relational.Postgres

	if cfg.Host == "" {
		return errors.New("postgres host is required")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("postgres port must be between 1 and 65535")
	}
	if cfg.User == "" {
		return errors.New("postgres user is required")
	}
	if cfg.Password == "" {
		return errors.New("postgres password is required")
	}
	if cfg.DBName == "" {
		return errors.New("postgres dbname is required")
	}

	return nil
}

// validateRedis validates Redis-specific configuration (pure validation, no side effects)
func (c *DatabaseConfig) validateRedis() error {
	if c.Redis.Addr == "" {
		return errors.New("redis addr is required")
	}

	return nil
}

// BuildMySQLDSN builds a MySQL DSN from the configuration
func (c *DatabaseConfig) BuildMySQLDSN() string {
	cfg := c.Relational.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	if cfg.Params != "" {
		dsn += "?" + cfg.Params
	}

	return dsn
}

// BuildPostgresDSN builds a PostgreSQL DSN from the configuration
func (c *DatabaseConfig) BuildPostgresDSN() string {
	cfg := c.Relational.Postgres

	// Build the connection string manually for better control
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	return connStr
}

// SelectedDriver returns the currently selected relational database driver
func (c *DatabaseConfig) SelectedDriver() string {
	return c.Relational.Driver
}
