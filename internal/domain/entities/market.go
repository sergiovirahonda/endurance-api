package entities

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/domain/events"
)

type Market struct {
	ID            uuid.UUID `json:"id"`
	Symbol        string    `json:"symbol"`
	Enabled       bool      `json:"enabled"`
	AverageValue  float64   `json:"average_value"`
	AverageVolume float64   `json:"average_volume"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Markets []Market

type MarketData struct {
	ID            uuid.UUID `json:"id"`
	CorrelationID uuid.UUID `json:"correlation_id"`
	Symbol        string    `json:"symbol"`
	Timestamp     time.Time `json:"timestamp"`
	// OHLCV data
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	// Technical indicators
	MACD       *float64 `json:"macd"`
	MACDSignal *float64 `json:"macd_signal"`
	MACDHist   *float64 `json:"macd_hist"`
	RSI6       *float64 `json:"rsi6"`
	RSI12      *float64 `json:"rsi12"`
	RSI24      *float64 `json:"rsi24"`
	SMA20      *float64 `json:"sma20"`
	SMA50      *float64 `json:"sma50"`
	SMA200     *float64 `json:"sma200"`
	// Volatility indicators
	ATR                 *float64 `json:"atr"`
	BollingerBands      *float64 `json:"bollinger_bands"`
	BollingerBandsWidth *float64 `json:"bollinger_bands_width"`
	BollingerBandsUpper *float64 `json:"bollinger_bands_upper"`
	BollingerBandsLower *float64 `json:"bollinger_bands_lower"`
	// Volume indicators
	OBV         *float64 `json:"obv"`
	ADX         *float64 `json:"adx"`
	ADXIndex    *float64 `json:"adx_index"`
	ADXPositive *float64 `json:"adx_positive"`
	ADXNegative *float64 `json:"adx_negative"`
	// Meta
	// Score
	Score     *float64  `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MarketDatas []MarketData

// Validations

func (m *MarketData) Validate() error {
	if !strings.HasSuffix(m.Symbol, "USDT") {
		return errors.ErrInvalidMarketSymbol
	}
	if m.Timestamp.IsZero() {
		return errors.ErrInvalidMarketTimestamp
	}
	if m.Open < 0 {
		return errors.ErrInvalidMarketOpen
	}
	if m.High < 0 {
		return errors.ErrInvalidMarketHigh
	}
	if m.Low < 0 {
		return errors.ErrInvalidMarketLow
	}
	if m.Close < 0 {
		return errors.ErrInvalidMarketClose
	}
	if m.Volume < 0 {
		return errors.ErrInvalidMarketVolume
	}
	return nil
}

func (m *MarketData) FromEvent(event *events.MarketDataEvent) {
	factory := MarketDataFactory{}
	marketData := factory.NewRawMarketData(
		event.ID,
		event.Symbol,
		event.DataTimestamp,
		event.Open,
		event.High,
		event.Low,
		event.Close,
		event.Volume,
	)
	*m = *marketData
}

// Factories

type MarketDataFactory struct{}

func (f *MarketDataFactory) NewUnparsedMarketData(
	correlationID uuid.UUID,
	symbol string,
	timestamp time.Time,
	open string,
	high string,
	low string,
	close string,
	volume string,
) (*MarketData, error) {
	openFloat, err := strconv.ParseFloat(open, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketOpen
	}
	highFloat, err := strconv.ParseFloat(high, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketHigh
	}
	lowFloat, err := strconv.ParseFloat(low, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketLow
	}
	closeFloat, err := strconv.ParseFloat(close, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketClose
	}
	volumeFloat, err := strconv.ParseFloat(volume, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketVolume
	}
	return &MarketData{
		ID:                  uuid.New(),
		CorrelationID:       correlationID,
		Symbol:              symbol,
		Timestamp:           timestamp,
		Open:                openFloat,
		High:                highFloat,
		Low:                 lowFloat,
		Close:               closeFloat,
		Volume:              volumeFloat,
		MACD:                nil,
		MACDSignal:          nil,
		MACDHist:            nil,
		RSI6:                nil,
		RSI12:               nil,
		RSI24:               nil,
		SMA20:               nil,
		SMA50:               nil,
		SMA200:              nil,
		ATR:                 nil,
		BollingerBands:      nil,
		BollingerBandsWidth: nil,
		BollingerBandsUpper: nil,
		BollingerBandsLower: nil,
		OBV:                 nil,
		ADX:                 nil,
		ADXIndex:            nil,
		ADXPositive:         nil,
		ADXNegative:         nil,
		Score:               nil,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}, nil
}

func (f *MarketDataFactory) NewRawMarketData(
	correlationID uuid.UUID,
	symbol string,
	timestamp time.Time,
	open float64,
	high float64,
	low float64,
	close float64,
	volume float64,
) *MarketData {
	return f.NewMarketData(
		correlationID,
		symbol,
		timestamp,
		open,
		high,
		low,
		close,
		volume,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func (f *MarketDataFactory) NewMarketData(
	correlationID uuid.UUID,
	symbol string,
	timestamp time.Time,
	open float64,
	high float64,
	low float64,
	close float64,
	volume float64,
	macd *float64,
	macdSignal *float64,
	macdHist *float64,
	rsi6 *float64,
	rsi12 *float64,
	rsi24 *float64,
	sma20 *float64,
	sma50 *float64,
	sma200 *float64,
	atr *float64,
	bollingerBands *float64,
	bollingerBandsWidth *float64,
	bollingerBandsUpper *float64,
	bollingerBandsLower *float64,
	obv *float64,
	adx *float64,
	adxIndex *float64,
	adxPositive *float64,
	adxNegative *float64,
	score *float64,
) *MarketData {
	return &MarketData{
		ID:                  uuid.New(),
		CorrelationID:       correlationID,
		Symbol:              symbol,
		Timestamp:           timestamp,
		Open:                open,
		High:                high,
		Low:                 low,
		Close:               close,
		Volume:              volume,
		MACD:                macd,
		MACDSignal:          macdSignal,
		MACDHist:            macdHist,
		RSI6:                rsi6,
		RSI12:               rsi12,
		RSI24:               rsi24,
		SMA20:               sma20,
		SMA50:               sma50,
		SMA200:              sma200,
		ATR:                 atr,
		BollingerBands:      bollingerBands,
		BollingerBandsWidth: bollingerBandsWidth,
		BollingerBandsUpper: bollingerBandsUpper,
		BollingerBandsLower: bollingerBandsLower,
		OBV:                 obv,
		ADX:                 adx,
		ADXIndex:            adxIndex,
		ADXPositive:         adxPositive,
		ADXNegative:         adxNegative,
		Score:               score,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}
}

func (f *MarketDataFactory) Clone(
	marketData *MarketData,
	timestamp time.Time,
	open float64,
	high float64,
	low float64,
	close float64,
	volume float64,
	macd *float64,
	macdSignal *float64,
	macdHist *float64,
	rsi6 *float64,
	rsi12 *float64,
	rsi24 *float64,
	sma20 *float64,
	sma50 *float64,
	sma200 *float64,
	atr *float64,
	bollingerBands *float64,
	bollingerBandsWidth *float64,
	bollingerBandsUpper *float64,
	bollingerBandsLower *float64,
	obv *float64,
	adx *float64,
	adxIndex *float64,
	adxPositive *float64,
	adxNegative *float64,
	score *float64,
) *MarketData {
	return &MarketData{
		ID:                  marketData.ID,
		CorrelationID:       marketData.CorrelationID,
		Symbol:              marketData.Symbol,
		Timestamp:           timestamp,
		Open:                open,
		High:                high,
		Low:                 low,
		Close:               close,
		Volume:              volume,
		MACD:                macd,
		MACDSignal:          macdSignal,
		MACDHist:            macdHist,
		RSI6:                rsi6,
		RSI12:               rsi12,
		RSI24:               rsi24,
		SMA20:               sma20,
		SMA50:               sma50,
		SMA200:              sma200,
		ATR:                 atr,
		BollingerBands:      bollingerBands,
		BollingerBandsWidth: bollingerBandsWidth,
		BollingerBandsUpper: bollingerBandsUpper,
		BollingerBandsLower: bollingerBandsLower,
		OBV:                 obv,
		ADX:                 adx,
		ADXIndex:            adxIndex,
		ADXPositive:         adxPositive,
		ADXNegative:         adxNegative,
		Score:               score,
		CreatedAt:           marketData.CreatedAt,
		UpdatedAt:           time.Now().UTC(),
	}
}

func (f *MarketDataFactory) NewMarketDataFromEvent(
	correlationID uuid.UUID,
	symbol string,
	timestamp time.Time,
	open float64,
	high float64,
	low float64,
	close float64,
	volume float64,
) *MarketData {
	return f.NewRawMarketData(
		correlationID,
		symbol,
		timestamp,
		open,
		high,
		low,
		close,
		volume,
	)
}
