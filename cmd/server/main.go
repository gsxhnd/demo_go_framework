package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"go_sample_code/internal/database"
	healthhandler "go_sample_code/internal/handler/health"
	"go_sample_code/internal/middleware"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"go_sample_code/internal/ent"
)

var cfgPathFlag = flag.String("c", "config.yaml", "")

func main() {
	flag.Parse()
	var cfgPath ConfigPath = ConfigPath(*cfgPathFlag)

	fx.New(
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(30*time.Second),
		fx.Supply(cfgPath),
		fx.Provide(
			NewLogger,
			NewConfig,
			newEntClients,
			database.NewRedisClient,
			newHealthCheckerWithConfig,
			NewFiberApp,
			healthhandler.NewHandler,
		),
		fx.Invoke(RegisterHooks),
	).Run()
}

// newEntClients creates both the sql.DB and ent.Client
func newEntClients(cfg *database.DatabaseConfig, log logger.Logger) (*sql.DB, *ent.Client, error) {
	return database.NewEntClient(cfg, log)
}

// newHealthCheckerWithConfig creates a health checker with configuration parameters
func newHealthCheckerWithConfig(
	db *sql.DB,
	redisClient *redis.Client,
	log logger.Logger,
	cfg *database.DatabaseConfig,
) database.HealthChecker {
	return database.NewHealthChecker(
		db,
		redisClient,
		cfg.SelectedDriver(),
		log,
		cfg.GetHealthCheckTimeout(),
	)
}

func NewConfig(cfgPath ConfigPath) (*database.DatabaseConfig, error) {
	var cfg database.DatabaseConfig

	data, err := os.ReadFile(string(cfgPath))
	if err != nil {
		// If config file doesn't exist, use defaults
		cfg = database.DatabaseConfig{
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
				MySQL: database.MySQLConfig{
					Host:     "localhost",
					Port:     3306,
					User:     "root",
					Password: "root",
					DBName:   "demo",
					Pool:     database.DefaultPoolConfig(),
				},
			},
			Redis: database.DefaultRedisConfig(),
		}
		cfg.Redis.Addr = "localhost:6379"
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Apply default values to configuration
	cfg.ApplyDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func NewLogger(cfgPath ConfigPath) (logger.Logger, error) {
	cfg := logger.DefaultConfig()
	return logger.NewLogger(cfg)
}

func NewFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
		Prefork:               false,
	})
}

func RegisterHooks(
	lifecycle fx.Lifecycle,
	app *fiber.App,
	log logger.Logger,
	db *sql.DB,
	entClient *ent.Client,
	redisClient *redis.Client,
	healthChecker database.HealthChecker,
	healthHandler healthhandler.Handler,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			app.Use(middleware.Recovery(log))
			app.Use(middleware.Logger(log))

			app.Get("/api/health", healthHandler.Check)

			go func() {
				if err := app.Listen(":8080"); err != nil {
					log.Error("failed to start server", zap.Error(err))
				}
			}()

			log.Info("server started on :8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down server")

			// Close Redis client
			database.CloseRedisClient(redisClient, log)

			// Close ent client and sql.DB
			database.CloseEntClient(db, entClient, log)

			// Shutdown Fiber
			if err := app.Shutdown(); err != nil {
				log.Error("failed to shutdown fiber app", zap.Error(err))
			}

			return nil
		},
	})
}
