package models

type TaxBracket struct {
	Min  float64  `json:"min"`
	Max  *float64 `json:"max,omitempty"`
	Rate float64  `json:"rate"`
}

type TaxResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}
