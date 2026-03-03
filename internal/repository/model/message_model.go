package model

import "time"

type MessageModel struct {
	ID        uint      `gorm:"primaryKey"`
	RoomID    uint      `gorm:"not null;index"`
	SenderID  uint      `gorm:"not null;index"`
	Sender    UserModel `gorm:"foreignKey:SenderID"`
	Content   string    `gorm:"not null"`
	Type      string    `gorm:"not null;default:'text'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (MessageModel) TableName() string {
	return "messages"
}

type ReadStatusModel struct {
	ID         uint      `gorm:"primaryKey"`
	RoomID     uint      `gorm:"not null;index"`
	UserID     uint      `gorm:"not null;index"`
	LastReadAt time.Time `gorm:"not null"`
}

func (ReadStatusModel) TableName() string {
	return "read_statuses"
}
