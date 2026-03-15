package domain

import "time"

// ─── Enum ─────────────────────────────────────────────────────────────────────

type RoomType string
type MemberRole string

const (
	RoomTypePrivate RoomType = "private"
	RoomTypeGroup   RoomType = "group"

	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

// ─── Enum ─────────────────────────────────────────────────────────────────────

type Room struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Type      RoomType      `json:"type"`
	CreatedBy uint          `json:"created_by"`
	Members   []*RoomMember `json:"members"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type RoomMember struct {
	ID       uint       `json:"id"`
	RoomID   uint       `json:"room_id"`
	UserID   uint       `json:"user_id"`
	User     *User      `json:"user"`
	Role     MemberRole `json:"role"`
	JoinedAt time.Time  `json:"joinedAt"`
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type RoomRepository interface {
	Create(room *Room) error
	FindByID(id uint) (*Room, error)
	FindByUserID(userID uint) ([]*Room, error)
	Update(room *Room) error
	Delete(id uint) error
	AddMember(member *RoomMember) error
	RemoveMember(roomID uint, userID uint) error
	FindMember(roomID uint, userID uint) (*RoomMember, error)
}

// ─── Usecase Interface ─────────────────────────────────────────────────────

type RoomUsecase interface {
	CreateRoom(userID uint, req *CreateRoomRequest) (*Room, error)
	GetRooms(userID uint) ([]*Room, error)
	GetRoomByID(userID uint, roomID uint) (*Room, error)
	UpdateRoom(userID uint, roomID uint, req *UpdateRoomRequest) (*Room, error)
	DeleteRoom(userID uint, roomID uint) error
	AddMember(userID uint, roomID uint, req *AddMemberRequest) error
	RemoveMember(userID uint, roomID uint, targetUserID uint) error
}

// ─── Request dan Response ────────────────────────────────────────────────────

type CreateRoomRequest struct {
	Name    string   `json:"name" binding:"required_if=Type group,omitempty,min=3,max=50"`
	Type    RoomType `json:"type" binding:"required,oneof=private group"`
	Members []uint   `json:"members" binding:"required,min=1"`
}

type UpdateRoomRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

type AddMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
