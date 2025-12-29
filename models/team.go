package models

import "time"

type Team struct {
	ID          uint      `gorm:"column:id;primaryKey;type:int unsigned"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description *string   `gorm:"column:description;type:text"`
	LeaderID    uint      `gorm:"column:leader_id;type:int unsigned;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	Leader      User         `gorm:"foreignKey:LeaderID;references:ID"`
	Members     []User       `gorm:"many2many:team_members;foreignKey:ID;joinForeignKey:TeamID;references:ID;joinReferences:UserID"`
	Projects    []Project    `gorm:"foreignKey:TeamID;references:ID"`
	TeamMembers []TeamMember `gorm:"foreignKey:TeamID;references:ID"`
}
