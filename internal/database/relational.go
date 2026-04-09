package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go_sample_code/internal/ent"
	"go_sample_code/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// NewEntClient creates a new ent client based on the configuration.
// It supports MySQL and PostgreSQL drivers dynamically.
func NewEntClient(cfg *DatabaseConfig, log logger.Logger) (*ent.Client, error) {
	driver := cfg.SelectedDriver()
	dsn := cfg.buildRelationalDSN()

	log.Info("initializing relational database",
		zap.String("driver", driver),
		zap.String("host", cfg.getRelationalHost()),
		zap.Int("port", cfg.getRelationalPort()),
		zap.String("database", cfg.getRelationalDBName()),
	)

	client, err := ent.Open(driver, dsn)
	if err != nil {
		log.Error("failed to open relational database",
			zap.String("driver", driver),
			zap.String("host", cfg.getRelationalHost()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to open relational database: %w", err)
	}

	// Perform startup health check
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := PingRelational(ctx, client); err != nil {
		client.Close()
		log.Error("failed to ping relational database during initialization",
			zap.String("driver", driver),
			zap.String("host", cfg.getRelationalHost()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to ping relational database: %w", err)
	}

	log.Info("relational database initialized successfully",
		zap.String("driver", driver),
		zap.String("host", cfg.getRelationalHost()),
	)

	return client, nil
}

// PingRelational performs a health check on the relational database
// by executing a simple query through ent
func PingRelational(ctx context.Context, client *ent.Client) error {
	// Use ent's User query to verify the database connection
	// This will execute a SELECT query which tests the connection
	_, err := client.User.Query().Limit(1).All(ctx)
	if err != nil {
		// If the error is "sql: no rows in result set", the connection is actually working
		// We need to check if it's a connection error vs a "no rows" error
		return fmt.Errorf("failed to query relational database: %w", err)
	}
	return nil
}

// CloseEntClient closes the ent client gracefully
func CloseEntClient(client *ent.Client, log logger.Logger) {
	if client != nil {
		if err := client.Close(); err != nil {
			log.Error("failed to close relational database", zap.Error(err))
		} else {
			log.Info("relational database connection closed")
		}
	}
}

// buildRelationalDSN builds the appropriate DSN based on the selected driver
func (c *DatabaseConfig) buildRelationalDSN() string {
	switch c.Relational.Driver {
	case DriverMySQL:
		return c.BuildMySQLDSN()
	case DriverPostgres:
		return c.BuildPostgresDSN()
	default:
		return ""
	}
}

// getRelationalHost returns the host of the selected relational database
func (c *DatabaseConfig) getRelationalHost() string {
	switch c.Relational.Driver {
	case DriverMySQL:
		return c.Relational.MySQL.Host
	case DriverPostgres:
		return c.Relational.Postgres.Host
	default:
		return ""
	}
}

// getRelationalPort returns the port of the selected relational database
func (c *DatabaseConfig) getRelationalPort() int {
	switch c.Relational.Driver {
	case DriverMySQL:
		return c.Relational.MySQL.Port
	case DriverPostgres:
		return c.Relational.Postgres.Port
	default:
		return 0
	}
}

// getRelationalDBName returns the database name of the selected relational database
func (c *DatabaseConfig) getRelationalDBName() string {
	switch c.Relational.Driver {
	case DriverMySQL:
		return c.Relational.MySQL.DBName
	case DriverPostgres:
		return c.Relational.Postgres.DBName
	default:
		return ""
	}
}

// Ensure sql.DB interface is available
var _ = (*sql.DB)(nil)
