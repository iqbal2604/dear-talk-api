package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	UserID   uint
	Username string
	hub      *Hub
	conn     *websocket.Conn
	send     chan Event
	log      *zap.Logger
}

func NewClient(userID uint, username string, hub *Hub, conn *websocket.Conn, log *zap.Logger) *Client {
	return &Client{
		UserID:   userID,
		Username: username,
		hub:      hub,
		conn:     conn,
		send:     make(chan Event, 256),
		log:      log,
	}
}

// ─── Read Pump ────────────────────────────────────────────────────────────────
// Membaca pesan dari client (browser/app)

func (c *Client) ReadPump(onMessage func(client *Client, msg []byte)) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				c.log.Error("websocket read error",
					zap.Uint("user_id", c.UserID),
					zap.Error(err),
				)
			}
			break
		}
		// Teruskan pesan ke handler
		onMessage(c, message)
	}
}

// ─── Write Pump ───────────────────────────────────────────────────────────────
// Mengirim pesan dari server ke client

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub menutup channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(event)
			if err != nil {
				c.log.Error("failed to marshal event", zap.Error(err))
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.log.Error("websocket write error",
					zap.Uint("user_id", c.UserID),
					zap.Error(err),
				)
				return
			}

		case <-ticker.C:
			// Kirim ping secara berkala
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
