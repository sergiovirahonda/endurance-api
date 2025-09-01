package dtos

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestTradingPreference_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()
	watchlist := []string{"BTCUSDT", "ETHUSDT"}

	dto := &TradingPreference{
		ID:                  id,
		UserID:              userId,
		Algorithm:           "test-algorithm",
		Watchlist:           watchlist,
		Operate:             true,
		StopLossEnabled:     true,
		StopLossExitEnabled: true,
		RiskLevel:           constants.TradingPreferenceRiskLevelLow,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, userId, entity.UserID)
	assert.Equal(t, "test-algorithm", entity.Algorithm)
	assert.Equal(t, watchlist, entity.Watchlist)
	assert.True(t, entity.Operate)
	assert.True(t, entity.StopLossEnabled)
	assert.True(t, entity.StopLossExitEnabled)
	assert.Equal(t, constants.TradingPreferenceRiskLevelLow, entity.RiskLevel)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestTradingPreference_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()
	watchlist := []string{"BTCUSDT", "ETHUSDT"}

	entity := &entities.TradingPreference{
		ID:                  id,
		UserID:              userId,
		Algorithm:           "test-algorithm",
		Watchlist:           watchlist,
		Operate:             true,
		StopLossEnabled:     true,
		StopLossExitEnabled: true,
		RiskLevel:           constants.TradingPreferenceRiskLevelLow,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	dto := &TradingPreference{}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, userId, dto.UserID)
	assert.Equal(t, "test-algorithm", dto.Algorithm)
	// assert.Equal(t, watchlist, dto.Watchlist)
	assert.True(t, dto.Operate)
	assert.True(t, dto.StopLossEnabled)
	assert.True(t, dto.StopLossExitEnabled)
	assert.Equal(t, constants.TradingPreferenceRiskLevelLow, dto.RiskLevel)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestTradingPreferences_ToEntities(t *testing.T) {
	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()
	userId := uuid.New()
	now := time.Now()
	watchlist1 := []string{"BTCUSDT"}
	watchlist2 := []string{"ETHUSDT"}

	dtos := TradingPreferences{
		{
			ID:                  id1,
			UserID:              userId,
			Algorithm:           "algorithm-1",
			Watchlist:           watchlist1,
			Operate:             true,
			StopLossEnabled:     true,
			StopLossExitEnabled: false,
			RiskLevel:           constants.TradingPreferenceRiskLevelLow,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  id2,
			UserID:              userId,
			Algorithm:           "algorithm-2",
			Watchlist:           watchlist2,
			Operate:             false,
			StopLossEnabled:     true,
			StopLossExitEnabled: true,
			RiskLevel:           constants.TradingPreferenceRiskLevelLow,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	// Act
	entities := dtos.ToEntities()

	// Assert
	assert.NotNil(t, entities)
	assert.Len(t, *entities, 2)

	// Check first entity
	assert.Equal(t, id1, (*entities)[0].ID)
	assert.Equal(t, userId, (*entities)[0].UserID)
	assert.Equal(t, "algorithm-1", (*entities)[0].Algorithm)
	assert.Equal(t, watchlist1, (*entities)[0].Watchlist)
	assert.True(t, (*entities)[0].Operate)
	assert.True(t, (*entities)[0].StopLossEnabled)
	assert.False(t, (*entities)[0].StopLossExitEnabled)
	assert.Equal(t, constants.TradingPreferenceRiskLevelLow, (*entities)[0].RiskLevel)
	assert.Equal(t, now, (*entities)[0].CreatedAt)
	assert.Equal(t, now, (*entities)[0].UpdatedAt)

	// Check second entity
	assert.Equal(t, id2, (*entities)[1].ID)
	assert.Equal(t, userId, (*entities)[1].UserID)
	assert.Equal(t, "algorithm-2", (*entities)[1].Algorithm)
	assert.Equal(t, watchlist2, (*entities)[1].Watchlist)
	assert.False(t, (*entities)[1].Operate)
	assert.True(t, (*entities)[1].StopLossEnabled)
	assert.True(t, (*entities)[1].StopLossExitEnabled)
	assert.Equal(t, constants.TradingPreferenceRiskLevelLow, (*entities)[1].RiskLevel)
	assert.Equal(t, now, (*entities)[1].CreatedAt)
	assert.Equal(t, now, (*entities)[1].UpdatedAt)
}

func TestHolding_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()

	dto := &Holding{
		ID:         id,
		UserID:     userId,
		Symbol:     "BTCUSDT",
		Quantity:   1.5,
		EntryPrice: 50000.0,
		ExitPrice:  51000.0,
		Profit:     1500.0,
		Status:     "closed",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, userId, entity.UserID)
	assert.Equal(t, "BTCUSDT", entity.Symbol)
	assert.Equal(t, 1.5, entity.Quantity)
	assert.Equal(t, 50000.0, entity.EntryPrice)
	assert.Equal(t, 51000.0, entity.ExitPrice)
	assert.Equal(t, 1500.0, entity.Profit)
	assert.Equal(t, "closed", entity.Status)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestHolding_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()

	entity := &entities.Holding{
		ID:         id,
		UserID:     userId,
		Symbol:     "BTCUSDT",
		Quantity:   1.5,
		EntryPrice: 50000.0,
		ExitPrice:  51000.0,
		Profit:     1500.0,
		Status:     "closed",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	dto := &Holding{}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, userId, dto.UserID)
	assert.Equal(t, "BTCUSDT", dto.Symbol)
	assert.Equal(t, 1.5, dto.Quantity)
	assert.Equal(t, 50000.0, dto.EntryPrice)
	assert.Equal(t, 51000.0, dto.ExitPrice)
	assert.Equal(t, 1500.0, dto.Profit)
	assert.Equal(t, "closed", dto.Status)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestHoldings_ToEntities(t *testing.T) {
	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()
	userId := uuid.New()
	now := time.Now()

	dtos := Holdings{
		{
			ID:         id1,
			UserID:     userId,
			Symbol:     "BTCUSDT",
			Quantity:   1.5,
			EntryPrice: 50000.0,
			ExitPrice:  51000.0,
			Profit:     1500.0,
			Status:     "closed",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         id2,
			UserID:     userId,
			Symbol:     "ETHUSDT",
			Quantity:   10.0,
			EntryPrice: 3000.0,
			ExitPrice:  0.0,
			Profit:     0.0,
			Status:     "open",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	// Act
	entities := dtos.ToEntities()

	// Assert
	assert.NotNil(t, entities)
	assert.Len(t, *entities, 2)

	// Check first entity
	assert.Equal(t, id1, (*entities)[0].ID)
	assert.Equal(t, userId, (*entities)[0].UserID)
	assert.Equal(t, "BTCUSDT", (*entities)[0].Symbol)
	assert.Equal(t, 1.5, (*entities)[0].Quantity)
	assert.Equal(t, 50000.0, (*entities)[0].EntryPrice)
	assert.Equal(t, 51000.0, (*entities)[0].ExitPrice)
	assert.Equal(t, 1500.0, (*entities)[0].Profit)
	assert.Equal(t, "closed", (*entities)[0].Status)
	assert.Equal(t, now, (*entities)[0].CreatedAt)
	assert.Equal(t, now, (*entities)[0].UpdatedAt)

	// Check second entity
	assert.Equal(t, id2, (*entities)[1].ID)
	assert.Equal(t, userId, (*entities)[1].UserID)
	assert.Equal(t, "ETHUSDT", (*entities)[1].Symbol)
	assert.Equal(t, 10.0, (*entities)[1].Quantity)
	assert.Equal(t, 3000.0, (*entities)[1].EntryPrice)
	assert.Equal(t, 0.0, (*entities)[1].ExitPrice)
	assert.Equal(t, 0.0, (*entities)[1].Profit)
	assert.Equal(t, "open", (*entities)[1].Status)
	assert.Equal(t, now, (*entities)[1].CreatedAt)
	assert.Equal(t, now, (*entities)[1].UpdatedAt)
}

func TestOrder_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()

	dto := &Order{
		ID:        id,
		UserID:    userId,
		Symbol:    "BTCUSDT",
		Quantity:  1.5,
		Price:     50000.0,
		Status:    "filled",
		TradeType: "buy",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, userId, entity.UserID)
	assert.Equal(t, "BTCUSDT", entity.Symbol)
	assert.Equal(t, 1.5, entity.Quantity)
	assert.Equal(t, 50000.0, entity.Price)
	assert.Equal(t, "filled", entity.Status)
	assert.Equal(t, "buy", entity.TradeType)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestOrder_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()

	entity := &entities.Order{
		ID:        id,
		UserID:    userId,
		Symbol:    "BTCUSDT",
		Quantity:  1.5,
		Price:     50000.0,
		Status:    "filled",
		TradeType: "buy",
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := &Order{}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, userId, dto.UserID)
	assert.Equal(t, "BTCUSDT", dto.Symbol)
	assert.Equal(t, 1.5, dto.Quantity)
	assert.Equal(t, 50000.0, dto.Price)
	assert.Equal(t, "filled", dto.Status)
	assert.Equal(t, "buy", dto.TradeType)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestOrders_ToEntities(t *testing.T) {
	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()
	userId := uuid.New()
	now := time.Now()

	dtos := Orders{
		{
			ID:        id1,
			UserID:    userId,
			Symbol:    "BTCUSDT",
			Quantity:  1.5,
			Price:     50000.0,
			Status:    "filled",
			TradeType: "buy",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        id2,
			UserID:    userId,
			Symbol:    "ETHUSDT",
			Quantity:  10.0,
			Price:     3000.0,
			Status:    "pending",
			TradeType: "sell",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Act
	entities := dtos.ToEntities()

	// Assert
	assert.NotNil(t, entities)
	assert.Len(t, *entities, 2)

	// Check first entity
	assert.Equal(t, id1, (*entities)[0].ID)
	assert.Equal(t, userId, (*entities)[0].UserID)
	assert.Equal(t, "BTCUSDT", (*entities)[0].Symbol)
	assert.Equal(t, 1.5, (*entities)[0].Quantity)
	assert.Equal(t, 50000.0, (*entities)[0].Price)
	assert.Equal(t, "filled", (*entities)[0].Status)
	assert.Equal(t, "buy", (*entities)[0].TradeType)
	assert.Equal(t, now, (*entities)[0].CreatedAt)
	assert.Equal(t, now, (*entities)[0].UpdatedAt)

	// Check second entity
	assert.Equal(t, id2, (*entities)[1].ID)
	assert.Equal(t, userId, (*entities)[1].UserID)
	assert.Equal(t, "ETHUSDT", (*entities)[1].Symbol)
	assert.Equal(t, 10.0, (*entities)[1].Quantity)
	assert.Equal(t, 3000.0, (*entities)[1].Price)
	assert.Equal(t, "pending", (*entities)[1].Status)
	assert.Equal(t, "sell", (*entities)[1].TradeType)
	assert.Equal(t, now, (*entities)[1].CreatedAt)
	assert.Equal(t, now, (*entities)[1].UpdatedAt)
}
