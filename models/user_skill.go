package models

import "time"

type UserSkill struct {
	UserID         uint      `gorm:"column:user_id;primaryKey;type:int unsigned;not null"`
	SkillID        uint      `gorm:"column:skill_id;primaryKey;type:int unsigned;not null"`
	Level          int       `gorm:"column:level;type:int;not null"`
	UsedYearNumber int       `gorm:"column:used_year_number;type:int;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Skill Skill `gorm:"foreignKey:SkillID;references:ID"`
}
