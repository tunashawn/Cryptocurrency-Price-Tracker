package db

import (
	"backend/internal/config"
	"backend/price-tracker/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewSqliteDB(cfg config.SqliteConfig) (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, cfg.SqliteURL)
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	err = createNewTableIfNotExist(db)
	if err != nil {
		return nil, fmt.Errorf("init table: %w", err)
	}

	err = setIndexes(db)
	if err != nil {
		return nil, fmt.Errorf("init indexes: %w", err)
	}

	err = setMaxDbSize(db)
	if err != nil {
		return nil, fmt.Errorf("init max db size: %w", err)
	}

	return db, nil
}

func createNewTableIfNotExist(db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model(&models.PriceDatum{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	return nil
}

func setIndexes(db *bun.DB) error {
	_, err := db.NewCreateIndex().
		Model(&models.PriceDatum{}).
		Index("idx_price_data_symbol_currency_timestamp").
		Column("symbol", "currency", "timestamp").
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	return nil
}

func setMaxDbSize(db *bun.DB) error {
	_, err := db.Exec("PRAGMA max_page_count = 262144;")
	if err != nil {
		return err
	}
	return nil
}
