package rarible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tabernol/inforce-go-task/internal/config"
	"io"
	"net/http"
	"strings"
)

type HttpClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewClient(cfg config.Config) Client {
	return &HttpClient{
		baseURL: cfg.RaribleBaseURL,
		apiKey:  cfg.RaribleAPIKey,
		client: &http.Client{
			Timeout: cfg.RaribleTimeout,
		},
	}
}

func (c *HttpClient) doRequest(ctx context.Context, method, path string, body any, out any) error {
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	base := strings.TrimRight(c.baseURL, "/")
	req, err := http.NewRequestWithContext(ctx, method, base+path, reqBody)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
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

func (c *HttpClient) GetOwnershipByID(ctx context.Context, ownershipID string) (*Ownership, error) {
	var out Ownership
	path := fmt.Sprintf("/v0.1/ownerships/%s", ownershipID)

	if err := c.doRequest(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *HttpClient) QueryTraitsWithRarity(ctx context.Context, req TraitRarityRequest) (*TraitRarityResponse, error) {
	var out TraitRarityResponse
	path := "/v0.1/items/traits/rarity"

	if err := c.doRequest(ctx, http.MethodPost, path, req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
