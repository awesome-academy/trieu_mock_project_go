package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Mail     MailConfig
}

type DatabaseConfig struct {
	Driver       string
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	MaxIdleConns int
	MaxOpenConns int
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (c RabbitMQConfig) GetURL() string {
	return "amqp://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/"
}

type MailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SenderEmail  string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		godotenv.Load()

		maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
		maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
		smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "1025"))

		cfg = &Config{
			Database: DatabaseConfig{
				Driver:       getEnv("DB_DRIVER", "mysql"),
				Host:         getEnv("DB_HOST", "localhost"),
				Port:         getEnv("DB_PORT", "3306"),
				User:         getEnv("DB_USER", "root"),
				Password:     getEnv("DB_PASSWORD", "password"),
				Database:     getEnv("DB_NAME", "trieu_mock_project_go"),
				MaxIdleConns: maxIdleConns,
				MaxOpenConns: maxOpenConns,
			},
			RabbitMQ: RabbitMQConfig{
				Host:     getEnv("RABBITMQ_HOST", "localhost"),
				Port:     getEnv("RABBITMQ_PORT", "5672"),
				User:     getEnv("RABBITMQ_USER", "guest"),
				Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			},
			Mail: MailConfig{
				SMTPHost:     getEnv("SMTP_HOST", "localhost"),
				SMTPPort:     smtpPort,
				SMTPUser:     getEnv("SMTP_USER", ""),
				SMTPPassword: getEnv("SMTP_PASSWORD", ""),
				SenderEmail:  getEnv("SENDER_EMAIL", "no-reply@trieu-mock-project-go.com"),
			},
		}
	})
	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
