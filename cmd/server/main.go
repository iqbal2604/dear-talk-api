package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
	"github.com/iqbal2604/dear-talk-api.git/pkg/database"
	"github.com/iqbal2604/dear-talk-api.git/pkg/logger"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Starting server...",
		zap.String("app", cfg.App.Env),
		zap.String("env", cfg.App.Name),
		zap.String("port", cfg.App.Port),
	)

	//Connect Database
	db, err := database.NewPostgresConnection(&cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect Database", zap.Error(err))
	}

	_ = db

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, "server is running", gin.H{
			"status":  "ok",
			"service": cfg.App.Name,
			"env":     cfg.App.Env,
		})
	})

	r.Run(":" + cfg.App.Port)

}
