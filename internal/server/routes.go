package server

import (
	"github.com/Tabernol/inforce-go-task/internal/config"
	"github.com/Tabernol/inforce-go-task/pkg/rarible"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, cfg config.Config) {
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Rarible handler
	rh := NewRaribleHandler(rarible.NewClient(cfg))
	mux.HandleFunc("/api/v1/rarible/ownership", rh.GetOwnership)
	mux.HandleFunc("/api/v1/rarible/traits", rh.GetTraitRarity)
}
