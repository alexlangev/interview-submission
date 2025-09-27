// internal/api/router.go
package api

import (
	"net/http"

	v1 "github.com/alexlangev/interview-submission/internal/api/v1/income-tax"
	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Top-level health probe (not versioned)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Mount versioned API
	r.Mount("/v1", v1.New())

	return r
}
