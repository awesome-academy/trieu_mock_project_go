package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"trieu_mock_project_go/internal/config"
	appErrors "trieu_mock_project_go/internal/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailQueue = "email_queue"

type RabbitMQService struct {
	conn *amqp.Connection
	mu   sync.RWMutex
}

func NewRabbitMQService() *RabbitMQService {
	s := &RabbitMQService{}
	s.connect()
	return s
}

func (s *RabbitMQService) connect() error {
	url := config.LoadConfig().RabbitMQ.GetURL()
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
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
	conn, err := s.getConnection()
	if err != nil {
		log.Printf("Error getting RabbitMQ connection: %v", err)
		return appErrors.ErrInternalServerError
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Error creating RabbitMQ channel: %v", err)
		return appErrors.ErrInternalServerError
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		EmailQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Printf("Error declaring RabbitMQ queue: %v", err)
		return appErrors.ErrInternalServerError
	}

	body, err := json.Marshal(job)
	if err != nil {
		log.Printf("Error marshaling job to JSON: %v", err)
		return appErrors.ErrInternalServerError
	}

	err = ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
		return appErrors.ErrFailedToPublishMessage
	}

	log.Printf(" [x] Sent email job to queue: %s", body)
	return nil
}

func (s *RabbitMQService) ConsumeEmailJobs(handler func(body []byte) error) error {
	return s.consumeWithRecovery(handler)
}

func (s *RabbitMQService) consumeWithRecovery(handler func(body []byte) error) error {
	for {
		err := s.startConsumer(handler)
		if err != nil {
			log.Printf("CRITICAL: Email consumer stopped with error: %v", err)
			log.Printf("Attempting to reconnect in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (s *RabbitMQService) startConsumer(handler func(body []byte) error) error {
	conn, err := s.getConnection()
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global (applies to this channel only)
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	q, err := ch.QueueDeclare(
		EmailQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf(" [*] Waiting for messages in %s.", EmailQueue)

	// Monitor channel closure
	closeChan := make(chan *amqp.Error)
	ch.NotifyClose(closeChan)

	for {
		select {
		case err := <-closeChan:
			if err != nil {
				return fmt.Errorf("channel closed with error: %w", err)
			}
			return fmt.Errorf("channel closed unexpectedly")
		case d, ok := <-msgs:
			if !ok {
				log.Printf("CRITICAL: Messages channel closed, consumer stopping")
				return fmt.Errorf("messages channel closed")
			}
			log.Printf("Received a message: %s", d.Body)
			err := handler(d.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				_ = d.Nack(false, true)
			} else {
				_ = d.Ack(false)
			}
		}
	}
}
