package models

import "time"

type TeamMember struct {
	ID        uint       `gorm:"column:id;primaryKey;type:int unsigned"`
	UserID    uint       `gorm:"column:user_id;type:int unsigned;not null"`
	TeamID    uint       `gorm:"column:team_id;type:int unsigned;not null"`
	JoinedAt  time.Time  `gorm:"column:joined_at;type:timestamp;autoCreateTime;not null"`
	LeftAt    *time.Time `gorm:"column:left_at;type:timestamp"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`

	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID"`
	Team Team `gorm:"foreignKey:TeamID;references:ID"`
}
