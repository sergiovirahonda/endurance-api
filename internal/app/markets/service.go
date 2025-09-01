package markets

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type MarketDataService interface {
	GetByID(ctx echo.Context, id uuid.UUID) (*entities.MarketData, error)
	GetAll(ctx echo.Context, filters filtering.ComplexFilters) (*entities.MarketDatas, error)
	Create(ctx echo.Context, marketData *entities.MarketData) (*entities.MarketData, error)
	Update(ctx echo.Context, marketData *entities.MarketData) (*entities.MarketData, error)
	Delete(ctx echo.Context, id uuid.UUID) error
	// Technical indicators
	GetSymbolScore(ctx echo.Context, symbol string) (valueobjects.SymbolScore, error)
	CalculateMACD(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateRSI(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateSMA(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateATR(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateBollingerBands(marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateOBV(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateADX(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateGeneralTechnicalIndicators(ctx echo.Context, marketDatas *entities.MarketDatas) (*entities.MarketData, error)
	CalculateOpportunityScore(ctx echo.Context, marketData *entities.MarketData) (*entities.MarketData, error)
	// Opportunity score helper methods
	CalculateMACDScore(macd, signal, histogram float64) float64
	CalculateRSIScore(rsi6, rsi12, rsi24 float64) float64
	CalculateSMAScore(close, sma20, sma50, sma200 float64) float64
	CalculateBollingerBandsScore(close, upper, lower, width float64) float64
	CalculateVolumeScore(obv, volume float64) float64
	CalculateTrendScore(adx, adxPositive, adxNegative float64) float64
	CalculateVolatilityScore(atr, close float64) float64
}
