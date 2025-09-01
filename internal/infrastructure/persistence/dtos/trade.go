package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"gorm.io/gorm"
)

type TradingPreference struct {
	gorm.Model
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;"`
	UserID              uuid.UUID      `gorm:"type:uuid;not null;"`
	Algorithm           string         `gorm:"type:varchar(20);not null;"`
	Watchlist           pq.StringArray `gorm:"type:text[]"`
	Operate             bool           `gorm:"type:boolean;not null;default:false;"`
	StopLossEnabled     bool           `gorm:"type:boolean;not null;default:false;"`
	StopLossExitEnabled bool           `gorm:"type:boolean;not null;default:false;"`
	RiskLevel           string         `gorm:"type:varchar(20);not null;default:'low';"`
	CreatedAt           time.Time      `gorm:"type:timestamp;not null;"`
	UpdatedAt           time.Time      `gorm:"type:timestamp;not null;"`
}

type TradingPreferences []TradingPreference

type Holding struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;"`
	Symbol     string    `gorm:"type:varchar(20);not null;"`
	Quantity   float64   `gorm:"type:decimal(10,2);not null;"`
	EntryPrice float64   `gorm:"type:decimal(10,2);not null;"`
	ExitPrice  float64   `gorm:"type:decimal(10,2);"`
	Profit     float64   `gorm:"type:decimal(10,2);"`
	EntryScore float64   `gorm:"type:decimal(10,2);"`
	Status     string    `gorm:"type:varchar(20);not null;default:'open';"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null;"`
}

type Holdings []Holding

type Order struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;"`
	Symbol    string    `gorm:"type:varchar(20);not null;"`
	Quantity  float64   `gorm:"type:decimal(10,2);not null;"`
	Price     float64   `gorm:"type:decimal(10,2);not null;"`
	Status    string    `gorm:"type:varchar(20);not null;"`
	TradeType string    `gorm:"type:varchar(20);not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;"`
}

type Orders []Order

// Receivers

func (t *TradingPreference) ToEntity() *entities.TradingPreference {
	return &entities.TradingPreference{
		ID:                  t.ID,
		UserID:              t.UserID,
		Algorithm:           t.Algorithm,
		Watchlist:           t.Watchlist,
		Operate:             t.Operate,
		StopLossEnabled:     t.StopLossEnabled,
		StopLossExitEnabled: t.StopLossExitEnabled,
		RiskLevel:           t.RiskLevel,
		CreatedAt:           t.CreatedAt,
		UpdatedAt:           t.UpdatedAt,
	}
}

func (t *TradingPreference) FromEntity(tradingPreference *entities.TradingPreference) {
	t.ID = tradingPreference.ID
	t.UserID = tradingPreference.UserID
	t.Algorithm = tradingPreference.Algorithm
	t.Watchlist = tradingPreference.Watchlist
	t.Operate = tradingPreference.Operate
	t.StopLossEnabled = tradingPreference.StopLossEnabled
	t.StopLossExitEnabled = tradingPreference.StopLossExitEnabled
	t.RiskLevel = tradingPreference.RiskLevel
	t.CreatedAt = tradingPreference.CreatedAt
	t.UpdatedAt = tradingPreference.UpdatedAt
}

func (t *TradingPreferences) ToEntities() *entities.TradingPreferences {
	entities := make(entities.TradingPreferences, len(*t))
	for i, tradingPreference := range *t {
		entities[i] = *tradingPreference.ToEntity()
	}
	return &entities
}

func (h *Holding) ToEntity() *entities.Holding {
	return &entities.Holding{
		ID:         h.ID,
		UserID:     h.UserID,
		Symbol:     h.Symbol,
		Quantity:   h.Quantity,
		EntryPrice: h.EntryPrice,
		ExitPrice:  h.ExitPrice,
		Profit:     h.Profit,
		EntryScore: h.EntryScore,
		Status:     h.Status,
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}
}

func (h *Holding) FromEntity(holding *entities.Holding) {
	h.ID = holding.ID
	h.UserID = holding.UserID
	h.Symbol = holding.Symbol
	h.Quantity = holding.Quantity
	h.EntryPrice = holding.EntryPrice
	h.ExitPrice = holding.ExitPrice
	h.Profit = holding.Profit
	h.EntryScore = holding.EntryScore
	h.Status = holding.Status
	h.CreatedAt = holding.CreatedAt
	h.UpdatedAt = holding.UpdatedAt
}

func (h *Holdings) ToEntities() *entities.Holdings {
	entities := make(entities.Holdings, len(*h))
	for i, holding := range *h {
		entities[i] = *holding.ToEntity()
	}
	return &entities
}

func (o *Order) ToEntity() *entities.Order {
	return &entities.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		Symbol:    o.Symbol,
		Quantity:  o.Quantity,
		Price:     o.Price,
		Status:    o.Status,
		TradeType: o.TradeType,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func (o *Order) FromEntity(order *entities.Order) {
	o.ID = order.ID
	o.UserID = order.UserID
	o.Symbol = order.Symbol
	o.Quantity = order.Quantity
	o.Price = order.Price
	o.Status = order.Status
	o.TradeType = order.TradeType
	o.CreatedAt = order.CreatedAt
	o.UpdatedAt = order.UpdatedAt
}

func (o *Orders) ToEntities() *entities.Orders {
	entities := make(entities.Orders, len(*o))
	for i, order := range *o {
		entities[i] = *order.ToEntity()
	}
	return &entities
}
