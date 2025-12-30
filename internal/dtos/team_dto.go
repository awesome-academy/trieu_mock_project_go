package dtos

import "time"

type PaginationRequestQuery struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}

type ListTeamsResponse struct {
	Teams []Team             `json:"teams"`
	Page  PaginationResponse `json:"page"`
}

type Team struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Leader   UserSummary      `json:"leader"`
	Members  []UserSummary    `json:"members"`
	Projects []ProjectSummary `json:"projects"`
}

type UserSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TeamMemberSummary struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	JoinedAt time.Time `json:"joined_at"`
}

type ListTeamMembersResponse struct {
	Members []TeamMemberSummary `json:"members"`
	Page    PaginationResponse  `json:"page"`
}

type CreateOrUpdateTeamRequest struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description *string `json:"description" binding:"max=1000"`
	LeaderID    uint    `json:"leader_id" binding:"required"`
}

type AddMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
