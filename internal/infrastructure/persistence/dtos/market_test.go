package dtos

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestMarket_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()

	dto := &Market{
		ID:            id,
		Symbol:        "BTCUSDT",
		Enabled:       true,
		AverageValue:  10000.0,
		AverageVolume: 1000000.0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, "BTCUSDT", entity.Symbol)
	assert.True(t, entity.Enabled)
	assert.Equal(t, 10000.0, entity.AverageValue)
	assert.Equal(t, 1000000.0, entity.AverageVolume)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestMarket_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()

	dto := &Market{}
	entity := &entities.Market{
		ID:            id,
		Symbol:        "BTCUSDT",
		Enabled:       true,
		AverageValue:  10000.0,
		AverageVolume: 1000000.0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, "BTCUSDT", dto.Symbol)
	assert.True(t, dto.Enabled)
	assert.Equal(t, 10000.0, dto.AverageValue)
	assert.Equal(t, 1000000.0, dto.AverageVolume)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestMarkets_ToEntities(t *testing.T) {
	// Arrange
	markets := Markets{
		{
			ID:            uuid.New(),
			Symbol:        "BTCUSDT",
			Enabled:       true,
			AverageValue:  10000.0,
			AverageVolume: 1000000.0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			Symbol:        "ETHUSDT",
			Enabled:       false,
			AverageValue:  20000.0,
			AverageVolume: 2000000.0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	entities := markets.ToEntities()

	// Assert
	assert.NotNil(t, entities)
	assert.Equal(t, 2, len(*entities))
	assert.Equal(t, "BTCUSDT", (*entities)[0].Symbol)
	assert.True(t, (*entities)[0].Enabled)
	assert.Equal(t, 10000.0, (*entities)[0].AverageValue)
	assert.Equal(t, 1000000.0, (*entities)[0].AverageVolume)
	assert.Equal(t, "ETHUSDT", (*entities)[1].Symbol)
	assert.False(t, (*entities)[1].Enabled)
	assert.Equal(t, 20000.0, (*entities)[1].AverageValue)
	assert.Equal(t, 2000000.0, (*entities)[1].AverageVolume)
}

func TestMarketData_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()
	macd := 1.23
	rsi6 := 45.67
	sma20 := 100.50
	atr := 2.34
	bb := 3.45
	obv := 1000.0
	adx := 25.0

	dto := &MarketData{
		ID:                  id,
		Symbol:              "BTCUSDT",
		Timestamp:           now,
		Open:                50000.0,
		High:                51000.0,
		Low:                 49000.0,
		Close:               50500.0,
		Volume:              100.0,
		MACD:                &macd,
		MACDSignal:          &macd,
		MACDHist:            &macd,
		RSI6:                &rsi6,
		RSI12:               &rsi6,
		RSI24:               &rsi6,
		SMA20:               &sma20,
		SMA50:               &sma20,
		SMA200:              &sma20,
		ATR:                 &atr,
		BollingerBands:      &bb,
		BollingerBandsWidth: &bb,
		BollingerBandsUpper: &bb,
		BollingerBandsLower: &bb,
		OBV:                 &obv,
		ADX:                 &adx,
		ADXIndex:            &adx,
		ADXPositive:         &adx,
		ADXNegative:         &adx,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, "BTCUSDT", entity.Symbol)
	assert.Equal(t, now, entity.Timestamp)
	assert.Equal(t, 50000.0, entity.Open)
	assert.Equal(t, 51000.0, entity.High)
	assert.Equal(t, 49000.0, entity.Low)
	assert.Equal(t, 50500.0, entity.Close)
	assert.Equal(t, 100.0, entity.Volume)
	assert.Equal(t, &macd, entity.MACD)
	assert.Equal(t, &macd, entity.MACDSignal)
	assert.Equal(t, &macd, entity.MACDHist)
	assert.Equal(t, &rsi6, entity.RSI6)
	assert.Equal(t, &rsi6, entity.RSI12)
	assert.Equal(t, &rsi6, entity.RSI24)
	assert.Equal(t, &sma20, entity.SMA20)
	assert.Equal(t, &sma20, entity.SMA50)
	assert.Equal(t, &sma20, entity.SMA200)
	assert.Equal(t, &atr, entity.ATR)
	assert.Equal(t, &bb, entity.BollingerBands)
	assert.Equal(t, &bb, entity.BollingerBandsWidth)
	assert.Equal(t, &bb, entity.BollingerBandsUpper)
	assert.Equal(t, &bb, entity.BollingerBandsLower)
	assert.Equal(t, &obv, entity.OBV)
	assert.Equal(t, &adx, entity.ADX)
	assert.Equal(t, &adx, entity.ADXIndex)
	assert.Equal(t, &adx, entity.ADXPositive)
	assert.Equal(t, &adx, entity.ADXNegative)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestMarketData_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()
	macd := 1.23
	rsi6 := 45.67
	sma20 := 100.50
	atr := 2.34
	bb := 3.45
	obv := 1000.0
	adx := 25.0

	entity := &entities.MarketData{
		ID:                  id,
		Symbol:              "BTCUSDT",
		Timestamp:           now,
		Open:                50000.0,
		High:                51000.0,
		Low:                 49000.0,
		Close:               50500.0,
		Volume:              100.0,
		MACD:                &macd,
		MACDSignal:          &macd,
		MACDHist:            &macd,
		RSI6:                &rsi6,
		RSI12:               &rsi6,
		RSI24:               &rsi6,
		SMA20:               &sma20,
		SMA50:               &sma20,
		SMA200:              &sma20,
		ATR:                 &atr,
		BollingerBands:      &bb,
		BollingerBandsWidth: &bb,
		BollingerBandsUpper: &bb,
		BollingerBandsLower: &bb,
		OBV:                 &obv,
		ADX:                 &adx,
		ADXIndex:            &adx,
		ADXPositive:         &adx,
		ADXNegative:         &adx,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	dto := &MarketData{}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, "BTCUSDT", dto.Symbol)
	assert.Equal(t, now, dto.Timestamp)
	assert.Equal(t, 50000.0, dto.Open)
	assert.Equal(t, 51000.0, dto.High)
	assert.Equal(t, 49000.0, dto.Low)
	assert.Equal(t, 50500.0, dto.Close)
	assert.Equal(t, 100.0, dto.Volume)
	assert.Equal(t, &macd, dto.MACD)
	assert.Equal(t, &macd, dto.MACDSignal)
	assert.Equal(t, &macd, dto.MACDHist)
	assert.Equal(t, &rsi6, dto.RSI6)
	assert.Equal(t, &rsi6, dto.RSI12)
	assert.Equal(t, &rsi6, dto.RSI24)
	assert.Equal(t, &sma20, dto.SMA20)
	assert.Equal(t, &sma20, dto.SMA50)
	assert.Equal(t, &sma20, dto.SMA200)
	assert.Equal(t, &atr, dto.ATR)
	assert.Equal(t, &bb, dto.BollingerBands)
	assert.Equal(t, &bb, dto.BollingerBandsWidth)
	assert.Equal(t, &bb, dto.BollingerBandsUpper)
	assert.Equal(t, &bb, dto.BollingerBandsLower)
	assert.Equal(t, &obv, dto.OBV)
	assert.Equal(t, &adx, dto.ADX)
	assert.Equal(t, &adx, dto.ADXIndex)
	assert.Equal(t, &adx, dto.ADXPositive)
	assert.Equal(t, &adx, dto.ADXNegative)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestMarketData_WithNilPointers(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()

	dto := &MarketData{
		ID:        id,
		Symbol:    "BTCUSDT",
		Timestamp: now,
		Open:      50000.0,
		High:      51000.0,
		Low:       49000.0,
		Close:     50500.0,
		Volume:    100.0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, "BTCUSDT", entity.Symbol)
	assert.Equal(t, now, entity.Timestamp)
	assert.Equal(t, 50000.0, entity.Open)
	assert.Equal(t, 51000.0, entity.High)
	assert.Equal(t, 49000.0, entity.Low)
	assert.Equal(t, 50500.0, entity.Close)
	assert.Equal(t, 100.0, entity.Volume)
	assert.Nil(t, entity.MACD)
	assert.Nil(t, entity.MACDSignal)
	assert.Nil(t, entity.MACDHist)
	assert.Nil(t, entity.RSI6)
	assert.Nil(t, entity.RSI12)
	assert.Nil(t, entity.RSI24)
	assert.Nil(t, entity.SMA20)
	assert.Nil(t, entity.SMA50)
	assert.Nil(t, entity.SMA200)
	assert.Nil(t, entity.ATR)
	assert.Nil(t, entity.BollingerBands)
	assert.Nil(t, entity.BollingerBandsWidth)
	assert.Nil(t, entity.BollingerBandsUpper)
	assert.Nil(t, entity.BollingerBandsLower)
	assert.Nil(t, entity.OBV)
	assert.Nil(t, entity.ADX)
	assert.Nil(t, entity.ADXIndex)
	assert.Nil(t, entity.ADXPositive)
	assert.Nil(t, entity.ADXNegative)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}
