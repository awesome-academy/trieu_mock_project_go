package services

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
	"trieu_mock_project_go/internal/config"
	"trieu_mock_project_go/internal/dtos"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config          config.MailConfig
	rabbitMQService *RabbitMQService
}

func NewEmailService(rabbitMQService *RabbitMQService) *EmailService {
	cfg := config.LoadConfig()
	return &EmailService{
		config:          cfg.Mail,
		rabbitMQService: rabbitMQService,
	}
}

func (s *EmailService) enqueue(to string, subject string, templateName string, data interface{}) {
	job := dtos.EmailJobDTO{
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

func (s *EmailService) SendEmail(to string, subject string, templateName string, data interface{}) error {
	tmplPath := filepath.Join("templates", "emails", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUser, s.config.SMTPPassword)

	return d.DialAndSend(m)
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

func (s *EmailService) SendProjectDeadlineReminderEmail(data dtos.ProjectDeadlineReminderEmailDTO) {
	s.enqueue(data.To, "Project Reminder: "+data.ProjectName+" is due soon", "project_deadline_reminder.html", data)
}
