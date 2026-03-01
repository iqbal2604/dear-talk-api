package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	//Bungkus gin ke http server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//Jalankan Server di Goroutine
	go func() {
		log.Info("Server is running", zap.String("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start Server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	//Beri waktu 10 detik untuk request yang sedang berjalan selesai

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to Shutdown", zap.Error(err))
	}

	log.Info("Server exited gratefully")

}
