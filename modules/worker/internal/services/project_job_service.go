package services

import (
	"context"
	"log"
	"trieu_mock_project_go_worker/internal/dtos"
	"trieu_mock_project_go_worker/internal/repositories"

	"gorm.io/gorm"
)

type ProjectJobService struct {
	db                *gorm.DB
	projectRepository *repositories.ProjectRepository
	emailService      *EmailService
}

func NewProjectJobService(db *gorm.DB, projectRepository *repositories.ProjectRepository, emailService *EmailService) *ProjectJobService {
	return &ProjectJobService{
		db:                db,
		projectRepository: projectRepository,
		emailService:      emailService,
	}
}

func (s *ProjectJobService) RemindProjectDeadlines(c context.Context) error {
	projects, err := s.projectRepository.FindProjectsNearDeadline(s.db.WithContext(c), 3)
	if err != nil {
		log.Printf("Error finding projects near deadline: %v", err)
		return err
	}

	log.Printf("Found %d projects nearing deadlines for reminders\n", len(projects))
	for _, p := range projects {
		if p.EndDate == nil {
			continue
		}

		dueDate := p.EndDate.Format("2006-01-02")

		for _, m := range p.Members {
			s.emailService.SendProjectDeadlineReminderEmail(dtos.ProjectDeadlineReminderEmailDTO{
				To:          m.Email,
				UserName:    m.Name,
				ProjectName: p.Name,
				DueDate:     dueDate,
			})
		}
	}

	return nil
}
