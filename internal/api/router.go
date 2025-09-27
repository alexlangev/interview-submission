// internal/api/router.go
package api

import (
	"net/http"

	v1 "github.com/alexlangev/interview-submission/internal/api/v1/income-tax"
	"github.com/alexlangev/interview-submission/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(calc *core.Calculator) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Mount versioned API
	r.Mount("/v1", v1.New(calc))

	return r
}
