package usecase

import (
	"context"
	"errors"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo       domain.UserRepository
	jwtUtil        *jwt.JWTUtil
	tokenBlacklist domain.TokenBlacklist
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtUtil *jwt.JWTUtil, tokenBlacklist domain.TokenBlacklist) domain.UserUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		jwtUtil:        jwtUtil,
		tokenBlacklist: tokenBlacklist,
	}
}

// ─── Register ─────────────────────────────────────────────────────────────────

func (u *authUsecase) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// Cek apakah email sudah dipakai
	existing, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Cek apakah username sudah dipakai
	existingUsername, err := u.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUsername != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Buat user baru
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ─── Login ────────────────────────────────────────────────────────────────────

func (u *authUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Cari user by email
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := u.jwtUtil.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwtUtil.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, token string) error {
	//Validasi token terlebih dahulu untuk ambil expiry
	claims, err := u.jwtUtil.ValidateToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	//Hitung sisa waktu token
	expiry := claims.ExpiresAt.Sub(claims.IssuedAt.Time)

	//Masukkan ke blaclist
	return u.tokenBlacklist.Add(ctx, token, expiry)
}
