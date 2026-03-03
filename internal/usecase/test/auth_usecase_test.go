package test

import (
	"context"
	"testing"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/mocks"
	"github.com/iqbal2604/dear-talk-api.git/internal/usecase"
	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
	"github.com/iqbal2604/dear-talk-api.git/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ─── Helper ───────────────────────────────────────────────────────────────────

func newJWTUtil() *jwt.JWTUtil {
	return jwt.NewJWTUtil(&config.JWTConfig{
		Secret:        "test-secret",
		AccessExpire:  15 * 60 * 1000000000,       // 15 menit
		RefreshExpire: 168 * 60 * 60 * 1000000000, // 7 hari
	})
}

func hashedPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// ─── Register Tests ───────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	req := &domain.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	// Mock — email belum ada
	userRepo.On("FindByEmail", req.Email).Return(nil, nil)
	// Mock — username belum ada
	userRepo.On("FindByUsername", req.Username).Return(nil, nil)
	// Mock — create user berhasil
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := authUsecase.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	userRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	req := &domain.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	// Mock — email sudah ada
	userRepo.On("FindByEmail", req.Email).Return(&domain.User{
		ID:    1,
		Email: req.Email,
	}, nil)

	user, err := authUsecase.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already registered", err.Error())
	userRepo.AssertExpectations(t)
}

func TestRegister_UsernameAlreadyTaken(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	req := &domain.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	// Mock — email belum ada
	userRepo.On("FindByEmail", req.Email).Return(nil, nil)
	// Mock — username sudah ada
	userRepo.On("FindByUsername", req.Username).Return(&domain.User{
		ID:       1,
		Username: req.Username,
	}, nil)

	user, err := authUsecase.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username already taken", err.Error())
	userRepo.AssertExpectations(t)
}

// ─── Login Tests ──────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	req := &domain.LoginRequest{
		Email:    "john@example.com",
		Password: "secret123",
	}

	// Mock — user ditemukan dengan password ter-hash
	userRepo.On("FindByEmail", req.Email).Return(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    req.Email,
		Password: hashedPassword(req.Password),
	}, nil)

	result, err := authUsecase.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, req.Email, result.User.Email)
	userRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	req := &domain.LoginRequest{
		Email:    "notfound@example.com",
		Password: "secret123",
	}

	// Mock — user tidak ditemukan
	userRepo.On("FindByEmail", req.Email).Return(nil, nil)

	result, err := authUsecase.Login(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid email or password", err.Error())
	userRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	req := &domain.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	// Mock — user ditemukan tapi password salah
	userRepo.On("FindByEmail", req.Email).Return(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    req.Email,
		Password: hashedPassword("secret123"),
	}, nil)

	result, err := authUsecase.Login(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid email or password", err.Error())
	userRepo.AssertExpectations(t)
}

// ─── Logout Tests ─────────────────────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	// Generate token valid dulu
	token, _ := jwtUtil.GenerateAccessToken(1, "johndoe")

	// Mock — blacklist berhasil
	tokenBlacklist.On("Add", context.Background(), token, mock.Anything).Return(nil)

	err := authUsecase.Logout(context.Background(), token)

	assert.NoError(t, err)
	tokenBlacklist.AssertExpectations(t)
}

func TestLogout_InvalidToken(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	tokenBlacklist := new(mocks.TokenBlacklist)
	jwtUtil := newJWTUtil()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, tokenBlacklist)

	err := authUsecase.Logout(context.Background(), "invalid-token")

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
}
