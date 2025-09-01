package exchanges

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	binance "github.com/adshao/go-binance/v2"
	binanceSapiConnector "github.com/binance/binance-connector-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/domain/events"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/pubsub"
)

// Structs

type DefaultExchangeService struct {
	sapiClient    *binanceSapiConnector.Client
	generalClient *binance.Client
}

type DefaultExchangeDataService struct {
	sapiClient    *binanceSapiConnector.Client
	generalClient *binance.Client
}

type DefaultExchangeWebSocketService struct {
	ps *pubsub.EventsPubSub
}

// Factories

func NewDefaultExchangeService(
	sapiClient *binanceSapiConnector.Client,
	generalClient *binance.Client,
) *DefaultExchangeService {
	return &DefaultExchangeService{
		sapiClient:    sapiClient,
		generalClient: generalClient,
	}
}

func NewDefaultExchangeDataService(
	sapiClient *binanceSapiConnector.Client,
	generalClient *binance.Client,
) *DefaultExchangeDataService {
	return &DefaultExchangeDataService{
		sapiClient:    sapiClient,
		generalClient: generalClient,
	}
}

func NewDefaultExchangeWebSocketService(
	ps *pubsub.EventsPubSub,
) *DefaultExchangeWebSocketService {
	return &DefaultExchangeWebSocketService{
		ps: ps,
	}
}

// ExchangeService implementation

func (s *DefaultExchangeService) GetBalance(
	ctx echo.Context,
	asset string,
) (*valueobjects.ExchangeBalance, error) {
	balances, err := s.GetBalances(ctx)
	if err != nil {
		return nil, err
	}
	for _, balance := range *balances {
		if balance.Asset == asset {
			return &balance, nil
		}
	}
	return nil, errors.ErrBalanceNotFound
}

func (s *DefaultExchangeService) GetBalances(
	ctx echo.Context,
) (*[]valueobjects.ExchangeBalance, error) {
	logger := config.GetLoggerFromContext(ctx)
	account, err := s.generalClient.
		NewGetAccountService().
		Do(ctx.Request().Context())
	logger.Infof("Account: %+v", account)
	if err != nil {
		logger.Errorf("Error getting account: %s", err)
		return nil, errors.ErrAccountNotAvailable
	}
	if len(account.Balances) == 0 {
		logger.Errorf("No balances available")
		return nil, errors.ErrNoBalance
	}
	balances := make([]valueobjects.ExchangeBalance, 0)
	for _, balance := range account.Balances {
		if balance.Free == "0" && balance.Locked == "0" {
			continue
		}
		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			logger.Errorf("Error parsing free balance: %s", err)
			return nil, errors.ErrInvalidBalance
		}
		locked, err := strconv.ParseFloat(balance.Locked, 64)
		if err != nil {
			logger.Errorf("Error parsing locked balance: %s", err)
			return nil, errors.ErrInvalidBalance
		}
		balances = append(balances, valueobjects.ExchangeBalance{
			Asset:  balance.Asset,
			Free:   free,
			Locked: locked,
		})
	}
	return &balances, nil
}

func (s *DefaultExchangeService) GetTicker(
	ctx echo.Context,
	symbol string,
) (*valueobjects.ExchangeTicker, error) {
	logger := config.GetLoggerFromContext(ctx)
	tickerService := s.sapiClient.NewTicker24hrService()
	tickers, err := tickerService.Symbol(symbol).Do(ctx.Request().Context())
	if err != nil {
		logger.Errorf("Error getting ticker: %s", err)
		return nil, errors.ErrTickerNotAvailable
	}
	if len(tickers) == 0 {
		logger.Errorf("No ticker found for symbol: %s", symbol)
		return nil, errors.ErrTickerNotAvailable
	}
	ticker := tickers[0]
	price, err := strconv.ParseFloat(ticker.LastPrice, 64)
	if err != nil {
		logger.Errorf("Error parsing price: %s", err)
		return nil, errors.ErrInvalidPrice
	}
	volume, err := strconv.ParseFloat(ticker.Volume, 64)
	if err != nil {
		logger.Errorf("Error parsing volume: %s", err)
		return nil, errors.ErrInvalidVolume
	}
	pricePercentageChange, err := strconv.ParseFloat(ticker.PriceChangePercent, 64)
	if err != nil {
		logger.Errorf("Error parsing price percentage change: %s", err)
		return nil, errors.ErrInvalidPricePercentageChange
	}
	return &valueobjects.ExchangeTicker{
		Symbol:                ticker.Symbol,
		Price:                 price,
		Volume:                volume,
		PricePercentageChange: pricePercentageChange,
	}, nil
}

func (s *DefaultExchangeService) GetAvailableSymbols(
	ctx echo.Context,
) (*[]valueobjects.ExchangeAvailableSymbol, error) {
	logger := config.GetLoggerFromContext(ctx)
	exchangeInfo, err := s.sapiClient.NewExchangeInfoService().Do(ctx.Request().Context())
	if err != nil {
		logger.Errorf("Error getting exchange info: %s", err)
		return nil, errors.ErrExchangeInfoNotAvailable
	}
	var activeSymbols []valueobjects.ExchangeAvailableSymbol
	for _, symbol := range exchangeInfo.Symbols {
		if symbol.Status == "TRADING" {
			as := valueobjects.ExchangeAvailableSymbol{
				Symbol:         symbol.Symbol,
				BaseAsset:      symbol.BaseAsset,
				QuoteAsset:     symbol.QuoteAsset,
				Status:         symbol.Status,
				QuotePrecision: symbol.QuotePrecision,
			}
			activeSymbols = append(activeSymbols, as)
		}
	}
	return &activeSymbols, nil
}

func (s *DefaultExchangeService) GetConversionQuote(
	ctx echo.Context,
	fromAsset string,
	toAsset string,
	fromAmount float64,
	walletType string,
) (*entities.ExchangeConversionQuote, error) {
	logger := config.GetLoggerFromContext(ctx)
	fromAmountStr := strconv.FormatFloat(fromAmount, 'f', -1, 64)
	quote, err := s.generalClient.NewConvertQuoteService().
		FromAsset(fromAsset).
		ToAsset(toAsset).
		FromAmount(fromAmountStr).
		WalletType(walletType).
		Do(ctx.Request().Context())
	if err != nil {
		if strings.Contains(err.Error(), "insufficient balance") {
			return nil, errors.ErrInsufficientBalance
		}
		if strings.Contains(err.Error(), "invalid symbol") {
			return nil, errors.ErrInvalidSymbol
		}
		if strings.Contains(err.Error(), "invalid amount") {
			return nil, errors.ErrInvalidQuoteAmount
		}
		logger.Errorf("Error getting conversion quote: %s", err)
		return nil, errors.ErrConversionQuoteNotAvailable
	}
	ratio, err := strconv.ParseFloat(quote.Ratio, 64)
	if err != nil {
		logger.Errorf("Error parsing ratio: %s", err)
		return nil, errors.ErrInvalidRatio
	}
	inverseRatio, err := strconv.ParseFloat(quote.InverseRatio, 64)
	if err != nil {
		logger.Errorf("Error parsing inverse ratio: %s", err)
		return nil, errors.ErrInvalidInverseRatio
	}
	toAmount, err := strconv.ParseFloat(quote.ToAmount, 64)
	if err != nil {
		logger.Errorf("Error parsing to amount: %s", err)
		return nil, errors.ErrInvalidToAmount
	}
	conversionQuote := entities.ExchangeConversionQuote{
		ID:           quote.QuoteId,
		FromAsset:    fromAsset,
		ToAsset:      toAsset,
		FromAmount:   fromAmount,
		ToAmount:     toAmount,
		Ratio:        ratio,
		InverseRatio: inverseRatio,
		ValidTime:    quote.ValidTime,
	}
	err = conversionQuote.Validate()
	if err != nil {
		return nil, err
	}
	return &conversionQuote, nil
}

func (s *DefaultExchangeService) AcceptConversionQuote(
	ctx echo.Context,
	id string,
) (*entities.ExchangeConversionOrder, error) {
	logger := config.GetLoggerFromContext(ctx)
	quote, err := s.generalClient.NewConvertAcceptQuoteService().
		QuoteId(id).
		Do(ctx.Request().Context())
	if err != nil {
		if strings.Contains(err.Error(), "quote expired") {
			return nil, errors.ErrQuoteExpired
		}
		logger.Errorf("Error accepting conversion quote: %s", err)
		return nil, errors.ErrConversionQuoteNotAvailable
	}
	return &entities.ExchangeConversionOrder{
		ID:        quote.OrderId,
		CreatedAt: time.Unix(quote.CreateTime, 0),
		Status:    quote.OrderStatus,
	}, nil
}

func (s *DefaultExchangeService) ConvertAsset(
	ctx echo.Context,
	userID string,
	fromAsset string,
	toAsset string,
	fromAmount float64,
	walletType string,
) (*entities.ExchangeConversionOrder, error) {
	logger := config.GetLoggerFromContext(ctx)
	quote, err := s.GetConversionQuote(
		ctx,
		fromAsset,
		toAsset,
		fromAmount,
		walletType,
	)
	logger.Infof("Quote for user %s: %+v", userID, quote)
	if err != nil {
		return nil, err
	}
	order, err := s.AcceptConversionQuote(ctx, quote.ID)
	logger.Infof("Order for user %s: %+v", userID, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// ExchangeDataService implementation

func (s *DefaultExchangeDataService) GetKlines(
	ctx echo.Context,
	symbol string,
	interval string,
	from time.Time,
	to time.Time,
) (*[]valueobjects.ExchangeKline, error) {
	logger := config.GetLoggerFromContext(ctx)
	klines, err := s.generalClient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(from.UnixMilli()).
		EndTime(to.UnixMilli()).
		Do(ctx.Request().Context())
	if err != nil {
		logger.Errorf("Error getting klines: %s", err)
		return nil, errors.ErrKlinesNotAvailable
	}
	kls := make([]valueobjects.ExchangeKline, 0)
	for _, kline := range klines {
		open, err := strconv.ParseFloat(kline.Open, 64)
		if err != nil {
			logger.Errorf("Error parsing open: %s", err)
			return nil, errors.ErrInvalidOpen
		}
		high, err := strconv.ParseFloat(kline.High, 64)
		if err != nil {
			logger.Errorf("Error parsing high: %s", err)
			return nil, errors.ErrInvalidHigh
		}
		low, err := strconv.ParseFloat(kline.Low, 64)
		if err != nil {
			logger.Errorf("Error parsing low: %s", err)
			return nil, errors.ErrInvalidLow
		}
		close, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			logger.Errorf("Error parsing close: %s", err)
			return nil, errors.ErrInvalidClose
		}
		volume, err := strconv.ParseFloat(kline.Volume, 64)
		if err != nil {
			logger.Errorf("Error parsing volume: %s", err)
			return nil, errors.ErrInvalidVolume
		}
		kl := valueobjects.ExchangeKline{
			OpenTime: time.Unix(kline.OpenTime, 0),
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			Volume:   volume,
			Trades:   kline.TradeNum,
		}
		kls = append(kls, kl)
	}
	return &kls, nil
}

// ExchangeWebSocketService implementation

func (s *DefaultExchangeWebSocketService) KlineHandler(
	event *binance.WsKlineEvent,
) {
	// Echo context
	ctx := echo.New().NewContext(nil, nil)
	logger := config.GetLogger()
	ctx.Set("logger", logger)
	marketDataEventFactory := events.MarketDataEventFactory{}
	ev, err := marketDataEventFactory.NewRawMarketDataEvent(
		uuid.Nil,
		event.Kline.Symbol,
		time.Unix(event.Kline.StartTime, 0),
		event.Kline.Open,
		event.Kline.High,
		event.Kline.Low,
		event.Kline.Close,
		event.Kline.Volume,
		event.Kline.IsFinal,
	)
	if err != nil {
		logger.Fatalf("Error creating market data: %s", err)
	}
	err = ev.Dispatch(s.ps)
	if err != nil {
		logger.Fatalf("Error dispatching market data event: %s", err)
	}
}

func (s *DefaultExchangeWebSocketService) ErrorHandler(
	err error,
) {
	fmt.Printf("Error: %s\n", err)
}

func (s *DefaultExchangeWebSocketService) Subscribe(
	ctx echo.Context,
	symbols []string,
	interval string,
) {
	logger := config.GetLoggerFromContext(ctx)
	for _, symbol := range symbols {
		_, _, err := binance.WsKlineServe(
			symbol,
			interval,
			s.KlineHandler,
			s.ErrorHandler,
		)
		if err != nil {
			logger.Errorf("Error subscribing to klines: %s", err)
		}
	}
}
