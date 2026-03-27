package main

import (
	"context"
	"flag"
	"time"

	healthhandler "go_sample_code/internal/handler/health"
	"go_sample_code/internal/middleware"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
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
			NewLogger,
			NewFiberApp,
			healthhandler.NewHandler,
		),
		fx.Invoke(RegisterHooks),
	).Run()
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
			return app.Shutdown()
		},
	})
}
