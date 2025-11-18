package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tabernol/inforce-go-task/internal/config"
	"github.com/Tabernol/inforce-go-task/pkg/rarible"
	"github.com/stretchr/testify/require"
)

func TestGetOwnershipByID(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/v0.1/ownerships/test-id", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)
		resp := rarible.Ownership{
			Id:       "test-id",
			Owner:    "0xabc",
			Contract: "0x123",
			TokenId:  "1",
			Value:    "100",
		}
		json.NewEncoder(w).Encode(resp)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := config.Config{
		RaribleBaseURL: server.URL,
		RaribleTimeout: 2 * time.Second,
		RaribleAPIKey:  "",
	}
	client := rarible.NewClient(cfg)

	ownership, err := client.GetOwnershipByID(context.Background(), "test-id")
	require.NoError(t, err)
	require.Equal(t, "test-id", ownership.Id)
	require.Equal(t, "0xabc", ownership.Owner)
}

func TestQueryTraitsWithRarity(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/v0.1/items/traits/rarity", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var req rarible.TraitRarityRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		require.Equal(t, "collection-1", req.CollectionId)

		resp := rarible.TraitRarityResponse{
			Traits: []rarible.TraitRarity{
				{Key: "Hat", Value: "Halo", Rarity: "0"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := config.Config{
		RaribleBaseURL: server.URL,
		RaribleTimeout: 2 * time.Second,
		RaribleAPIKey:  "",
	}
	client := rarible.NewClient(cfg)

	req := rarible.TraitRarityRequest{
		CollectionId: "collection-1",
		Properties:   []rarible.TraitProperty{{Key: "Hat", Value: "Halo"}},
		Limit:        10,
	}

	resp, err := client.QueryTraitsWithRarity(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, resp.Traits, 1)
	require.Equal(t, "Hat", resp.Traits[0].Key)
	require.Equal(t, "Halo", resp.Traits[0].Value)
}
