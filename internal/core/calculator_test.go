package core

import (
	"testing"

	"github.com/alexlangev/interview-submission/internal/models"
	u "github.com/alexlangev/interview-submission/internal/utils"
)

type fakeProvider struct {
	brackets []models.TaxBracket
	err      error
}

func (f *fakeProvider) GetTaxBrackets(year int) ([]models.TaxBracket, error) {
	return f.brackets, f.err
}

// 2022 brackets from the assignment readme
// max is pointer for nil value == no maximum
func brackets2022() []models.TaxBracket {
	max1 := 50197.0
	max2 := 100392.0
	max3 := 155625.0
	max4 := 221708.0

	return []models.TaxBracket{
		{Min: 0, Max: &max1, Rate: 0.15},
		{Min: 50197, Max: &max2, Rate: 0.205},
		{Min: 100392, Max: &max3, Rate: 0.26},
		{Min: 155625, Max: &max4, Rate: 0.29},
		{Min: 221708, Max: nil, Rate: 0.33},
	}
}

func TestCalculate_ZeroIncome(t *testing.T) {
	calc := NewCalculator(&fakeProvider{brackets: brackets2022()})

	res, err := calc.CalculateCents(2022, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalTaxCents != 0 {
		t.Fatalf("expected total tax 0, got %d", res.TotalTaxCents)
	}
	if res.EffectiveBps != 0 {
		t.Fatalf("expected effective 0 basis points, got %d", res.EffectiveBps)
	}
}

func TestCalculateCents_100kIncome_2022(t *testing.T) {
	calc := NewCalculator(&fakeProvider{brackets: brackets2022()})

	incomeCents := u.DollarsToCents(100000.00)
	res, err := calc.CalculateCents(2022, incomeCents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantTotal := u.DollarsToCents(17739.17)
	if res.TotalTaxCents != wantTotal {
		t.Fatalf("total tax: want %d, got %d", wantTotal, res.TotalTaxCents)
	}

	wantEffBps := int64(1774)
	if res.EffectiveBps != wantEffBps {
		t.Fatalf("effective bps: want %d, got %d", wantEffBps, res.EffectiveBps)
	}
}
