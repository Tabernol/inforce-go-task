package rarible

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	BaseURL     string
	HttpTimeout time.Duration
	ApiKey      string
}

func NewConfigFromEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return Config{
		BaseURL:     os.Getenv("RARIBLE_BASE_URL"),
		HttpTimeout: 10 * time.Second,
		ApiKey:      os.Getenv("RARIBLE_API_KEY"),
	}
}
