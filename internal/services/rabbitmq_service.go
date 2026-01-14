package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"trieu_mock_project_go/internal/config"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go_shared/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn  *amqp.Connection
	pubCh *amqp.Channel
	mu    sync.RWMutex
}

func NewRabbitMQService() *RabbitMQService {
	s := &RabbitMQService{}
	_ = s.connect()
	return s
}

func (s *RabbitMQService) connect() error {
	url := config.LoadConfig().RabbitMQ.GetURL()
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	s.conn = conn
	s.pubCh = nil
	return nil
}

func (s *RabbitMQService) getPublishChannel() (*amqp.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil || s.conn.IsClosed() {
		if err := s.connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
		}
	}

	if s.pubCh == nil || s.pubCh.IsClosed() {
		ch, err := s.conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to create channel: %w", err)
		}

		s.pubCh = ch
	}

	return s.pubCh, nil
}

func (s *RabbitMQService) getConnection() (*amqp.Connection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil || s.conn.IsClosed() {
		if err := s.connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
		}
	}

	return s.conn, nil
}

func (s *RabbitMQService) PublishEmailJob(job interface{}) *appErrors.AppError {
	body, err := json.Marshal(job)
	if err != nil {
		log.Printf("Error marshaling job to JSON: %v", err)
		return appErrors.ErrInternalServerError
	}

	ch, err := s.getPublishChannel()
	if err != nil {
		log.Printf("Error getting RabbitMQ channel: %v", err)
		return appErrors.ErrInternalServerError
	}

	// Lock when publishing to ensure thread-safety on shared channel
	s.mu.Lock()
	err = ch.PublishWithContext(context.Background(),
		"",                  // exchange
		rabbitmq.EmailQueue, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	s.mu.Unlock()

	if err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
		return appErrors.ErrFailedToPublishMessage
	}

	log.Printf(" [x] Sent email job to queue: %s", body)
	return nil
}

func (s *RabbitMQService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pubCh != nil && !s.pubCh.IsClosed() {
		_ = s.pubCh.Close()
	}
	if s.conn != nil && !s.conn.IsClosed() {
		_ = s.conn.Close()
	}
	log.Println("RabbitMQ connection and channels closed gracefully")
}
