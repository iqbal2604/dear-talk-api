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

// ─── Send Message ─────────────────────────────────────────────────────────────

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

// ─── Get Messages ─────────────────────────────────────────────────────────────

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

// ─── Edit Message ─────────────────────────────────────────────────────────────

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

// ─── Delete Message ───────────────────────────────────────────────────────────

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

// ─── Mark As Read ─────────────────────────────────────────────────────────────

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
