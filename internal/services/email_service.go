package services

import (
	"log"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go_shared/contracts"
)

type EmailService struct {
	rabbitMQService *RabbitMQService
}

func NewEmailService(rabbitMQService *RabbitMQService) *EmailService {
	return &EmailService{
		rabbitMQService: rabbitMQService,
	}
}

func (s *EmailService) enqueue(to string, subject string, templateName string, data interface{}) {
	job := contracts.EmailJobDTO{
		To:           to,
		Subject:      subject,
		TemplateName: templateName,
		Data:         data,
	}
	err := s.rabbitMQService.PublishEmailJob(job)
	if err != nil {
		log.Printf("Failed to enqueue email job: %v", err)
	}
}

func (s *EmailService) SendTeamJoinEmail(data dtos.TeamMembershipEmailDTO) {
	s.enqueue(data.To, "Welcome to Team "+data.TeamName, "team_join.html", data)
}

func (s *EmailService) SendTeamLeaveEmail(data dtos.TeamMembershipEmailDTO) {
	s.enqueue(data.To, "Leaving Team "+data.TeamName, "team_leave.html", data)
}

func (s *EmailService) SendProjectJoinEmail(data dtos.ProjectMembershipEmailDTO) {
	s.enqueue(data.To, "Assigned to Project "+data.ProjectName, "project_join.html", data)
}

func (s *EmailService) SendProjectLeaveEmail(data dtos.ProjectMembershipEmailDTO) {
	s.enqueue(data.To, "Removed from Project "+data.ProjectName, "project_leave.html", data)
}
