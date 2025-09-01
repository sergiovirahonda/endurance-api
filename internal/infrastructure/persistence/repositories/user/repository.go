package user

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type UserRepository interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx echo.Context, email string) (*entities.User, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Users, error)
	Create(ctx echo.Context, user *entities.User) (*entities.User, error)
	Update(ctx echo.Context, user *entities.User) (*entities.User, error)
	Delete(ctx echo.Context, id uuid.UUID) error
}
