package mock

import (
	"backend/price-tracker/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetSymbolList() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) Insert(datum models.PriceDatum) error {
	args := m.Called(datum)
	return args.Error(0)
}

func (m *MockRepository) GetLatestPrice(req models.PriceDatum) (*models.PriceDatum, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PriceDatum), args.Error(1)
}

func (m *MockRepository) GetPriceOfTheLast24h(req models.PriceDatum) ([]models.PriceDatum, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PriceDatum), args.Error(1)
}

func (m *MockRepository) BulkInsert(data []models.PriceDatum) error {
	args := m.Called(data)
	return args.Error(0)
}
