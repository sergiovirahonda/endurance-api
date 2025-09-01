package errors

import "errors"

var (
	ErrNoBalance                    = errors.New("no balances available")
	ErrAccountNotAvailable          = errors.New("account not available")
	ErrInvalidBalance               = errors.New("invalid balance")
	ErrInvalidPrice                 = errors.New("invalid price")
	ErrInvalidVolume                = errors.New("invalid volume")
	ErrInvalidPricePercentageChange = errors.New("invalid price percentage change")
	ErrTickerNotAvailable           = errors.New("ticker not available")
	ErrExchangeInfoNotAvailable     = errors.New("exchange info not available")
	ErrConversionQuoteNotAvailable  = errors.New("conversion quote not available")
	ErrInvalidRatio                 = errors.New("invalid ratio")
	ErrInvalidInverseRatio          = errors.New("invalid inverse ratio")
	ErrInvalidToAmount              = errors.New("invalid to amount")
	ErrInsufficientBalance          = errors.New("insufficient balance")
	ErrInvalidSymbol                = errors.New("invalid symbol")
	ErrQuoteExpired                 = errors.New("quote expired")
	ErrInvalidQuoteAmount           = errors.New("invalid quote amount")
	ErrBalanceNotFound              = errors.New("balance not found")
	// Validation errors - Exchange Conversion Quote
	ErrInvalidFromAmount      = errors.New("invalid from amount")
	ErrInvalidValidTime       = errors.New("invalid valid time")
	ErrInvalidFee             = errors.New("invalid fee")
	ErrInvalidConversionDrift = errors.New("invalid conversion drift")
	// Klines
	ErrKlinesNotAvailable = errors.New("klines not available")
	ErrInvalidOpen        = errors.New("invalid open")
	ErrInvalidHigh        = errors.New("invalid high")
	ErrInvalidLow         = errors.New("invalid low")
	ErrInvalidClose       = errors.New("invalid close")
	ErrInvalidTrades      = errors.New("invalid trades")
	// WebSocket
	ErrWebSocketConnectionFailed   = errors.New("websocket connection failed")
	ErrWebSocketSubscriptionFailed = errors.New("websocket subscription failed")
	ErrWebSocketMessageInvalid     = errors.New("invalid websocket message")
	ErrWebSocketStreamClosed       = errors.New("websocket stream closed")
	ErrWebSocketReconnectionFailed = errors.New("websocket reconnection failed")
)
