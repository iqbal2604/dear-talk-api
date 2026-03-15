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
	ID        uint        `json:"id"`
	RoomID    uint        `json:"room_id"`
	SenderID  uint        `json:"sender_id"`
	Sender    *User       `json:"sender"`
	Content   string      `json:"content"`
	Type      MessageType `json:"type"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt *time.Time  `json:"deletedAt"`
}

type ReadStatus struct {
	ID         uint      `json:"id"`
	RoomID     uint      `json:"room_id"`
	UserID     uint      `json:"user_id"`
	LastReadAt time.Time `json:"lastreadAt"`
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
