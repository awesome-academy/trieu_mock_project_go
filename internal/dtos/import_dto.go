package dtos

import "time"

type UserImportData struct {
	Name       string
	Email      string
	Birthday   *time.Time
	PositionID uint
	TeamID     *uint
	Skills     []SkillImportData
}

type SkillImportData struct {
	ID             uint
	Level          int
	UsedYearNumber int
}

type PositionImportData struct {
	Name         string
	Abbreviation string
}

type TeamImportData struct {
	Name        string
	Description *string
	LeaderID    uint
}

type ProjectImportData struct {
	Name         string
	Abbreviation string
	StartDate    *time.Time
	EndDate      *time.Time
	LeaderID     uint
	TeamID       uint
	MemberIDs    []uint
}
