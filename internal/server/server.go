package server

import (
	"fmt"
	"github.com/Tabernol/inforce-go-task/internal/config"
	"log"
	"net/http"
)

type Server struct {
	cfg config.Config
	mux *http.ServeMux
}

func NewServer(cfg config.Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}

	RegisterRoutes(s.mux, cfg)
	return s
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) Run() {
	addr := fmt.Sprintf(":%s", s.cfg.ServerPort)
	log.Printf("Server running on %s", addr)
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
