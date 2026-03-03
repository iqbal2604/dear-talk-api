package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type RoomHandler struct {
	roomUsecase domain.RoomUsecase
}

func NewRoomHandler(roomUsecase domain.RoomUsecase) *RoomHandler {
	return &RoomHandler{roomUsecase: roomUsecase}
}

// ─── Create Room ──────────────────────────────────────────────────────────────

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req domain.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	room, err := h.roomUsecase.CreateRoom(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "room created", room)
}

// ─── Get Rooms ────────────────────────────────────────────────────────────────

func (h *RoomHandler) GetRooms(c *gin.Context) {
	userID := c.GetUint("user_id")

	rooms, err := h.roomUsecase.GetRooms(userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "rooms fetched", rooms)
}

// ─── Get Room By ID ───────────────────────────────────────────────────────────

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	room, err := h.roomUsecase.GetRoomByID(userID, uint(roomID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "room fetched", room)
}

// ─── Update Room ──────────────────────────────────────────────────────────────

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	var req domain.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	room, err := h.roomUsecase.UpdateRoom(userID, uint(roomID), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "room updated", room)
}

// ─── Delete Room ──────────────────────────────────────────────────────────────

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	if err := h.roomUsecase.DeleteRoom(userID, uint(roomID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "room deleted", nil)
}

// ─── Add Member ───────────────────────────────────────────────────────────────

func (h *RoomHandler) AddMember(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	var req domain.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	if err := h.roomUsecase.AddMember(userID, uint(roomID), &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "member added", nil)
}

// ─── Remove Member ────────────────────────────────────────────────────────────

func (h *RoomHandler) RemoveMember(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id", nil)
		return
	}

	if err := h.roomUsecase.RemoveMember(userID, uint(roomID), uint(targetUserID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "member removed", nil)
}
