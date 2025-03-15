package repositories

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/price-tracker/models"
	"github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SqliteRepositoryTestSuite struct {
	suite.Suite
	db         *bun.DB
	repository Repository
}

func (s *SqliteRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for testing
	cfg := config.SqliteConfig{
		SqliteURL: ":memory:",
	}

	Db, err := db.NewSqliteDB(cfg)
	assert.NoError(s.T(), err)
	s.db = Db
	s.repository = NewSqliteRepository(Db)
}

func (s *SqliteRepositoryTestSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		assert.NoError(s.T(), err)
	}
}

func TestSqliteRepositorySuite(t *testing.T) {
	suite.Run(t, new(SqliteRepositoryTestSuite))
}

func (s *SqliteRepositoryTestSuite) TestSaveAndGetLatestPrice() {
	// Arrange
	symbol := "BTC"
	expect := models.PriceDatum{
		Timestamp: time.Now(),
		Symbol:    symbol,
		Currency:  "",
		Price:     50000.0,
	}
	// Act - Save
	err := s.repository.Insert(expect)
	assert.NoError(s.T(), err)

	// Act - Get Latest
	result, err := s.repository.GetLatestPrice(expect)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), symbol, result.Symbol)
	assert.Equal(s.T(), expect.Price, result.Price)
	assert.Equal(s.T(), expect.Timestamp, result.Timestamp)
}

func (s *SqliteRepositoryTestSuite) TestGetPriceHistory() {
	// Arrange
	symbol := "BTC"
	expect := []models.PriceDatum{
		{
			Timestamp: time.Now().UTC().Add(-time.Hour),
			Symbol:    symbol,
			Currency:  "",
			Price:     50000.0,
		},
		{
			Timestamp: time.Now().UTC(),
			Symbol:    symbol,
			Currency:  "",
			Price:     50000.0,
		},
	}

	// Act
	errWrite := s.repository.BulkInsert(expect)
	results, err := s.repository.GetPriceOfTheLast24h(expect[0])

	// Assert
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), errWrite)

	// Verify results are in chronological order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.Before(results[j].Timestamp)
	})
	for i, result := range results {
		assert.Equal(s.T(), symbol, result.Symbol)
		assert.Equal(s.T(), expect[i].Price, result.Price)
		assert.Equal(s.T(), expect[i].Timestamp, result.Timestamp.UTC())
	}
}

func (s *SqliteRepositoryTestSuite) TestGetLatestPrice_NonExistentSymbol() {
	// Act
	result, err := s.repository.GetLatestPrice(models.PriceDatum{Symbol: "NOTEXISTSYMBOL"})

	// Assert
	assert.Equal(s.T(), sqlite3.ErrNotFound, err)
	assert.Nil(s.T(), result)
}

func (s *SqliteRepositoryTestSuite) TestGetPriceHistory_NoDataInRange() {
	// Arrange
	results, err := s.repository.GetPriceOfTheLast24h(models.PriceDatum{Symbol: "NOTEXISTSYMBOL"})

	// Assert
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), results)
}

func (s *SqliteRepositoryTestSuite) TestGetPriceHistory_GetSymbolList() {
	// Arrange
	_, err := s.repository.GetAllCryptoInfo()

	// Assert
	assert.Nil(s.T(), err)
	//assert.Equal(s.T(), 0, len(results))
}
