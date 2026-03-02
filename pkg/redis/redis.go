package redis

import (
	"context"
	"fmt"

	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(cfg *config.RedisConfig, log *zap.Logger) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping untuk pastikan koneksi berhasil
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	log.Info("Redis connected successfully",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
	)

	return client, nil
}
