package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() chi.Router {
	r := chi.NewRouter()

	r.Get("/income-tax", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("smoke at v1/calculate"))
	})

	return r
}
