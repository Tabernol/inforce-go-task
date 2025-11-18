package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tabernol/inforce-go-task/pkg/rarible"
	"net/http"
)

type RaribleHandler struct {
	client rarible.Client
}

func NewRaribleHandler(client rarible.Client) *RaribleHandler {
	return &RaribleHandler{client: client}
}

func (h *RaribleHandler) GetOwnership(w http.ResponseWriter, r *http.Request) {
	fmt.Print("try to GET info by ownershipId")
	tokenID := r.URL.Query().Get("ownershipId")
	if tokenID == "" {
		http.Error(w, "ownershipId required", http.StatusBadRequest)
		return
	}

	data, err := h.client.GetOwnershipByID(context.Background(), tokenID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *RaribleHandler) GetTraitRarity(w http.ResponseWriter, r *http.Request) {
	// парсим body JSON
	var req rarible.TraitRarityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.client.QueryTraitsWithRarity(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
