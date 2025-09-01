package market

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Market factory for testing
type MarketFactory struct{}

func (f *MarketFactory) NewMarket(
	symbol string,
	enabled bool,
	averageValue float64,
	averageVolume float64,
) *entities.Market {
	return &entities.Market{
		ID:            uuid.New(),
		Symbol:        symbol,
		Enabled:       enabled,
		AverageValue:  averageValue,
		AverageVolume: averageVolume,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
}

// Market repository tests

func TestGetMarketByIDReturnsError(t *testing.T) {
	// Arrange
	id := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := marketRepository.GetByID(ctx, id)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetMarketByIDReturnsMarket(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketFactory := &MarketFactory{}
	market := marketFactory.NewMarket(
		"BTCUSDT",
		true,
		50000.0,
		1000000.0,
	)

	dto := dtos.Market{}
	dto.FromEntity(market)

	database.Create(&dto)

	// Act
	foundMarket, err := marketRepository.GetByID(ctx, market.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarket)
	assert.Equal(t, market.ID, foundMarket.ID)
	assert.Equal(t, market.Symbol, foundMarket.Symbol)
	assert.Equal(t, market.Enabled, foundMarket.Enabled)
	assert.Equal(t, market.AverageValue, foundMarket.AverageValue)
	assert.Equal(t, market.AverageVolume, foundMarket.AverageVolume)
}

func TestGetMarketBySymbolReturnsError(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := marketRepository.GetBySymbol(ctx, "NONEXISTENT")

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetMarketBySymbolReturnsMarket(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketFactory := &MarketFactory{}
	market := marketFactory.NewMarket(
		"ETHUSDT",
		false,
		3000.0,
		500000.0,
	)

	dto := dtos.Market{}
	dto.FromEntity(market)

	database.Create(&dto)

	// Act
	foundMarket, err := marketRepository.GetBySymbol(ctx, market.Symbol)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarket)
	assert.Equal(t, market.ID, foundMarket.ID)
	assert.Equal(t, market.Symbol, foundMarket.Symbol)
	assert.Equal(t, market.Enabled, foundMarket.Enabled)
	assert.Equal(t, market.AverageValue, foundMarket.AverageValue)
	assert.Equal(t, market.AverageVolume, foundMarket.AverageVolume)
}

func TestGetAllMarketsReturnsMarkets(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketFactory := &MarketFactory{}
	market1 := marketFactory.NewMarket(
		"LTCUSDT",
		true,
		200.0,
		200000.0,
	)
	market2 := marketFactory.NewMarket(
		"LTCUSDT",
		false,
		250.0,
		300000.0,
	)

	dto1 := dtos.Market{}
	dto1.FromEntity(market1)
	database.Create(&dto1)

	dto2 := dtos.Market{}
	dto2.FromEntity(market2)
	database.Create(&dto2)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "LTCUSDT",
	}, "created_at", "desc", 0, 10)
	foundMarkets, err := marketRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarkets)
	assert.Equal(t, 2, len(*foundMarkets))
	assert.Equal(t, market1.ID, (*foundMarkets)[0].ID)
	assert.Equal(t, market2.ID, (*foundMarkets)[1].ID)
}

func TestGetAllMarketsWithUnmatchingSymbolReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "XRPUSDT",
	}, "created_at", "desc", 0, 10)
	foundMarkets, err := marketRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarkets)
	assert.Equal(t, 0, len(*foundMarkets))
}

func TestCreateMarketReturnsMarket(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketFactory := &MarketFactory{}
	market := marketFactory.NewMarket(
		"CHRUSDT",
		true,
		0.5,
		50000.0,
	)

	// Act
	createdMarket, err := marketRepository.Create(ctx, market)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdMarket)
	assert.Equal(t, market.ID, createdMarket.ID)
	assert.Equal(t, market.Symbol, createdMarket.Symbol)
	assert.Equal(t, market.Enabled, createdMarket.Enabled)
	assert.Equal(t, market.AverageValue, createdMarket.AverageValue)
	assert.Equal(t, market.AverageVolume, createdMarket.AverageVolume)

	// Assert that the market was created in the database
	foundMarket, err := marketRepository.GetByID(ctx, market.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundMarket)
	assert.Equal(t, market.ID, foundMarket.ID)
}

func TestUpdateMarketReturnsMarket(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketFactory := &MarketFactory{}
	market := marketFactory.NewMarket(
		"DOGEUSDT",
		false,
		0.1,
		100000.0,
	)

	dto := dtos.Market{}
	dto.FromEntity(market)
	database.Create(&dto)

	// Update market values
	market.Enabled = true
	market.AverageValue = 0.15
	market.AverageVolume = 200000.0

	// Act
	updatedMarket, err := marketRepository.Update(ctx, market)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedMarket)
	assert.Equal(t, market.ID, updatedMarket.ID)
	assert.Equal(t, market.Symbol, updatedMarket.Symbol)
	assert.Equal(t, market.Enabled, updatedMarket.Enabled)
	assert.Equal(t, market.AverageValue, updatedMarket.AverageValue)
	assert.Equal(t, market.AverageVolume, updatedMarket.AverageVolume)

	// Verify the update was persisted
	foundMarket, err := marketRepository.GetByID(ctx, market.ID)
	assert.NoError(t, err)
	assert.True(t, foundMarket.Enabled)
	assert.Equal(t, 0.15, foundMarket.AverageValue)
	assert.Equal(t, 200000.0, foundMarket.AverageVolume)
}

func TestDeleteMarketReturnsErrorIfIDDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	err := marketRepository.Delete(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteMarketDeletesMarket(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketFactory := &MarketFactory{}
	market := marketFactory.NewMarket(
		"DOGEUSDT",
		true,
		0.1,
		100000.0,
	)

	dto := dtos.Market{}
	dto.FromEntity(market)
	database.Create(&dto)

	// Act
	err := marketRepository.Delete(ctx, market.ID)

	// Assert
	assert.NoError(t, err)

	// Verify the market was deleted
	foundMarket, err := marketRepository.GetByID(ctx, market.ID)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, foundMarket)
}

// Market data repository tests

func TestGetMarketDataByIDReturnsError(t *testing.T) {
	// Arrange
	id := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := marketDataRepository.GetByID(ctx, id)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetMarketDataByIDReturnsMarketData(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketDataFactory := &entities.MarketDataFactory{}
	value := float64(10000.1)
	marketData := marketDataFactory.NewMarketData(
		uuid.New(),
		"BTCUSDT",
		time.Now().UTC(),
		10000.1,
		10000.2,
		10000.3,
		10000.4,
		1000,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
	)

	dto := dtos.MarketData{}
	dto.FromEntity(marketData)

	database.Create(&dto)

	// Act
	foundMarketData, err := marketDataRepository.GetByID(ctx, marketData.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarketData)
	assert.Equal(t, marketData.ID, foundMarketData.ID)
}

func TestGetAllMarketDataReturnsMarketData(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	marketDataFactory := &entities.MarketDataFactory{}
	value := float64(10000.1)
	marketData1 := marketDataFactory.NewMarketData(
		uuid.New(),
		"LTCUSDT",
		time.Now().UTC(),
		10000.1,
		10000.2,
		10000.3,
		10000.4,
		1000,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
	)
	marketData2 := marketDataFactory.NewMarketData(
		uuid.New(),
		"LTCUSDT",
		time.Now().UTC(),
		10000.1,
		10000.2,
		10000.3,
		10000.4,
		1000,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
	)

	dto1 := dtos.MarketData{}
	dto1.FromEntity(marketData1)
	database.Create(&dto1)

	dto2 := dtos.MarketData{}
	dto2.FromEntity(marketData2)
	database.Create(&dto2)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "LTCUSDT",
	}, "created_at", "desc", 0, 10)
	foundMarketData, err := marketDataRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarketData)
	assert.Equal(t, 2, len(*foundMarketData))
	assert.Equal(t, marketData1.ID, (*foundMarketData)[0].ID)
	assert.Equal(t, marketData2.ID, (*foundMarketData)[1].ID)
}

func TestGetAllMarketDataWithUnmatchingSymbolReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "XRPUSDT",
	}, "created_at", "desc", 0, 10)
	foundMarketData, err := marketDataRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundMarketData)
	assert.Equal(t, 0, len(*foundMarketData))
}

func TestCreateMarketDataReturnsMarketData(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	value := float64(10000.1)

	marketDataFactory := &entities.MarketDataFactory{}
	marketData := marketDataFactory.NewMarketData(
		uuid.New(),
		"CHRUSDT",
		time.Now().UTC(),
		10000.1,
		10000.2,
		10000.3,
		10000.4,
		1000,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
	)

	dto := dtos.MarketData{}
	dto.FromEntity(marketData)

	// Act
	createdMarketData, err := marketDataRepository.Create(ctx, marketData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdMarketData)
	assert.Equal(t, marketData.ID, createdMarketData.ID)

	// Assert that the market data was created in the database
	md, err := marketDataRepository.GetByID(ctx, marketData.ID)
	assert.NoError(t, err)
	assert.NotNil(t, md)
	assert.Equal(t, marketData.ID, md.ID)
}

func TestUpdateMarketDataReturnsMarketData(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	value := float64(10000.1)

	marketDataFactory := &entities.MarketDataFactory{}
	marketData := marketDataFactory.NewMarketData(
		uuid.New(),
		"DOGEUSDT",
		time.Now().UTC(),
		10000.1,
		10000.2,
		10000.3,
		10000.4,
		1000,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
	)

	dto := dtos.MarketData{}
	dto.FromEntity(marketData)

	database.Create(&dto)

	marketData.Close = 10.0

	// Act
	updatedMarketData, err := marketDataRepository.Update(ctx, marketData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedMarketData)
	assert.Equal(t, marketData.ID, updatedMarketData.ID)
	assert.Equal(t, marketData.Close, updatedMarketData.Close)
}

func TestDeleteMarketDataReturnsErrorIfIDDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	err := marketDataRepository.Delete(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteMarketDataDeletesMarketData(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	value := float64(10000.1)

	marketDataFactory := &entities.MarketDataFactory{}
	marketData := marketDataFactory.NewMarketData(
		uuid.New(),
		"DOGEUSDT",
		time.Now().UTC(),
		10000.1,
		10000.2,
		10000.3,
		10000.4,
		1000,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
		&value,
	)

	dto := dtos.MarketData{}
	dto.FromEntity(marketData)

	database.Create(&dto)

	// Act
	err := marketDataRepository.Delete(ctx, marketData.ID)

	// Assert
	assert.NoError(t, err)

	foundMarketData, err := marketDataRepository.GetByID(ctx, marketData.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, foundMarketData)
}
