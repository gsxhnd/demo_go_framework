package database

import (
	"context"
	"database/sql"
	"fmt"

	"go_sample_code/internal/ent"
	"go_sample_code/pkg/logger"

	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// configurePool applies the connection pool configuration to the sql.DB instance
func configurePool(db *sql.DB, pool PoolConfig) {
	if pool.MaxOpenConns > 0 {
		db.SetMaxOpenConns(pool.MaxOpenConns)
	}
	if pool.MaxIdleConns > 0 {
		db.SetMaxIdleConns(pool.MaxIdleConns)
	}
	if pool.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(pool.ConnMaxLifetime)
	}
	if pool.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
	}
}

// NewEntClient creates a new ent client based on the configuration.
// It supports MySQL and PostgreSQL drivers dynamically.
func NewEntClient(cfg *DatabaseConfig, log logger.Logger) (*sql.DB, *ent.Client, error) {
	driver := cfg.SelectedDriver()
	dsn := cfg.buildRelationalDSN()

	log.Info("initializing relational database",
		zap.String("driver", driver),
		zap.String("host", cfg.getRelationalHost()),
		zap.Int("port", cfg.getRelationalPort()),
		zap.String("database", cfg.getRelationalDBName()),
	)

	// Open sql.DB directly to access connection pool configuration
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Error("failed to open database connection",
			zap.String("driver", driver),
			zap.String("host", cfg.getRelationalHost()),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Apply connection pool configuration
	configurePool(db, cfg.getPoolConfig())

	// Create ent client using the configured sql.DB
	entDriver := entsql.OpenDB(driver, db)
	client := ent.NewClient(ent.Driver(entDriver))

	// Perform startup health check
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetHealthCheckTimeout())
	defer cancel()

	if err := PingRelational(ctx, db); err != nil {
		client.Close()
		db.Close()
		log.Error("failed to ping relational database during initialization",
			zap.String("driver", driver),
			zap.String("host", cfg.getRelationalHost()),
			zap.Error(err),
		)
		return nil, nil, fmt.Errorf("failed to ping relational database: %w", err)
	}

	log.Info("relational database initialized successfully",
		zap.String("driver", driver),
		zap.String("host", cfg.getRelationalHost()),
	)

	return db, client, nil
}

// PingRelational performs a health check on the relational database
// using sql.DB.PingContext for a generic and reliable connection test
func PingRelational(ctx context.Context, db *sql.DB) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping relational database: %w", err)
	}
	return nil
}

// PingRelationalWithEnt performs health check using ent client
// (fallback when direct sql.DB is not available)
func PingRelationalWithEnt(ctx context.Context, client *ent.Client) error {
	_, err := client.User.Query().Limit(1).All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query relational database: %w", err)
	}
	return nil
}

// CloseEntClient closes the ent client and underlying sql.DB
func CloseEntClient(db *sql.DB, client *ent.Client, log logger.Logger) {
	if client != nil {
		if err := client.Close(); err != nil {
			log.Error("failed to close ent client", zap.Error(err))
		}
	}
	if db != nil {
		if err := db.Close(); err != nil {
			log.Error("failed to close database connection", zap.Error(err))
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
