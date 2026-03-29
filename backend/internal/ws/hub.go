package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message is the standard unit of real-time communication on Skillora.
type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// Client represents a single active user connection.
type Client struct {
	UserID string
	Conn   *websocket.Conn
	mu     sync.Mutex
}

// Hub maintains the set of active clients and handles broadcasting.
type Hub struct {
	clients map[string]*Client // key: UserID
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Register adds a new client to the hub.
func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// If client already exists, close the old connection first.
	if old, exists := h.clients[userID]; exists {
		old.Conn.Close()
	}

	h.clients[userID] = &Client{
		UserID: userID,
		Conn:   conn,
	}
	log.Printf("[ws] user %s connected", userID)
}

// Unregister removes a client.
func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c, exists := h.clients[userID]; exists {
		c.Conn.Close()
		delete(h.clients, userID)
		log.Printf("[ws] user %s disconnected", userID)
	}
}

// BroadcastToUser sends a message to a specific user if they are online.
func (h *Hub) BroadcastToUser(userID string, msg Message) {
	h.mu.RLock()
	client, exists := h.clients[userID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	data, _ := json.Marshal(msg)
	if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[ws] send error to %s: %v", userID, err)
		h.Unregister(userID)
	}
}
