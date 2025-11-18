package main

import (
	"github.com/Tabernol/inforce-go-task/internal/config"
	"github.com/Tabernol/inforce-go-task/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	srv := server.NewServer(cfg)
	srv.Run()
}
