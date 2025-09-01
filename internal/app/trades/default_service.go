package trades

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	exchanges "github.com/sergiovirahonda/endurance-api/internal/app/exchange"
	"github.com/sergiovirahonda/endurance-api/internal/app/markets"
	"github.com/sergiovirahonda/endurance-api/internal/app/notifications"
	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/aggregate"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/trade"
	"github.com/sergiovirahonda/endurance-api/internal/lib"
)

// Structs

type DefaultTradingService struct {
	TradingPreferenceService *DefaultTradingPreferenceService
	HoldingService           *DefaultHoldingService
	OrderService             *DefaultOrderService
	ExchangeService          *exchanges.DefaultExchangeService
	NotificationService      *notifications.DefaultNotificationService
	MarketDataService        *markets.DefaultMarketDataService
	UacService               uacs.UacService
}

type DefaultTradingPreferenceService struct {
	TradingPreferenceRepository trade.TradingPreferenceRepository
	UacService                  uacs.UacService
}

type DefaultHoldingService struct {
	HoldingRepository trade.HoldingRepository
	UacService        uacs.UacService
}

type DefaultOrderService struct {
	OrderRepository trade.OrderRepository
	UacService      uacs.UacService
}

// Factories

func NewDefaultTradingService(
	tradingPreferenceService *DefaultTradingPreferenceService,
	holdingService *DefaultHoldingService,
	orderService *DefaultOrderService,
	exchangeService *exchanges.DefaultExchangeService,
	notificationService *notifications.DefaultNotificationService,
	uacService uacs.UacService,
) *DefaultTradingService {
	return &DefaultTradingService{
		TradingPreferenceService: tradingPreferenceService,
		HoldingService:           holdingService,
		OrderService:             orderService,
		ExchangeService:          exchangeService,
		NotificationService:      notificationService,
		UacService:               uacService,
	}
}

func NewDefaultTradingPreferenceService(
	tradingPreferenceRepository trade.TradingPreferenceRepository,
	uacService uacs.UacService,
) *DefaultTradingPreferenceService {
	return &DefaultTradingPreferenceService{
		TradingPreferenceRepository: tradingPreferenceRepository,
		UacService:                  uacService,
	}
}

func NewDefaultHoldingService(
	holdingRepository trade.HoldingRepository,
	uacService uacs.UacService,
) *DefaultHoldingService {
	return &DefaultHoldingService{
		HoldingRepository: holdingRepository,
		UacService:        uacService,
	}
}

func NewDefaultOrderService(
	orderRepository trade.OrderRepository,
	uacService uacs.UacService,
) *DefaultOrderService {
	return &DefaultOrderService{
		OrderRepository: orderRepository,
		UacService:      uacService,
	}
}

// Receivers

// Trading Service

func (s *DefaultTradingService) PullBackTrade(
	ctx echo.Context,
	tradingPosition *aggregate.TradingPositionAggregate,
) error {
	logger := config.GetLoggerFromContext(ctx)
	signal, err := s.PullBackTradeSignal(ctx, tradingPosition)
	if err != nil {
		return err
	}
	if signal == constants.PullBackTradeSignalHold {
		logger.Info("No pull back trade signal. Holding position for user %s", tradingPosition.Holding.UserID)
		return nil
	}
	scores, err := s.MarketDataService.GetScores(
		ctx,
		tradingPosition.TradingPreference.Watchlist,
	)
	if err != nil {
		return err
	}
	if scores[0].Symbol == tradingPosition.Holding.Symbol {
		logger.Info("Current holding is the best performing asset in the watchlist")
		return nil
	}
	attractiveSymbol := scores[0].Symbol
	attractiveMarketData, err := s.MarketDataService.GetLatest(ctx, attractiveSymbol)
	if err != nil {
		return err
	}
	attractive, err := s.IsAttractiveSymbol(ctx, tradingPosition, attractiveMarketData)
	if err != nil {
		return err
	}
	if !attractive {
		logger.Info("Attractive symbol does not meet the risk level criteria.")
		// Stop loss if enabled
		if tradingPosition.TradingPreference.StopLossEnabled {
			s.ExecuteStopLoss(
				ctx,
				tradingPosition.Holding,
				"spot", // TODO: Get wallet type from trading preference
			)
		}
		return nil
	}
	// Execute trade
	s.ExecuteTrade(
		ctx,
		tradingPosition.Holding,
		attractiveSymbol,
		attractiveMarketData,
		"spot", // TODO: Get wallet type from trading preference
	)
	return nil
}

func (s *DefaultTradingService) IsAttractiveSymbol(
	ctx echo.Context,
	tradingPosition *aggregate.TradingPositionAggregate,
	toMarketData *entities.MarketData,
) (bool, error) {
	if tradingPosition.TradingPreference.RiskLevel == constants.TradingPreferenceRiskLevelLow {
		if *toMarketData.Score > 0.8 {
			return true, nil
		}
	} else if tradingPosition.TradingPreference.RiskLevel == constants.TradingPreferenceRiskLevelMedium {
		if *toMarketData.Score > 0.7 {
			return true, nil
		}
	} else if tradingPosition.TradingPreference.RiskLevel == constants.TradingPreferenceRiskLevelHigh {
		if *toMarketData.Score > 0.5 {
			return true, nil
		}
	}
	return false, nil
}

func (s *DefaultTradingService) PullBackTradeSignal(
	ctx echo.Context,
	tradingPosition *aggregate.TradingPositionAggregate,
) (string, error) {
	// The higher the score, the slower is the risk, and the higher the chance of profit
	holdingCurrentScore, err := s.MarketDataService.GetSymbolScore(ctx, tradingPosition.Holding.Symbol)
	if err != nil {
		return "", err
	}
	holdingTicker, err := s.ExchangeService.GetTicker(ctx, tradingPosition.Holding.Symbol)
	if err != nil {
		return "", err
	}
	profit := (holdingTicker.Price - tradingPosition.Holding.EntryPrice) * tradingPosition.Holding.Quantity
	profitPercentage := profit / (holdingTicker.Price * tradingPosition.Holding.Quantity) * 100
	if holdingCurrentScore.Score >= tradingPosition.Holding.EntryScore {
		return constants.PullBackTradeSignalHold, nil
	}
	if tradingPosition.TradingPreference.RiskLevel == constants.TradingPreferenceRiskLevelLow {
		if profitPercentage > 5 {
			return constants.PullBackTradeSignalSell, nil
		}
	} else if tradingPosition.TradingPreference.RiskLevel == constants.TradingPreferenceRiskLevelMedium {
		if profitPercentage > 10 {
			return constants.PullBackTradeSignalSell, nil
		}
	} else if tradingPosition.TradingPreference.RiskLevel == constants.TradingPreferenceRiskLevelHigh {
		if profitPercentage > 13 {
			return constants.PullBackTradeSignalSell, nil
		}
	}
	return constants.PullBackTradeSignalHold, nil
}

func (s *DefaultTradingService) GetOpenPositionsForSymbol(
	ctx echo.Context,
	symbol string,
) (*aggregate.TradingPositionAggregates, error) {
	if err := s.UacService.IsRegularUser(ctx); err != nil {
		return nil, err
	}
	tradingPreferencesFilters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"operate": true,
		},
		"created_at",
		"desc",
		1,
		10000,
	)
	tradingPreferences, err := s.TradingPreferenceService.GetAll(ctx, tradingPreferencesFilters)
	if err != nil {
		return nil, err
	}
	userIDs := make([]uuid.UUID, len(*tradingPreferences))
	for i, tp := range *tradingPreferences {
		userIDs[i] = tp.UserID
	}
	holdingFilters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol":      symbol,
			"status":      constants.HoldingStatusOpen,
			"user_id__in": userIDs,
		},
		"created_at",
		"desc",
		1,
		10000,
	)
	holdings, err := s.HoldingService.GetAll(ctx, holdingFilters)
	if err != nil {
		return nil, err
	}
	tradingPositionAggregates := make(aggregate.TradingPositionAggregates, len(*holdings))
	// TODO: Optimize this
	for i, holding := range *holdings {
		for _, tp := range *tradingPreferences {
			if tp.UserID == holding.UserID {
				tradingPositionAggregates[i] = aggregate.TradingPositionAggregate{
					Holding:           &holding,
					TradingPreference: &tp,
				}
			}
		}
	}
	return &tradingPositionAggregates, nil
}

func (s *DefaultTradingService) ExecuteTrade(
	ctx echo.Context,
	holding *entities.Holding,
	toAsset string,
	toMarketData *entities.MarketData,
	walletType string,
) error {
	logger := config.GetLoggerFromContext(ctx)
	user := s.UacService.GetUser(ctx)
	tradingPreference, err := s.TradingPreferenceService.GetByUserID(ctx, user.ID)
	if err != nil {
		return err
	}
	if !tradingPreference.Operate {
		logger.Info("Trading preference is not active")
		return nil
	}
	balance, err := s.ExchangeService.GetBalance(ctx, holding.GetAsset())
	if err != nil {
		return err
	}
	newAssetSymbol := fmt.Sprintf("%s%s", toAsset, "USDT")
	toTicker, err := s.ExchangeService.GetTicker(ctx, newAssetSymbol)
	if err != nil {
		return err
	}
	fromTicker, err := s.ExchangeService.GetTicker(ctx, holding.GetAsset())
	if err != nil {
		return err
	}
	orderFactory := entities.OrderFactory{}
	conversionQuote, err := s.ExchangeService.GetConversionQuote(
		ctx,
		holding.GetAsset(),
		toAsset,
		balance.Free,
		walletType,
	)
	if err != nil {
		return err
	}
	err = conversionQuote.ValidateConversionDrift(
		fromTicker.Price*balance.Free,
		toTicker.Price,
	)
	if err != nil {
		return err
	}
	order := orderFactory.NewOrder(
		user.ID,
		newAssetSymbol,
		conversionQuote.ToAmount,
		conversionQuote.ToAmount*toTicker.Price,
		constants.OrderTypeTakeProfit,
	)
	_, err = s.OrderService.Create(ctx, order)
	if err != nil {
		return err
	}
	_, err = s.ExchangeService.AcceptConversionQuote(ctx, conversionQuote.ID)
	if err != nil {
		return err
	}
	holding.ExitPrice = fromTicker.Price
	profit := (holding.ExitPrice - holding.EntryPrice) * holding.Quantity
	holding.Status = constants.HoldingStatusClosed
	holding.Profit = profit
	s.HoldingService.Update(ctx, holding)
	holdingFactory := entities.HoldingFactory{}
	newHolding := holdingFactory.NewHolding(
		user.ID,
		newAssetSymbol,
		conversionQuote.ToAmount,
		toTicker.Price,
		0,
		0,
		*toMarketData.Score,
		constants.HoldingStatusOpen,
	)
	s.HoldingService.Create(ctx, newHolding)
	order.Status = constants.OrderStatusFilled
	s.OrderService.Update(ctx, order)
	profit = (holding.ExitPrice * holding.Quantity) - (holding.ExitPrice * holding.Quantity)
	profitPercentage := profit / (holding.ExitPrice * holding.Quantity) * 100
	// Send notification
	s.NotificationService.SendTradeNotification(
		ctx,
		holding.Symbol,
		newAssetSymbol,
		fromTicker.Price,
		profit,
		profitPercentage,
	)
	return nil
}

func (s *DefaultTradingService) ExecuteStopLoss(
	ctx echo.Context,
	holding *entities.Holding,
	walletType string,
) error {
	logger := config.GetLoggerFromContext(ctx)
	user := s.UacService.GetUser(ctx)
	tradingPreference, err := s.TradingPreferenceService.GetByUserID(ctx, user.ID)
	if err != nil {
		return err
	}
	if !tradingPreference.Operate {
		logger.Info("Trading preference is not active")
		return nil
	}
	balance, err := s.ExchangeService.GetBalance(ctx, holding.GetAsset())
	if err != nil {
		return err
	}
	fromTicker, err := s.ExchangeService.GetTicker(ctx, holding.GetAsset())
	if err != nil {
		return err
	}
	toTicker, err := s.ExchangeService.GetTicker(ctx, "USDT")
	if err != nil {
		return err
	}
	orderFactory := entities.OrderFactory{}
	conversionQuote, err := s.ExchangeService.GetConversionQuote(
		ctx,
		holding.GetAsset(),
		"USDT",
		balance.Free,
		walletType,
	)
	if err != nil {
		return err
	}
	order := orderFactory.NewOrder(
		user.ID,
		"USDT",
		conversionQuote.ToAmount,
		conversionQuote.ToAmount*toTicker.Price,
		constants.OrderTypeStopLoss,
	)
	_, err = s.OrderService.Create(ctx, order)
	if err != nil {
		return err
	}
	_, err = s.ExchangeService.AcceptConversionQuote(ctx, conversionQuote.ID)
	if err != nil {
		return err
	}
	holding.ExitPrice = fromTicker.Price
	profit := (holding.ExitPrice * holding.Quantity) - (holding.EntryPrice * holding.Quantity)
	holding.Status = constants.HoldingStatusClosed
	holding.Profit = profit
	s.HoldingService.Update(ctx, holding)
	holdingFactory := entities.HoldingFactory{}
	newHolding := holdingFactory.NewHolding(
		user.ID,
		holding.GetAsset(),
		conversionQuote.ToAmount,
		fromTicker.Price,
		0,
		0,
		holding.EntryScore,
		constants.HoldingStatusOpen,
	)
	s.HoldingService.Create(ctx, newHolding)
	order.Status = constants.OrderStatusFilled
	s.OrderService.Update(ctx, order)
	profit = (holding.ExitPrice * holding.Quantity) - (holding.ExitPrice * holding.Quantity)
	profitPercentage := profit / (holding.ExitPrice * holding.Quantity) * 100
	// Send notification
	s.NotificationService.SendStopLossNotification(
		ctx,
		holding.Symbol,
		fromTicker.Price,
		profit,
		profitPercentage,
	)
	return nil
}

// Trading Preference Service

func (s *DefaultTradingPreferenceService) GetByUserID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.TradingPreference, error) {
	tp, err := s.TradingPreferenceRepository.GetByUserID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, tp.UserID); err != nil {
		return nil, err
	}
	return tp, nil
}

func (s *DefaultTradingPreferenceService) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.TradingPreference, error) {
	tp, err := s.TradingPreferenceRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, tp.UserID); err != nil {
		return nil, err
	}
	return tp, nil
}

func (s *DefaultTradingPreferenceService) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.TradingPreferences, error) {
	filters.SetMetaParameters()
	filters.NarrowUserFilters("user_id")
	return s.TradingPreferenceRepository.GetAll(ctx, filters)
}

func (s *DefaultTradingPreferenceService) Create(
	ctx echo.Context,
	tp *entities.TradingPreference,
) (*entities.TradingPreference, error) {
	_, err := s.GetByUserID(ctx, tp.UserID)
	if err == nil {
		return nil, errors.ErrTradingPreferenceAlreadyExists
	}
	if err := s.UacService.IsResourceOwner(ctx, tp.UserID); err != nil {
		return nil, err
	}
	err = tp.Validate()
	if err != nil {
		return nil, err
	}
	return s.TradingPreferenceRepository.Create(ctx, tp)
}

func (s *DefaultTradingPreferenceService) Update(
	ctx echo.Context,
	entity *entities.TradingPreference,
) (*entities.TradingPreference, error) {
	_, err := s.GetByID(ctx, entity.ID)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, entity.UserID); err != nil {
		return nil, err
	}
	err = entity.Validate()
	if err != nil {
		return nil, err
	}
	return s.TradingPreferenceRepository.Update(ctx, entity)
}

func (s *DefaultTradingPreferenceService) Delete(
	ctx echo.Context,
	id uuid.UUID,
) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.TradingPreferenceRepository.Delete(ctx, id)
}

// Holding Service

func (s *DefaultHoldingService) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.Holding, error) {
	holding, err := s.HoldingRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, holding.UserID); err != nil {
		return nil, err
	}
	return holding, nil
}

func (s *DefaultHoldingService) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.Holdings, error) {
	filters.SetMetaParameters()
	filters.NarrowUserFilters("user_id")
	return s.HoldingRepository.GetAll(ctx, filters)
}

func (s *DefaultHoldingService) Create(
	ctx echo.Context,
	holding *entities.Holding,
) (*entities.Holding, error) {
	err := holding.Validate()
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, holding.UserID); err != nil {
		return nil, err
	}
	return s.HoldingRepository.Create(ctx, holding)
}

func (s *DefaultHoldingService) Update(
	ctx echo.Context,
	entity *entities.Holding,
) (*entities.Holding, error) {
	_, err := s.GetByID(ctx, entity.ID)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, entity.UserID); err != nil {
		return nil, err
	}
	err = entity.Validate()
	if err != nil {
		return nil, err
	}
	return s.HoldingRepository.Update(ctx, entity)
}

func (s *DefaultHoldingService) Delete(
	ctx echo.Context,
	id uuid.UUID,
) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.HoldingRepository.Delete(ctx, id)
}

func (s *DefaultHoldingService) GetBySymbol(
	ctx echo.Context,
	symbol string,
) (*entities.Holdings, error) {
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol": symbol,
		},
		"created_at",
		"desc",
		1,
		100,
	)
	holdings, err := s.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(*holdings) == 0 {
		return nil, errors.ErrHoldingNotFound
	}
	return holdings, nil
}

func (s *DefaultHoldingService) GetBySymbolAndStatus(
	ctx echo.Context,
	symbol string,
	status string,
) (*entities.Holdings, error) {
	if !lib.SliceContains(constants.HoldingStatuses, status) {
		return nil, errors.ErrInvalidHoldingStatus
	}
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol": symbol,
			"status": status,
		},
		"created_at",
		"desc",
		1,
		1000,
	)
	return s.GetAll(ctx, filters)
}

// Order Service

func (s *DefaultOrderService) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.Order, error) {
	order, err := s.OrderRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, order.UserID); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *DefaultOrderService) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.Orders, error) {
	filters.SetMetaParameters()
	filters.NarrowUserFilters("user_id")
	return s.OrderRepository.GetAll(ctx, filters)
}

func (s *DefaultOrderService) Create(
	ctx echo.Context,
	order *entities.Order,
) (*entities.Order, error) {
	err := order.Validate()
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, order.UserID); err != nil {
		return nil, err
	}
	return s.OrderRepository.Create(ctx, order)
}

func (s *DefaultOrderService) Update(
	ctx echo.Context,
	entity *entities.Order,
) (*entities.Order, error) {
	_, err := s.GetByID(ctx, entity.ID)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsResourceOwner(ctx, entity.UserID); err != nil {
		return nil, err
	}
	err = entity.Validate()
	if err != nil {
		return nil, err
	}
	return s.OrderRepository.Update(ctx, entity)
}

func (s *DefaultOrderService) Delete(
	ctx echo.Context,
	id uuid.UUID,
) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.OrderRepository.Delete(ctx, id)
}
