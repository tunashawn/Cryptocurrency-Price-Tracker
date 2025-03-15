package ws

import (
	"backend/price-tracker/models"
	"backend/price-tracker/repositories/mock"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BinanceWebsocketTestSuite struct {
	suite.Suite
	mockRepo *mock.MockRepository
	cache    *sync.Map
	ws       BinanceWebsocketImpl
}

func (s *BinanceWebsocketTestSuite) SetupTest() {
	s.mockRepo = new(mock.MockRepository)
	s.cache = &sync.Map{}
	s.ws = BinanceWebsocketImpl{
		cache:      s.cache,
		repository: s.mockRepo,
	}
}

func TestBinanceWebsocketSuite(t *testing.T) {
	suite.Run(t, new(BinanceWebsocketTestSuite))
}

func (s *BinanceWebsocketTestSuite) TestUpdateLatestPrice() {
	// Arrange
	expect := []models.PriceDatum{
		{
			Timestamp: time.Now().UTC(),
			Symbol:    "BTC",
			Currency:  "USDT",
			Price:     100.0,
		},
		{
			Timestamp: time.Now().UTC(),
			Symbol:    "ABC",
			Currency:  "USDT",
			Price:     200.0,
		},
	}
	priceTicker := models.BinancePriceTicker{
		Status: 200,
		BinanceResult: []models.BinanceResult{
			{
				Symbol: "BTCUSDT",
				Price:  "100.0",
			},
			{
				Symbol: "ABCUSDT",
				Price:  "200.0",
			},
		},
	}

	// Act
	s.mockRepo.On("BulkInsert", m.MatchedBy(func(data []models.PriceDatum) bool {
		return len(data) == 2 &&
			data[0].Symbol == "BTC" && data[0].Price == 100.0 &&
			data[1].Symbol == "ABC" && data[1].Price == 200.0
	})).Return(nil)

	s.ws.updateLatestPrice(priceTicker)

	// Assert
	if value, ok := s.cache.Load("BTC"); ok {
		assert.Equal(s.T(), value.(models.PriceDatum).Symbol, expect[0].Symbol)
		assert.Equal(s.T(), value.(models.PriceDatum).Price, expect[0].Price)
	}

	if value, ok := s.cache.Load("ABC"); ok {
		assert.Equal(s.T(), value.(models.PriceDatum).Symbol, expect[1].Symbol)
		assert.Equal(s.T(), value.(models.PriceDatum).Price, expect[1].Price)
	}
}
