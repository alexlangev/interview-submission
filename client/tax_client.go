package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexlangev/interview-submission/internal/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	// retries logic?
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

func (c *Client) GetTaxBrackets(year int) ([]models.TaxBracket, error) {
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%d", c.BaseURL, year)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}

	var payload models.TaxResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}

	return payload.TaxBrackets, nil
}
