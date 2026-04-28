package realtime

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type           string `json:"type"`
	ConversationID int64  `json:"conversation_id"`
	UserID         int64  `json:"user_id,omitempty"`
	MessageID      int64  `json:"message_id,omitempty"`
	IsTyping       *bool  `json:"is_typing,omitempty"`
	Payload        any    `json:"payload,omitempty"`
}

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan Event
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]bool),
	}
}

func (h *Hub) AddClient(userID int64, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]bool)
	}
	h.clients[userID][c] = true
}

func (h *Hub) RemoveClient(userID int64, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] != nil {
		delete(h.clients[userID], c)
		close(c.Send)

		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
	}
}

func (h *Hub) SendToUser(userID int64, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients[userID] {
		select {
		case c.Send <- event:
		default:
		}
	}
}

func (h *Hub) BroadcastToUsers(userIDs []int64, event Event) {
	for _, userID := range userIDs {
		h.SendToUser(userID, event)
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for event := range c.Send {
		if err := c.Conn.WriteJSON(event); err != nil {
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, userID int64) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Hub:    hub,
		Send:   make(chan Event, 32),
	}

	hub.AddClient(userID, client)

	go client.WritePump()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			hub.RemoveClient(userID, client)
			_ = conn.Close()
			return
		}
	}
}

func ParseUserIDFromQuery(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
}

func MustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
