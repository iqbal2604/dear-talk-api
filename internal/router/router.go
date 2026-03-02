package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
	"github.com/iqbal2604/dear-talk-api.git/internal/middleware"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type Handlers struct {
	AuthHandler    *handler.AuthHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func Setup(r *gin.Engine, h *Handlers) {
	r.GET("/health", func(c *gin.Context) {
		response.OK(c, "server is running", gin.H{
			"service": "ChatApp",
		})
	})

	v1 := r.Group("/api/v1")

	// Public routes
	registerAuthRoutes(v1, h.AuthHandler)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(h.AuthMiddleware.Authenticate())
	registerProtectedRoutes(protected)
}

func registerProtectedRoutes(protected *gin.RouterGroup) {
	// Sementara hanya /me, nanti akan diisi seiring phase berkembang
	protected.GET("/me", func(c *gin.Context) {
		response.OK(c, "ok", gin.H{
			"user_id":  c.GetUint("user_id"),
			"username": c.GetString("username"),
		})
	})
}
