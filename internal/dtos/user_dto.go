package dtos

import (
	"trieu_mock_project_go/types"
)

type UserProfile struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Birthday *types.Date `json:"birthday"`

	CurrentTeam *PositionSummary   `json:"current_team,omitempty"`
	Position    Position           `json:"position"`
	Projects    []ProjectSummary   `json:"projects"`
	Skills      []UserSkillSummary `json:"skills"`
}

type UserSearchRequest struct {
	TeamId *uint `form:"team_id"`
	Limit  int   `form:"limit" binding:"min=1,max=100"`
	Offset int   `form:"offset" binding:"min=0"`
}

type UserSearchResponse struct {
	Users []UserDataForSearch `json:"users"`
	Page  PaginationResponse  `json:"page"`
}

type UserDataForSearch struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateOrUpdateUserRequest struct {
	Name       string            `json:"name" binding:"required"`
	Email      string            `json:"email" binding:"required,email"`
	Birthday   *types.Date       `json:"birthday"`
	PositionID uint              `json:"position_id" binding:"required"`
	TeamID     *uint             `json:"team_id"`
	Skills     []UpdateUserSkill `json:"skills"`
}

type UpdateUserSkill struct {
	ID             uint `json:"id" binding:"required"`
	Level          int  `json:"level" binding:"required,min=1,max=10"`
	UsedYearNumber int  `json:"used_year_number" binding:"min=0,max=100"`
}
