package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
	"github.com/iqbal2604/dear-talk-api.git/internal/middleware"
	"github.com/iqbal2604/dear-talk-api.git/internal/websocket"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type Handlers struct {
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	RoomHandler    *handler.RoomHandler
	AuthMiddleware *middleware.AuthMiddleware
	WSHandler      *websocket.WSHandler
	MessageHandler *handler.MessageHandler
}

func Setup(r *gin.Engine, h *Handlers) {
	//Global Middlewares
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, "server is running", gin.H{
			"service": "DearTalk",
		})
	})

	// WebSocket route (auth via query param token)
	r.GET("/ws", h.WSHandler.ServeWS)

	v1 := r.Group("/api/v1")

	// Public routes
	registerAuthRoutes(v1, h.AuthHandler)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(h.AuthMiddleware.Authenticate())
	registerUserRoutes(protected, h.UserHandler)
	registerRoomRoutes(protected, h.RoomHandler)
	registerMessageRoutes(protected, h.MessageHandler)
}
