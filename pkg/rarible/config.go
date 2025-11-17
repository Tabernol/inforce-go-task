package rarible

import "time"

type Config struct {
	BaseURL     string
	HttpTimeout time.Duration
	ApiKey      string // maybe key
}
