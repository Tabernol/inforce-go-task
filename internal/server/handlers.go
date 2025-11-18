package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Tabernol/inforce-go-task/pkg/rarible"
	"net/http"
	"regexp"
	"time"
)

type errorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type RaribleHandler struct {
	client    rarible.Client
	timeout   time.Duration
	ethRegexp *regexp.Regexp
}

func NewRaribleHandler(client rarible.Client) *RaribleHandler {
	return &RaribleHandler{
		client:    client,
		timeout:   10 * time.Second, // timeout for third-party API
		ethRegexp: regexp.MustCompile(`^ETHEREUM:.+:.+:.+$`),
	}
}

// GetOwnership handles GET /ownership?ownershipId=...
func (h *RaribleHandler) GetOwnership(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ownershipID := r.URL.Query().Get("ownershipId")
	if ownershipID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "ownershipId required",
			Details: "Query param ownershipId is missing",
		})
		return
	}

	if !h.ethRegexp.MatchString(ownershipID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "invalid ownershipId",
			Details: "ownershipId must match format ETHEREUM:contract:tokenId:owner",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	data, err := h.client.GetOwnershipByID(ctx, ownershipID)
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusGatewayTimeout
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "failed to fetch ownership",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// GetTraitRarity handles POST /traits
func (h *RaribleHandler) GetTraitRarity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req rarible.TraitRarityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "invalid JSON",
			Details: err.Error(),
		})
		return
	}

	if req.CollectionId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "collectionId required",
			Details: "collectionId must not be empty",
		})
		return
	}

	if len(req.Properties) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "properties required",
			Details: "properties array must not be empty",
		})
		return
	}

	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 50
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	data, err := h.client.QueryTraitsWithRarity(ctx, req)
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusGatewayTimeout
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(errorResponse{
			Error:   "failed to fetch trait rarity",
			Details: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
