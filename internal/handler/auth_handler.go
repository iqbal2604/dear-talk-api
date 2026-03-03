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

// Register godoc
// @Summary      Register akun baru
// @Description  Membuat akun user baru dengan username, email, dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      domain.RegisterRequest  true  "Register Request"
// @Success      201      {object}  response.Response{data=domain.User}
// @Failure      400      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Login user
// @Description  Login dengan email dan password, mendapatkan JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      domain.LoginRequest  true  "Login Request"
// @Success      200      {object}  response.Response{data=domain.LoginResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /auth/login [post]
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

// Logout godoc
// @Summary      Logout user
// @Description  Logout dan invalidate JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := h.authUsecase.Logout(c.Request.Context(), token); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "logout success", nil)
}
