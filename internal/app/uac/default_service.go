package uacs

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
)

type DefaultUacService struct{}

func NewDefaultUacService() *DefaultUacService {
	return &DefaultUacService{}
}

func (s DefaultUacService) GetUser(ctx echo.Context) *entities.User {
	return ctx.Get("user").(*entities.User)
}

func (s DefaultUacService) IsResourceOwner(ctx echo.Context, resourceUserID uuid.UUID) error {
	user := s.GetUser(ctx)
	if user.ID != resourceUserID {
		return errors.ErrForbidden
	}
	return nil
}

func (s DefaultUacService) IsAdminUser(ctx echo.Context) error {
	user := s.GetUser(ctx)
	if user.Role != constants.RoleAdmin {
		return errors.ErrForbidden
	}
	return nil
}

func (s DefaultUacService) IsFunctionalUser(ctx echo.Context) error {
	user := s.GetUser(ctx)
	if user.Role != constants.RoleFunctional {
		return errors.ErrForbidden
	}
	return nil
}

func (s DefaultUacService) IsRegularUser(ctx echo.Context) error {
	user := s.GetUser(ctx)
	if user.Role != constants.RoleUser {
		return errors.ErrForbidden
	}
	return nil
}
