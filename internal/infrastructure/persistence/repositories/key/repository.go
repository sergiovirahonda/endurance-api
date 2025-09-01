package key

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type KeyRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.ApiKey, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.ApiKeys, error)
	Create(ctx echo.Context, apiKey *entities.ApiKey) (*entities.ApiKey, error)
	Update(ctx echo.Context, apiKey *entities.ApiKey) (*entities.ApiKey, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}
