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

// ─── CreateRoom Tests ─────────────────────────────────────────────────────────

func TestCreateRoom_GroupSuccess(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.CreateRoomRequest{
		Name:    "General",
		Type:    domain.RoomTypeGroup,
		Members: []uint{2},
	}

	// Mock — create room berhasil
	roomRepo.On("Create", mock.AnythingOfType("*domain.Room")).Return(nil)
	// Mock — tambah creator sebagai admin
	roomRepo.On("AddMember", mock.AnythingOfType("*domain.RoomMember")).Return(nil)
	// Mock — cek user member ada
	userRepo.On("FindByID", uint(2)).Return(&domain.User{ID: 2, Username: "janedoe"}, nil)
	// Mock — ambil room lengkap
	roomRepo.On("FindByID", uint(0)).Return(&domain.Room{
		ID:        0,
		Name:      req.Name,
		Type:      req.Type,
		CreatedBy: 1,
		Members: []*domain.RoomMember{
			{UserID: 1, Role: domain.MemberRoleAdmin},
			{UserID: 2, Role: domain.MemberRoleMember},
		},
	}, nil)

	room, err := roomUsecase.CreateRoom(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, req.Name, room.Name)
	assert.Equal(t, domain.RoomTypeGroup, room.Type)
	assert.Len(t, room.Members, 2)
	roomRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCreateRoom_PrivateSuccess(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.CreateRoomRequest{
		Type:    domain.RoomTypePrivate,
		Members: []uint{2},
	}

	roomRepo.On("Create", mock.AnythingOfType("*domain.Room")).Return(nil)
	roomRepo.On("AddMember", mock.AnythingOfType("*domain.RoomMember")).Return(nil)
	userRepo.On("FindByID", uint(2)).Return(&domain.User{ID: 2, Username: "janedoe"}, nil)
	roomRepo.On("FindByID", uint(0)).Return(&domain.Room{
		ID:        0,
		Name:      "private",
		Type:      domain.RoomTypePrivate,
		CreatedBy: 1,
		Members: []*domain.RoomMember{
			{UserID: 1, Role: domain.MemberRoleAdmin},
			{UserID: 2, Role: domain.MemberRoleMember},
		},
	}, nil)

	room, err := roomUsecase.CreateRoom(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, domain.RoomTypePrivate, room.Type)
	roomRepo.AssertExpectations(t)
}

func TestCreateRoom_PrivateInvalidMembers(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	// Private room dengan lebih dari 1 member
	req := &domain.CreateRoomRequest{
		Type:    domain.RoomTypePrivate,
		Members: []uint{2, 3},
	}

	room, err := roomUsecase.CreateRoom(1, req)

	assert.Error(t, err)
	assert.Nil(t, room)
	assert.Equal(t, "private room must have exactly 1 other member", err.Error())
}

func TestCreateRoom_MemberNotFound(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.CreateRoomRequest{
		Name:    "General",
		Type:    domain.RoomTypeGroup,
		Members: []uint{99},
	}

	roomRepo.On("Create", mock.AnythingOfType("*domain.Room")).Return(nil)
	roomRepo.On("AddMember", mock.AnythingOfType("*domain.RoomMember")).Return(nil)
	// Mock — member tidak ditemukan
	userRepo.On("FindByID", uint(99)).Return(nil, nil)

	room, err := roomUsecase.CreateRoom(1, req)

	assert.Error(t, err)
	assert.Nil(t, room)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

// ─── GetRooms Tests ───────────────────────────────────────────────────────────

func TestGetRooms_Success(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	expectedRooms := []*domain.Room{
		{ID: 1, Name: "General", Type: domain.RoomTypeGroup},
		{ID: 2, Name: "private", Type: domain.RoomTypePrivate},
	}

	roomRepo.On("FindByUserID", uint(1)).Return(expectedRooms, nil)

	rooms, err := roomUsecase.GetRooms(1)

	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	roomRepo.AssertExpectations(t)
}

func TestGetRooms_Empty(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	roomRepo.On("FindByUserID", uint(1)).Return([]*domain.Room{}, nil)

	rooms, err := roomUsecase.GetRooms(1)

	assert.NoError(t, err)
	assert.Empty(t, rooms)
	roomRepo.AssertExpectations(t)
}

// ─── UpdateRoom Tests ─────────────────────────────────────────────────────────

func TestUpdateRoom_Success(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.UpdateRoomRequest{
		Name: "General Updated",
	}

	// Mock — user adalah admin
	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		RoomID: 1,
		Role:   domain.MemberRoleAdmin,
	}, nil)

	roomRepo.On("FindByID", uint(1)).Return(&domain.Room{
		ID:   1,
		Name: "General",
		Type: domain.RoomTypeGroup,
	}, nil)

	roomRepo.On("Update", mock.AnythingOfType("*domain.Room")).Return(nil)

	roomRepo.On("FindByID", uint(1)).Return(&domain.Room{
		ID:   1,
		Name: req.Name,
		Type: domain.RoomTypeGroup,
	}, nil)

	room, err := roomUsecase.UpdateRoom(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, req.Name, room.Name)
	roomRepo.AssertExpectations(t)
}

func TestUpdateRoom_NotAdmin(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.UpdateRoomRequest{
		Name: "General Updated",
	}

	// Mock — user bukan admin
	roomRepo.On("FindMember", uint(1), uint(2)).Return(&domain.RoomMember{
		UserID: 2,
		RoomID: 1,
		Role:   domain.MemberRoleMember,
	}, nil)

	room, err := roomUsecase.UpdateRoom(2, 1, req)

	assert.Error(t, err)
	assert.Nil(t, room)
	assert.Equal(t, "only admin can update room", err.Error())
	roomRepo.AssertExpectations(t)
}

// ─── AddMember Tests ──────────────────────────────────────────────────────────

func TestAddMember_Success(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.AddMemberRequest{
		UserID: 3,
	}

	// Mock — user adalah admin
	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		Role:   domain.MemberRoleAdmin,
	}, nil)

	// Mock — user yang ditambah ada
	userRepo.On("FindByID", uint(3)).Return(&domain.User{
		ID:       3,
		Username: "newmember",
	}, nil)

	// Mock — user belum jadi member
	roomRepo.On("FindMember", uint(1), uint(3)).Return(nil, nil)

	// Mock — tambah member berhasil
	roomRepo.On("AddMember", mock.MatchedBy(func(m *domain.RoomMember) bool {
		return m.UserID == 3 && m.RoomID == 1
	})).Return(nil)

	err := roomUsecase.AddMember(1, 1, req)

	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestAddMember_AlreadyMember(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	req := &domain.AddMemberRequest{
		UserID: 2,
	}

	// Mock — user adalah admin
	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		Role:   domain.MemberRoleAdmin,
	}, nil)

	// Mock — user yang ditambah ada
	userRepo.On("FindByID", uint(2)).Return(&domain.User{
		ID:       2,
		Username: "janedoe",
	}, nil)

	// Mock — user sudah jadi member
	roomRepo.On("FindMember", uint(1), uint(2)).Return(&domain.RoomMember{
		UserID: 2,
		Role:   domain.MemberRoleMember,
	}, nil)

	err := roomUsecase.AddMember(1, 1, req)

	assert.Error(t, err)
	assert.Equal(t, "user is already a member", err.Error())
	roomRepo.AssertExpectations(t)
}

// ─── RemoveMember Tests ───────────────────────────────────────────────────────

func TestRemoveMember_Success(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	// Mock — user adalah admin
	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		Role:   domain.MemberRoleAdmin,
	}, nil)

	// Mock — target adalah member
	roomRepo.On("FindMember", uint(1), uint(2)).Return(&domain.RoomMember{
		UserID: 2,
		Role:   domain.MemberRoleMember,
	}, nil)

	// Mock — hapus member berhasil
	roomRepo.On("RemoveMember", uint(1), uint(2)).Return(nil)

	err := roomUsecase.RemoveMember(1, 1, 2)

	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
}

func TestRemoveMember_AdminRemoveSelf(t *testing.T) {
	roomRepo := new(mocks.RoomRepository)
	userRepo := new(mocks.UserRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, userRepo)

	// Mock — user adalah admin
	roomRepo.On("FindMember", uint(1), uint(1)).Return(&domain.RoomMember{
		UserID: 1,
		Role:   domain.MemberRoleAdmin,
	}, nil)

	// Admin coba hapus dirinya sendiri
	err := roomUsecase.RemoveMember(1, 1, 1)

	assert.Error(t, err)
	assert.Equal(t, "admin cannot remove themselves", err.Error())
	roomRepo.AssertExpectations(t)
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func newRoomMember(roomID, userID uint, role domain.MemberRole) *domain.RoomMember {
	return &domain.RoomMember{
		RoomID:   roomID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}
