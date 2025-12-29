package models

import "time"

type Notification struct {
	ID        uint      `gorm:"column:id;primaryKey;type:int unsigned"`
	UserID    uint      `gorm:"column:user_id;type:int unsigned;not null"`
	Title     string    `gorm:"column:title;type:varchar(255);not null"`
	Content   string    `gorm:"column:content;type:text;not null"`
	IsRead    bool      `gorm:"column:is_read;type:boolean;default:false;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID"`
}
