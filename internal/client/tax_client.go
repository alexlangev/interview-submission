package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexlangev/interview-submission/internal/models"
)

type Client struct {
	BaseURL      string
	HTTPClient   *http.Client
	MaxRetries   int
	RetriesDelay time.Duration
	TotalTimeout time.Duration
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	return &Client{
		BaseURL:      baseURL,
		HTTPClient:   httpClient,
		MaxRetries:   5,
		RetriesDelay: 250 * time.Millisecond,
		TotalTimeout: 10 * time.Second,
	}
}

func (c *Client) GetTaxBrackets(year int) ([]models.TaxBracket, error) {
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%d", c.BaseURL, year)

	start := time.Now() // first query
	attempts := c.MaxRetries + 1

	for attempt := 0; attempt < attempts; attempt++ {
		fmt.Println("client attempt: ", attempt)
		// sleep between retries
		if attempt > 0 {
			if time.Since(start)+c.RetriesDelay > c.TotalTimeout {
				return nil, fmt.Errorf("Timeout: exceeded total timeout (%s)", c.TotalTimeout)
			}
			time.Sleep(c.RetriesDelay)
		}

		resp, err := c.HTTPClient.Get(url)
		if err != nil {
			if attempt < c.MaxRetries && time.Since(start) < c.TotalTimeout {
				continue
			}
			return nil, fmt.Errorf("request failed: %w", err)
		}

		// Status 200
		if resp.StatusCode == http.StatusOK {
			var responseData models.TaxResponse

			err := json.NewDecoder(resp.Body).Decode(&responseData)
			if err != nil {
				resp.Body.Close()
				return nil, fmt.Errorf("decode JSON: %w", err)
			}
			resp.Body.Close()

			brackets := []models.TaxBracket{}
			for _, v := range responseData.TaxBrackets {
				brackets = append(brackets, models.TaxBracket{Min: v.Min, Max: v.Max, Rate: v.Rate})
			}
			return brackets, nil
		}

		resp.Body.Close()

		// Status 500
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			if attempt < c.MaxRetries && time.Since(start) < c.TotalTimeout {
				continue
			}
			return nil, fmt.Errorf("Server error after retries: %d", resp.StatusCode)
		}

		// Other non-retryable errors
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}

	return nil, fmt.Errorf("exhausted retries fetching tax brackets")
}
