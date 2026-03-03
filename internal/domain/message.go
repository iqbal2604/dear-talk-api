package domain

import "time"

// ─── Enum ─────────────────────────────────────────────────────────────────────

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

// ─── Entity ───────────────────────────────────────────────────────────────────

type Message struct {
	ID        uint
	RoomID    uint
	SenderID  uint
	Sender    *User
	Content   string
	Type      MessageType
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ReadStatus struct {
	ID         uint
	RoomID     uint
	UserID     uint
	LastReadAt time.Time
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type MessageRepository interface {
	Create(message *Message) error
	FindByID(id uint) (*Message, error)
	FindByRoomID(roomID uint, page int, limit int) ([]*Message, int64, error)
	Update(message *Message) error
	Delete(id uint) error
	MarkAsRead(roomID uint, userID uint) error
}

// ─── Usecase Interface ────────────────────────────────────────────────────────

type MessageUsecase interface {
	SendMessage(userID uint, roomID uint, req *SendMessageRequest) (*Message, error)
	GetMessages(userID uint, roomID uint, page int, limit int) ([]*Message, int64, error)
	EditMessage(userID uint, messageID uint, req *EditMessageRequest) (*Message, error)
	DeleteMessage(userID uint, messageID uint) error
	MarkAsRead(userID uint, roomID uint) error
}

// ─── Request & Response ───────────────────────────────────────────────────────

type SendMessageRequest struct {
	Content string      `json:"content" binding:"required,min=1"`
	Type    MessageType `json:"type" binding:"omitempty,oneof=text image file"`
}

type EditMessageRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}
