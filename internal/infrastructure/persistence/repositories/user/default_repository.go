package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"gorm.io/gorm"
)

// Structs

type DefaultUserRepository struct {
	Connection *gorm.DB
}

// Factories

func NewDefaultUserRepository(connection *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{Connection: connection}
}

// Repositories

func (dur *DefaultUserRepository) GetByID(ctx echo.Context, id uuid.UUID) (*entities.User, error) {
	var user dtos.User
	result := dur.Connection.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return user.ToEntity(), nil
}

func (dur *DefaultUserRepository) GetByEmail(
	ctx echo.Context,
	email string,
) (*entities.User, error) {
	var user dtos.User
	result := dur.Connection.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return user.ToEntity(), nil
}

func (dur *DefaultUserRepository) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.Users, error) {
	instances := dtos.Users{}
	query := filters.QueryFromFilter(dur.Connection)
	result := query.Find(&instances).
		Order(filters.GetOrdering()).
		Offset(filters.GetPagination().Page).
		Limit(filters.GetPagination().PageSize)
	if result.Error != nil {
		return nil, result.Error
	}
	return instances.ToEntities(), nil
}

func (dur *DefaultUserRepository) Create(
	ctx echo.Context,
	user *entities.User,
) (*entities.User, error) {
	instance := dtos.User{}
	instance.FromEntity(user)
	result := dur.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dur *DefaultUserRepository) Update(
	ctx echo.Context,
	user *entities.User,
) (*entities.User, error) {
	instance := dtos.User{}
	instance.FromEntity(user)
	result := dur.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dur *DefaultUserRepository) Delete(
	ctx echo.Context,
	id uuid.UUID,
) error {
	_, err := dur.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return dur.Connection.Delete(&dtos.User{}, id).Error
}
