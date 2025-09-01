package trades

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type TradingService interface {
	ExecuteTrade(ctx echo.Context, holding *entities.Holding, toAsset string, toMarketData *entities.MarketData, walletType string) error
	ExecuteStopLoss(ctx echo.Context, holding *entities.Holding, walletType string) error
}

type TradingPreferenceService interface {
	GetByUserID(ctx echo.Context, id uuid.UUID) (*entities.TradingPreference, error)
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.TradingPreference, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.TradingPreferences, error)
	Create(ctx echo.Context, tp *entities.TradingPreference) (*entities.TradingPreference, error)
	Update(ctx echo.Context, entity *entities.TradingPreference) (*entities.TradingPreference, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}

type HoldingService interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.Holding, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Holdings, error)
	Create(ctx echo.Context, holding *entities.Holding) (*entities.Holding, error)
	Update(ctx echo.Context, entity *entities.Holding) (*entities.Holding, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}

type OrderService interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.Order, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Orders, error)
	Create(ctx echo.Context, order *entities.Order) (*entities.Order, error)
	Update(ctx echo.Context, entity *entities.Order) (*entities.Order, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}
