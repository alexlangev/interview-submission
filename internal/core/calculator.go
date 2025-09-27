package core

import (
	"github.com/alexlangev/interview-submission/internal/providers"
	u "github.com/alexlangev/interview-submission/internal/utils"
)

type Calculator struct {
	provider providers.TaxBracketProvider
}

func NewCalculator(p providers.TaxBracketProvider) *Calculator {
	return &Calculator{provider: p}
}

type BracketResultCents struct {
	MinCents int64
	MaxCents *int64
	RateBps  int64
	TaxCents int64
}

type ResultCents struct {
	IncomeCents   int64
	Year          int
	TotalTaxCents int64
	EffectiveBps  int64
	PerBracket    []BracketResultCents
}

func (c *Calculator) CalculateCents(year int, incomeCents int64) (ResultCents, error) {
	brackets, err := c.provider.GetTaxBrackets(year)
	if err != nil {
		return ResultCents{}, err
	}

	perBracket := []BracketResultCents{}
	total := int64(0)

	for _, b := range brackets {
		min := u.DollarsToCents(b.Min)

		var maxPtr *int64
		var upper int64
		if b.Max != nil {
			m := u.DollarsToCents(*b.Max)
			maxPtr = &m
			upper = u.MinInt64(incomeCents, m)
		} else {
			upper = incomeCents
		}

		taxable := upper - min
		if taxable < 0 {
			taxable = 0
		}

		bps := u.RateToBasisPoint(b.Rate)
		tax := u.DivRoundHalfUp(taxable*bps, 10000)

		total += tax
		perBracket = append(perBracket, BracketResultCents{
			MinCents: min,
			MaxCents: maxPtr,
			RateBps:  bps,
			TaxCents: tax,
		})
	}

	var effBps int64
	if incomeCents > 0 {
		effBps = u.DivRoundHalfUp(total*10000, incomeCents)
	}

	return ResultCents{
		IncomeCents:   incomeCents,
		Year:          year,
		TotalTaxCents: total,
		EffectiveBps:  effBps,
		PerBracket:    perBracket,
	}, nil
}
