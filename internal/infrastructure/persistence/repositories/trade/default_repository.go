package trade

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

type DefaultTradingPreferenceRepository struct {
	Connection *gorm.DB
}

type DefaultHoldingRepository struct {
	Connection *gorm.DB
}

type DefaultOrderRepository struct {
	Connection *gorm.DB
}

// Factories

func NewDefaultTradingPreferenceRepository(connection *gorm.DB) *DefaultTradingPreferenceRepository {
	return &DefaultTradingPreferenceRepository{Connection: connection}
}

func NewDefaultHoldingRepository(connection *gorm.DB) *DefaultHoldingRepository {
	return &DefaultHoldingRepository{Connection: connection}
}

func NewDefaultOrderRepository(connection *gorm.DB) *DefaultOrderRepository {
	return &DefaultOrderRepository{Connection: connection}
}

// TradingPreferenceRepository implementation

func (dtr *DefaultTradingPreferenceRepository) GetByID(ctx echo.Context, id uuid.UUID) (*entities.TradingPreference, error) {
	var preference dtos.TradingPreference
	result := dtr.Connection.Where("id = ?", id).First(&preference)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return preference.ToEntity(), nil
}

func (dtr *DefaultTradingPreferenceRepository) GetByUserID(ctx echo.Context, userID uuid.UUID) (*entities.TradingPreference, error) {
	var preference dtos.TradingPreference
	result := dtr.Connection.Where("user_id = ?", userID).First(&preference)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return preference.ToEntity(), nil
}

func (dtr *DefaultTradingPreferenceRepository) GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.TradingPreferences, error) {
	instances := dtos.TradingPreferences{}
	query := filters.QueryFromFilter(dtr.Connection)
	result := query.Find(&instances).
		Order(filters.GetOrdering()).
		Offset(filters.GetPagination().Page).
		Limit(filters.GetPagination().PageSize)
	if result.Error != nil {
		return nil, result.Error
	}
	return instances.ToEntities(), nil
}

func (dtr *DefaultTradingPreferenceRepository) Create(ctx echo.Context, preference *entities.TradingPreference) (*entities.TradingPreference, error) {
	instance := dtos.TradingPreference{}
	instance.FromEntity(preference)
	result := dtr.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dtr *DefaultTradingPreferenceRepository) Update(ctx echo.Context, preference *entities.TradingPreference) (*entities.TradingPreference, error) {
	instance := dtos.TradingPreference{}
	instance.FromEntity(preference)
	result := dtr.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dtr *DefaultTradingPreferenceRepository) Delete(ctx echo.Context, id uuid.UUID) error {
	_, err := dtr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return dtr.Connection.Delete(&dtos.TradingPreference{}, id).Error
}

// HoldingRepository implementation

func (dhr *DefaultHoldingRepository) GetByID(ctx echo.Context, id uuid.UUID) (*entities.Holding, error) {
	var holding dtos.Holding
	result := dhr.Connection.Where("id = ?", id).First(&holding)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return holding.ToEntity(), nil
}

func (dhr *DefaultHoldingRepository) GetByUserID(ctx echo.Context, userID uuid.UUID) (*entities.Holdings, error) {
	var holdings dtos.Holdings
	result := dhr.Connection.Where("user_id = ?", userID).Find(&holdings)
	if result.Error != nil {
		return nil, result.Error
	}
	return holdings.ToEntities(), nil
}

func (dhr *DefaultHoldingRepository) GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Holdings, error) {
	instances := dtos.Holdings{}
	query := filters.QueryFromFilter(dhr.Connection)
	result := query.Find(&instances).
		Order(filters.GetOrdering()).
		Offset(filters.GetPagination().Page).
		Limit(filters.GetPagination().PageSize)
	if result.Error != nil {
		return nil, result.Error
	}
	return instances.ToEntities(), nil
}

func (dhr *DefaultHoldingRepository) Create(ctx echo.Context, holding *entities.Holding) (*entities.Holding, error) {
	instance := dtos.Holding{}
	instance.FromEntity(holding)
	result := dhr.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dhr *DefaultHoldingRepository) Update(ctx echo.Context, holding *entities.Holding) (*entities.Holding, error) {
	instance := dtos.Holding{}
	instance.FromEntity(holding)
	result := dhr.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dhr *DefaultHoldingRepository) Delete(ctx echo.Context, id uuid.UUID) error {
	_, err := dhr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return dhr.Connection.Delete(&dtos.Holding{}, id).Error
}

// OrderRepository implementation

func (dor *DefaultOrderRepository) GetByID(ctx echo.Context, id uuid.UUID) (*entities.Order, error) {
	var order dtos.Order
	result := dor.Connection.Where("id = ?", id).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return order.ToEntity(), nil
}

func (dor *DefaultOrderRepository) GetByUserID(ctx echo.Context, userID uuid.UUID) (*entities.Orders, error) {
	var orders dtos.Orders
	result := dor.Connection.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders.ToEntities(), nil
}

func (dor *DefaultOrderRepository) GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Orders, error) {
	instances := dtos.Orders{}
	query := filters.QueryFromFilter(dor.Connection)
	result := query.Find(&instances).
		Order(filters.GetOrdering()).
		Offset(filters.GetPagination().Page).
		Limit(filters.GetPagination().PageSize)
	if result.Error != nil {
		return nil, result.Error
	}
	return instances.ToEntities(), nil
}

func (dor *DefaultOrderRepository) Create(ctx echo.Context, order *entities.Order) (*entities.Order, error) {
	instance := dtos.Order{}
	instance.FromEntity(order)
	result := dor.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dor *DefaultOrderRepository) Update(ctx echo.Context, order *entities.Order) (*entities.Order, error) {
	instance := dtos.Order{}
	instance.FromEntity(order)
	result := dor.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (dor *DefaultOrderRepository) Delete(ctx echo.Context, id uuid.UUID) error {
	_, err := dor.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return dor.Connection.Delete(&dtos.Order{}, id).Error
}
