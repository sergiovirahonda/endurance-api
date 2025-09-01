package markets

import (
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/events"
)

type MarketDataEventHandler interface {
	HandleMarketDataPushed(ctx echo.Context, event events.MarketDataEvent) error
	HandlePartialMarketData(ctx echo.Context, event events.MarketDataEvent) error
}

type MarketDataEventRegistry interface {
	HandleEvent(ctx echo.Context, msg []byte) error
}
