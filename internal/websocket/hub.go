package websocket

import "sync"

// Event types
const (
	EventNewMessage  = "new_message"
	EventUserOnline  = "user_online"
	EventUserOffline = "user_offline"
	EventTyping      = "typing"
)

// Event yang dikirim ke client
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Typing payload
type TypingPayload struct {
	RoomID   uint   `json:"room_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IsTyping bool   `json:"is_typing"`
}

// Hub mengelola semua koneksi aktif
type Hub struct {
	mu      sync.RWMutex
	clients map[uint]*Client // userID -> Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]*Client),
	}
}

// Register client baru
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.UserID] = client
}

// Unregister client yang disconnect
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client.UserID]; ok {
		delete(h.clients, client.UserID)
		close(client.send)
	}
}

// Kirim event ke user tertentu
func (h *Hub) SendToUser(userID uint, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[userID]; ok {
		select {
		case client.send <- event:
		default:
			// Channel penuh, skip
		}
	}
}

// Broadcast event ke semua member dalam room
func (h *Hub) BroadcastToRoom(memberIDs []uint, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userID := range memberIDs {
		if client, ok := h.clients[userID]; ok {
			select {
			case client.send <- event:
			default:
				// Channel penuh, skip
			}
		}
	}
}

// Cek apakah user sedang online
func (h *Hub) IsOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
