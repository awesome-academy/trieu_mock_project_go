package dtos

import "time"

type ActivityLogSummary struct {
	ID          uint      `json:"id"`
	Action      string    `json:"action"`
	UserID      uint      `json:"user_id"`
	UserName    string    `json:"user_name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ActivityLogSearchResponse struct {
	Logs []ActivityLogSummary `json:"logs"`
	Page PaginationResponse   `json:"page"`
}
