package test

import (
	"testing"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/mocks"
	"github.com/iqbal2604/dear-talk-api.git/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── GetProfile Tests ─────────────────────────────────────────────────────────

func TestGetProfile_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	expectedUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	userRepo.On("FindByID", uint(1)).Return(expectedUser, nil)

	user, err := userUsecase.GetProfile(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	userRepo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	// Mock — user tidak ditemukan
	userRepo.On("FindByID", uint(99)).Return(nil, nil)

	user, err := userUsecase.GetProfile(99)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

// ─── UpdateProfile Tests ──────────────────────────────────────────────────────

func TestUpdateProfile_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	existingUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	req := &domain.UpdateProfileRequest{
		Username: "johnnew",
	}

	// Mock — user ditemukan
	userRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	// Mock — username baru belum dipakai
	userRepo.On("FindByUsername", req.Username).Return(nil, nil)
	// Mock — update berhasil
	userRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := userUsecase.UpdateProfile(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	userRepo.AssertExpectations(t)
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	req := &domain.UpdateProfileRequest{
		Username: "johnnew",
	}

	// Mock — user tidak ditemukan
	userRepo.On("FindByID", uint(99)).Return(nil, nil)

	user, err := userUsecase.UpdateProfile(99, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

func TestUpdateProfile_UsernameTaken(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	existingUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	req := &domain.UpdateProfileRequest{
		Username: "janedoe",
	}

	// Mock — user ditemukan
	userRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	// Mock — username sudah dipakai orang lain
	userRepo.On("FindByUsername", req.Username).Return(&domain.User{
		ID:       2,
		Username: "janedoe",
	}, nil)

	user, err := userUsecase.UpdateProfile(1, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username already taken", err.Error())
	userRepo.AssertExpectations(t)
}

func TestUpdateProfile_SameUsername(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	existingUser := &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	// Update dengan username yang sama — tidak perlu cek
	req := &domain.UpdateProfileRequest{
		Username: "johndoe",
	}

	// Mock — user ditemukan
	userRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	// Mock — update berhasil tanpa cek username
	userRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := userUsecase.UpdateProfile(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "johndoe", user.Username)
	userRepo.AssertExpectations(t)
}

// ─── SearchUsers Tests ────────────────────────────────────────────────────────

func TestSearchUsers_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	expectedUsers := []*domain.User{
		{ID: 1, Username: "johndoe", Email: "john@example.com"},
		{ID: 2, Username: "johnjr", Email: "johnjr@example.com"},
	}

	userRepo.On("Search", "john").Return(expectedUsers, nil)

	users, err := userUsecase.SearchUsers("john")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	userRepo.AssertExpectations(t)
}

func TestSearchUsers_EmptyQuery(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	users, err := userUsecase.SearchUsers("")

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, "search query is required", err.Error())
}

func TestSearchUsers_NoResults(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userUsecase := usecase.NewUserManagementUsecase(userRepo)

	userRepo.On("Search", "unknown").Return([]*domain.User{}, nil)

	users, err := userUsecase.SearchUsers("unknown")

	assert.NoError(t, err)
	assert.Empty(t, users)
	userRepo.AssertExpectations(t)
}
