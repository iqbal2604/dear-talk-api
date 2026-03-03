package repository

import (
	"errors"
	"time"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/repository/model"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db: db}
}

// ─── Converter ────────────────────────────────────────────────────────────────

func toMessageDomain(m *model.MessageModel) *domain.Message {
	msg := &domain.Message{
		ID:        m.ID,
		RoomID:    m.RoomID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		Type:      domain.MessageType(m.Type),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}

	// Map sender jika ada
	if m.Sender.ID != 0 {
		msg.Sender = toUserDomain(&m.Sender)
	}

	return msg
}

// ─── Implementations ──────────────────────────────────────────────────────────

func (r *messageRepository) Create(message *domain.Message) error {
	m := &model.MessageModel{
		RoomID:   message.RoomID,
		SenderID: message.SenderID,
		Content:  message.Content,
		Type:     string(message.Type),
	}

	result := r.db.Create(m)
	if result.Error != nil {
		return result.Error
	}

	message.ID = m.ID
	message.CreatedAt = m.CreatedAt
	return nil
}

func (r *messageRepository) FindByID(id uint) (*domain.Message, error) {
	var m model.MessageModel
	result := r.db.
		Preload("Sender").
		First(&m, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return toMessageDomain(&m), nil
}

func (r *messageRepository) FindByRoomID(roomID uint, page int, limit int) ([]*domain.Message, int64, error) {
	var messages []model.MessageModel
	var total int64

	offset := (page - 1) * limit

	// Hitung total pesan
	r.db.Model(&model.MessageModel{}).
		Where("room_id = ? AND deleted_at IS NULL", roomID).
		Count(&total)

	// Ambil pesan dengan pagination, diurutkan terbaru
	result := r.db.
		Preload("Sender").
		Where("room_id = ? AND deleted_at IS NULL", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	domainMessages := make([]*domain.Message, len(messages))
	for i, m := range messages {
		domainMessages[i] = toMessageDomain(&m)
	}

	return domainMessages, total, nil
}

func (r *messageRepository) Update(message *domain.Message) error {
	return r.db.Model(&model.MessageModel{}).
		Where("id = ?", message.ID).
		Updates(map[string]interface{}{
			"content": message.Content,
		}).Error
}

func (r *messageRepository) Delete(id uint) error {
	now := time.Now()
	return r.db.Model(&model.MessageModel{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

func (r *messageRepository) MarkAsRead(roomID uint, userID uint) error {
	// Upsert — update jika ada, insert jika belum ada
	return r.db.
		Where(model.ReadStatusModel{RoomID: roomID, UserID: userID}).
		Assign(model.ReadStatusModel{LastReadAt: time.Now()}).
		FirstOrCreate(&model.ReadStatusModel{}).Error
}
