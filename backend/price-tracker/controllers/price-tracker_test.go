package controllers

import (
	gin_test_setup "backend/internal/gin-test-setup"
	"backend/internal/response"
	"backend/price-tracker/models"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPriceTrackingService implements services.PriceTrackingService
type MockPriceTrackingService struct {
	mock.Mock
}

func (m *MockPriceTrackingService) GetCryptoList() ([]models.PriceDatum, error) {
	args := m.Called()
	return args.Get(0).([]models.PriceDatum), args.Error(1)
}

func (m *MockPriceTrackingService) IsSymbolValid(symbol string) bool {
	args := m.Called(symbol)
	return args.Bool(0)
}

func (m *MockPriceTrackingService) GetLatestPrice(req models.PriceDatum) (*models.PriceDatum, error) {
	args := m.Called(req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.PriceDatum), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPriceTrackingService) GetPriceOfTheLast24h(req models.PriceDatum) ([]models.PriceDatum, error) {
	args := m.Called(req)
	if args.Get(0) != nil {
		return args.Get(0).([]models.PriceDatum), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestPriceTrackerControllerImpl_GetLatestPrice(t *testing.T) {
	type TestResponseData struct {
		Meta response.Meta      `json:"meta"`
		Data *models.PriceDatum `json:"data"`
	}
	type args struct {
		path  string
		valid bool
	}
	type want struct {
		res     TestResponseData
		mockErr error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "missing symbol",
			args: args{
				path:  "/?symbol=&currency=USD",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    400,
						Message: "missing symbol value ''",
					},
					Data: nil,
				},
			},
		},
		{
			name: "invalid symbol",
			args: args{
				path:  "/?symbol=HUHU&currency=USD",
				valid: false,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    400,
						Message: "invalid symbol value 'HUHU'",
					},
					Data: nil,
				},
			},
		},
		{
			name: "success",
			args: args{
				path:  "/?symbol=BTC&currency=USDT",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    200,
						Message: "ok",
					},
					Data: &models.PriceDatum{
						Symbol:   "BTC",
						Currency: "USDT",
						Price:    1,
					},
				},
			},
		},
		{
			name: "internal error",
			args: args{
				path:  "/?symbol=BTC&currency=USDT",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    500,
						Message: "error",
					},
					//Data:
				},
				mockErr: errors.New("error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := gin_test_setup.NewGinTestContext("GET", tt.args.path)

			mockPriceTrackingService := &MockPriceTrackingService{}

			p := &PriceTrackerControllerImpl{
				priceTrackingService: mockPriceTrackingService,
			}

			mockPriceTrackingService.On("GetLatestPrice", mock.Anything).Return(tt.want.res.Data, tt.want.mockErr)
			mockPriceTrackingService.On("IsSymbolValid", mock.Anything).Return(tt.args.valid)
			p.GetLatestPrice(ctx.Context)

			var get TestResponseData
			json.Unmarshal([]byte(ctx.GetResponseBody()), &get)

			assert.Equal(t, tt.want.res.Meta, get.Meta)
			assert.Equal(t, tt.want.res.Data, get.Data)

		})
	}
}

func TestNewPriceTrackerController(t *testing.T) {
	mockPriceTrackingService := &MockPriceTrackingService{}
	NewPriceTrackerController(mockPriceTrackingService)
}

func TestPriceTrackerControllerImpl_GetPriceHistory(t *testing.T) {
	type TestResponseData struct {
		Meta response.Meta       `json:"meta"`
		Data []models.PriceDatum `json:"data"`
	}
	type args struct {
		path  string
		valid bool
	}
	type want struct {
		res     TestResponseData
		mockErr error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "missing symbol",
			args: args{
				path:  "/?symbol=&currency=USD&interval=",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    400,
						Message: "missing symbol value ''",
					},
					Data: nil,
				},
			},
		},
		{
			name: "invalid symbol",
			args: args{
				path:  "/?symbol=AAA&currency=USD&interval=",
				valid: false,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    400,
						Message: "invalid symbol value 'AAA'",
					},
					Data: nil,
				},
			},
		},
		{
			name: "wrong/missing interval",
			args: args{
				path:  "/?symbol=BTC&currency=USDT&interval=",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    400,
						Message: "interval only support '24' hours",
					},
					Data: nil,
				},
			},
		},
		{
			name: "success",
			args: args{
				path:  "/?symbol=BTC&currency=USDT&interval=24",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    200,
						Message: "ok",
					},
					Data: []models.PriceDatum{
						{
							Symbol:   "BTC",
							Currency: "USDT",
							Price:    1,
						},
					},
				},
			},
		},
		{
			name: "internal error",
			args: args{
				path:  "/?symbol=BTC&currency=USDT&interval=24",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    500,
						Message: "error",
					},
					//Data:
				},
				mockErr: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := gin_test_setup.NewGinTestContext("GET", tt.args.path)

			mockPriceTrackingService := &MockPriceTrackingService{}

			p := &PriceTrackerControllerImpl{
				priceTrackingService: mockPriceTrackingService,
			}

			mockPriceTrackingService.On("GetPriceOfTheLast24h", mock.Anything).Return(tt.want.res.Data, tt.want.mockErr)
			mockPriceTrackingService.On("IsSymbolValid", mock.Anything).Return(tt.args.valid)

			p.GetPriceHistory(ctx.Context)

			var get TestResponseData
			json.Unmarshal([]byte(ctx.GetResponseBody()), &get)

			assert.Equal(t, tt.want.res.Meta, get.Meta)
			assert.Equal(t, tt.want.res.Data, get.Data)
		})
	}
}

func TestPriceTrackerControllerImpl_GetCryptoList(t *testing.T) {
	type TestResponseData struct {
		Meta response.Meta       `json:"meta"`
		Data []models.PriceDatum `json:"data"`
	}
	type args struct {
		path  string
		valid bool
	}
	type want struct {
		res     TestResponseData
		mockErr error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success",
			args: args{
				path:  "/",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    200,
						Message: "ok",
					},
					Data: []models.PriceDatum{
						{
							Symbol: "ABC",
						},
						{
							Symbol: "BCD",
						},
					},
				},
				mockErr: nil,
			},
		},
		{
			name: "internal error",
			args: args{
				path:  "/",
				valid: true,
			},
			want: want{
				res: TestResponseData{
					Meta: response.Meta{
						Code:    500,
						Message: "error",
					},
					Data: nil,
				},
				mockErr: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := gin_test_setup.NewGinTestContext("GET", tt.args.path)

			mockPriceTrackingService := &MockPriceTrackingService{}

			p := &PriceTrackerControllerImpl{
				priceTrackingService: mockPriceTrackingService,
			}

			mockPriceTrackingService.On("GetCryptoList", mock.Anything).Return(tt.want.res.Data, tt.want.mockErr)

			p.GetCryptoList(ctx.Context)

			var get TestResponseData
			json.Unmarshal([]byte(ctx.GetResponseBody()), &get)

			assert.Equal(t, tt.want.res.Meta, get.Meta)
			assert.Equal(t, tt.want.res.Data, get.Data)
		})
	}
}
