package models

import (
	"fmt"
	"github.com/uptrace/bun"
	"strconv"
	"time"
)

type PriceDatum struct {
	bun.BaseModel `json:"-" bun:"table:crypto_price"`
	Timestamp     time.Time `json:"timestamp,omitzero" bun:"timestamp"`
	Symbol        string    `json:"symbol,omitempty" bun:"symbol"`
	Currency      string    `json:"currency,omitempty" bun:"currency"`
	Price         float64   `json:"price,omitempty" bun:"price"`
}

func NewPriceDatumFromBinanceResult(res BinanceResult, timestamp time.Time) (PriceDatum, error) {
	price, err := strconv.ParseFloat(res.Price, 64)
	if err != nil {
		return PriceDatum{}, fmt.Errorf("parse float: %w", err)
	}
	return PriceDatum{
		Timestamp: timestamp,
		Symbol:    res.Symbol[:len(res.Symbol)-4],
		Currency:  res.Symbol[len(res.Symbol)-4:],
		Price:     price,
	}, nil
}
