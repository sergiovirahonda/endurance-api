package market

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type MarketRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.Market, error)
	GetBySymbol(ctx echo.Context, symbol string) (*entities.Market, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Markets, error)
	Create(ctx echo.Context, market *entities.Market) (*entities.Market, error)
	Update(ctx echo.Context, market *entities.Market) (*entities.Market, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}

type MarketDataRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.MarketData, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.MarketDatas, error)
	Create(ctx echo.Context, market *entities.MarketData) (*entities.MarketData, error)
	Update(ctx echo.Context, market *entities.MarketData) (*entities.MarketData, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}
