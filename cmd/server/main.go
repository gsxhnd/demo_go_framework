package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"time"

	"go_sample_code/internal/database"
	"go_sample_code/internal/ent"
	healthhandler "go_sample_code/internal/handler/health"
	userhandler "go_sample_code/internal/handler/user"
	"go_sample_code/internal/middleware"
	userrepo "go_sample_code/internal/repo/user"
	userservice "go_sample_code/internal/service/user"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/metrics"
	"go_sample_code/pkg/trace"
	"go_sample_code/pkg/validator"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
			NewAppConfig,
			NewDatabaseConfig,
			NewLoggerConfig,
			NewTraceConfig,
			NewMetricsConfig,
			NewLogger,
			newEntClients,
			database.NewRedisClient,
			newHealthCheckerWithConfig,
			newTracerProvider,
			trace.NewTracer,
			newMeterProvider,
			metrics.NewHTTPRecorder,
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
func newTracerProvider(cfg *trace.TraceConfig) (*trace.TraceConfig, *sdktrace.TracerProvider, error) {
	tp, err := trace.NewTracerProvider(cfg)
	if err != nil {
		return cfg, nil, fmt.Errorf("failed to create tracer provider: %w", err)
	}
	return cfg, tp, nil
}

// newMeterProvider 创建 metrics provider
func newMeterProvider(cfg *metrics.MetricsConfig) (*metrics.MeterProvider, error) {
	return metrics.NewMeterProvider(cfg)
}

// NewLogger 创建日志实例
func NewLogger(cfg *logger.LoggerConfig) (logger.Logger, error) {
	return logger.NewLogger(cfg)
}

func NewFiberApp(cfg *AppConfig) *fiber.App {
	return fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
		Prefork:               false,
	})
}

// NewValidator 创建全局 validator 实例
func NewValidator() *validator.Validate {
	// 注册 UpdateUserRequest 结构级校验：至少传一个可更新字段
	v := validator.New()
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
	httpMetrics *metrics.HTTPRecorder,
	meterProvider *metrics.MeterProvider,
	tracerProvider *sdktrace.TracerProvider,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Middleware 顺序: Recovery -> RateLimit -> Metrics -> Logger
			app.Use(middleware.Recovery(log))
			app.Use(middleware.RateLimit(middleware.DefaultRateLimitConfig(log)))
			app.Use(middleware.Metrics(httpMetrics))
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

			// 1. 先关闭 HTTP server
			if err := app.Shutdown(); err != nil {
				log.Error("failed to shutdown fiber app", zap.Error(err))
			}

			// 2. 关闭 Redis client
			database.CloseRedisClient(redisClient, log)

			// 3. 关闭 Ent/sql
			database.CloseEntClient(db, entClient, log)

			// 4. Shutdown meter provider (flush metrics)
			if meterProvider != nil {
				shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				if err := meterProvider.Shutdown(shutdownCtx); err != nil {
					log.Error("failed to shutdown meter provider", zap.Error(err))
				}
			}

			// 5. Shutdown tracer provider
			if tracerProvider != nil {
				if err := tracerProvider.Shutdown(ctx); err != nil {
					log.Error("failed to shutdown tracer provider", zap.Error(err))
				}
			}

			// 6. Shutdown logger
			log.Shutdown(ctx)

			return nil
		},
	})
}
