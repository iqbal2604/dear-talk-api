package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	cloudinarypkg "github.com/iqbal2604/dear-talk-api.git/pkg/cloudinary"
)

type userManagementUsecase struct {
	userRepo   domain.UserRepository
	cloudinary *cloudinarypkg.CloudinaryClient
}

func NewUserManagementUsecase(userRepo domain.UserRepository, cloudinary *cloudinarypkg.CloudinaryClient) domain.UserManagementUsecase {
	return &userManagementUsecase{
		userRepo:   userRepo,
		cloudinary: cloudinary,
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

func (u *userManagementUsecase) UploadAvatar(
	userID uint,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (*domain.User, error) {
	// Validasi ukuran file max 2MB
	if fileHeader.Size > 2*1024*1024 {
		return nil, errors.New("ukuran file maksimal 2MB")
	}

	// Validasi tipe file
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, errors.New("tipe file tidak didukung, gunakan JPEG/PNG/WebP")
	}

	// Upload ke Cloudinary
	ctx := context.Background()
	avatarURL, err := u.cloudinary.UploadAvatar(ctx, file, userID)
	if err != nil {
		return nil, errors.New("gagal upload avatar")
	}

	// Update database
	if err := u.userRepo.UpdateAvatar(userID, avatarURL); err != nil {
		return nil, errors.New("gagal menyimpan avatar")
	}

	// Return updated user
	return u.userRepo.FindByID(userID)
}
