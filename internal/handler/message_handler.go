package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/response"
)

type MessageHandler struct {
	messageUsecase domain.MessageUsecase
}

func NewMessageHandler(messageUsecase domain.MessageUsecase) *MessageHandler {
	return &MessageHandler{messageUsecase: messageUsecase}
}

// SendMessage godoc
// @Summary      Kirim pesan
// @Description  Mengirim pesan baru ke sebuah room
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "Room ID"
// @Param        request  body      domain.SendMessageRequest  true  "Send Message Request"
// @Success      201      {object}  response.Response{data=domain.Message}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /rooms/{id}/messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	var req domain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	message, err := h.messageUsecase.SendMessage(userID, uint(roomID), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "message sent", message)
}

// GetMessages godoc
// @Summary      Ambil riwayat pesan
// @Description  Mendapatkan riwayat pesan dalam sebuah room dengan pagination
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int  true   "Room ID"
// @Param        page   query     int  false  "Page number (default: 1)"
// @Param        limit  query     int  false  "Items per page (default: 20, max: 50)"
// @Success      200    {object}  response.Response{data=[]domain.Message}
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /rooms/{id}/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	// Ambil query pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	messages, total, err := h.messageUsecase.GetMessages(userID, uint(roomID), page, limit)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Hitung total page
	totalPage := int(total) / limit
	if int(total)%limit != 0 {
		totalPage++
	}

	response.OKWithMeta(c, "messages fetched", messages, &response.Meta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	})
}

// EditMessage godoc
// @Summary      Edit pesan
// @Description  Mengedit isi pesan, hanya pengirim yang bisa melakukan ini
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                       true  "Message ID"
// @Param        request  body      domain.EditMessageRequest  true  "Edit Message Request"
// @Success      200      {object}  response.Response{data=domain.Message}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /messages/{id} [put]
func (h *MessageHandler) EditMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid message id", nil)
		return
	}

	var req domain.EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err.Error())
		return
	}

	message, err := h.messageUsecase.EditMessage(userID, uint(messageID), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "message updated", message)
}

// DeleteMessage godoc
// @Summary      Hapus pesan
// @Description  Menghapus pesan secara soft delete, hanya pengirim yang bisa melakukan ini
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Message ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid message id", nil)
		return
	}

	if err := h.messageUsecase.DeleteMessage(userID, uint(messageID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "message deleted", nil)
}

// MarkAsRead godoc
// @Summary      Tandai pesan sudah dibaca
// @Description  Menandai semua pesan dalam room sudah dibaca oleh user
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Room ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /rooms/{id}/read [post]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid room id", nil)
		return
	}

	if err := h.messageUsecase.MarkAsRead(userID, uint(roomID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "messages marked as read", nil)
}
