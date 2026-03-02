package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.UserUsecase
}

func NewAuthHandler(authUsecase domain.UserUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// ─── Register ─────────────────────────────────────────────────────────────────

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid Request", err.Error())
		return
	}

	user, err := h.authUsecase.Register(&req)
	if err != nil {
		response.Conflict(c, err.Error())
		return
	}

	response.Created(c, "Register Success", user)
}

// ─── Login ────────────────────────────────────────────────────────────────────

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid Request", err.Error())
		return
	}

	result, err := h.authUsecase.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "Login Success", result)
}

// ─── Logout ────────────────────────────────────────────────────────────────────
func (h *AuthHandler) Logout(c *gin.Context) {
	//Ambil token dari header
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := h.authUsecase.Logout(c.Request.Context(), token); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "logout success", nil)
}
