package database_test

import (
	"testing"
	"time"

	"go_sample_code/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *database.DatabaseConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid MySQL config",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverMySQL,
					MySQL: database.MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
						Pool:     database.DefaultPoolConfig(),
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: false,
		},
		{
			name: "valid PostgreSQL config",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverPostgres,
					Postgres: database.PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
						Pool:     database.DefaultPoolConfig(),
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid driver",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: "oracle",
					MySQL: database.MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "invalid relational driver",
		},
		{
			name: "MySQL missing host",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverMySQL,
					MySQL: database.MySQLConfig{
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "mysql host is required",
		},
		{
			name: "MySQL missing user",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverMySQL,
					MySQL: database.MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "mysql user is required",
		},
		{
			name: "MySQL invalid port",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverMySQL,
					MySQL: database.MySQLConfig{
						Host:     "localhost",
						Port:     0,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "mysql port must be between 1 and 65535",
		},
		{
			name: "Postgres missing host",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverPostgres,
					Postgres: database.PostgresConfig{
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "postgres host is required",
		},
		{
			name: "Postgres missing user",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverPostgres,
					Postgres: database.PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
					},
				},
				Redis: database.RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "postgres user is required",
		},
		{
			name: "Redis missing addr",
			cfg: &database.DatabaseConfig{
				Relational: database.RelationalConfig{
					Driver: database.DriverPostgres,
					Postgres: database.PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
						Pool:     database.DefaultPoolConfig(),
					},
				},
				Redis: database.RedisConfig{},
			},
			wantErr: true,
			errMsg:  "redis addr is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultPoolConfig(t *testing.T) {
	cfg := database.DefaultPoolConfig()
	assert.Equal(t, 25, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxLifetime)
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxIdleTime)
}

func TestBuildMySQLDSN(t *testing.T) {
	cfg := &database.DatabaseConfig{
		Relational: database.RelationalConfig{
			Driver: database.DriverMySQL,
			MySQL: database.MySQLConfig{
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "password",
				DBName:   "testdb",
			},
		},
		Redis: database.RedisConfig{
			Addr: "localhost:6379",
		},
	}
	dsn := cfg.BuildMySQLDSN()
	assert.Contains(t, dsn, "root:password@tcp(localhost:3306)/testdb")
}

func TestBuildMySQLDSN_WithParams(t *testing.T) {
	cfg := &database.DatabaseConfig{
		Relational: database.RelationalConfig{
			Driver: database.DriverMySQL,
			MySQL: database.MySQLConfig{
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "password",
				DBName:   "testdb",
				Params:   "parseTime=true&charset=utf8mb4",
			},
		},
		Redis: database.RedisConfig{
			Addr: "localhost:6379",
		},
	}
	dsn := cfg.BuildMySQLDSN()
	assert.Contains(t, dsn, "root:password@tcp(localhost:3306)/testdb")
	assert.Contains(t, dsn, "parseTime=true")
	assert.Contains(t, dsn, "charset=utf8mb4")
}

func TestBuildPostgresDSN(t *testing.T) {
	cfg := &database.DatabaseConfig{
		Relational: database.RelationalConfig{
			Driver: database.DriverPostgres,
			Postgres: database.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "password",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
		},
		Redis: database.RedisConfig{
			Addr: "localhost:6379",
		},
	}
	dsn := cfg.BuildPostgresDSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=postgres")
	assert.Contains(t, dsn, "dbname=testdb")
	assert.Contains(t, dsn, "sslmode=disable")
}
