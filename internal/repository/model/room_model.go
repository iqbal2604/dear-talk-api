package model

import "time"

type RoomModel struct {
	ID        uint              `gorm:"primaryKey"`
	Name      string            `gorm:"not null"`
	Type      string            `gorm:"not null"`
	CreatedBy uint              `gorm:"not null"`
	Members   []RoomMemberModel `gorm:"foreignKey:RoomID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RoomModel) TableName() string {
	return "rooms"
}

type RoomMemberModel struct {
	ID       uint      `gorm:"primaryKey"`
	RoomID   uint      `gorm:"not null;index"`
	UserID   uint      `gorm:"not null;index"`
	User     UserModel `gorm:"foreignKey:UserID"`
	Role     string    `gorm:"not null;default:'member'"`
	JoinedAt time.Time
}

func (RoomMemberModel) TableName() string {
	return "room_members"
}
