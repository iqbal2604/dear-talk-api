package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	accessExpire, err := time.ParseDuration(viper.GetString("JWT_ACCESS_EXPIRE"))
	if err != nil {
		accessExpire = 15 * time.Minute
	}

	refreshExpire, err := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRE"))
	if err != nil {
		refreshExpire = 168 * time.Hour
	}

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "Deartalk"),
			Env:  getEnv("APP_ENV", "production"),
			Port: getEnv("PORT", viper.GetString("APP_PORT")), // Railway pakai PORT
		},
		Database: DatabaseConfig{
			Host:     getEnv("PGHOST", viper.GetString("DB_HOST")),
			Port:     getEnv("PGPORT", viper.GetString("DB_PORT")),
			User:     getEnv("PGUSER", viper.GetString("DB_USER")),
			Password: getEnv("PGPASSWORD", viper.GetString("DB_PASSWORD")),
			Name:     getEnv("PGDATABASE", viper.GetString("DB_NAME")),
			SSLMode:  getEnv("DB_SSL_MODE", "require"), // Railway butuh SSL
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", viper.GetString("JWT_SECRET")),
			AccessExpire:  accessExpire,
			RefreshExpire: refreshExpire,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", viper.GetString("REDIS_HOST")),
			Port:     getEnv("REDIS_PORT", viper.GetString("REDIS_PORT")),
			Password: getEnv("REDIS_PASSWORD", viper.GetString("REDIS_PASSWORD")),
			DB:       0,
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
