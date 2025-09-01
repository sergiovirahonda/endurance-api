package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/lib"
)

type TradingPreference struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"user_id"`
	Algorithm           string    `json:"algorithm"`
	Watchlist           []string  `json:"watchlist"`
	Operate             bool      `json:"operate"`
	StopLossEnabled     bool      `json:"stop_loss_enabled"`
	StopLossExitEnabled bool      `json:"stop_loss_exit"`
	RiskLevel           string    `json:"risk_level"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type TradingPreferences []TradingPreference

type Holding struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Symbol     string    `json:"symbol"`
	Quantity   float64   `json:"quantity"`
	EntryPrice float64   `json:"entry_price"`
	ExitPrice  float64   `json:"exit_price"`
	Profit     float64   `json:"profit"`
	EntryScore float64   `json:"entry_score"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Holdings []Holding

type Order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Symbol    string    `json:"symbol"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	TradeType string    `json:"trade_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Orders []Order

// Validations

func (tp *TradingPreference) Validate() error {
	if !lib.SliceContains(constants.TradingPreferenceAlgorithms, tp.Algorithm) {
		return errors.ErrInvalidAlgorithm
	}
	if !lib.SliceContains(constants.TradingPreferenceRiskLevels, tp.RiskLevel) {
		return errors.ErrInvalidRiskLevel
	}
	for _, symbol := range tp.Watchlist {
		if !strings.HasSuffix(symbol, "USDT") {
			return errors.ErrInvalidWatchlistElement
		}
	}
	return nil
}

func (h *Holding) Validate() error {
	if !lib.SliceContains(constants.HoldingStatuses, h.Status) {
		return errors.ErrInvalidHoldingStatus
	}
	if !strings.HasSuffix(h.Symbol, "USDT") {
		return errors.ErrInvalidHoldingSymbol
	}
	if h.Quantity < 0 {
		return errors.ErrInvalidHoldingQuantity
	}
	if h.EntryPrice < 0 {
		return errors.ErrInvalidHoldingEntryPrice
	}
	if h.ExitPrice < 0 {
		return errors.ErrInvalidHoldingExitPrice
	}
	if h.EntryScore < 0 {
		return errors.ErrInvalidHoldingEntryScore
	}
	return nil
}

func (h *Holding) GetAsset() string {
	return strings.Split(h.Symbol, "USDT")[0]
}

func (o *Order) Validate() error {
	if !lib.SliceContains(constants.OrderStatuses, o.Status) {
		return errors.ErrInvalidOrderStatus
	}
	if !lib.SliceContains(constants.OrderTypes, o.TradeType) {
		return errors.ErrInvalidOrderType
	}
	if o.Quantity < 0 {
		return errors.ErrInvalidOrderQuantity
	}
	if o.Price < 0 {
		return errors.ErrInvalidOrderPrice
	}
	if !strings.HasSuffix(o.Symbol, "USDT") {
		return errors.ErrInvalidOrderSymbol
	}
	return nil
}

// Factories

type HoldingFactory struct{}

func (f *HoldingFactory) NewHolding(
	userID uuid.UUID,
	symbol string,
	quantity float64,
	entryPrice float64,
	exitPrice float64,
	profit float64,
	entryScore float64,
	status string,
) *Holding {
	return &Holding{
		ID:         uuid.New(),
		UserID:     userID,
		Symbol:     symbol,
		Quantity:   quantity,
		EntryPrice: entryPrice,
		ExitPrice:  exitPrice,
		Profit:     profit,
		EntryScore: entryScore,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (f *HoldingFactory) Clone(
	holding *Holding,
	symbol string,
	quantity float64,
	entryPrice float64,
	exitPrice float64,
	profit float64,
	entryScore float64,
	status string,
) *Holding {
	return &Holding{
		ID:         holding.ID,
		UserID:     holding.UserID,
		Symbol:     symbol,
		Quantity:   quantity,
		EntryPrice: entryPrice,
		ExitPrice:  exitPrice,
		Profit:     profit,
		EntryScore: entryScore,
		Status:     status,
		CreatedAt:  holding.CreatedAt,
		UpdatedAt:  holding.UpdatedAt,
	}
}

type TradingPreferenceFactory struct{}

func (f *TradingPreferenceFactory) NewTradingPreference(
	userID uuid.UUID,
	algorithm string,
	watchlist []string,
	operate bool,
	stopLossEnabled bool,
	StopLossExitEnabled bool,
	riskLevel string,
) *TradingPreference {
	return &TradingPreference{
		ID:                  uuid.New(),
		UserID:              userID,
		Algorithm:           algorithm,
		Watchlist:           watchlist,
		Operate:             operate,
		StopLossEnabled:     stopLossEnabled,
		StopLossExitEnabled: StopLossExitEnabled,
		RiskLevel:           riskLevel,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}
}

func (f *TradingPreferenceFactory) Clone(
	tradingPreference *TradingPreference,
	algorithm string,
	watchlist []string,
	operate bool,
	stopLossEnabled bool,
	StopLossExitEnabled bool,
	riskLevel string,
) *TradingPreference {
	return &TradingPreference{
		ID:                  tradingPreference.ID,
		UserID:              tradingPreference.UserID,
		Algorithm:           algorithm,
		Watchlist:           watchlist,
		Operate:             operate,
		StopLossEnabled:     stopLossEnabled,
		StopLossExitEnabled: StopLossExitEnabled,
		RiskLevel:           riskLevel,
		CreatedAt:           tradingPreference.CreatedAt,
		UpdatedAt:           tradingPreference.UpdatedAt,
	}
}

type OrderFactory struct{}

func (f *OrderFactory) NewOrder(
	userID uuid.UUID,
	symbol string,
	quantity float64,
	price float64,
	tradeType string,
) *Order {
	return &Order{
		ID:        uuid.New(),
		UserID:    userID,
		Symbol:    symbol,
		Quantity:  quantity,
		Price:     price,
		Status:    constants.OrderStatusOpen,
		TradeType: tradeType,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (f *OrderFactory) Clone(
	order *Order,
	symbol string,
	quantity float64,
	price float64,
	tradeType string,
) *Order {
	return &Order{
		ID:        order.ID,
		UserID:    order.UserID,
		Symbol:    symbol,
		Quantity:  quantity,
		Price:     price,
		Status:    order.Status,
		TradeType: tradeType,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
