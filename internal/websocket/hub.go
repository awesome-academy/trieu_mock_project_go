package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

const RedisNotificationChannel = "user_notifications"

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
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

func (h *Hub) SendNotification(userID uint, msg *NotificationMessage) {
	h.mu.RLock()
	clientsMap := h.clients[userID]

	clients := make([]*Client, 0, len(clientsMap))
	for c := range clientsMap {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered in SendNotification: %v", r)
				}
			}()
			select {
			case c.Send <- msg:
			default:
				go c.Close(h)
			}
		}()
	}
}

func (h *Hub) SubscribeToRedis(ctx context.Context, redisClient *redis.Client, channel string) {
	pubsub := redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			var notification NotificationMessage
			if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
				fmt.Printf("Error unmarshaling redis notification: %v\n", err)
				continue
			}
			h.SendNotification(notification.UserID, &notification)
		}
	}
}
