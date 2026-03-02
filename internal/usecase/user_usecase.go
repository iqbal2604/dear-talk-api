package usecase

import (
	"errors"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
)

type userManagementUsecase struct {
	userRepo domain.UserRepository
}

func NewUserManagementUsecase(userRepo domain.UserRepository) domain.UserManagementUsecase {
	return &userManagementUsecase{
		userRepo: userRepo,
	}
}

// ─── Get Profile ──────────────────────────────────────────────────────────────

func (u *userManagementUsecase) GetProfile(id uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// ─── Update Profile ──────────────────────────────────────────────────────────────
func (u *userManagementUsecase) UpdateProfile(id uint, req *domain.UpdateProfileRequest) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	//Cek kalo username sudah dipakai oleh orang lain
	if req.Username != "" && req.Username != user.Username {
		existing, err := u.userRepo.FindByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("username already taken")
		}
		user.Username = req.Username
	}

	//Update Avatar jika ada
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ─── Search Users ─────────────────────────────────────────────────────────────
func (u *userManagementUsecase) SearchUsers(query string) ([]*domain.User, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}

	users, err := u.userRepo.Search(query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// ─── Get User By ID ─────────────────────────────────────────────────────────────

func (u *userManagementUsecase) GetUserByID(id uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
