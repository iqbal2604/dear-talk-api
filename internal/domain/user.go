package domain

import (
	"context"
	"time"
)

//───Entity──────────────────────────────────────────────────────────────

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Avatar    string    `json:"avatar"`
	IsOnline  bool      `json:"isonline"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//───Repository Interface─────────────────────────────────────────────────

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Search(query string) ([]*User, error)
	Update(user *User) error
}

//───Token Blacklist Interface───────────────────────────────────────────

type TokenBlacklist interface {
	Add(ctx context.Context, token string, expiry time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

//───Usecase Interface─────────────────────────────────────────────────

type UserUsecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
}

type UserManagementUsecase interface {
	GetProfile(id uint) (*User, error)
	UpdateProfile(id uint, req *UpdateProfileRequest) (*User, error)
	SearchUsers(query string) ([]*User, error)
	GetUserByID(id uint) (*User, error)
}

//───Request dan Response─────────────────────────────────────────────────

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Avatar   string `json:"avatar" binfing:"omitempty,url"`
}
