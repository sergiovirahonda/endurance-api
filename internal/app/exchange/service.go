package exchanges

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
)

type ExchangeService interface {
	GetBalance(ctx echo.Context, asset string) (*valueobjects.ExchangeBalance, error)
	GetBalances(ctx echo.Context) (*[]valueobjects.ExchangeBalance, error)
	GetTicker(ctx echo.Context, symbol string) (*valueobjects.ExchangeTicker, error)
	GetAvailableSymbols(ctx echo.Context) (*[]valueobjects.ExchangeAvailableSymbol, error)
	GetConversionQuote(ctx echo.Context, fromAsset string, toAsset string, fromAmount float64, walletType string) (*entities.ExchangeConversionQuote, error)
	AcceptConversionQuote(ctx echo.Context, id string) (*entities.ExchangeConversionOrder, error)
	ConvertAsset(ctx echo.Context, userID string, fromAsset string, toAsset string, fromAmount float64, walletType string) (*entities.ExchangeConversionOrder, error)
}

type ExchangeDataService interface {
	GetKlines(ctx echo.Context, symbol string, interval string, from time.Time, to time.Time) (*[]valueobjects.ExchangeKline, error)
}

type ExchangeWebSocketService interface {
	SubscribeToKlines(ctx context.Context, subscription valueobjects.WebSocketSubscription) (<-chan valueobjects.WebSocketKline, error)
	UnsubscribeFromKlines(ctx context.Context, subscription valueobjects.WebSocketSubscription) error
	IsConnected() bool
	Close() error
}
