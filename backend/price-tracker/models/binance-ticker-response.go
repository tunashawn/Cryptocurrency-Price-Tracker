package models

import "backend/internal/constant"

type BinancePriceTicker struct {
	Status        int             `json:"status"`
	BinanceResult []BinanceResult `json:"result"`
}

type BinanceResult struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (b *BinanceResult) IsCurrencyUSDT() bool {
	return b.Symbol[len(b.Symbol)-4:] == constant.USDT
}
