package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"trieu_mock_project_go_shared/contracts"
	"trieu_mock_project_go_worker/internal/config"
	"trieu_mock_project_go_worker/internal/repositories"
	"trieu_mock_project_go_worker/internal/services"

	"github.com/robfig/cron/v3"
)

type WorkerContainer struct {
	// Repositories
	ProjectRepo *repositories.ProjectRepository

	// Services
	RabbitMQService   *services.RabbitMQService
	EmailService      *services.EmailService
	ProjectJobService *services.ProjectJobService

	// Cron
	CronScheduler *cron.Cron
}

func NewWorkerContainer() *WorkerContainer {
	// Initialize repositories
	projectRepo := repositories.NewProjectRepository()

	// Initialize services
	rabbitMQService := services.NewRabbitMQService()
	emailService := services.NewEmailService(rabbitMQService)
	projectJobService := services.NewProjectJobService(config.DB, projectRepo, emailService)

	return &WorkerContainer{
		ProjectRepo:       projectRepo,
		RabbitMQService:   rabbitMQService,
		EmailService:      emailService,
		ProjectJobService: projectJobService,
		CronScheduler:     cron.New(),
	}
}

func (c *WorkerContainer) StartEmailWorker() error {
	go func() {
		// ConsumeEmailJobs has built-in automatic recovery and will retry forever
		_ = c.RabbitMQService.ConsumeEmailJobs(func(body []byte) error {
			var job contracts.EmailJobDTO

			if err := json.Unmarshal(body, &job); err != nil {
				return err
			}

			return c.EmailService.SendEmail(job.To, job.Subject, job.TemplateName, job.Data)
		})
	}()

	log.Println("Email worker started with automatic recovery")
	return nil
}

func (c *WorkerContainer) InitializeApp() error {
	// Start RabbitMQ email worker
	if err := c.StartEmailWorker(); err != nil {
		return fmt.Errorf("failed to start email worker: %w", err)
	}

	// Start scheduled cron jobs
	c.StartCronJobs()

	return nil
}

func (c *WorkerContainer) Shutdown() {
	log.Println("Shutting down application services...")

	// Stop cron jobs
	if c.CronScheduler != nil {
		c.CronScheduler.Stop()
		log.Println("Cron jobs stopped")
	}

	if c.RabbitMQService != nil {
		c.RabbitMQService.Close()
	}
}

func (c *WorkerContainer) StartCronJobs() {
	cr := cron.New()

	// Schedule project deadline reminders every day at 8 AM
	_, err := cr.AddFunc("0 8 * * *", func() {
		log.Println("Running scheduled job: RemindProjectDeadlines")
		c.ProjectJobService.RemindProjectDeadlines(context.Background())
		log.Println("Completed scheduled job: RemindProjectDeadlines")
	})
	if err != nil {
		log.Printf("Error adding cron job for project deadline reminders: %v", err)
		return
	}

	cr.Start()
	c.CronScheduler = cr
	log.Println("Cron jobs started")
}
