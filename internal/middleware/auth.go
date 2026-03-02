package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/jwt"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type AuthMiddleware struct {
	jwtUtil        *jwt.JWTUtil
	tokenBlacklist domain.TokenBlacklist
}

func NewAuthMiddleware(jwtUtil *jwt.JWTUtil, tokenBlacklist domain.TokenBlacklist) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtil:        jwtUtil,
		tokenBlacklist: tokenBlacklist,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		//Cek apakah token sudah di blacklist
		blacklisted, err := m.tokenBlacklist.IsBlacklisted(c.Request.Context(), tokenStr)
		if err != nil || blacklisted {
			response.Unauthorized(c, "token has been invalidated")
		}

		claims, err := m.jwtUtil.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
