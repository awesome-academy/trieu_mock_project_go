package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Server           ServerConfig
	Database         DatabaseConfig
	Redis            RedisConfig
	RabbitMQ         RabbitMQConfig
	SessionConfig    SessionConfig
	JWT              JWTConfig
	RequestRateLimit int
}

type ServerConfig struct {
	Host string
	Port string
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

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
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

type SessionConfig struct {
	Secret string
	MaxAge int
	Secure bool
}

type JWTConfig struct {
	Secret string
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads config once and returns the singleton instance
func LoadConfig() *Config {
	once.Do(func() {
		// Load .env file
		godotenv.Load()

		sessionMaxAge, err := strconv.Atoi(getEnv("SESSION_MAX_AGE", "3600"))
		if err != nil {
			sessionMaxAge = 3600
		}
		maxIdleConns, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
		if err != nil {
			maxIdleConns = 10
		}
		maxOpenConns, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
		if err != nil {
			maxOpenConns = 100
		}
		redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
		if err != nil {
			redisDB = 0
		}
		requestRateLimit, err := strconv.Atoi(getEnv("REQUEST_RATE_LIMIT", "10"))
		if err != nil {
			requestRateLimit = 10
		}
		cfg = &Config{
			Server: ServerConfig{
				Host: getEnv("SERVER_HOST", "localhost"),
				Port: getEnv("SERVER_PORT", "8080"),
			},
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
			Redis: RedisConfig{
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     getEnv("REDIS_PORT", "6379"),
				Username: getEnv("REDIS_USERNAME", ""),
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       redisDB,
			},
			RabbitMQ: RabbitMQConfig{
				Host:     getEnv("RABBITMQ_HOST", "localhost"),
				Port:     getEnv("RABBITMQ_PORT", "5672"),
				User:     getEnv("RABBITMQ_USER", "guest"),
				Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			},
			SessionConfig: SessionConfig{
				Secret: getEnv("SESSION_SECRET", "trieu-mock-project-go-secret"),
				MaxAge: sessionMaxAge,
				Secure: getEnv("SESSION_SECURE", "false") == "true",
			},
			JWT: JWTConfig{
				Secret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			},
			RequestRateLimit: requestRateLimit,
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
