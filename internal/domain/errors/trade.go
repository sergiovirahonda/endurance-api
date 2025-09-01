package errors

import "errors"

var (
	ErrTradingPreferenceAlreadyExists = errors.New("trading preference already exists")
	// Validation errors - Trading Preference
	ErrInvalidAlgorithm        = errors.New("invalid algorithm")
	ErrInvalidRiskLevel        = errors.New("invalid risk level")
	ErrInvalidWatchlistElement = errors.New("invalid watchlist element")
	// Validation errors - Holding
	ErrInvalidHoldingStatus     = errors.New("invalid holding status")
	ErrInvalidHoldingSymbol     = errors.New("invalid holding symbol")
	ErrInvalidHoldingQuantity   = errors.New("invalid holding quantity")
	ErrInvalidHoldingEntryPrice = errors.New("invalid holding entry price")
	ErrInvalidHoldingExitPrice  = errors.New("invalid holding exit price")
	ErrInvalidHoldingEntryScore = errors.New("invalid holding entry score")
	ErrHoldingNotFound          = errors.New("holding not found")
	// Validation errors - Order
	ErrInvalidOrderStatus   = errors.New("invalid order status")
	ErrInvalidOrderType     = errors.New("invalid order type")
	ErrInvalidOrderPrice    = errors.New("invalid order price")
	ErrInvalidOrderQuantity = errors.New("invalid order quantity")
	ErrInvalidOrderSymbol   = errors.New("invalid order symbol")
)
