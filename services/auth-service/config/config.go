package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	// Load .env if it exists.
	_ = godotenv.Load()

	return &Config{
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://postgres:localhost@localhost:5432/shopsphere_db?sslmode=disable",
		),
		Port: getEnv("PORT", "8081"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
