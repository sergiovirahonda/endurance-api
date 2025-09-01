package key

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"gorm.io/gorm"
)

// Structs

type DefaultKeyRepository struct {
	Connection *gorm.DB
}

// Factories

func NewDefaultKeyRepository(connection *gorm.DB) *DefaultKeyRepository {
	return &DefaultKeyRepository{Connection: connection}
}

// KeyRepository implementation

func (d *DefaultKeyRepository) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.ApiKey, error) {
	var key dtos.ApiKey
	result := d.Connection.Where("id = ?", id).First(&key)
	if result.Error != nil {
		return nil, result.Error
	}
	return key.ToEntity(), nil
}

func (d *DefaultKeyRepository) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.ApiKeys, error) {
	instances := dtos.ApiKeys{}
	query := filters.QueryFromFilter(d.Connection)
	result := query.Find(&instances).
		Order(filters.GetOrdering()).
		Offset(filters.GetPagination().Page).
		Limit(filters.GetPagination().PageSize)
	if result.Error != nil {
		return nil, result.Error
	}
	return instances.ToEntities(), nil
}

func (d *DefaultKeyRepository) Create(
	ctx echo.Context,
	apiKey *entities.ApiKey,
) (*entities.ApiKey, error) {
	instance := dtos.ApiKey{}
	instance.FromEntity(apiKey)
	result := d.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (d *DefaultKeyRepository) Update(
	ctx echo.Context,
	apiKey *entities.ApiKey,
) (*entities.ApiKey, error) {
	instance := dtos.ApiKey{}
	instance.FromEntity(apiKey)
	result := d.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (d *DefaultKeyRepository) Delete(ctx echo.Context, id uuid.UUID) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}
	result := d.Connection.Delete(&dtos.ApiKey{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
