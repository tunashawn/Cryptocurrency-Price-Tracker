package repositories

import (
	"backend/price-tracker/models"
	"context"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
)

type SqliteRepositoryImpl struct {
	db *bun.DB
}

func NewSqliteRepository(db *bun.DB) Repository {
	return &SqliteRepositoryImpl{
		db: db,
	}
}

func (s *SqliteRepositoryImpl) Insert(datum models.PriceDatum) error {
	_, err := s.db.NewInsert().Model(&datum).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteRepositoryImpl) BulkInsert(data []models.PriceDatum) error {
	_, err := s.db.NewInsert().Model(&data).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteRepositoryImpl) GetLatestPrice(datum models.PriceDatum) (*models.PriceDatum, error) {
	exists, err := s.db.NewSelect().
		Model(&datum).
		Order("timestamp DESC").
		Limit(1).
		Exists(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if !exists {
		return nil, sqlite3.ErrNotFound
	}

	return &datum, nil
}

func (s *SqliteRepositoryImpl) GetPriceOfTheLast24h(req models.PriceDatum) ([]models.PriceDatum, error) {
	var res []models.PriceDatum

	query := `
			WITH Numbered AS (
				SELECT symbol, timestamp, price,
				ROW_NUMBER() OVER () AS rn,
				COUNT(*) OVER () AS total_count
			FROM crypto_price
			WHERE symbol = ?
			AND timestamp >= DATETIME('now', '-1 day')
			ORDER BY timestamp ASC
			)
			SELECT timestamp, symbol, price
			FROM Numbered
			WHERE (rn - 1) % (total_count / 30 + 1) = 0;`

	err := s.db.NewRaw(query, req.Symbol).
		Scan(context.Background(), &res)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return res, nil
}

func (s *SqliteRepositoryImpl) GetAllCryptoInfo() ([]models.PriceDatum, error) {
	res := make([]models.PriceDatum, 0)

	err := s.db.
		NewRaw("SELECT DISTINCT symbol FROM ?", bun.Ident("crypto_price")).
		Scan(context.Background(), &res)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return res, nil
}
