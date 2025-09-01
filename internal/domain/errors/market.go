package errors

import "errors"

// Validation errors

var (
	ErrMarketDataInsufficient             = errors.New("insufficient market data")
	ErrInvalidMarketSymbol                = errors.New("invalid market data symbol")
	ErrMarketDataTooOld                   = errors.New("market data too old")
	ErrInvalidMarketTimestamp             = errors.New("invalid market data timestamp")
	ErrInvalidMarketOpen                  = errors.New("invalid market data open")
	ErrInvalidMarketHigh                  = errors.New("invalid market data high")
	ErrInvalidMarketLow                   = errors.New("invalid market data low")
	ErrInvalidMarketClose                 = errors.New("invalid market data close")
	ErrInvalidMarketVolume                = errors.New("invalid market data volume")
	ErrInsufficientDataForMACDCalculation = errors.New("insufficient data for MACD calculation")
	ErrInsufficientDataForRSI             = errors.New("insufficient data for RSI calculation")
	ErrInsufficientDataForSMA             = errors.New("insufficient data for SMA calculation")
	ErrInsufficientDataForATR             = errors.New("insufficient data for ATR calculation")
	ErrInsufficientDataForBollingerBands  = errors.New("insufficient data for Bollinger Bands calculation")
	ErrInsufficientDataForOBV             = errors.New("insufficient data for OBV calculation")
	ErrInsufficientDataForADX             = errors.New("insufficient data for ADX calculation")
)
