package market

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"gorm.io/gorm"
)

// Structs

type DefaultMarketRepository struct {
	Connection *gorm.DB
}

type DefaultMarketDataRepository struct {
	Connection *gorm.DB
}

// Factories

func NewDefaultMarketRepository(connection *gorm.DB) *DefaultMarketRepository {
	return &DefaultMarketRepository{Connection: connection}
}

func NewDefaultMarketDataRepository(connection *gorm.DB) *DefaultMarketDataRepository {
	return &DefaultMarketDataRepository{Connection: connection}
}

// MarketRepository implementation

func (d *DefaultMarketRepository) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.Market, error) {
	var market dtos.Market
	result := d.Connection.Where("id = ?", id).First(&market)
	if result.Error != nil {
		return nil, result.Error
	}
	return market.ToEntity(), nil
}

func (d *DefaultMarketRepository) GetBySymbol(
	ctx echo.Context,
	symbol string,
) (*entities.Market, error) {
	var market dtos.Market
	result := d.Connection.Where("symbol = ?", symbol).First(&market)
	if result.Error != nil {
		return nil, result.Error
	}
	return market.ToEntity(), nil
}

func (d *DefaultMarketRepository) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.Markets, error) {
	markets := dtos.Markets{}
	query := filters.QueryFromFilter(d.Connection)
	result := query.Find(&markets).
		Order(filters.GetOrdering()).
		Offset(filters.GetPagination().Page).
		Limit(filters.GetPagination().PageSize)
	if result.Error != nil {
		return nil, result.Error
	}
	return markets.ToEntities(), nil
}

func (d *DefaultMarketRepository) Create(
	ctx echo.Context,
	market *entities.Market,
) (*entities.Market, error) {
	instance := dtos.Market{}
	instance.FromEntity(market)
	result := d.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (d *DefaultMarketRepository) Update(
	ctx echo.Context,
	market *entities.Market,
) (*entities.Market, error) {
	instance := dtos.Market{}
	instance.FromEntity(market)
	result := d.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (d *DefaultMarketRepository) Delete(ctx echo.Context, id uuid.UUID) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}
	result := d.Connection.Delete(&dtos.Market{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// MarketDataRepository implementation

func (d *DefaultMarketDataRepository) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.MarketData, error) {
	var market dtos.MarketData
	result := d.Connection.Where("id = ?", id).First(&market)
	if result.Error != nil {
		return nil, result.Error
	}
	return market.ToEntity(), nil
}

func (d *DefaultMarketDataRepository) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.MarketDatas, error) {
	instances := dtos.MarketDatas{}
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

func (d *DefaultMarketDataRepository) Create(
	ctx echo.Context,
	market *entities.MarketData,
) (*entities.MarketData, error) {
	instance := dtos.MarketData{}
	instance.FromEntity(market)
	result := d.Connection.Create(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (d *DefaultMarketDataRepository) Update(
	ctx echo.Context,
	market *entities.MarketData,
) (*entities.MarketData, error) {
	instance := dtos.MarketData{}
	instance.FromEntity(market)
	result := d.Connection.Save(&instance)
	if result.Error != nil {
		return nil, result.Error
	}
	return instance.ToEntity(), nil
}

func (d *DefaultMarketDataRepository) Delete(ctx echo.Context, id uuid.UUID) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}
	result := d.Connection.Delete(&dtos.MarketData{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
