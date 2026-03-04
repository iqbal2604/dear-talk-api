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

// CreateRoom godoc
// @Summary      Buat room baru
// @Description  Membuat room baru bertipe private atau group
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      domain.CreateRoomRequest  true  "Create Room Request"
// @Success      201      {object}  response.Response{data=domain.Room}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /rooms [post]
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

// GetRooms godoc
// @Summary      Lihat semua room
// @Description  Mendapatkan semua room yang dimiliki user yang sedang login
// @Tags         Rooms
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=[]domain.Room}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /rooms [get]
func (h *RoomHandler) GetRooms(c *gin.Context) {
	userID := c.GetUint("user_id")

	rooms, err := h.roomUsecase.GetRooms(userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "rooms fetched", rooms)
}

// GetRoomByID godoc
// @Summary      Lihat detail room
// @Description  Mendapatkan detail room beserta semua member
// @Tags         Rooms
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Room ID"
// @Success      200  {object}  response.Response{data=domain.Room}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /rooms/{id} [get]
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

// UpdateRoom godoc
// @Summary      Update room
// @Description  Update nama room, hanya admin yang bisa melakukan ini
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                       true  "Room ID"
// @Param        request  body      domain.UpdateRoomRequest  true  "Update Room Request"
// @Success      200      {object}  response.Response{data=domain.Room}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /rooms/{id} [put]
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

// DeleteRoom godoc
// @Summary      Hapus room
// @Description  Menghapus room, hanya admin yang bisa melakukan ini
// @Tags         Rooms
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Room ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /rooms/{id} [delete]
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

// AddMember godoc
// @Summary      Tambah member
// @Description  Menambah member baru ke room, hanya admin yang bisa melakukan ini
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                      true  "Room ID"
// @Param        request  body      domain.AddMemberRequest  true  "Add Member Request"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /rooms/{id}/members [post]
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

// RemoveMember godoc
// @Summary      Hapus member
// @Description  Menghapus member dari room, hanya admin yang bisa melakukan ini
// @Tags         Rooms
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int  true  "Room ID"
// @Param        userId  path      int  true  "User ID yang akan dihapus"
// @Success      200     {object}  response.Response
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Router       /rooms/{id}/members/{userId} [delete]
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
