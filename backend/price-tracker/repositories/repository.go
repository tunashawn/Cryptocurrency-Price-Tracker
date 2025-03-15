package repositories

import "backend/price-tracker/models"

type Repository interface {
	// Insert writes a single datum
	Insert(datum models.PriceDatum) error
	// BulkInsert writes a bulket of data
	BulkInsert(data []models.PriceDatum) error
	// GetLatestPrice returns the latest price of a symbol
	GetLatestPrice(req models.PriceDatum) (*models.PriceDatum, error)
	// GetPriceOfTheLast24h returns all prices of a symbol in the last 24h
	GetPriceOfTheLast24h(req models.PriceDatum) ([]models.PriceDatum, error)
	// GetSymbolList returns a list of all symbols
	GetSymbolList() ([]string, error)
}
