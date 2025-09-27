package providers

import "github.com/alexlangev/interview-submission/internal/models"

// Abstracts where the tax data comes from.
// Implementations may fetch a external API, a database, a file,
type TaxBracketProvider interface {
	GetTaxBrackets(year int) ([]models.TaxBracket, error)
}
