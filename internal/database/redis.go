package database

import (
	"context"
	"fmt"
	"time"

	"go_sample_code/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewRedisClient creates a new Redis client based on the configuration
func NewRedisClient(cfg *DatabaseConfig, log logger.Logger) (*redis.Client, error) {
	log.Info("initializing Redis client",
		zap.String("addr", cfg.Redis.Addr),
		zap.Int("db", cfg.Redis.DB),
	)

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Username:     cfg.Redis.Username,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	})

	// Perform startup health check
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := PingRedis(ctx, client); err != nil {
		client.Close()
		log.Error("failed to ping Redis during initialization",
			zap.String("addr", cfg.Redis.Addr),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info("Redis client initialized successfully",
		zap.String("addr", cfg.Redis.Addr),
	)

	return client, nil
}

// PingRedis performs a health check on the Redis client
func PingRedis(ctx context.Context, client *redis.Client) error {
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}
	return nil
}

// CloseRedisClient closes the Redis client gracefully
func CloseRedisClient(client *redis.Client, log logger.Logger) {
	if client != nil {
		if err := client.Close(); err != nil {
			log.Error("failed to close Redis client", zap.Error(err))
		} else {
			log.Info("Redis client connection closed")
		}
	}
}
