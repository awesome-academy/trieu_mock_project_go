package dtos

import "time"

type UserProfile struct {
	ID       uint       `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Birthday *time.Time `json:"birthday"`

	CurrentTeam *TeamSummary       `json:"current_team,omitempty"`
	Position    Position           `json:"position"`
	Projects    []ProjectSummary   `json:"projects"`
	Skills      []UserSkillSummary `json:"skills"`
}

type TeamSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Position struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type ProjectSummary struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Abbreviation string     `json:"abbreviation"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
}

type UserSkillSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
