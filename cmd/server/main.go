package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
)

func main() {
	cfg := config.Load()
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "deartalk-backend",
		})
	})

	r.Run(":" + cfg.App.Port)

}
