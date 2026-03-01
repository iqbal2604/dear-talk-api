package database

import (
	"fmt"

	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(cfg *config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect database: %w", err)
	}

	log.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.String("db", cfg.Name),
	)

	return db, nil
}
