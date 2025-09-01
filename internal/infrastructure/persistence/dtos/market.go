package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"gorm.io/gorm"
)

type Market struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	Symbol        string    `gorm:"type:varchar(10);not null;"`
	Enabled       bool      `gorm:"type:boolean;not null;default:false;"`
	AverageValue  float64   `gorm:"type:decimal(10,2);"`
	AverageVolume float64   `gorm:"type:decimal(10,2);"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;"`
	UpdatedAt     time.Time `gorm:"type:timestamp;not null;"`
}

type Markets []Market

type MarketData struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	CorrelationID uuid.UUID `gorm:"type:uuid;not null;"`
	Symbol        string    `gorm:"type:varchar(10);not null;"`
	Timestamp     time.Time `gorm:"type:timestamp;not null;"`
	// OHLCV data
	Open   float64 `gorm:"type:decimal(10,2);not null;"`
	High   float64 `gorm:"type:decimal(10,2);not null;"`
	Low    float64 `gorm:"type:decimal(10,2);not null;"`
	Close  float64 `gorm:"type:decimal(10,2);not null;"`
	Volume float64 `gorm:"type:decimal(10,2);not null;"`
	// Technical indicators
	MACD       *float64 `gorm:"type:decimal(10,2);"`
	MACDSignal *float64 `gorm:"type:decimal(10,2);"`
	MACDHist   *float64 `gorm:"type:decimal(10,2);"`
	RSI6       *float64 `gorm:"type:decimal(10,2);"`
	RSI12      *float64 `gorm:"type:decimal(10,2);"`
	RSI24      *float64 `gorm:"type:decimal(10,2);"`
	SMA20      *float64 `gorm:"type:decimal(10,2);"`
	SMA50      *float64 `gorm:"type:decimal(10,2);"`
	SMA200     *float64 `gorm:"type:decimal(10,2);"`
	// Volatility indicators
	ATR                 *float64 `gorm:"type:decimal(10,2);"`
	BollingerBands      *float64 `gorm:"type:decimal(10,2);"`
	BollingerBandsWidth *float64 `gorm:"type:decimal(10,2);"`
	BollingerBandsUpper *float64 `gorm:"type:decimal(10,2);"`
	BollingerBandsLower *float64 `gorm:"type:decimal(10,2);"`
	// Volume indicators
	OBV         *float64 `gorm:"type:decimal(10,2);"`
	ADX         *float64 `gorm:"type:decimal(10,2);"`
	ADXIndex    *float64 `gorm:"type:decimal(10,2);"`
	ADXPositive *float64 `gorm:"type:decimal(10,2);"`
	ADXNegative *float64 `gorm:"type:decimal(10,2);"`
	// Meta
	CreatedAt time.Time `gorm:"type:timestamp;not null;"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;"`
}

type MarketDatas []MarketData

// Receivers

func (m *Market) ToEntity() *entities.Market {
	return &entities.Market{
		ID:            m.ID,
		Symbol:        m.Symbol,
		Enabled:       m.Enabled,
		AverageValue:  m.AverageValue,
		AverageVolume: m.AverageVolume,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (m *Market) FromEntity(market *entities.Market) {
	m.ID = market.ID
	m.Symbol = market.Symbol
	m.Enabled = market.Enabled
	m.AverageValue = market.AverageValue
	m.AverageVolume = market.AverageVolume
	m.CreatedAt = market.CreatedAt
	m.UpdatedAt = market.UpdatedAt
}

func (m *Markets) ToEntities() *entities.Markets {
	entities := make(entities.Markets, len(*m))
	for i, market := range *m {
		entities[i] = *market.ToEntity()
	}
	return &entities
}

func (m *MarketData) ToEntity() *entities.MarketData {
	return &entities.MarketData{
		ID:                  m.ID,
		CorrelationID:       m.CorrelationID,
		Symbol:              m.Symbol,
		Timestamp:           m.Timestamp,
		Open:                m.Open,
		High:                m.High,
		Low:                 m.Low,
		Close:               m.Close,
		Volume:              m.Volume,
		MACD:                m.MACD,
		MACDSignal:          m.MACDSignal,
		MACDHist:            m.MACDHist,
		RSI6:                m.RSI6,
		RSI12:               m.RSI12,
		RSI24:               m.RSI24,
		SMA20:               m.SMA20,
		SMA50:               m.SMA50,
		SMA200:              m.SMA200,
		ATR:                 m.ATR,
		BollingerBands:      m.BollingerBands,
		BollingerBandsWidth: m.BollingerBandsWidth,
		BollingerBandsUpper: m.BollingerBandsUpper,
		BollingerBandsLower: m.BollingerBandsLower,
		OBV:                 m.OBV,
		ADX:                 m.ADX,
		ADXIndex:            m.ADXIndex,
		ADXPositive:         m.ADXPositive,
		ADXNegative:         m.ADXNegative,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func (m *MarketData) FromEntity(marketData *entities.MarketData) {
	m.ID = marketData.ID
	m.CorrelationID = marketData.CorrelationID
	m.Symbol = marketData.Symbol
	m.Timestamp = marketData.Timestamp
	m.Open = marketData.Open
	m.High = marketData.High
	m.Low = marketData.Low
	m.Close = marketData.Close
	m.Volume = marketData.Volume
	m.MACD = marketData.MACD
	m.MACDSignal = marketData.MACDSignal
	m.MACDHist = marketData.MACDHist
	m.RSI6 = marketData.RSI6
	m.RSI12 = marketData.RSI12
	m.RSI24 = marketData.RSI24
	m.SMA20 = marketData.SMA20
	m.SMA50 = marketData.SMA50
	m.SMA200 = marketData.SMA200
	m.ATR = marketData.ATR
	m.BollingerBands = marketData.BollingerBands
	m.BollingerBandsWidth = marketData.BollingerBandsWidth
	m.BollingerBandsUpper = marketData.BollingerBandsUpper
	m.BollingerBandsLower = marketData.BollingerBandsLower
	m.OBV = marketData.OBV
	m.ADX = marketData.ADX
	m.ADXIndex = marketData.ADXIndex
	m.ADXPositive = marketData.ADXPositive
	m.ADXNegative = marketData.ADXNegative
	m.CreatedAt = marketData.CreatedAt
	m.UpdatedAt = marketData.UpdatedAt
}

func (m *MarketDatas) ToEntities() *entities.MarketDatas {
	entities := make(entities.MarketDatas, len(*m))
	for i, marketData := range *m {
		entities[i] = *marketData.ToEntity()
	}
	return &entities
}
