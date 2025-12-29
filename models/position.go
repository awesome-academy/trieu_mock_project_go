package models

import "time"

type Position struct {
	ID           uint      `gorm:"column:id;primaryKey;type:int unsigned"`
	Name         string    `gorm:"column:name;type:varchar(255);not null"`
	Abbreviation string    `gorm:"column:abbreviation;type:varchar(50);not null"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	Users []User `gorm:"foreignKey:PositionID;references:ID"`
}
