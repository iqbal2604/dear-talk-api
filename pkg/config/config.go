package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
		},
	}
}
