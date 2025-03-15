package services

import (
	"backend/price-tracker/models"
	"backend/price-tracker/repositories"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type PriceTrackingService interface {
	// GetLatestPrice returns the latest price of a symbol.
	// First, it tries to fetch from cache.
	// Second, it tries to fetch from API.
	// Finally, it tries to fetch from database
	GetLatestPrice(req models.PriceDatum) (*models.PriceDatum, error)
	// GetPriceOfTheLast24h returns all prices of a symbol in the last 24h
	GetPriceOfTheLast24h(req models.PriceDatum) ([]models.PriceDatum, error)
	// IsSymbolValid returns true if them symbol is valid and exist
	IsSymbolValid(symbol string) bool
}

type PriceTrackingServiceImpl struct {
	cache      *sync.Map
	repository repositories.Repository
	httpClient *http.Client
}

func NewPriceTrackingService(cache *sync.Map, repository repositories.Repository) PriceTrackingService {
	return &PriceTrackingServiceImpl{
		cache:      cache,
		repository: repository,
		httpClient: &http.Client{
			Timeout: time.Second * 2,
		},
	}
}

// GetLatestPrice returns the latest price of a symbol.
// First, it tries to fetch from cache. Next if not found
// Second, it tries to fetch from binance API. Next if not found
// Finally, it tries to fetch from database
func (p *PriceTrackingServiceImpl) GetLatestPrice(req models.PriceDatum) (*models.PriceDatum, error) {
	if val, ok := p.cache.Load(req.Symbol); ok {
		res := val.(models.PriceDatum)
		return &res, nil
	} else if res, err := p.fetchPriceFromBinanceAPI(req); err == nil {
		return res, nil
	} else {
		slog.Error("fetchPriceFromBinanceAPI", "err", err)

		res, err := p.repository.GetLatestPrice(req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

// fetchPriceFromBinanceAPI returns the latest price of a symbol from Binance's REST API
func (p *PriceTrackingServiceImpl) fetchPriceFromBinanceAPI(req models.PriceDatum) (*models.PriceDatum, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", req.Symbol+req.Currency)

	resp, err := p.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	var binanceResult models.BinanceResult
	err = json.NewDecoder(resp.Body).Decode(&binanceResult)
	if err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	res, err := models.NewPriceDatumFromBinanceResult(binanceResult, time.Now())
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &res, nil
}

// GetPriceOfTheLast24h returns all prices of a symbol in the last 24h
func (p *PriceTrackingServiceImpl) GetPriceOfTheLast24h(req models.PriceDatum) ([]models.PriceDatum, error) {
	res, err := p.repository.GetPriceOfTheLast24h(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// IsSymbolValid returns true if symbol is valid and exist
func (p *PriceTrackingServiceImpl) IsSymbolValid(symbol string) bool {
	if _, ok := p.cache.Load(symbol); ok {
		return true
	}

	return false
}
