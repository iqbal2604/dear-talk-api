package domain

import "time"

//───Entity──────────────────────────────────────────────────────────────

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	Avatar    string
	IsOnline  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

//───Repository Interface─────────────────────────────────────────────────

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *User) error
}

//───Usecase Interface─────────────────────────────────────────────────

type UserUsecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
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
