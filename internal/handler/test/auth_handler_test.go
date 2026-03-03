package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
	"github.com/iqbal2604/dear-talk-api.git/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── Setup ────────────────────────────────────────────────────────────────────

func setupRouter(authHandler *handler.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)
	r.POST("/api/v1/auth/logout", authHandler.Logout)
	return r
}

func makeRequest(r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ─── Register Tests ───────────────────────────────────────────────────────────

func TestRegisterHandler_Success(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	req := domain.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	authUsecase.On("Register", mock.AnythingOfType("*domain.RegisterRequest")).
		Return(&domain.User{
			ID:       1,
			Username: req.Username,
			Email:    req.Email,
		}, nil)

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/register", req, "")

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "register success", response["message"])
	authUsecase.AssertExpectations(t)
}

func TestRegisterHandler_InvalidRequest(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	// Request tanpa email
	req := map[string]interface{}{
		"username": "johndoe",
		"password": "secret123",
	}

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/register", req, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, false, response["success"])
}

func TestRegisterHandler_EmailAlreadyExists(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	req := domain.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	authUsecase.On("Register", mock.AnythingOfType("*domain.RegisterRequest")).
		Return(nil, assert.AnError)

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/register", req, "")

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, false, response["success"])
	authUsecase.AssertExpectations(t)
}

// ─── Login Tests ──────────────────────────────────────────────────────────────

func TestLoginHandler_Success(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	req := domain.LoginRequest{
		Email:    "john@example.com",
		Password: "secret123",
	}

	authUsecase.On("Login", mock.AnythingOfType("*domain.LoginRequest")).
		Return(&domain.LoginResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			User: &domain.User{
				ID:       1,
				Username: "johndoe",
				Email:    req.Email,
			},
		}, nil)

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/login", req, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "login success", response["message"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	authUsecase.AssertExpectations(t)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	req := domain.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	authUsecase.On("Login", mock.AnythingOfType("*domain.LoginRequest")).
		Return(nil, assert.AnError)

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/login", req, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, false, response["success"])
	authUsecase.AssertExpectations(t)
}

func TestLoginHandler_InvalidRequest(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	// Request tanpa password
	req := map[string]interface{}{
		"email": "john@example.com",
	}

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/login", req, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, false, response["success"])
}

// ─── Logout Tests ─────────────────────────────────────────────────────────────

func TestLogoutHandler_Success(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	authUsecase.On("Logout", mock.Anything, mock.AnythingOfType("string")).
		Return(nil)

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/logout", nil, "valid-token")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "logout success", response["message"])
	authUsecase.AssertExpectations(t)
}

func TestLogoutHandler_InvalidToken(t *testing.T) {
	authUsecase := new(mocks.UserUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := setupRouter(authHandler)

	authUsecase.On("Logout", mock.Anything, mock.AnythingOfType("string")).
		Return(assert.AnError)

	w := makeRequest(r, http.MethodPost, "/api/v1/auth/logout", nil, "invalid-token")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, false, response["success"])
	authUsecase.AssertExpectations(t)
}
