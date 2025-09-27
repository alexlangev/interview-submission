package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const mockJSON = `{
  "tax_brackets": [
    {
        "min": 0,
        "max": 50197,
        "rate": 0.15
    },
    {
        "min": 50197,
        "max": 100392,
        "rate": 0.205
    },
    {
        "min": 100392,
        "max": 155625,
        "rate": 0.26
    },
    {
        "min": 155625,
        "max": 221708,
        "rate": 0.29
    },
    {
        "min": 221708,
        "rate": 0.33
    }
  ]
}`

func newTestClient(baseURL string) *Client {
	return &Client{
		BaseURL:      baseURL,
		HTTPClient:   &http.Client{Timeout: 2 * time.Second},
		MaxRetries:   0,
		RetriesDelay: 0,
		TotalTimeout: 2 * time.Second,
	}
}

// 200
func TestTaxBrackets_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/tax-calculator/tax-year/2022"; got != want {
			t.Fatalf("expected path %q, got %q", want, got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	t.Cleanup(srv.Close)

	c := newTestClient(srv.URL)
	brackets, err := c.GetTaxBrackets(2022)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Basic shape checks
	if len(brackets) != 5 {
		t.Fatalf("expected 5 brackets, got %d", len(brackets))
	}
}

func TestTaxBrackets_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tax_brackets":"12345"}`)) // not an array

	}))
	t.Cleanup(srv.Close)

	c := newTestClient(srv.URL)
	if _, err := c.GetTaxBrackets(2022); err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}
}

func TestTaxBrackets_500_WithRetries(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits += 1
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	c := newTestClient(srv.URL)
	c.MaxRetries = 2

	_, err := c.GetTaxBrackets(2022)
	if err == nil {
		t.Fatal("expected error on repeated 500s")
	}
	if hits != 3 {
		t.Fatalf("expected 3 attempts (1+2 retries), got %d", hits)
	}
}

func TestTaxBrackets_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockJSON))
	}))
	t.Cleanup(srv.Close)

	c := newTestClient(srv.URL)
	c.HTTPClient.Timeout = 50 * time.Millisecond

	if _, err := c.GetTaxBrackets(2022); err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}
