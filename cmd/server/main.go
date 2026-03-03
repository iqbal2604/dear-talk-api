package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/repository/model"
	"github.com/iqbal2604/dear-talk-api.git/internal/router"
	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
	"github.com/iqbal2604/dear-talk-api.git/pkg/logger"
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
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// Wire
	app, err := InitializeApp(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize app", zap.Error(err))
	}

	// Auto migrate
	if err := app.DB.AutoMigrate(&model.UserModel{}, &model.RoomModel{}, &model.RoomMemberModel{}, &model.MessageModel{}, &model.ReadStatusModel{}); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}
	log.Info("Database migrated successfully")

	// Setup router
	r := gin.Default()
	router.Setup(r, &router.Handlers{
		AuthHandler:    app.AuthHandler,
		AuthMiddleware: app.AuthMiddleware,
		UserHandler:    app.UserHandler,
		MessageHandler: app.MessageHandler,
		WSHandler:      app.WSHandler,
		RoomHandler:    app.RoomHandler,
	})

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Server is running", zap.String("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}
