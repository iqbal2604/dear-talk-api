package usecase

import (
	"errors"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
)

type messageUsecase struct {
	messageRepo domain.MessageRepository
	roomRepo    domain.RoomRepository
}

func NewMessageUsecase(
	messageRepo domain.MessageRepository,
	roomRepo domain.RoomRepository,
) domain.MessageUsecase {
	return &messageUsecase{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

// ─── Send Message ─────────────────────────────────────────────────────────────

func (u *messageUsecase) SendMessage(userID uint, roomID uint, req *domain.SendMessageRequest) (*domain.Message, error) {
	// Pastikan user adalah member room
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("you are not a member of this room")
	}

	// Set default type jika kosong
	msgType := req.Type
	if msgType == "" {
		msgType = domain.MessageTypeText
	}

	message := &domain.Message{
		RoomID:   roomID,
		SenderID: userID,
		Content:  req.Content,
		Type:     msgType,
	}

	if err := u.messageRepo.Create(message); err != nil {
		return nil, err
	}

	// Ambil pesan lengkap dengan sender
	return u.messageRepo.FindByID(message.ID)
}

// ─── Get Messages ─────────────────────────────────────────────────────────────

func (u *messageUsecase) GetMessages(userID uint, roomID uint, page int, limit int) ([]*domain.Message, int64, error) {
	// Pastikan user adalah member room
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return nil, 0, err
	}
	if member == nil {
		return nil, 0, errors.New("you are not a member of this room")
	}

	// Default pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	return u.messageRepo.FindByRoomID(roomID, page, limit)
}

// ─── Edit Message

func (u *messageUsecase) EditMessage(userID uint, messageID uint, req *domain.EditMessageRequest) (*domain.Message, error) {
	// Cari pesan
	message, err := u.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, errors.New("message not found")
	}

	// Hanya pengirim yang boleh edit
	if message.SenderID != userID {
		return nil, errors.New("you can only edit your own message")
	}

	// Pastikan pesan belum dihapus
	if message.DeletedAt != nil {
		return nil, errors.New("message has been deleted")
	}

	message.Content = req.Content

	if err := u.messageRepo.Update(message); err != nil {
		return nil, err
	}

	return u.messageRepo.FindByID(messageID)
}

// ─── Delete Message ───────────────────────────────────────────────────────────

func (u *messageUsecase) DeleteMessage(userID uint, messageID uint) error {
	// Cari pesan
	message, err := u.messageRepo.FindByID(messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("message not found")
	}

	// Hanya pengirim yang boleh hapus
	if message.SenderID != userID {
		return errors.New("you can only delete your own message")
	}

	// Pastikan pesan belum dihapus
	if message.DeletedAt != nil {
		return errors.New("message already deleted")
	}

	return u.messageRepo.Delete(messageID)
}

// ─── Mark As Read ─────────────────────────────────────────────────────────────

func (u *messageUsecase) MarkAsRead(userID uint, roomID uint) error {
	// Pastikan user adalah member room
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("you are not a member of this room")
	}

	return u.messageRepo.MarkAsRead(roomID, userID)
}
