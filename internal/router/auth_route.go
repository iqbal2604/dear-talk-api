package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
)

func registerAuthRoutes(
	v1 *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	strictLimiter gin.HandlerFunc,
) {
	auth := v1.Group("/auth")
	auth.Use(strictLimiter)
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}
}
