package repository

import (
	"errors"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/repository/model"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &roomRepository{db: db}
}

// ─── Converter ────────────────────────────────────────────────────────────────

func toRoomDomain(m *model.RoomModel) *domain.Room {
	room := &domain.Room{
		ID:        m.ID,
		Name:      m.Name,
		Type:      domain.RoomType(m.Type),
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	for _, member := range m.Members {
		room.Members = append(room.Members, toRoomMemberDomain(&member))
	}

	return room
}

func toRoomMemberDomain(m *model.RoomMemberModel) *domain.RoomMember {
	member := &domain.RoomMember{
		ID:       m.ID,
		RoomID:   m.RoomID,
		UserID:   m.UserID,
		Role:     domain.MemberRole(m.Role),
		JoinedAt: m.JoinedAt,
	}

	// Map user jika ada
	if m.User.ID != 0 {
		member.User = toUserDomain(&m.User)
	}

	return member
}

// ─── Implementations ──────────────────────────────────────────────────────────

func (r *roomRepository) Create(room *domain.Room) error {
	m := &model.RoomModel{
		Name:      room.Name,
		Type:      string(room.Type),
		CreatedBy: room.CreatedBy,
	}

	result := r.db.Create(m)
	if result.Error != nil {
		return result.Error
	}

	room.ID = m.ID
	return nil
}

func (r *roomRepository) FindByID(id uint) (*domain.Room, error) {
	var m model.RoomModel
	result := r.db.
		Preload("Members").
		Preload("Members.User").
		First(&m, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return toRoomDomain(&m), nil
}

func (r *roomRepository) FindByUserID(userID uint) ([]*domain.Room, error) {
	var members []model.RoomMemberModel
	result := r.db.Where("user_id = ?", userID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}

	// Ambil semua room ID yang dimiliki user
	roomIDs := make([]uint, len(members))
	for i, m := range members {
		roomIDs[i] = m.RoomID
	}

	if len(roomIDs) == 0 {
		return []*domain.Room{}, nil
	}

	// Ambil detail room beserta members
	var rooms []model.RoomModel
	result = r.db.
		Preload("Members").
		Preload("Members.User").
		Where("id IN ?", roomIDs).
		Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}

	domainRooms := make([]*domain.Room, len(rooms))
	for i, room := range rooms {
		domainRooms[i] = toRoomDomain(&room)
	}

	return domainRooms, nil
}

func (r *roomRepository) Update(room *domain.Room) error {
	return r.db.Model(&model.RoomModel{}).
		Where("id = ?", room.ID).
		Updates(map[string]interface{}{
			"name": room.Name,
		}).Error
}

func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&model.RoomModel{}, id).Error
}

func (r *roomRepository) AddMember(member *domain.RoomMember) error {
	m := &model.RoomMemberModel{
		RoomID:   member.RoomID,
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
	}
	result := r.db.Create(m)
	if result.Error != nil {
		return result.Error
	}
	member.ID = m.ID
	return nil
}

func (r *roomRepository) RemoveMember(roomID uint, userID uint) error {
	return r.db.
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.RoomMemberModel{}).Error
}

func (r *roomRepository) FindMember(roomID uint, userID uint) (*domain.RoomMember, error) {
	var m model.RoomMemberModel
	result := r.db.
		Where("room_id = ? AND user_id = ?", roomID, userID).
		First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return toRoomMemberDomain(&m), nil
}
