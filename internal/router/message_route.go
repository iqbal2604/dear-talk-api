package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
)

func registerMessageRoutes(protected *gin.RouterGroup, messageHandler *handler.MessageHandler) {
	// Routes di bawah rooms
	rooms := protected.Group("/rooms")
	{
		rooms.POST("/:id/messages", messageHandler.SendMessage)
		rooms.GET("/:id/messages", messageHandler.GetMessages)
		rooms.POST("/:id/read", messageHandler.MarkAsRead)
	}

	// Routes khusus message
	messages := protected.Group("/messages")
	{
		messages.PUT("/:id", messageHandler.EditMessage)
		messages.DELETE("/:id", messageHandler.DeleteMessage)
	}
}
