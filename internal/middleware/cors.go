package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Di production ganti dengan domain spesifik
		// contoh: []string{"https://yourdomain.com"}
		AllowOrigins: []string{"*"},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},

		ExposeHeaders: []string{
			"Content-Length",
		},

		// Izinkan credentials (cookie, auth header)
		AllowCredentials: true,

		// Cache preflight request selama 12 jam
		MaxAge: 12 * time.Hour,
	})
}
