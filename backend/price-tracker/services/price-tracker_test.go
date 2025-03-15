package services

import (
	"backend/price-tracker/models"
	"backend/price-tracker/repositories/mock"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PriceTrackingServiceTestSuite struct {
	suite.Suite
	mockRepo   *mock.MockRepository
	cache      *sync.Map
	service    PriceTrackingService
	httpClient *http.Client
}

func (s *PriceTrackingServiceTestSuite) SetupTest() {
	s.mockRepo = new(mock.MockRepository)
	s.cache = &sync.Map{}
	s.httpClient = &http.Client{}
	s.service = NewPriceTrackingService(s.cache, s.mockRepo)
}

func TestPriceTrackingServiceSuite(t *testing.T) {
	suite.Run(t, new(PriceTrackingServiceTestSuite))
}

func (s *PriceTrackingServiceTestSuite) TestGetLatestPrice() {
	// Arrange
	symbol := "BTC"
	expectedPrice := models.PriceDatum{
		Symbol:   "BTC",
		Currency: "USDT",
		Price:    50000.0,
	}
	s.cache.Store(symbol, expectedPrice)

	// Act
	price, err := s.service.GetLatestPrice(expectedPrice)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedPrice.Price, price.Price)
}

func (s *PriceTrackingServiceTestSuite) TestGetPriceHistory() {
	// Arrange
	symbol := "BTC"
	now := time.Now().UTC()
	from := now.Add(-1 * time.Hour)
	to := now
	req := models.PriceDatum{
		Symbol:   symbol,
		Currency: "USDT",
	}
	expectedPrices := []models.PriceDatum{
		{
			Symbol:    symbol,
			Currency:  "USDT",
			Price:     50000.0,
			Timestamp: from.Add(30 * time.Minute),
		},
		{
			Symbol:    symbol,
			Currency:  "USDT",
			Price:     51000.0,
			Timestamp: to,
		},
	}
	s.mockRepo.On("GetPriceOfTheLast24h", req).Return(expectedPrices, nil)

	// Act
	prices, err := s.service.GetPriceOfTheLast24h(req)

	// Assert
	assert.NoError(s.T(), err)
	for i, price := range prices {
		assert.Equal(s.T(), expectedPrices[i].Symbol, price.Symbol)
		assert.Equal(s.T(), expectedPrices[i].Currency, price.Currency)
		assert.Equal(s.T(), expectedPrices[i].Price, price.Price)
		assert.Equal(s.T(), expectedPrices[i].Timestamp.UTC(), price.Timestamp.UTC())
	}
	s.mockRepo.AssertExpectations(s.T())
}

func TestPriceTrackingServiceImpl_fetchPriceFromBinanceAPI(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name           string
		request        models.PriceDatum
		mockResponse   string
		mockStatusCode int
		mockError      error
		expectedResult *models.PriceDatum
		expectedError  bool
	}{
		{
			name: "successful response",
			request: models.PriceDatum{
				Symbol:   "BTC",
				Currency: "USDT",
			},
			mockResponse:   `{"symbol": "BTCUSDT", "price": "64235.12000000"}`,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectedResult: &models.PriceDatum{
				Symbol:    "BTC",
				Currency:  "USDT",
				Price:     64235.12,
				Timestamp: time.Time{}, // Will be compared separately
			},
			expectedError: false,
		},
		{
			name: "http error",
			request: models.PriceDatum{
				Symbol:   "BTC",
				Currency: "USDT",
			},
			mockResponse:   "",
			mockStatusCode: 0,
			mockError:      errors.New("connection refused"),
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "invalid json response",
			request: models.PriceDatum{
				Symbol:   "BTC",
				Currency: "USDT",
			},
			mockResponse:   `{invalid json`,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectedResult: nil,
			expectedError:  true,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock transport
			mockTransport := &mockRoundTripper{
				response: &http.Response{
					StatusCode: tc.mockStatusCode,
					Body:       io.NopCloser(strings.NewReader(tc.mockResponse)),
				},
				err: tc.mockError,
			}

			// Create a client with the mock transport
			mockClient := &http.Client{
				Transport: mockTransport,
			}

			// Create the service with the mock client
			service := &PriceTrackingServiceImpl{
				httpClient: mockClient,
			}

			// Call the method being tested
			result, err := service.fetchPriceFromBinanceAPI(tc.request)

			// Check error
			if tc.expectedError && err == nil {
				t.Error("Expected an error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check result
			if tc.expectedResult == nil && result != nil {
				t.Error("Expected nil result but got a value")
			}
			if tc.expectedResult != nil {
				if result == nil {
					t.Error("Expected a result but got nil")
				} else {
					// Compare relevant fields
					if result.Symbol != tc.expectedResult.Symbol {
						t.Errorf("Expected Symbol %s, got %s", tc.expectedResult.Symbol, result.Symbol)
					}
					if result.Currency != tc.expectedResult.Currency {
						t.Errorf("Expected Currency %s, got %s", tc.expectedResult.Currency, result.Currency)
					}
					if result.Price != tc.expectedResult.Price {
						t.Errorf("Expected Price %f, got %f", tc.expectedResult.Price, result.Price)
					}
					// Check that timestamp is recent
					timeNow := time.Now()
					if result.Timestamp.Sub(timeNow) > 2*time.Second {
						t.Errorf("Expected recent timestamp, got %v", result.Timestamp)
					}
				}
			}

			// Verify that the correct URL was called
			expectedURL := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s",
				tc.request.Symbol+tc.request.Currency)
			if mockTransport.requestURL != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, mockTransport.requestURL)
			}
		})
	}
}

// Define the mock transport
type mockRoundTripper struct {
	response   *http.Response
	err        error
	requestURL string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.requestURL = req.URL.String()
	return m.response, m.err
}

func TestPriceTrackingServiceImpl_IsSymbolValid(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid",
			args: args{symbol: "ABC"},
			want: true,
		},
		{
			name: "invalid",
			args: args{symbol: "BCD"},
			want: false,
		},
	}

	m := sync.Map{}
	m.Store("ABC", true)

	p := &PriceTrackingServiceImpl{
		cache: &m,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, p.IsSymbolValid(tt.args.symbol), "IsSymbolValid(%v)", tt.args.symbol)
		})
	}
}
