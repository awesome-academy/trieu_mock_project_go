package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	AppSecret string
	JWT       JWTConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Database string
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

		cfg = &Config{
			Server: ServerConfig{
				Host: getEnv("SERVER_HOST", "localhost"),
				Port: getEnv("SERVER_PORT", "8080"),
			},
			Database: DatabaseConfig{
				Driver:   getEnv("DB_DRIVER", "mysql"),
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "3306"),
				User:     getEnv("DB_USER", "root"),
				Password: getEnv("DB_PASSWORD", "password"),
				Database: getEnv("DB_NAME", "trieu_mock_project_go"),
			},
			AppSecret: getEnv("SESSION_SECRET", "trieu-mock-project-go-secret"),
			JWT: JWTConfig{
				Secret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
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
