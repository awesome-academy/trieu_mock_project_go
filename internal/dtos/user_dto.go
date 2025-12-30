package dtos

import (
	"trieu_mock_project_go/types"
)

type UserProfile struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Birthday *types.Date `json:"birthday"`

	CurrentTeam *TeamSummary       `json:"current_team,omitempty"`
	Position    Position           `json:"position"`
	Projects    []ProjectSummary   `json:"projects"`
	Skills      []UserSkillSummary `json:"skills"`
}
