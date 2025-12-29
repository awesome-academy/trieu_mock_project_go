package models

import "time"

type Skill struct {
	ID        uint      `gorm:"column:id;primaryKey;type:int unsigned"`
	Name      string    `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_name"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null"`

	// Relationships
	Users      []User      `gorm:"many2many:user_skills;foreignKey:ID;joinForeignKey:SkillID;references:ID;joinReferences:UserID"`
	UserSkills []UserSkill `gorm:"foreignKey:SkillID;references:ID"`
}
