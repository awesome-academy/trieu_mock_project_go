package dtos

import (
	"trieu_mock_project_go/types"

	"github.com/go-playground/validator/v10"
)

type ProjectSearchRequest struct {
	PaginationRequestQuery
	TeamID *uint `form:"team_id"`
}

type ProjectSearchResponse struct {
	Projects []ProjectListItem  `json:"projects"`
	Page     PaginationResponse `json:"page"`
}

type ProjectListItem struct {
	ID           uint        `json:"id"`
	Name         string      `json:"name"`
	Abbreviation string      `json:"abbreviation"`
	StartDate    *types.Date `json:"start_date"`
	EndDate      *types.Date `json:"end_date"`
	LeaderName   string      `json:"leader_name"`
	TeamName     string      `json:"team_name"`
}

type ProjectDetail struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Abbreviation string        `json:"abbreviation"`
	StartDate    *types.Date   `json:"start_date"`
	EndDate      *types.Date   `json:"end_date"`
	Leader       UserSummary   `json:"leader"`
	Team         TeamSummary   `json:"team"`
	Members      []UserSummary `json:"members"`
}

type CreateOrUpdateProjectRequest struct {
	Name         string      `json:"name" binding:"required,max=255"`
	Abbreviation string      `json:"abbreviation" binding:"required,max=50"`
	StartDate    *types.Date `json:"start_date" binding:"omitempty"`
	EndDate      *types.Date `json:"end_date" binding:"omitempty"`
	LeaderID     uint        `json:"leader_id" binding:"required"`
	TeamID       uint        `json:"team_id" binding:"required"`
	MemberIDs    []uint      `json:"member_ids"`
}

func ProjectRequestStructLevelValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(CreateOrUpdateProjectRequest)

	// If EndDate is provided, StartDate must also be provided
	if req.EndDate != nil && req.StartDate == nil {
		sl.ReportError(req.StartDate, "StartDate", "start_date", "required_with_end_date", "")
	}

	// If both are provided, EndDate must be after StartDate
	if req.StartDate != nil && req.EndDate != nil {
		if !req.EndDate.After(req.StartDate.Time) {
			sl.ReportError(req.EndDate, "EndDate", "end_date", "gt_start_date", "")
		}
	}
}
