package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexlangev/interview-submission/internal/core"
	u "github.com/alexlangev/interview-submission/internal/utils"
)

type CalculateResponse struct {
	Year          int                  `json:"year"`
	SalaryInput   string               `json:"salary_input"`
	SalaryCents   int64                `json:"salary_cents"`
	TotalTaxCents int64                `json:"total_tax_cents,omitempty"`
	EffectiveRate float64              `json:"effective_rate,omitempty"`
	PerBracket    []PerBracketResponse `json:"per_bracket"`
	Message       string               `json:"message"`
}

type PerBracketResponse struct {
	MinCents        int64  `json:"min_cents"`
	MaxCents        *int64 `json:"max_cents,omitempty"`
	RateBasisPoints int64  `json:"rate_basis_points"`
	TaxCents        int64  `json:"tax_cents"`
}

func Calculate(calc *core.Calculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		yearStr := q.Get("year")
		salaryStr := q.Get("salary")

		// missing params
		if yearStr == "" || salaryStr == "" {
			http.Error(w, "missing query params: year and salary are required", http.StatusBadRequest)
			return
		}

		// unsupported years
		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 2019 || year > 2022 {
			http.Error(w, "invalid or unsupported year (2019–2022)", http.StatusBadRequest)
			return
		}

		// invalid or negative salary
		salaryFloat, err := strconv.ParseFloat(salaryStr, 64)
		if err != nil {
			http.Error(w, "invalid salary format", http.StatusBadRequest)
			return
		}
		if salaryFloat < 0 {
			http.Error(w, "salary must be non-negative", http.StatusBadRequest)
			return
		}

		salaryCents := u.DollarsToCents(salaryFloat)
		calcResult, err := calc.CalculateCents(year, salaryCents)
		if err != nil {
			http.Error(w, "error calculating tax", http.StatusInternalServerError)
			return
		}

		resp := CalculateResponse{
			Year:          year,
			SalaryInput:   salaryStr,
			SalaryCents:   salaryCents,
			TotalTaxCents: calcResult.TotalTaxCents,
			EffectiveRate: u.BasisPointToRate(calcResult.EffectiveBps),
			PerBracket:    []PerBracketResponse{},
			Message:       "OK",
		}

		// build per bracket response
		for _, bracket := range calcResult.PerBracket {
			if bracket.TaxCents == 0 {
				continue
			}
			resp.PerBracket = append(resp.PerBracket, PerBracketResponse{
				MinCents:        bracket.MinCents,
				MaxCents:        bracket.MaxCents,
				RateBasisPoints: bracket.RateBps,
				TaxCents:        bracket.TaxCents,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
