package test

import (
	"testing"
	"time"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/internal/mocks"
	"github.com/iqbal2604/dear-talk-api.git/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── SendMessage Tests ────────────────────────────────────────────────────────

func TestSendMessage_Success(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	req := &domain.SendMessageRequest{
		Content: "Hello everyone!",
		Type:    domain.MessageTypeText,
	}

	// Mock — user adalah member room
	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		RoomID: 1,
		Role:   domain.MemberRoleMember,
	}, nil)

	// Mock — create message berhasil
	messageRepo.On("Create", mock.AnythingOfType("*domain.Message")).Return(nil)

	// Mock — ambil message lengkap
	messageRepo.On("FindByID", uint(0)).Return(&domain.Message{
		ID:       0,
		RoomID:   1,
		SenderID: 1,
		Content:  req.Content,
		Type:     domain.MessageTypeText,
		Sender: &domain.User{
			ID:       1,
			Username: "johndoe",
		},
	}, nil)

	message, err := messageUsecase.SendMessage(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, req.Content, message.Content)
	assert.Equal(t, uint(1), message.SenderID)
	assert.NotNil(t, message.Sender)
	roomRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestSendMessage_NotMember(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	req := &domain.SendMessageRequest{
		Content: "Hello everyone!",
		Type:    domain.MessageTypeText,
	}

	// Mock — user bukan member room
	roomRepo.On("FindMember", uint(1), uint(99)).Return(nil, nil)

	message, err := messageUsecase.SendMessage(99, 1, req)

	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Equal(t, "you are not a member of this room", err.Error())
	roomRepo.AssertExpectations(t)
}

func TestSendMessage_DefaultTypeText(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	// Tidak set type — default ke text
	req := &domain.SendMessageRequest{
		Content: "Hello!",
	}

	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		RoomID: 1,
	}, nil)

	messageRepo.On("Create", mock.MatchedBy(func(m *domain.Message) bool {
		// Pastikan type default ke text
		return m.Type == domain.MessageTypeText
	})).Return(nil)

	messageRepo.On("FindByID", uint(0)).Return(&domain.Message{
		ID:      0,
		Content: req.Content,
		Type:    domain.MessageTypeText,
	}, nil)

	message, err := messageUsecase.SendMessage(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, domain.MessageTypeText, message.Type)
	roomRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

// ─── GetMessages Tests ────────────────────────────────────────────────────────

func TestGetMessages_Success(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	expectedMessages := []*domain.Message{
		{ID: 1, Content: "Hello!", SenderID: 1},
		{ID: 2, Content: "Hi there!", SenderID: 2},
	}

	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		RoomID: 1,
	}, nil)

	messageRepo.On("FindByRoomID", uint(1), 1, 20).
		Return(expectedMessages, int64(2), nil)

	messages, total, err := messageUsecase.GetMessages(1, 1, 1, 20)

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, int64(2), total)
	roomRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestGetMessages_DefaultPagination(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		RoomID: 1,
	}, nil)

	// Kalau page dan limit 0 — default ke page 1 limit 20
	messageRepo.On("FindByRoomID", uint(1), 1, 20).
		Return([]*domain.Message{}, int64(0), nil)

	messages, total, err := messageUsecase.GetMessages(1, 1, 0, 0)

	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.Equal(t, int64(0), total)
	roomRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestGetMessages_NotMember(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	roomRepo.On("FindMember", uint(1), uint(99)).Return(nil, nil)

	messages, total, err := messageUsecase.GetMessages(99, 1, 1, 20)

	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, "you are not a member of this room", err.Error())
	roomRepo.AssertExpectations(t)
}

// ─── EditMessage Tests ────────────────────────────────────────────────────────

func TestEditMessage_Success(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	req := &domain.EditMessageRequest{
		Content: "Hello everyone! (edited)",
	}

	// Mock — pesan ditemukan milik user
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:        1,
		SenderID:  1,
		Content:   "Hello everyone!",
		DeletedAt: nil,
	}, nil).Once()

	// Mock — update berhasil
	messageRepo.On("Update", mock.AnythingOfType("*domain.Message")).Return(nil)

	// Mock — ambil pesan terupdate
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:       1,
		SenderID: 1,
		Content:  req.Content,
	}, nil).Once()

	message, err := messageUsecase.EditMessage(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, req.Content, message.Content)
	messageRepo.AssertExpectations(t)
}

func TestEditMessage_NotOwner(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	req := &domain.EditMessageRequest{
		Content: "Hello everyone! (edited)",
	}

	// Mock — pesan milik user lain
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:       1,
		SenderID: 2, // bukan user 1
		Content:  "Hello everyone!",
	}, nil)

	message, err := messageUsecase.EditMessage(1, 1, req)

	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Equal(t, "you can only edit your own message", err.Error())
	messageRepo.AssertExpectations(t)
}

func TestEditMessage_AlreadyDeleted(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	req := &domain.EditMessageRequest{
		Content: "edited",
	}

	deletedAt := time.Now()

	// Mock — pesan sudah dihapus
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:        1,
		SenderID:  1,
		DeletedAt: &deletedAt,
	}, nil)

	message, err := messageUsecase.EditMessage(1, 1, req)

	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Equal(t, "message has been deleted", err.Error())
	messageRepo.AssertExpectations(t)
}

// ─── DeleteMessage Tests ──────────────────────────────────────────────────────

func TestDeleteMessage_Success(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	// Mock — pesan ditemukan milik user
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:        1,
		SenderID:  1,
		DeletedAt: nil,
	}, nil)

	// Mock — delete berhasil
	messageRepo.On("Delete", uint(1)).Return(nil)

	err := messageUsecase.DeleteMessage(1, 1)

	assert.NoError(t, err)
	messageRepo.AssertExpectations(t)
}

func TestDeleteMessage_NotOwner(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	// Mock — pesan milik user lain
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:       1,
		SenderID: 2,
	}, nil)

	err := messageUsecase.DeleteMessage(1, 1)

	assert.Error(t, err)
	assert.Equal(t, "you can only delete your own message", err.Error())
	messageRepo.AssertExpectations(t)
}

func TestDeleteMessage_AlreadyDeleted(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	deletedAt := time.Now()

	// Mock — pesan sudah dihapus
	messageRepo.On("FindByID", uint(1)).Return(&domain.Message{
		ID:        1,
		SenderID:  1,
		DeletedAt: &deletedAt,
	}, nil)

	err := messageUsecase.DeleteMessage(1, 1)

	assert.Error(t, err)
	assert.Equal(t, "message already deleted", err.Error())
	messageRepo.AssertExpectations(t)
}

// ─── MarkAsRead Tests ─────────────────────────────────────────────────────────

func TestMarkAsRead_Success(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		RoomID: 1,
	}, nil)

	messageRepo.On("MarkAsRead", uint(1), uint(1)).Return(nil)

	err := messageUsecase.MarkAsRead(1, 1)

	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestMarkAsRead_NotMember(t *testing.T) {
	messageRepo := new(mocks.MessageRepository)
	roomRepo := new(mocks.RoomRepository)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, roomRepo)

	roomRepo.On("FindMember", uint(1), uint(99)).Return(nil, nil)

	err := messageUsecase.MarkAsRead(99, 1)

	assert.Error(t, err)
	assert.Equal(t, "you are not a member of this room", err.Error())
	roomRepo.AssertExpectations(t)
}
