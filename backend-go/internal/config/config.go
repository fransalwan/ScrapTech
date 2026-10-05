package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/scrapflow_db?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-fintech-key-change-in-production"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		Environment: env,
	}
}
