package models

import "time"

type Project struct {
	ID           uint       `gorm:"column:id;primaryKey;type:int unsigned"`
	Name         string     `gorm:"column:name;type:varchar(255);not null"`
	Abbreviation string     `gorm:"column:abbreviation;type:varchar(50);not null"`
	StartDate    *time.Time `gorm:"column:start_date;type:date"`
	EndDate      *time.Time `gorm:"column:end_date;type:date"`
	LeaderID     uint       `gorm:"column:leader_id;type:int unsigned;not null"`
	TeamID       uint       `gorm:"column:team_id;type:int unsigned;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	Leader  User   `gorm:"foreignKey:LeaderID;references:ID"`
	Team    Team   `gorm:"foreignKey:TeamID;references:ID"`
	Members []User `gorm:"many2many:project_members;foreignKey:ID;joinForeignKey:ProjectID;references:ID;joinReferences:UserID"`
}
