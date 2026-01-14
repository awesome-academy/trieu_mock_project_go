package services

import (
	"context"
	"encoding/json"
	"log"
	"trieu_mock_project_go/internal/config"
	appErrors "trieu_mock_project_go/internal/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailQueue = "email_queue"

type RabbitMQService struct {
	conn *amqp.Connection
}

func NewRabbitMQService() *RabbitMQService {
	return &RabbitMQService{
		conn: config.RabbitMQConn,
	}
}

func (s *RabbitMQService) PublishEmailJob(job interface{}) *appErrors.AppError {
	ch, err := s.conn.Channel()
	if err != nil {
		return appErrors.ErrFailedToPublishMessage
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
		return appErrors.ErrFailedToPublishMessage
	}

	body, err := json.Marshal(job)
	if err != nil {
		return appErrors.ErrFailedToPublishMessage
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
		return appErrors.ErrFailedToPublishMessage
	}

	log.Printf(" [x] Sent email job to queue: %s", body)
	return nil
}

func (s *RabbitMQService) ConsumeEmailJobs(handler func(body []byte) error) error {
	ch, err := s.conn.Channel()
	if err != nil {
		return err
	}
	// We don't close the channel here because we want to keep consuming

	q, err := ch.QueueDeclare(
		EmailQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
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
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := handler(d.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				_ = d.Nack(false, true)
			} else {
				_ = d.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s.", EmailQueue)
	return nil
}
