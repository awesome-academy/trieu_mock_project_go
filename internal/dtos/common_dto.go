package dtos

import "time"

type PaginationResponse struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type PositionSummary struct {
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

type SkillSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserSkillSummary struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Level          int    `json:"level"`
	UsedYearNumber int    `json:"used_year_number"`
}
