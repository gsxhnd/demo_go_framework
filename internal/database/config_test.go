package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *DatabaseConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid MySQL config",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverMySQL,
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
						Pool:     DefaultPoolConfig(),
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: false,
		},
		{
			name: "valid PostgreSQL config",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverPostgres,
					Postgres: PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
						Pool:     DefaultPoolConfig(),
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid driver",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: "oracle",
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "invalid relational driver",
		},
		{
			name: "MySQL missing host",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverMySQL,
					MySQL: MySQLConfig{
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "mysql host is required",
		},
		{
			name: "MySQL missing user",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverMySQL,
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "mysql user is required",
		},
		{
			name: "MySQL invalid port",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverMySQL,
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     0,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "mysql port must be between 1 and 65535",
		},
		{
			name: "Postgres missing host",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverPostgres,
					Postgres: PostgresConfig{
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "postgres host is required",
		},
		{
			name: "Postgres missing user",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverPostgres,
					Postgres: PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
					},
				},
				Redis: RedisConfig{
					Addr: "localhost:6379",
				},
			},
			wantErr: true,
			errMsg:  "postgres user is required",
		},
		{
			name: "Redis missing addr",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Driver: DriverPostgres,
					Postgres: PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
						Pool:     DefaultPoolConfig(),
					},
				},
				Redis: RedisConfig{},
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
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDatabaseConfig_BuildMySQLDSN(t *testing.T) {
	tests := []struct {
		name string
		cfg  *DatabaseConfig
		want string
	}{
		{
			name: "basic DSN",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
					},
				},
			},
			want: "root:password@tcp(localhost:3306)/testdb",
		},
		{
			name: "DSN with params",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					MySQL: MySQLConfig{
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "password",
						DBName:   "testdb",
						Params:   "charset=utf8mb4&parseTime=True",
					},
				},
			},
			want: "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True",
		},
		{
			name: "DSN with custom port",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					MySQL: MySQLConfig{
						Host:     "192.168.1.100",
						Port:     3307,
						User:     "admin",
						Password: "secret",
						DBName:   "mydb",
					},
				},
			},
			want: "admin:secret@tcp(192.168.1.100:3307)/mydb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.BuildMySQLDSN()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDatabaseConfig_BuildPostgresDSN(t *testing.T) {
	tests := []struct {
		name string
		cfg  *DatabaseConfig
		want string
	}{
		{
			name: "basic DSN",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Postgres: PostgresConfig{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "password",
						DBName:   "testdb",
						SSLMode:  "disable",
					},
				},
			},
			want: "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable",
		},
		{
			name: "DSN with require ssl",
			cfg: &DatabaseConfig{
				Relational: RelationalConfig{
					Postgres: PostgresConfig{
						Host:     "192.168.1.100",
						Port:     5433,
						User:     "admin",
						Password: "secret",
						DBName:   "mydb",
						SSLMode:  "require",
					},
				},
			},
			want: "host=192.168.1.100 port=5433 user=admin password=secret dbname=mydb sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.BuildPostgresDSN()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultPoolConfig(t *testing.T) {
	cfg := DefaultPoolConfig()

	assert.Equal(t, 25, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxLifetime)
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxIdleTime)
}

func TestDefaultRedisConfig(t *testing.T) {
	cfg := DefaultRedisConfig()

	assert.Equal(t, 10, cfg.PoolSize)
	assert.Equal(t, 5, cfg.MinIdleConns)
	assert.Equal(t, 5*time.Second, cfg.DialTimeout)
	assert.Equal(t, 3*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 3*time.Second, cfg.WriteTimeout)
}

func TestIsValidDriver(t *testing.T) {
	tests := []struct {
		driver string
		want   bool
	}{
		{DriverMySQL, true},
		{DriverPostgres, true},
		{"oracle", false},
		{"sqlite", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.driver, func(t *testing.T) {
			got := isValidDriver(tt.driver)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectedDriver(t *testing.T) {
	cfg := &DatabaseConfig{
		Relational: RelationalConfig{
			Driver: DriverMySQL,
		},
	}

	got := cfg.SelectedDriver()
	assert.Equal(t, DriverMySQL, got)

	cfg.Relational.Driver = DriverPostgres
	got = cfg.SelectedDriver()
	assert.Equal(t, DriverPostgres, got)
}
