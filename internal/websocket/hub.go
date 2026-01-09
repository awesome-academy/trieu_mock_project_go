package websocket

import (
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[uint]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c.UserID]; !ok {
		h.clients[c.UserID] = make(map[*Client]bool)
	}
	h.clients[c.UserID][c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[c.UserID]; ok {
		if _, exists := clients[c]; exists {
			delete(clients, c)
			close(c.Send)
		}
		if len(clients) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

func (h *Hub) SendNotification(userID uint, notificationMsg *NotificationMessage) {
	h.mu.RLock()
	clients := h.clients[userID]
	h.mu.RUnlock()

	for c := range clients {
		select {
		case c.Send <- notificationMsg:
		default:
			go h.Unregister(c)
		}
	}
}
