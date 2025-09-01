package trade

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TradePreferenceRepository tests

func TestGetTradingPreferenceByIDReturnsError(t *testing.T) {
	// Arrange
	id := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := tradePreferenceRepository.GetByID(ctx, id)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetTradingPreferenceByIDReturnsTradingPreference(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	tradingPreferenceFactory := &entities.TradingPreferenceFactory{}
	tradingPreference := tradingPreferenceFactory.NewTradingPreference(
		uuid.New(),
		"test-algorithm",
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)

	dto := dtos.TradingPreference{}
	dto.FromEntity(tradingPreference)

	database.Create(&dto)

	// Act
	foundTradingPreference, err := tradePreferenceRepository.GetByID(ctx, tradingPreference.ID)

	fmt.Println(foundTradingPreference)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundTradingPreference)
	assert.Equal(t, tradingPreference.ID, foundTradingPreference.ID)
}

func TestGetTradingPreferenceByUserIDReturnsError(t *testing.T) {
	// Arrange
	userID := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := tradePreferenceRepository.GetByUserID(ctx, userID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetTradingPreferenceByUserIDReturnsTradingPreference(t *testing.T) {
	// Arrange
	userID := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	tradingPreferenceFactory := &entities.TradingPreferenceFactory{}
	tradingPreference := tradingPreferenceFactory.NewTradingPreference(
		userID,
		"test-algorithm",
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)

	dto := dtos.TradingPreference{}
	dto.FromEntity(tradingPreference)

	database.Create(&dto)

	// Act
	foundTradingPreference, err := tradePreferenceRepository.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundTradingPreference)
	assert.Equal(t, dto.ID, foundTradingPreference.ID)
}
func TestGetAllTradingPreferencesWithNotMatchingFiltersReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"risk_level": "something_wrong",
	}, "created_at", "desc", 0, 10)
	tradingPreferences, err := tradePreferenceRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tradingPreferences)
	assert.Equal(t, 0, len(*tradingPreferences))
}

func TestGetAllTradingPreferencesReturnsTradingPreferences(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	tradingPreferenceFactory := &entities.TradingPreferenceFactory{}
	tradingPreference := tradingPreferenceFactory.NewTradingPreference(
		uuid.New(),
		"fake-algorithm",
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)

	dto := dtos.TradingPreference{}
	dto.FromEntity(tradingPreference)

	database.Create(&dto)

	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"algorithm": "fake-algorithm",
	}, "created_at", "desc", 0, 10)

	tradingPreferences, err := tradePreferenceRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tradingPreferences)
	assert.Equal(t, 1, len(*tradingPreferences))
	assert.Equal(t, dto.ID, (*tradingPreferences)[0].ID)
}

func TestCreateTradingPreferenceReturnsTradingPreference(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	tradingPreferenceFactory := &entities.TradingPreferenceFactory{}
	tradingPreference := tradingPreferenceFactory.NewTradingPreference(
		uuid.New(),
		"swing_trading",
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)

	// Act
	createdTradingPreference, err := tradePreferenceRepository.Create(ctx, tradingPreference)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdTradingPreference)
	assert.Equal(t, tradingPreference.ID, createdTradingPreference.ID)

	foundTradingPreference, err := tradePreferenceRepository.GetByID(ctx, createdTradingPreference.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundTradingPreference)
	assert.Equal(t, createdTradingPreference.ID, foundTradingPreference.ID)
}

func TestUpdateTradingPreferenceReturnsTradingPreference(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	tradingPreferenceFactory := &entities.TradingPreferenceFactory{}
	tradingPreference := tradingPreferenceFactory.NewTradingPreference(
		uuid.New(),
		"swing_trading",
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)

	dto := dtos.TradingPreference{}
	dto.FromEntity(tradingPreference)

	database.Create(&dto)

	tradingPreference.Algorithm = "scalping"
	tradingPreference.RiskLevel = constants.TradingPreferenceRiskLevelMedium

	// Act
	updatedTradingPreference, err := tradePreferenceRepository.Update(ctx, tradingPreference)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedTradingPreference)
	assert.Equal(t, tradingPreference.ID, updatedTradingPreference.ID)
	assert.Equal(t, tradingPreference.Algorithm, updatedTradingPreference.Algorithm)
	assert.Equal(t, tradingPreference.RiskLevel, updatedTradingPreference.RiskLevel)

	foundTradingPreference, err := tradePreferenceRepository.GetByID(ctx, updatedTradingPreference.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundTradingPreference)
	assert.Equal(t, updatedTradingPreference.ID, foundTradingPreference.ID)
	assert.Equal(t, updatedTradingPreference.Algorithm, foundTradingPreference.Algorithm)
	assert.Equal(t, updatedTradingPreference.RiskLevel, foundTradingPreference.RiskLevel)
}

func TestDeleteTradingPreferenceReturnsErrorIfIDDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	err := tradePreferenceRepository.Delete(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteTradingPreferencesDeletesTradingPreference(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	tradingPreferenceFactory := &entities.TradingPreferenceFactory{}
	tradingPreference := tradingPreferenceFactory.NewTradingPreference(
		uuid.New(),
		"swing_trading",
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)

	dto := dtos.TradingPreference{}
	dto.FromEntity(tradingPreference)

	database.Create(&dto)

	// Act
	err := tradePreferenceRepository.Delete(ctx, dto.ID)

	// Assert
	assert.NoError(t, err)

	foundTradingPreference, err := tradePreferenceRepository.GetByID(ctx, dto.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, foundTradingPreference)
}

// HoldingRepository tests

func TestGetHoldingByIDReturnsError(t *testing.T) {
	// Arrange
	id := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := holdingRepository.GetByID(ctx, id)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetHoldingByIDReturnsHolding(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	holdingFactory := &entities.HoldingFactory{}
	holding := holdingFactory.NewHolding(
		uuid.New(),
		"BTCUSDT",
		1.5,
		50000.0,
		51000.0,
		1500.0,
		100.0,
		"closed",
	)

	dto := dtos.Holding{}
	dto.FromEntity(holding)

	database.Create(&dto)

	// Act
	foundHolding, err := holdingRepository.GetByID(ctx, holding.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundHolding)
	assert.Equal(t, holding.ID, foundHolding.ID)
}

func TestGetHoldingByUserIDReturnsHoldings(t *testing.T) {
	// Arrange
	userID := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	holdingFactory := &entities.HoldingFactory{}
	holding := holdingFactory.NewHolding(
		userID,
		"BTCUSDT",
		1.5,
		50000.0,
		51000.0,
		1500.0,
		100.0,
		"closed",
	)

	dto := dtos.Holding{}
	dto.FromEntity(holding)

	database.Create(&dto)

	// Act
	foundHoldings, err := holdingRepository.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundHoldings)
	assert.Equal(t, 1, len(*foundHoldings))
	assert.Equal(t, dto.ID, (*foundHoldings)[0].ID)
}

func TestGetAllHoldingsWithNotMatchingFiltersReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "NONEXISTENT",
	}, "created_at", "desc", 0, 10)
	holdings, err := holdingRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, holdings)
	assert.Equal(t, 0, len(*holdings))
}

func TestGetAllHoldingsReturnsHoldings(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	holdingFactory := &entities.HoldingFactory{}
	holding := holdingFactory.NewHolding(
		uuid.New(),
		"LTCUSDT",
		1.5,
		50000.0,
		51000.0,
		1500.0,
		100.0,
		"closed",
	)

	dto := dtos.Holding{}
	dto.FromEntity(holding)

	database.Create(&dto)

	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "LTCUSDT",
	}, "created_at", "desc", 0, 10)

	// Act
	holdings, err := holdingRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, holdings)
	assert.Equal(t, 1, len(*holdings))
	assert.Equal(t, dto.ID, (*holdings)[0].ID)
}

func TestCreateHoldingReturnsHolding(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	holdingFactory := &entities.HoldingFactory{}
	holding := holdingFactory.NewHolding(
		uuid.New(),
		"ETHUSDT",
		10.0,
		3000.0,
		0.0,
		0.0,
		100.0,
		"open",
	)

	// Act
	createdHolding, err := holdingRepository.Create(ctx, holding)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdHolding)
	assert.Equal(t, holding.ID, createdHolding.ID)

	foundHolding, err := holdingRepository.GetByID(ctx, createdHolding.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundHolding)
	assert.Equal(t, createdHolding.ID, foundHolding.ID)
}

func TestUpdateHoldingReturnsHolding(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	holdingFactory := &entities.HoldingFactory{}
	holding := holdingFactory.NewHolding(
		uuid.New(),
		"BTCUSDT",
		1.5,
		50000.0,
		0.0,
		0.0,
		100.0,
		"open",
	)

	dto := dtos.Holding{}
	dto.FromEntity(holding)

	database.Create(&dto)

	holding.Quantity = 2.0
	holding.EntryPrice = 55000.0
	holding.ExitPrice = 56000.0
	holding.Profit = 2000.0
	holding.Status = "closed"

	// Act
	updatedHolding, err := holdingRepository.Update(ctx, holding)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedHolding)
	assert.Equal(t, holding.ID, updatedHolding.ID)
	assert.Equal(t, holding.Quantity, updatedHolding.Quantity)
	assert.Equal(t, holding.EntryPrice, updatedHolding.EntryPrice)
	assert.Equal(t, holding.ExitPrice, updatedHolding.ExitPrice)
	assert.Equal(t, holding.Profit, updatedHolding.Profit)
	assert.Equal(t, holding.Status, updatedHolding.Status)

	foundHolding, err := holdingRepository.GetByID(ctx, updatedHolding.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundHolding)
	assert.Equal(t, updatedHolding.ID, foundHolding.ID)
	assert.Equal(t, updatedHolding.Quantity, foundHolding.Quantity)
	assert.Equal(t, updatedHolding.EntryPrice, foundHolding.EntryPrice)
	assert.Equal(t, updatedHolding.ExitPrice, foundHolding.ExitPrice)
	assert.Equal(t, updatedHolding.Profit, foundHolding.Profit)
	assert.Equal(t, updatedHolding.Status, foundHolding.Status)
}

func TestDeleteHoldingReturnsErrorIfIDDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	err := holdingRepository.Delete(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteHoldingDeletesHolding(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	holdingFactory := &entities.HoldingFactory{}
	holding := holdingFactory.NewHolding(
		uuid.New(),
		"BTCUSDT",
		1.5,
		50000.0,
		51000.0,
		1500.0,
		100.0,
		"closed",
	)

	dto := dtos.Holding{}
	dto.FromEntity(holding)

	database.Create(&dto)

	// Act
	err := holdingRepository.Delete(ctx, dto.ID)

	// Assert
	assert.NoError(t, err)

	foundHolding, err := holdingRepository.GetByID(ctx, dto.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, foundHolding)
}

// OrderRepository tests

func TestGetOrderByIDReturnsError(t *testing.T) {
	// Arrange
	id := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := orderRepository.GetByID(ctx, id)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetOrderByIDReturnsOrder(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	orderFactory := &entities.OrderFactory{}
	order := orderFactory.NewOrder(
		uuid.New(),
		"BTCUSDT",
		1.5,
		50000.0,
		constants.OrderTypeStopLoss,
	)

	dto := dtos.Order{}
	dto.FromEntity(order)

	database.Create(&dto)

	// Act
	foundOrder, err := orderRepository.GetByID(ctx, order.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundOrder)
	assert.Equal(t, order.ID, foundOrder.ID)
}

func TestGetOrderByUserIDReturnsOrders(t *testing.T) {
	// Arrange
	userID := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	orderFactory := &entities.OrderFactory{}
	order := orderFactory.NewOrder(
		userID,
		"BTCUSDT",
		1.5,
		50000.0,
		constants.OrderTypeStopLoss,
	)

	dto := dtos.Order{}
	dto.FromEntity(order)

	database.Create(&dto)

	// Act
	foundOrders, err := orderRepository.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundOrders)
	assert.Equal(t, 1, len(*foundOrders))
	assert.Equal(t, dto.ID, (*foundOrders)[0].ID)
}

func TestGetAllOrdersWithNotMatchingFiltersReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "NONEXISTENT",
	}, "created_at", "desc", 0, 10)
	orders, err := orderRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 0, len(*orders))
}

func TestGetAllOrdersReturnsOrders(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	orderFactory := &entities.OrderFactory{}
	order := orderFactory.NewOrder(
		uuid.New(),
		"LTCUSDT",
		1.5,
		50000.0,
		constants.OrderTypeStopLoss,
	)

	dto := dtos.Order{}
	dto.FromEntity(order)

	database.Create(&dto)

	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"symbol": "LTCUSDT",
	}, "created_at", "desc", 0, 10)

	// Act
	orders, err := orderRepository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 1, len(*orders))
	assert.Equal(t, dto.ID, (*orders)[0].ID)
}

func TestCreateOrderReturnsOrder(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	orderFactory := &entities.OrderFactory{}
	order := orderFactory.NewOrder(
		uuid.New(),
		"ETHUSDT",
		10.0,
		3000.0,
		constants.OrderTypeStopLoss,
	)

	// Act
	createdOrder, err := orderRepository.Create(ctx, order)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.Equal(t, order.ID, createdOrder.ID)

	foundOrder, err := orderRepository.GetByID(ctx, createdOrder.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundOrder)
	assert.Equal(t, createdOrder.ID, foundOrder.ID)
}

func TestUpdateOrderReturnsOrder(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	orderFactory := &entities.OrderFactory{}
	order := orderFactory.NewOrder(
		uuid.New(),
		"BTCUSDT",
		1.5,
		50000.0,
		constants.OrderTypeStopLoss,
	)

	dto := dtos.Order{}
	dto.FromEntity(order)

	database.Create(&dto)

	order.Quantity = 2.0
	order.Price = 55000.0
	order.Status = constants.OrderStatusFilled
	order.TradeType = constants.OrderTypeStopLoss

	// Act
	updatedOrder, err := orderRepository.Update(ctx, order)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, order.ID, updatedOrder.ID)
	assert.Equal(t, order.Quantity, updatedOrder.Quantity)
	assert.Equal(t, order.Price, updatedOrder.Price)
	assert.Equal(t, order.Status, updatedOrder.Status)
	assert.Equal(t, order.TradeType, updatedOrder.TradeType)

	foundOrder, err := orderRepository.GetByID(ctx, updatedOrder.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundOrder)
	assert.Equal(t, updatedOrder.ID, foundOrder.ID)
	assert.Equal(t, updatedOrder.Quantity, foundOrder.Quantity)
	assert.Equal(t, updatedOrder.Price, foundOrder.Price)
	assert.Equal(t, updatedOrder.Status, foundOrder.Status)
	assert.Equal(t, updatedOrder.TradeType, foundOrder.TradeType)
}

func TestDeleteOrderReturnsErrorIfIDDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	err := orderRepository.Delete(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteOrderDeletesOrder(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	orderFactory := &entities.OrderFactory{}
	order := orderFactory.NewOrder(
		uuid.New(),
		"BTCUSDT",
		1.5,
		50000.0,
		constants.OrderTypeStopLoss,
	)

	dto := dtos.Order{}
	dto.FromEntity(order)

	database.Create(&dto)

	// Act
	err := orderRepository.Delete(ctx, dto.ID)

	// Assert
	assert.NoError(t, err)

	foundOrder, err := orderRepository.GetByID(ctx, dto.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, foundOrder)
}
