package trade

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type TradingPreferenceRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.TradingPreference, error)
	GetByUserID(ctx echo.Context, userID uuid.UUID) (*entities.TradingPreference, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.TradingPreferences, error)
	Create(ctx echo.Context, preference *entities.TradingPreference) (*entities.TradingPreference, error)
	Update(ctx echo.Context, preference *entities.TradingPreference) (*entities.TradingPreference, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}

type HoldingRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.Holding, error)
	GetByUserID(ctx echo.Context, userID uuid.UUID) (*entities.Holdings, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Holdings, error)
	Create(ctx echo.Context, holding *entities.Holding) (*entities.Holding, error)
	Update(ctx echo.Context, holding *entities.Holding) (*entities.Holding, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}

type OrderRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.Order, error)
	GetByUserID(ctx echo.Context, userID uuid.UUID) (*entities.Orders, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Orders, error)
	Create(ctx echo.Context, order *entities.Order) (*entities.Order, error)
	Update(ctx echo.Context, order *entities.Order) (*entities.Order, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}
