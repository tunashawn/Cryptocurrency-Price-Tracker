package ws

import (
	"backend/price-tracker/models"
	"backend/price-tracker/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	"time"
)

type BinanceWebsocketImpl struct {
	url           string
	cache         *sync.Map
	repository    repositories.Repository
	fetchInterval time.Duration
}

func NewBinanceWebsocket(
	cache *sync.Map,
	url string,
	fetchInterval int64,
	repository repositories.Repository,
) WebSocketFetcher {
	return &BinanceWebsocketImpl{
		cache:         cache,
		url:           url,
		repository:    repository,
		fetchInterval: time.Duration(fetchInterval),
	}
}

// Connect returns the connection to websocket
func (b *BinanceWebsocketImpl) Connect() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(b.url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	slog.Info("connected to binance websocket")

	return conn, nil
}

// Fetch keeps running every few seconds to retrieve latest price from given connection
func (b *BinanceWebsocketImpl) Fetch(conn *websocket.Conn) error {
	defer conn.Close()

	slog.Info("fetching price from binance")

	for {
		// Request latest price
		err := conn.WriteJSON(map[string]interface{}{
			"id":     "043a7cf2-bde3-4888-9604-c8ac41fcba4d",
			"method": "ticker.price",
		})
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}

		// Read the response
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		// Parse response
		var priceTicker models.BinancePriceTicker
		if err := json.Unmarshal(message, &priceTicker); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		if priceTicker.Status != 200 {
			return errors.New("price not available")
		}

		// Update latest price to cache and db
		go b.updateLatestPrice(priceTicker)

		time.Sleep(b.fetchInterval * time.Second)
	}
}

func (b *BinanceWebsocketImpl) updateLatestPrice(priceTicker models.BinancePriceTicker) {
	var bulk []models.PriceDatum

	timestamp := time.Now() // set fixed timestamp for this batch

	for _, res := range priceTicker.BinanceResult {
		if res.IsCurrencyUSDT() {

			datum, err := models.NewPriceDatumFromBinanceResult(res, timestamp)
			if err != nil {
				slog.Error("create price datum from binance", "err", err)
				continue
			}

			b.cache.Store(datum.Symbol, datum)

			bulk = append(bulk, datum)
		}
	}

	if len(bulk) > 0 {
		err := b.repository.BulkInsert(bulk)
		if err != nil {
			slog.Error("bulk insert", "err", err)
		}

		slog.Info("bulk insert", "count", len(bulk))
	}
}
