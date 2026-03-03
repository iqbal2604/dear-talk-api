package usecase

import (
	"errors"
	"time"

	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
)

type roomUsecase struct {
	roomRepo domain.RoomRepository
	userRepo domain.UserRepository
}

func NewRoomUsecase(roomRepo domain.RoomRepository, userRepo domain.UserRepository) domain.RoomUsecase {
	return &roomUsecase{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

// ─── Create Room ──────────────────────────────────────────────────────────────

func (u *roomUsecase) CreateRoom(userID uint, req *domain.CreateRoomRequest) (*domain.Room, error) {
	// Validasi untuk private room hanya boleh 1 member tambahan
	if req.Type == domain.RoomTypePrivate && len(req.Members) != 1 {
		return nil, errors.New("private room must have exactly 1 other member")
	}

	room := &domain.Room{
		Name:      req.Name,
		Type:      req.Type,
		CreatedBy: userID,
	}

	// Untuk private room, nama diisi otomatis nanti dari username lawan
	if req.Type == domain.RoomTypePrivate {
		room.Name = "private"
	}

	if err := u.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// Tambah creator sebagai admin
	if err := u.roomRepo.AddMember(&domain.RoomMember{
		RoomID:   room.ID,
		UserID:   userID,
		Role:     domain.MemberRoleAdmin,
		JoinedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	// Tambah member lainnya
	for _, memberID := range req.Members {
		// Pastikan user ada
		user, err := u.userRepo.FindByID(memberID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}

		if err := u.roomRepo.AddMember(&domain.RoomMember{
			RoomID:   room.ID,
			UserID:   memberID,
			Role:     domain.MemberRoleMember,
			JoinedAt: time.Now(),
		}); err != nil {
			return nil, err
		}
	}

	// Ambil room lengkap dengan members
	return u.roomRepo.FindByID(room.ID)
}

// ─── Get Rooms ────────────────────────────────────────────────────────────────

func (u *roomUsecase) GetRooms(userID uint) ([]*domain.Room, error) {
	rooms, err := u.roomRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// ─── Get Room By ID ───────────────────────────────────────────────────────────

func (u *roomUsecase) GetRoomByID(userID uint, roomID uint) (*domain.Room, error) {
	// Pastikan user adalah member room ini
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("room not found")
	}

	room, err := u.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}

	return room, nil
}

// ─── Update Room ──────────────────────────────────────────────────────────────

func (u *roomUsecase) UpdateRoom(userID uint, roomID uint, req *domain.UpdateRoomRequest) (*domain.Room, error) {
	// Hanya admin yang boleh update
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("room not found")
	}
	if member.Role != domain.MemberRoleAdmin {
		return nil, errors.New("only admin can update room")
	}

	room, err := u.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	room.Name = req.Name

	if err := u.roomRepo.Update(room); err != nil {
		return nil, err
	}

	return u.roomRepo.FindByID(roomID)
}

// ─── Delete Room ──────────────────────────────────────────────────────────────

func (u *roomUsecase) DeleteRoom(userID uint, roomID uint) error {
	// Hanya admin yang boleh hapus
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("room not found")
	}
	if member.Role != domain.MemberRoleAdmin {
		return errors.New("only admin can delete room")
	}

	return u.roomRepo.Delete(roomID)
}

// ─── Add Member ───────────────────────────────────────────────────────────────

func (u *roomUsecase) AddMember(userID uint, roomID uint, req *domain.AddMemberRequest) error {
	// Hanya admin yang boleh tambah member
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("room not found")
	}
	if member.Role != domain.MemberRoleAdmin {
		return errors.New("only admin can add member")
	}

	// Cek apakah user yang akan ditambah ada
	user, err := u.userRepo.FindByID(req.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Cek apakah sudah menjadi member
	existing, err := u.roomRepo.FindMember(roomID, req.UserID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("user is already a member")
	}

	return u.roomRepo.AddMember(&domain.RoomMember{
		RoomID:   roomID,
		UserID:   req.UserID,
		Role:     domain.MemberRoleMember,
		JoinedAt: time.Now(),
	})
}

// ─── Remove Member ────────────────────────────────────────────────────────────

func (u *roomUsecase) RemoveMember(userID uint, roomID uint, targetUserID uint) error {
	// Hanya admin yang boleh hapus member
	member, err := u.roomRepo.FindMember(roomID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("room not found")
	}
	if member.Role != domain.MemberRoleAdmin {
		return errors.New("only admin can remove member")
	}

	// Admin tidak bisa hapus dirinya sendiri
	if targetUserID == userID {
		return errors.New("admin cannot remove themselves")
	}

	// Cek apakah target adalah member
	target, err := u.roomRepo.FindMember(roomID, targetUserID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("user is not a member")
	}

	return u.roomRepo.RemoveMember(roomID, targetUserID)
}
