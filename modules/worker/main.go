package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"trieu_mock_project_go_worker/internal/bootstrap"
	"trieu_mock_project_go_worker/internal/config"
)

func main() {
	// Initialize config
	config.LoadConfig()

	// Initialize Database
	if err := config.ConnectToMySQL(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize worker container
	container := bootstrap.NewWorkerContainer()

	// Initialize Application Services (Email worker, Cron jobs)
	if err := container.InitializeApp(); err != nil {
		log.Fatalf("Failed to initialize worker services: %v", err)
	}

	log.Println("Worker service is running.")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker service...")
	container.Shutdown()

	log.Println("Worker service exited")
}
