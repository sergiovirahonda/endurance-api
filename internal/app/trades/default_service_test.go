package trades

import (
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

// --- TradingPreferenceService Tests ---

func TestGetTradingPreferenceByIDReturnsErrorIfNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	_, err := tradingPreferenceService.GetByID(ctx, uuid.New())
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCreateTradingPreferenceAndGetByID(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	tpFactory := &entities.TradingPreferenceFactory{}
	tp := tpFactory.NewTradingPreference(
		userID,
		constants.TradingAlgorithmSwingTrading,
		[]string{"BTCUSDT", "ETHUSDT"},
		true,
		true,
		true,
		constants.TradingPreferenceRiskLevelLow,
	)
	ctx.Set("user", &entities.User{ID: userID})
	created, err := tradingPreferenceService.Create(ctx, tp)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	fetched, err := tradingPreferenceService.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestGetTradingPreferenceByUserIDReturnsErrorIfNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	_, err := tradingPreferenceService.GetByUserID(ctx, uuid.New())
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetTradingPreferenceByUserIDReturnsPreference(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	tpFactory := &entities.TradingPreferenceFactory{}
	tp := tpFactory.NewTradingPreference(
		userID,
		constants.TradingAlgorithmDayTrading,
		[]string{"BTCUSDT"},
		false,
		false,
		false,
		constants.TradingPreferenceRiskLevelMedium,
	)
	dto := dtos.TradingPreference{}
	dto.FromEntity(tp)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	fetched, err := tradingPreferenceService.GetByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, tp.ID, fetched.ID)
}

func TestGetAllTradingPreferencesReturnsEmptyIfNoMatch(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", &entities.User{ID: uuid.New()})
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{"algorithm": "nonexistent"}, "created_at", "desc", 0, 10)
	prefs, err := tradingPreferenceService.GetAll(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*prefs))
}

func TestGetAllTradingPreferencesReturnsPreferences(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	tpFactory := &entities.TradingPreferenceFactory{}
	tp := tpFactory.NewTradingPreference(
		userID,
		constants.TradingAlgorithmScalping,
		[]string{"BTCUSDT"},
		true,
		false,
		false,
		constants.TradingPreferenceRiskLevelHigh,
	)
	dto := dtos.TradingPreference{}
	dto.FromEntity(tp)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{"algorithm": constants.TradingAlgorithmScalping}, "created_at", "desc", 0, 10)
	prefs, err := tradingPreferenceService.GetAll(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*prefs))
	assert.Equal(t, tp.ID, (*prefs)[0].ID)
}

func TestUpdateTradingPreference(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	tpFactory := &entities.TradingPreferenceFactory{}
	tp := tpFactory.NewTradingPreference(
		userID,
		constants.TradingAlgorithmSwingTrading,
		[]string{"BTCUSDT"},
		true,
		false,
		false,
		constants.TradingPreferenceRiskLevelLow,
	)
	dto := dtos.TradingPreference{}
	dto.FromEntity(tp)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	tp.Algorithm = constants.TradingAlgorithmDayTrading
	tp.RiskLevel = constants.TradingPreferenceRiskLevelHigh
	updated, err := tradingPreferenceService.Update(ctx, tp)
	assert.NoError(t, err)
	assert.Equal(t, tp.ID, updated.ID)
	assert.Equal(t, constants.TradingAlgorithmDayTrading, updated.Algorithm)
	assert.Equal(t, constants.TradingPreferenceRiskLevelHigh, updated.RiskLevel)
}

func TestDeleteTradingPreference(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	tpFactory := &entities.TradingPreferenceFactory{}
	tp := tpFactory.NewTradingPreference(
		userID,
		constants.TradingAlgorithmSwingTrading,
		[]string{"BTCUSDT"},
		true,
		false,
		false,
		constants.TradingPreferenceRiskLevelLow,
	)
	dto := dtos.TradingPreference{}
	dto.FromEntity(tp)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	err := tradingPreferenceService.Delete(ctx, tp.ID)
	assert.NoError(t, err)
	_, err = tradingPreferenceService.GetByID(ctx, tp.ID)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// --- HoldingService Tests ---

func TestGetHoldingByIDReturnsErrorIfNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	_, err := holdingService.GetByID(ctx, uuid.New())
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCreateHoldingAndGetByID(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	hFactory := &entities.HoldingFactory{}
	h := hFactory.NewHolding(
		userID,
		"BTCUSDT",
		1.5,
		50000.0,
		51000.0,
		1500.0,
		100.0,
		constants.HoldingStatusOpen,
	)
	ctx.Set("user", &entities.User{ID: userID})
	created, err := holdingService.Create(ctx, h)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	fetched, err := holdingService.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestGetAllHoldingsReturnsEmptyIfNoMatch(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", &entities.User{ID: uuid.New()})
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{"symbol": "NONEXISTENT"}, "created_at", "desc", 0, 10)
	holdings, err := holdingService.GetAll(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*holdings))
}

func TestGetAllHoldingsReturnsHoldings(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	hFactory := &entities.HoldingFactory{}
	h := hFactory.NewHolding(
		userID,
		"LTCUSDT",
		2.0,
		100.0,
		120.0,
		40.0,
		100.0,
		constants.HoldingStatusClosed,
	)
	dto := dtos.Holding{}
	dto.FromEntity(h)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{"symbol": "LTCUSDT"}, "created_at", "desc", 0, 10)
	holdings, err := holdingService.GetAll(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*holdings))
	assert.Equal(t, h.ID, (*holdings)[0].ID)
}

func TestUpdateHolding(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	hFactory := &entities.HoldingFactory{}
	h := hFactory.NewHolding(
		userID,
		"BTCUSDT",
		1.0,
		50000.0,
		51000.0,
		1000.0,
		100.0,
		constants.HoldingStatusOpen,
	)
	dto := dtos.Holding{}
	dto.FromEntity(h)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	h.Quantity = 2.0
	h.Status = constants.HoldingStatusClosed
	updated, err := holdingService.Update(ctx, h)
	assert.NoError(t, err)
	assert.Equal(t, h.ID, updated.ID)
	assert.Equal(t, 2.0, updated.Quantity)
	assert.Equal(t, constants.HoldingStatusClosed, updated.Status)
}

func TestDeleteHolding(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	hFactory := &entities.HoldingFactory{}
	h := hFactory.NewHolding(
		userID,
		"BTCUSDT",
		1.0,
		50000.0,
		51000.0,
		1000.0,
		100.0,
		constants.HoldingStatusOpen,
	)
	dto := dtos.Holding{}
	dto.FromEntity(h)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	err := holdingService.Delete(ctx, h.ID)
	assert.NoError(t, err)
	_, err = holdingService.GetByID(ctx, h.ID)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// --- OrderService Tests ---

func TestGetOrderByIDReturnsErrorIfNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	_, err := orderService.GetByID(ctx, uuid.New())
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCreateOrderAndGetByID(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	oFactory := &entities.OrderFactory{}
	o := oFactory.NewOrder(
		userID,
		"BTCUSDT",
		1.0,
		50000.0,
		constants.OrderTypeStopLoss,
	)
	ctx.Set("user", &entities.User{ID: userID})
	created, err := orderService.Create(ctx, o)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	fetched, err := orderService.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestGetAllOrdersReturnsEmptyIfNoMatch(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", &entities.User{ID: uuid.New()})
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{"symbol": "NONEXISTENT"}, "created_at", "desc", 0, 10)
	orders, err := orderService.GetAll(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*orders))
}

func TestGetAllOrdersReturnsOrders(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	oFactory := &entities.OrderFactory{}
	o := oFactory.NewOrder(
		userID,
		"ETHUSDT",
		2.0,
		2000.0,
		constants.OrderTypeTakeProfit,
	)
	dto := dtos.Order{}
	dto.FromEntity(o)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{"symbol": "ETHUSDT"}, "created_at", "desc", 0, 10)
	orders, err := orderService.GetAll(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*orders))
	assert.Equal(t, o.ID, (*orders)[0].ID)
}

func TestUpdateOrder(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	oFactory := &entities.OrderFactory{}
	o := oFactory.NewOrder(
		userID,
		"BTCUSDT",
		1.0,
		50000.0,
		constants.OrderTypeStopLoss,
	)
	dto := dtos.Order{}
	dto.FromEntity(o)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	o.Quantity = 3.0
	o.Status = constants.OrderStatusFilled
	updated, err := orderService.Update(ctx, o)
	assert.NoError(t, err)
	assert.Equal(t, o.ID, updated.ID)
	assert.Equal(t, 3.0, updated.Quantity)
	assert.Equal(t, constants.OrderStatusFilled, updated.Status)
}

func TestDeleteOrder(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	oFactory := &entities.OrderFactory{}
	o := oFactory.NewOrder(
		userID,
		"BTCUSDT",
		1.0,
		50000.0,
		constants.OrderTypeStopLoss,
	)
	dto := dtos.Order{}
	dto.FromEntity(o)
	database.Create(&dto)
	ctx.Set("user", &entities.User{ID: userID})
	err := orderService.Delete(ctx, o.ID)
	assert.NoError(t, err)
	_, err = orderService.GetByID(ctx, o.ID)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
