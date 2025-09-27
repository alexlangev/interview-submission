package v1

import (
	"github.com/alexlangev/interview-submission/internal/core"
	"github.com/go-chi/chi/v5"
)

func New(calc *core.Calculator) chi.Router {
	r := chi.NewRouter()

	r.Get("/income-tax", Calculate(calc))

	return r
}
