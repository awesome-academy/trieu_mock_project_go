package services

import (
	"context"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"

	"gorm.io/gorm"
)

type SkillService struct {
	db              *gorm.DB
	skillRepository *repositories.SkillRepository
}

func NewSkillService(db *gorm.DB, skillRepository *repositories.SkillRepository) *SkillService {
	return &SkillService{db: db, skillRepository: skillRepository}
}

func (s *SkillService) GetAllSkillsSummary(c context.Context) []dtos.SkillSummary {
	skills, err := s.skillRepository.FindAllSkillSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.SkillSummary{}
	}

	skillDtos := make([]dtos.SkillSummary, 0, len(skills))
	for _, skill := range skills {
		skillDtos = append(skillDtos, dtos.SkillSummary{
			ID:   skill.ID,
			Name: skill.Name,
		})
	}
	return skillDtos
}
