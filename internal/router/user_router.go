package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
)

func registerUserRoutes(protected *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := protected.Group("/users")
	{
		users.GET("/me", userHandler.GetMyProfile)
		users.PUT("/me", userHandler.UpdateMyProfile)
		users.GET("/search", userHandler.SearchUsers)
		users.GET("/:id", userHandler.GetUserByID)
	}
}
