package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
	"github.com/iqbal2604/dear-talk-api.git/pkg/validator"
)

type AuthHandler struct {
	authUsecase domain.UserUsecase
}

func NewAuthHandler(authUsecase domain.UserUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	// Validasi input
	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation error", errs)
		return
	}

	user, err := h.authUsecase.Register(&req)
	if err != nil {
		response.Conflict(c, err.Error())
		return
	}

	response.Created(c, "register success", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	// Validasi input
	if errs := validator.Validate(req); len(errs) > 0 {
		response.UnprocessableEntity(c, "validation error", errs)
		return
	}

	result, err := h.authUsecase.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "login success", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := h.authUsecase.Logout(c.Request.Context(), token); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "logout success", nil)
}
