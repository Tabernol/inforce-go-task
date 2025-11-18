package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	ServerPort     string
	RaribleAPIKey  string
	RaribleBaseURL string
	RaribleTimeout time.Duration
}

func LoadConfig() Config {
	// Load from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("RARIBLE_API_KEY")
	baseURL := os.Getenv("RARIBLE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.rarible.org"
	}

	return Config{
		ServerPort:     port,
		RaribleAPIKey:  apiKey,
		RaribleBaseURL: baseURL,
		RaribleTimeout: 10 * time.Second,
	}
}
