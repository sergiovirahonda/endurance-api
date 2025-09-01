package markets

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/domain/events"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/pubsub"
)

// Structs

type DefaultMarketDataEventHandler struct {
	marketDataService MarketDataService
	eventsPubSub      *pubsub.EventsPubSub
}

type DefaultMarketDataEventRegistry struct {
	handlers map[string]func(ctx echo.Context, event events.MarketDataEvent) error
}

// Factories

func NewDefaultMarketDataEventHandler(
	marketDataService MarketDataService,
	eventsPubSub *pubsub.EventsPubSub,
) *DefaultMarketDataEventHandler {
	return &DefaultMarketDataEventHandler{
		marketDataService: marketDataService,
		eventsPubSub:      eventsPubSub,
	}
}

func NewDefaultMarketDataEventRegistry(
	marketDataEventHandler MarketDataEventHandler,
) *DefaultMarketDataEventRegistry {
	return &DefaultMarketDataEventRegistry{
		handlers: map[string]func(ctx echo.Context, event events.MarketDataEvent) error{
			constants.MarketDataPushedEvent:  marketDataEventHandler.HandleMarketDataPushed,
			constants.PartialMarketDataEvent: marketDataEventHandler.HandlePartialMarketData,
		},
	}
}

// Market data event handlers

func (h *DefaultMarketDataEventHandler) HandleMarketDataPushed(
	ctx echo.Context,
	event events.MarketDataEvent,
) error {
	logger := config.GetLoggerFromContext(ctx)
	logger.Info("Handling market data pushed event...")
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"correlation_id": event.ID,
		},
		"created_at",
		"desc",
		0,
		1,
	)
	mks, err := h.marketDataService.GetAll(ctx, filters)
	if err != nil {
		logger.Error("Error getting market data by correlation ID: %s", err)
		return err
	}
	if len(*mks) > 0 {
		// Idempotency check
		logger.Info("Event already processed. Skipping.")
		return nil
	}
	if !event.CandleClose {
		// TODO: Handle partial candles
		logger.Info("Candle not closed. Skipping.")
		return nil
	}
	// Create market data
	factory := entities.MarketDataFactory{}
	marketData := factory.NewMarketDataFromEvent(
		event.ID,
		event.Symbol,
		event.DataTimestamp,
		event.Open,
		event.High,
		event.Low,
		event.Close,
		event.Volume,
	)
	err = marketData.Validate()
	if err != nil {
		logger.Error("Error validating market data: %s", err)
		return err
	}
	_, err = h.marketDataService.Create(ctx, marketData)
	if err != nil {
		logger.Error("Error creating market data: %s", err)
		return err
	}
	marketDataEventFactory := events.MarketDataEventFactory{}
	marketDataEvent := marketDataEventFactory.NewPartialMarketDataEvent(
		marketData.ID,
	)
	err = marketDataEvent.Dispatch(h.eventsPubSub)
	if err != nil {
		logger.Error("Error dispatching market data event: %s", err)
		return err
	}
	return nil
}

func (h *DefaultMarketDataEventHandler) HandlePartialMarketData(
	ctx echo.Context,
	event events.MarketDataEvent,
) error {
	logger := config.GetLoggerFromContext(ctx)
	logger.Info("Handling partial market data event...")
	datapoint, err := h.marketDataService.GetByID(ctx, event.DatapointID)
	if err != nil {
		logger.Error("Error getting datapoint: %s", err)
		return err
	}
	startTimestamp := datapoint.Timestamp.Add(-15 * 24 * time.Hour)
	endTimestamp := datapoint.Timestamp
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol":         datapoint.Symbol,
			"timestamp__gte": startTimestamp,
			"timestamp__lte": endTimestamp,
		},
		"timestamp",
		"asc",
		0,
		1000,
	)
	mks, err := h.marketDataService.GetAll(ctx, filters)
	if err != nil {
		logger.Error("Error getting market data: %s", err)
		return err
	}
	if len(*mks) == 0 {
		return errors.ErrMarketDataInsufficient
	}
	newMarketData, err := h.marketDataService.CalculateGeneralTechnicalIndicators(
		ctx,
		mks,
	)
	if err != nil {
		logger.Error("Error calculating general technical indicators: %s", err)
		return err
	}
	newMarketData, err = h.marketDataService.CalculateOpportunityScore(ctx, newMarketData)
	if err != nil {
		logger.Error("Error calculating opportunity score: %s", err)
		return err
	}
	_, err = h.marketDataService.Update(ctx, newMarketData)
	if err != nil {
		logger.Error("Error updating market data: %s", err)
		return err
	}
	return nil
}

// Main event handler

func (r DefaultMarketDataEventRegistry) HandleEvent(
	ctx echo.Context,
	msg []byte,
) error {
	logger := config.GetLoggerFromContext(ctx)
	marketDataEvent := events.MarketDataEvent{}
	err := json.Unmarshal(msg, &marketDataEvent)
	if err != nil {
		return nil
	}
	logger.Info("Market domain event received: %s", marketDataEvent.Type)
	handler, ok := r.handlers[marketDataEvent.Type]
	if !ok {
		logger.Error("No handler found for market domain event: %s", marketDataEvent.Type)
		return nil
	}
	return handler(ctx, marketDataEvent)
}
