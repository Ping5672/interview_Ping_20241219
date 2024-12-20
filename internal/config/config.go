package config

import (
	"fmt"
	"os"
)

var IsTestEnvironment bool

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetDatabaseConfig() DatabaseConfig {
	if IsTestEnvironment {
		return DatabaseConfig{
			Host:     "localhost", // Use localhost for tests
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "game_db",
		}
	}

	// Production/Docker configuration
	return DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "db"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("DB_NAME", "game_db"),
	}
}

func (c DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
