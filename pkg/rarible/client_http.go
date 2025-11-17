package rarible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewClient(cfg Config) Client {
	return &httpClient{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.ApiKey,
		client: &http.Client{
			Timeout: cfg.HttpTimeout,
		},
	}
}

func (c *httpClient) doRequest(ctx context.Context, method, path string, body any, out any) error {
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}

func (c *httpClient) GetOwnershipByID(ctx context.Context, ownershipID string) (*Ownership, error) {
	var out Ownership
	path := fmt.Sprintf("/v0.1/ownerships/%s", ownershipID)

	if err := c.doRequest(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *httpClient) QueryTraitsWithRarity(ctx context.Context, req TraitRarityRequest) (*TraitRarityResponse, error) {
	var out TraitRarityResponse
	path := "/v0.1/items/traits/rarity" // <- новий endpoint

	if err := c.doRequest(ctx, http.MethodPost, path, req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

//func (c *httpClient) QueryTraitsWithRarity(ctx context.Context, req TraitRarityRequest) (*TraitRarityResponse, error) {
//	var out TraitRarityResponse
//	path := "/v0.1/items/byTraitRarities"
//
//	// Викликаємо doRequest для POST + JSON
//	if err := c.doRequest(ctx, http.MethodPost, path, req, &out); err != nil {
//		return nil, err
//	}
//
//	return &out, nil
//}

//func (c *httpClient) QueryTraitsWithRarity(ctx context.Context, req TraitRarityRequest) (*TraitRarityResponse, error) {
//	var out TraitRarityResponse
//	path := "/v0.1/items/byTraitRarities"
//
//	if err := c.doRequest(ctx, http.MethodPost, path, req, &out); err != nil {
//		return nil, err
//	}
//
//	return &out, nil
//}
