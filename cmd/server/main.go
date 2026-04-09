package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"go_sample_code/internal/database"
	healthhandler "go_sample_code/internal/handler/health"
	userhandler "go_sample_code/internal/handler/user"
	"go_sample_code/internal/middleware"
	userrepo "go_sample_code/internal/repo/user"
	userservice "go_sample_code/internal/service/user"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/go-playground/validator/v10"
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
			newTracerProvider,
			trace.NewTracer,
			NewFiberApp,
			NewValidator,
			userrepo.NewUserRepo,
			userservice.NewUserService,
			healthhandler.NewHandler,
			userhandler.NewHandler,
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

// newTracerProvider creates a tracer provider for distributed tracing
func newTracerProvider() (*trace.TraceConfig, *sdktrace.TracerProvider) {
	cfg := &trace.TraceConfig{
		OtelEnable:         false, // 默认关闭，可通过配置启用
		OtelServiceName:    "demo-go-framework",
		OtelServiceVersion: "1.0.0",
	}
	tp, _ := trace.NewTracerProvider(cfg)
	return cfg, tp
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

// NewValidator 创建全局 validator 实例
func NewValidator() *validator.Validate {
	v := validator.New()
	// 设置 TagNameFunc，优先使用 json/query/params 标签名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Tag.Get("query")
		}
		if name == "" {
			name = fld.Tag.Get("params")
		}
		if name == "" {
			name = fld.Name
		}
		// 移除 omitempty,validate 等额外标签
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		return name
	})

	// 注册 UpdateUserRequest 结构级校验：至少传一个可更新字段
	v.RegisterStructValidation(func(sl validator.StructLevel) {
		req := sl.Current().Interface().(userhandler.UpdateUserRequest)
		if req.Email == nil && req.Password == nil && req.Nickname == nil &&
			req.Avatar == nil && req.Phone == nil && req.IsActive == nil {
			sl.ReportError(reflect.ValueOf(req), "UpdateUserRequest", "", "at_least_one_field", "")
		}
	}, userhandler.UpdateUserRequest{})

	return v
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
