package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Tabernol/inforce-go-task/pkg/rarible"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockClient implements rarible.Client for tests.
type mockClient struct {
	// control returned values
	ownershipResp *rarible.Ownership
	ownershipErr  error

	traitResp *rarible.TraitRarityResponse
	traitErr  error
}

func (m *mockClient) GetOwnershipByID(ctx context.Context, ownershipID string) (*rarible.Ownership, error) {
	// simulate context deadline exceeded if requested
	if m.ownershipErr != nil {
		return nil, m.ownershipErr
	}
	return m.ownershipResp, nil
}

func (m *mockClient) QueryTraitsWithRarity(ctx context.Context, req rarible.TraitRarityRequest) (*rarible.TraitRarityResponse, error) {
	if m.traitErr != nil {
		return nil, m.traitErr
	}
	return m.traitResp, nil
}

func TestGetOwnership_Handler(t *testing.T) {
	validID := "ETHEREUM:0xabc:1:0xdead"

	t.Run("missing ownershipId -> 400", func(t *testing.T) {
		m := &mockClient{}
		h := NewRaribleHandler(m)

		req := httptest.NewRequest(http.MethodGet, "/rarible/ownership", nil)
		rec := httptest.NewRecorder()

		h.GetOwnership(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "ownershipId required", er.Error)
	})

	t.Run("invalid ownershipId -> 400", func(t *testing.T) {
		m := &mockClient{}
		h := NewRaribleHandler(m)

		req := httptest.NewRequest(http.MethodGet, "/rarible/ownership?ownershipId=badid", nil)
		rec := httptest.NewRecorder()

		h.GetOwnership(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "invalid ownershipId", er.Error)
	})

	t.Run("client returns success -> 200 with ownership", func(t *testing.T) {
		m := &mockClient{
			ownershipResp: &rarible.Ownership{
				Id:       validID,
				Owner:    "ETHEREUM:0xowner",
				Contract: "ETHEREUM:0xabc",
				TokenId:  "1",
				Value:    "1",
			},
		}
		h := NewRaribleHandler(m)

		req := httptest.NewRequest(http.MethodGet, "/rarible/ownership?ownershipId="+validID, nil)
		rec := httptest.NewRecorder()

		h.GetOwnership(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		var out rarible.Ownership
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&out))
		require.Equal(t, validID, out.Id)
		require.Equal(t, "ETHEREUM:0xowner", out.Owner)
	})

	t.Run("client returns error -> 502", func(t *testing.T) {
		m := &mockClient{
			ownershipErr: errors.New("remote failure"),
		}
		h := NewRaribleHandler(m)

		req := httptest.NewRequest(http.MethodGet, "/rarible/ownership?ownershipId="+validID, nil)
		rec := httptest.NewRecorder()

		h.GetOwnership(rec, req)

		require.Equal(t, http.StatusBadGateway, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "failed to fetch ownership", er.Error)
	})

	t.Run("timeout -> 504", func(t *testing.T) {
		// simulate timeout by returning context.DeadlineExceeded error
		m := &mockClient{
			ownershipErr: context.DeadlineExceeded,
		}
		h := NewRaribleHandler(m)

		req := httptest.NewRequest(http.MethodGet, "/rarible/ownership?ownershipId="+validID, nil)
		rec := httptest.NewRecorder()

		h.GetOwnership(rec, req)

		require.Equal(t, http.StatusGatewayTimeout, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "failed to fetch ownership", er.Error)
	})
}

func TestGetTraitRarity_Handler(t *testing.T) {
	validReq := rarible.TraitRarityRequest{
		CollectionId: "ETHEREUM:0xabc",
		Properties:   []rarible.TraitProperty{{Key: "Hat", Value: "Halo"}},
		Limit:        10,
	}

	t.Run("invalid json -> 400", func(t *testing.T) {
		m := &mockClient{}
		h := NewRaribleHandler(m)

		body := bytes.NewBufferString("{invalid json")
		req := httptest.NewRequest(http.MethodPost, "/rarible/traits", body)
		rec := httptest.NewRecorder()

		h.GetTraitRarity(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "invalid JSON", er.Error)
	})

	t.Run("missing collectionId -> 400", func(t *testing.T) {
		m := &mockClient{}
		h := NewRaribleHandler(m)

		payload := rarible.TraitRarityRequest{
			CollectionId: "",
			Properties:   []rarible.TraitProperty{{Key: "a", Value: "b"}},
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/rarible/traits", bytes.NewBuffer(b))
		rec := httptest.NewRecorder()

		h.GetTraitRarity(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "collectionId required", er.Error)
	})

	t.Run("missing properties -> 400", func(t *testing.T) {
		m := &mockClient{}
		h := NewRaribleHandler(m)

		payload := rarible.TraitRarityRequest{
			CollectionId: "ETHEREUM:0xabc",
			Properties:   []rarible.TraitProperty{},
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/rarible/traits", bytes.NewBuffer(b))
		rec := httptest.NewRecorder()

		h.GetTraitRarity(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "properties required", er.Error)
	})

	t.Run("successful -> 200 and body", func(t *testing.T) {
		m := &mockClient{
			traitResp: &rarible.TraitRarityResponse{
				Traits: []rarible.TraitRarity{
					{Key: "Hat", Value: "Halo", Rarity: "0"},
				},
			},
		}
		h := NewRaribleHandler(m)

		b, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/rarible/traits", bytes.NewBuffer(b))
		rec := httptest.NewRecorder()

		h.GetTraitRarity(rec, req)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		var out rarible.TraitRarityResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&out))
		require.Len(t, out.Traits, 1)
		require.Equal(t, "Hat", out.Traits[0].Key)
	})

	t.Run("remote error -> 502", func(t *testing.T) {
		m := &mockClient{
			traitErr: errors.New("remote error"),
		}
		h := NewRaribleHandler(m)

		b, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/rarible/traits", bytes.NewBuffer(b))
		rec := httptest.NewRecorder()

		h.GetTraitRarity(rec, req)

		require.Equal(t, http.StatusBadGateway, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "failed to fetch trait rarity", er.Error)
	})

	t.Run("timeout -> 504", func(t *testing.T) {
		m := &mockClient{
			traitErr: context.DeadlineExceeded,
		}
		h := NewRaribleHandler(m)

		b, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/rarible/traits", bytes.NewBuffer(b))
		rec := httptest.NewRecorder()

		h.GetTraitRarity(rec, req)

		require.Equal(t, http.StatusGatewayTimeout, rec.Result().StatusCode)
		var er errorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&er))
		require.Equal(t, "failed to fetch trait rarity", er.Error)
	})
}
