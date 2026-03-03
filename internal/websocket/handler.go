package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/iqbal2604/dear-talk-api.git/internal/domain"
	"github.com/iqbal2604/dear-talk-api.git/pkg/jwt"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow semua origin untuk development
	// Di production ganti dengan cek origin yang spesifik
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub         *Hub
	jwtUtil     *jwt.JWTUtil
	roomRepo    domain.RoomRepository
	messageRepo domain.MessageRepository
	log         *zap.Logger
}

func NewWSHandler(
	hub *Hub,
	jwtUtil *jwt.JWTUtil,
	roomRepo domain.RoomRepository,
	messageRepo domain.MessageRepository,
	log *zap.Logger,
) *WSHandler {
	return &WSHandler{
		hub:         hub,
		jwtUtil:     jwtUtil,
		roomRepo:    roomRepo,
		messageRepo: messageRepo,
		log:         log,
	}
}

// ─── Upgrade HTTP ke WebSocket ────────────────────────────────────────────────

func (h *WSHandler) ServeWS(c *gin.Context) {
	// Ambil token dari query param
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	// Validasi token
	claims, err := h.jwtUtil.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Upgrade koneksi ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("failed to upgrade websocket", zap.Error(err))
		return
	}

	// Buat client baru
	client := NewClient(claims.UserID, claims.Username, h.hub, conn, h.log)

	// Register client ke hub
	h.hub.Register(client)

	h.log.Info("websocket client connected",
		zap.Uint("user_id", client.UserID),
		zap.String("username", client.Username),
	)

	// Broadcast user online ke semua client
	h.hub.BroadcastToRoom(h.getOnlineUserIDs(), Event{
		Type: EventUserOnline,
		Payload: gin.H{
			"user_id":  client.UserID,
			"username": client.Username,
		},
	})

	// Jalankan write pump di goroutine terpisah
	go client.WritePump()

	// Jalankan read pump (blocking)
	client.ReadPump(h.handleMessage)

	// Kalau sampai sini berarti client disconnect
	h.log.Info("websocket client disconnected",
		zap.Uint("user_id", client.UserID),
		zap.String("username", client.Username),
	)

	// Broadcast user offline
	h.hub.BroadcastToRoom(h.getOnlineUserIDs(), Event{
		Type: EventUserOffline,
		Payload: gin.H{
			"user_id":  client.UserID,
			"username": client.Username,
		},
	})
}

// ─── Handle Incoming Message dari Client ─────────────────────────────────────

type IncomingEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	RoomID  uint   `json:"room_id"`
	Content string `json:"content"`
}

func (h *WSHandler) handleMessage(client *Client, msg []byte) {
	var event IncomingEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		h.log.Error("failed to unmarshal event", zap.Error(err))
		return
	}

	switch event.Type {

	// Client kirim pesan baru
	case EventNewMessage:
		var payload SendMessagePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			h.log.Error("failed to unmarshal send message payload", zap.Error(err))
			return
		}
		h.handleSendMessage(client, payload)

	// Client sedang mengetik
	case EventTyping:
		var payload TypingPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			h.log.Error("failed to unmarshal typing payload", zap.Error(err))
			return
		}
		h.handleTyping(client, payload)
	}
}

// ─── Handle Send Message ──────────────────────────────────────────────────────

func (h *WSHandler) handleSendMessage(client *Client, payload SendMessagePayload) {
	// Pastikan user adalah member room
	member, err := h.roomRepo.FindMember(payload.RoomID, client.UserID)
	if err != nil || member == nil {
		return
	}

	// Simpan pesan ke DB
	message := &domain.Message{
		RoomID:   payload.RoomID,
		SenderID: client.UserID,
		Content:  payload.Content,
		Type:     domain.MessageTypeText,
	}

	if err := h.messageRepo.Create(message); err != nil {
		h.log.Error("failed to save message", zap.Error(err))
		return
	}

	// Ambil pesan lengkap dengan sender
	savedMessage, err := h.messageRepo.FindByID(message.ID)
	if err != nil {
		h.log.Error("failed to fetch message", zap.Error(err))
		return
	}

	// Ambil semua member room untuk broadcast
	room, err := h.roomRepo.FindByID(payload.RoomID)
	if err != nil || room == nil {
		return
	}

	memberIDs := make([]uint, len(room.Members))
	for i, m := range room.Members {
		memberIDs[i] = m.UserID
	}

	// Broadcast pesan ke semua member room
	h.hub.BroadcastToRoom(memberIDs, Event{
		Type:    EventNewMessage,
		Payload: savedMessage,
	})
}

// ─── Handle Typing ────────────────────────────────────────────────────────────

func (h *WSHandler) handleTyping(client *Client, payload TypingPayload) {
	// Pastikan user adalah member room
	member, err := h.roomRepo.FindMember(payload.RoomID, client.UserID)
	if err != nil || member == nil {
		return
	}

	// Ambil semua member room
	room, err := h.roomRepo.FindByID(payload.RoomID)
	if err != nil || room == nil {
		return
	}

	memberIDs := make([]uint, len(room.Members))
	for i, m := range room.Members {
		memberIDs[i] = m.UserID
	}

	// Broadcast typing ke semua member kecuali pengirim
	for _, memberID := range memberIDs {
		if memberID != client.UserID {
			h.hub.SendToUser(memberID, Event{
				Type: EventTyping,
				Payload: TypingPayload{
					RoomID:   payload.RoomID,
					UserID:   client.UserID,
					Username: client.Username,
					IsTyping: payload.IsTyping,
				},
			})
		}
	}
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func (h *WSHandler) getOnlineUserIDs() []uint {
	h.hub.mu.RLock()
	defer h.hub.mu.RUnlock()

	ids := make([]uint, 0, len(h.hub.clients))
	for id := range h.hub.clients {
		ids = append(ids, id)
	}
	return ids
}
