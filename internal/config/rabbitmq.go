package config

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection

func InitRabbitMQ() error {
	rabbitConfig := LoadConfig().RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitConfig.User,
		rabbitConfig.Password,
		rabbitConfig.Host,
		rabbitConfig.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	RabbitMQConn = conn
	log.Println("RabbitMQ connected successfully")
	return nil
}
