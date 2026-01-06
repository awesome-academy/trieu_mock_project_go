package dtos

import "time"

type NotificationSearchRequest struct {
	Limit  int `json:"limit" form:"limit,min=1,max=100"`
	Offset int `json:"offset" form:"offset,min=0"`
}

type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
}
