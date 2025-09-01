package uacs

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
)

type UacService interface {
	GetUser(ctx echo.Context) *entities.User
	IsResourceOwner(ctx echo.Context, resourceUserID uuid.UUID) error
	IsAdminUser(ctx echo.Context) error
	IsFunctionalUser(ctx echo.Context) error
	IsRegularUser(ctx echo.Context) error
}
