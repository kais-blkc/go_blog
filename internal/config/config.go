package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	JwtSecret   string
	Port        string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		DatabaseUrl: getEnv("DATABASE_URL", ""),
		JwtSecret:   getEnv("JWT_SECRET", ""),
		Port:        getEnv("PORT", "8080"),
	}

	return config, nil
}
