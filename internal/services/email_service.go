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
	config config.MailConfig
}

func NewEmailService() *EmailService {
	cfg := config.LoadConfig()
	return &EmailService{config: cfg.Mail}
}

func (s *EmailService) send(to string, subject string, templateName string, data interface{}) {
	// Send email in a goroutine to avoid blocking
	go func() {
		err := s.SendEmail(to, subject, templateName, data)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", to, err)
		}
	}()
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
	s.send(data.To, "Welcome to Team "+data.TeamName, "team_join.html", data)
}

func (s *EmailService) SendTeamLeaveEmail(data dtos.TeamMembershipEmailDTO) {
	s.send(data.To, "Leaving Team "+data.TeamName, "team_leave.html", data)
}

func (s *EmailService) SendProjectJoinEmail(data dtos.ProjectMembershipEmailDTO) {
	s.send(data.To, "Assigned to Project "+data.ProjectName, "project_join.html", data)
}

func (s *EmailService) SendProjectLeaveEmail(data dtos.ProjectMembershipEmailDTO) {
	s.send(data.To, "Removed from Project "+data.ProjectName, "project_leave.html", data)
}
