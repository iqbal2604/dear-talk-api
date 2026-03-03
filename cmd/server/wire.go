//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
	"github.com/iqbal2604/dear-talk-api.git/internal/middleware"
	"github.com/iqbal2604/dear-talk-api.git/internal/repository"
	"github.com/iqbal2604/dear-talk-api.git/internal/usecase"
	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
	"github.com/iqbal2604/dear-talk-api.git/pkg/database"
	"github.com/iqbal2604/dear-talk-api.git/pkg/jwt"
	"github.com/iqbal2604/dear-talk-api.git/pkg/redis"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ─── Provider Sets ────────────────────────────────────────────────────────────
var infrastructureSet = wire.NewSet(
	database.NewPostgresConnection,
	jwt.NewJWTUtil,
	redis.NewRedisClient,
	redis.NewTokenBlacklist,
	wire.Bind(new(domain.TokenBlacklist), new(*redis.TokenBlacklist)),
)

var repositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewRoomRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewAuthUsecase,
	usecase.NewUserManagementUsecase,
	usecase.NewRoomUsecase,
)

var handlerSet = wire.NewSet(
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewRoomHandler,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

// ─── App struct ───────────────────────────────────────────────────────────────

type App struct {
	DB             *gorm.DB
	Redis          *goredis.Client
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	RoomHandler    *handler.RoomHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func NewApp(
	db *gorm.DB,
	redisClient *goredis.Client,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roomHandler *handler.RoomHandler,
	authMiddleware *middleware.AuthMiddleware,
) *App {
	return &App{
		DB:             db,
		Redis:          redisClient,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		RoomHandler:    roomHandler,
		AuthMiddleware: authMiddleware,
	}
}

func InitializeApp(cfg *config.Config, log *zap.Logger) (*App, error) {
	wire.Build(
		wire.FieldsOf(new(*config.Config), "Database", "JWT", "Redis"),

		infrastructureSet,
		repositorySet,
		usecaseSet,
		handlerSet,
		middlewareSet,

		NewApp,
	)
	return nil, nil
}
