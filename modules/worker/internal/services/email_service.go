package services

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"path/filepath"
	"trieu_mock_project_go_shared/contracts"
	"trieu_mock_project_go_worker/internal/config"
	"trieu_mock_project_go_worker/internal/dtos"

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

func (s *EmailService) Enqueue(to string, subject string, templateName string, data interface{}) {
	job := contracts.EmailJobDTO{
		To:           to,
		Subject:      subject,
		TemplateName: templateName,
		Data:         data,
	}
	body, _ := json.Marshal(job)
	err := s.rabbitMQService.PublishEmailJob(body)
	if err != nil {
		log.Printf("Failed to enqueue email job: %v", err)
	}
}

func (s *EmailService) SendEmail(to string, subject string, templateName string, data interface{}) error {
	tmplPath := filepath.Join("templates", "emails", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Failed to parse template %s: %v", tmplPath, err)
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("Failed to execute template %s: %v", tmplPath, err)
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

func (s *EmailService) SendProjectDeadlineReminderEmail(data dtos.ProjectDeadlineReminderEmailDTO) {
	s.Enqueue(data.To, "Project Reminder: "+data.ProjectName+" is due soon", "project_deadline_reminder.html", data)
}
