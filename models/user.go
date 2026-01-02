package models

import (
	"time"
)

type User struct {
	ID            uint       `gorm:"column:id;primaryKey;type:int unsigned"`
	Name          string     `gorm:"column:name;type:varchar(255);not null"`
	Email         string     `gorm:"column:email;type:varchar(255);not null;uniqueIndex:idx_email"`
	Password      string     `gorm:"column:password;type:varchar(255);not null"`
	Birthday      *time.Time `gorm:"column:birthday;type:date"`
	CurrentTeamID *uint      `gorm:"column:current_team_id;type:int unsigned"`
	PositionID    uint       `gorm:"column:position_id;type:int unsigned;not null"`
	Role          string     `gorm:"column:role;type:enum('admin','user');default:'user';not null"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	CurrentTeam   *Team          `gorm:"foreignKey:CurrentTeamID;references:ID"`
	Position      Position       `gorm:"foreignKey:PositionID;references:ID"`
	Teams         []Team         `gorm:"many2many:team_members;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:TeamID"`
	Projects      []Project      `gorm:"many2many:project_members;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:ProjectID"`
	Skills        []Skill        `gorm:"many2many:user_skills;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:SkillID"`
	UserSkill     []UserSkill    `gorm:"foreignKey:UserID;references:ID"`
	LeadTeams     []Team         `gorm:"foreignKey:LeaderID;references:ID"`
	LeadProjects  []Project      `gorm:"foreignKey:LeaderID;references:ID"`
	ActivityLogs  []ActivityLog  `gorm:"foreignKey:UserID;references:ID"`
	Notifications []Notification `gorm:"foreignKey:UserID;references:ID"`
	TeamMembers   []TeamMember   `gorm:"foreignKey:UserID;references:ID"`
}
