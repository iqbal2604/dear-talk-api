package repository

import (
	"errors"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/repository/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// ─── Converter ────────────────────────────────────────────────────────────────

func toUserDomain(m *model.UserModel) *domain.User {
	return &domain.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Avatar:    m.Avatar,
		IsOnline:  m.IsOnline,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toUserModel(u *domain.User) *model.UserModel {
	return &model.UserModel{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Avatar:   u.Avatar,
		IsOnline: u.IsOnline,
	}
}

// ─── Implementations ──────────────────────────────────────────────────────────

func (r *userRepository) Create(user *domain.User) error {
	m := toUserModel(user)
	result := r.db.Create(m)
	if result.Error != nil {
		return result.Error
	}

	//Kembalikan ID yang di Generate ke Domain
	user.ID = m.ID
	return nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var m model.UserModel
	result := r.db.First(&m, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return toUserDomain(&m), nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var m model.UserModel
	result := r.db.Where("email = ?", email).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return toUserDomain(&m), nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var m model.UserModel
	result := r.db.Where("username = ?", username).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return toUserDomain(&m), nil
}

func (r *userRepository) Search(query string) ([]*domain.User, error) {
	var models []model.UserModel
	result := r.db.Where("username ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").Limit(20).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = toUserDomain(&m)
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	m := toUserModel(user)
	return r.db.Save(m).Error
}
