package markets

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/market"
)

// Structs

type DefaultMarketDataService struct {
	MarketDataRepository market.MarketDataRepository
	UacService           uacs.UacService
}

// Factories

func NewDefaultMarketDataService(
	MarketDataRepository market.MarketDataRepository,
	uacService uacs.UacService,
) *DefaultMarketDataService {
	return &DefaultMarketDataService{
		MarketDataRepository: MarketDataRepository,
		UacService:           uacService,
	}
}

// MarketService implementation

func (s *DefaultMarketDataService) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.MarketData, error) {
	market, err := s.MarketDataRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return market, nil
}

func (s *DefaultMarketDataService) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.MarketDatas, error) {
	filters.SetMetaParameters()
	return s.MarketDataRepository.GetAll(ctx, filters)
}

func (s *DefaultMarketDataService) GetLatest(
	ctx echo.Context,
	symbol string,
) (*entities.MarketData, error) {
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol": symbol,
		},
		"created_at",
		"desc",
		1,
		1,
	)
	marketDatas, err := s.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(*marketDatas) == 0 {
		return nil, errors.ErrMarketDataInsufficient
	}
	if (*marketDatas)[0].Timestamp.Before(time.Now().Add(-time.Minute * 2)) {
		return nil, errors.ErrMarketDataTooOld
	}
	return &(*marketDatas)[0], nil
}

func (s *DefaultMarketDataService) Create(
	ctx echo.Context,
	marketData *entities.MarketData,
) (*entities.MarketData, error) {
	err := marketData.Validate()
	if err != nil {
		return nil, err
	}
	timestampLowerBand := marketData.Timestamp.Truncate(time.Minute)
	timestampUpperBand := timestampLowerBand.Add(time.Minute)
	cf := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol":         marketData.Symbol,
			"timestamp__gte": timestampLowerBand,
			"timestamp__lt":  timestampUpperBand,
		},
		"created_at",
		"desc",
		1,
		100,
	)
	marketDatas, err := s.MarketDataRepository.GetAll(ctx, cf)
	if err != nil {
		return nil, err
	}
	if len(*marketDatas) > 0 {
		factory := entities.MarketDataFactory{}
		updatedMarketData := factory.Clone(
			&(*marketDatas)[0],
			marketData.Timestamp,
			marketData.Open,
			marketData.High,
			marketData.Low,
			marketData.Close,
			marketData.Volume,
			marketData.MACD,
			marketData.MACDSignal,
			marketData.MACDHist,
			marketData.RSI6,
			marketData.RSI12,
			marketData.RSI24,
			marketData.SMA20,
			marketData.SMA50,
			marketData.SMA200,
			marketData.ATR,
			marketData.BollingerBands,
			marketData.BollingerBandsWidth,
			marketData.BollingerBandsUpper,
			marketData.BollingerBandsLower,
			marketData.OBV,
			marketData.ADX,
			marketData.ADXIndex,
			marketData.ADXPositive,
			marketData.ADXNegative,
			marketData.Score,
		)
		return s.Update(ctx, updatedMarketData)
	}
	return s.MarketDataRepository.Create(ctx, marketData)
}

func (s *DefaultMarketDataService) Update(
	ctx echo.Context,
	marketData *entities.MarketData,
) (*entities.MarketData, error) {
	return s.MarketDataRepository.Update(ctx, marketData)
}

func (s *DefaultMarketDataService) Delete(
	ctx echo.Context,
	id uuid.UUID,
) error {
	return s.MarketDataRepository.Delete(ctx, id)
}

// Indicators implementation

// Opportunity score

func (s *DefaultMarketDataService) GetScores(
	ctx echo.Context,
	symbols []string,
) (valueobjects.SymbolScores, error) {
	scores := make(valueobjects.SymbolScores, len(symbols))
	for i, symbol := range symbols {
		score, err := s.GetSymbolScore(ctx, symbol)
		if err != nil {
			return nil, err
		}
		scores[i] = score
	}
	// Sort by score value in descending order and assign rankings
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	for i, score := range scores {
		score.Ranking = i + 1
	}
	return scores, nil
}

func (s *DefaultMarketDataService) GetSymbolScore(
	ctx echo.Context,
	symbol string,
) (valueobjects.SymbolScore, error) {
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"symbol": symbol,
		},
		"created_at",
		"desc",
		1,
		1,
	)
	marketDatas, err := s.MarketDataRepository.GetAll(ctx, filters)
	if err != nil {
		return valueobjects.SymbolScore{}, err
	}
	if len(*marketDatas) == 0 {
		return valueobjects.SymbolScore{}, errors.ErrMarketDataInsufficient
	}
	timestamp := time.Now().Add(-time.Minute * 2)
	if (*marketDatas)[0].Timestamp.Before(timestamp) {
		return valueobjects.SymbolScore{}, errors.ErrMarketDataTooOld
	}
	score := (*marketDatas)[0].Score
	if score == nil {
		return valueobjects.SymbolScore{}, errors.ErrMarketDataInsufficient
	}
	return valueobjects.SymbolScore{
		Symbol:  symbol,
		Score:   *score,
		Ranking: 0,
	}, nil
}

// Moving Average Convergence Divergence

func (s *DefaultMarketDataService) CalculateMACD(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 26 {

		return nil, errors.ErrInsufficientDataForMACDCalculation
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Create a time series from market data
	series := techan.NewTimeSeries()
	for _, data := range *marketDatas {
		candle := techan.NewCandle(techan.NewTimePeriod(data.Timestamp, time.Minute))
		candle.OpenPrice = big.NewDecimal(data.Open)
		candle.MaxPrice = big.NewDecimal(data.High)
		candle.MinPrice = big.NewDecimal(data.Low)
		candle.ClosePrice = big.NewDecimal(data.Close)
		candle.Volume = big.NewDecimal(data.Volume)
		series.AddCandle(candle)
	}

	// Create close price indicator
	closeIndicator := techan.NewClosePriceIndicator(series)

	// Calculate MACD using techan library
	// MACD parameters: fast period (12), slow period (26)
	macdIndicator := techan.NewMACDIndicator(closeIndicator, 12, 26)

	// Calculate signal line (EMA of MACD, typically 9 periods)
	signalIndicator := techan.NewEMAIndicator(macdIndicator, 9)

	// Get the last index
	lastIndex := len(*marketDatas) - 1

	// Calculate MACD values
	macdValue := macdIndicator.Calculate(lastIndex)
	signalValue := signalIndicator.Calculate(lastIndex)
	histogramValue := macdValue.Sub(signalValue)

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with MACD values (convert from Decimal to float64)
	if !macdValue.IsZero() {
		macdFloat := macdValue.Float()
		latestData.MACD = &macdFloat
	}

	if !signalValue.IsZero() {
		signalFloat := signalValue.Float()
		latestData.MACDSignal = &signalFloat
	}

	if !histogramValue.IsZero() {
		histogramFloat := histogramValue.Float()
		latestData.MACDHist = &histogramFloat
	}

	return &latestData, nil
}

// Relative Strength Index

func (s *DefaultMarketDataService) CalculateRSI(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 24 {
		return nil, errors.ErrInsufficientDataForRSI
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Create a time series from market data
	series := techan.NewTimeSeries()
	for _, data := range *marketDatas {
		candle := techan.NewCandle(techan.NewTimePeriod(data.Timestamp, time.Minute))
		candle.OpenPrice = big.NewDecimal(data.Open)
		candle.MaxPrice = big.NewDecimal(data.High)
		candle.MinPrice = big.NewDecimal(data.Low)
		candle.ClosePrice = big.NewDecimal(data.Close)
		candle.Volume = big.NewDecimal(data.Volume)
		series.AddCandle(candle)
	}

	// Create close price indicator
	closeIndicator := techan.NewClosePriceIndicator(series)

	// Calculate RSI for different periods using the correct techan method
	rsi6Indicator := techan.NewRelativeStrengthIndexIndicator(closeIndicator, 6)
	rsi12Indicator := techan.NewRelativeStrengthIndexIndicator(closeIndicator, 12)
	rsi24Indicator := techan.NewRelativeStrengthIndexIndicator(closeIndicator, 24)

	// Get the last index
	lastIndex := len(*marketDatas) - 1

	// Calculate RSI values
	rsi6Value := rsi6Indicator.Calculate(lastIndex)
	rsi12Value := rsi12Indicator.Calculate(lastIndex)
	rsi24Value := rsi24Indicator.Calculate(lastIndex)

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with RSI values (convert from Decimal to float64)
	if !rsi6Value.IsZero() {
		rsi6Float := rsi6Value.Float()
		latestData.RSI6 = &rsi6Float
	}

	if !rsi12Value.IsZero() {
		rsi12Float := rsi12Value.Float()
		latestData.RSI12 = &rsi12Float
	}

	if !rsi24Value.IsZero() {
		rsi24Float := rsi24Value.Float()
		latestData.RSI24 = &rsi24Float
	}

	return &latestData, nil
}

// Simple Moving Average

func (s *DefaultMarketDataService) CalculateSMA(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 200 {
		return nil, errors.ErrInsufficientDataForSMA
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Create a time series from market data
	series := techan.NewTimeSeries()
	for _, data := range *marketDatas {
		candle := techan.NewCandle(techan.NewTimePeriod(data.Timestamp, time.Minute))
		candle.OpenPrice = big.NewDecimal(data.Open)
		candle.MaxPrice = big.NewDecimal(data.High)
		candle.MinPrice = big.NewDecimal(data.Low)
		candle.ClosePrice = big.NewDecimal(data.Close)
		candle.Volume = big.NewDecimal(data.Volume)
		series.AddCandle(candle)
	}

	// Create close price indicator
	closeIndicator := techan.NewClosePriceIndicator(series)

	// Calculate SMA for different periods
	sma20Indicator := techan.NewSimpleMovingAverage(closeIndicator, 20)
	sma50Indicator := techan.NewSimpleMovingAverage(closeIndicator, 50)
	sma200Indicator := techan.NewSimpleMovingAverage(closeIndicator, 200)

	// Get the last index
	lastIndex := len(*marketDatas) - 1

	// Calculate SMA values
	sma20Value := sma20Indicator.Calculate(lastIndex)
	sma50Value := sma50Indicator.Calculate(lastIndex)
	sma200Value := sma200Indicator.Calculate(lastIndex)

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with SMA values (convert from Decimal to float64)
	if !sma20Value.IsZero() {
		sma20Float := sma20Value.Float()
		latestData.SMA20 = &sma20Float
	}

	if !sma50Value.IsZero() {
		sma50Float := sma50Value.Float()
		latestData.SMA50 = &sma50Float
	}

	if !sma200Value.IsZero() {
		sma200Float := sma200Value.Float()
		latestData.SMA200 = &sma200Float
	}

	return &latestData, nil
}

// Average True Range

func (s *DefaultMarketDataService) CalculateATR(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 14 {
		return nil, errors.ErrInsufficientDataForATR
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Create a time series from market data
	series := techan.NewTimeSeries()
	for _, data := range *marketDatas {
		candle := techan.NewCandle(techan.NewTimePeriod(data.Timestamp, time.Minute))
		candle.OpenPrice = big.NewDecimal(data.Open)
		candle.MaxPrice = big.NewDecimal(data.High)
		candle.MinPrice = big.NewDecimal(data.Low)
		candle.ClosePrice = big.NewDecimal(data.Close)
		candle.Volume = big.NewDecimal(data.Volume)
		series.AddCandle(candle)
	}

	// Calculate ATR (typically 14 periods)
	atrIndicator := techan.NewAverageTrueRangeIndicator(series, 14)

	// Get the last index
	lastIndex := len(*marketDatas) - 1

	// Calculate ATR value
	atrValue := atrIndicator.Calculate(lastIndex)

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with ATR value (convert from Decimal to float64)
	if !atrValue.IsZero() {
		atrFloat := atrValue.Float()
		latestData.ATR = &atrFloat
	}

	return &latestData, nil
}

// Bollinger Bands

func (s *DefaultMarketDataService) CalculateBollingerBands(
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 20 {
		return nil, errors.ErrInsufficientDataForBollingerBands
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Create a time series from market data
	series := techan.NewTimeSeries()
	for _, data := range *marketDatas {
		candle := techan.NewCandle(techan.NewTimePeriod(data.Timestamp, time.Minute))
		candle.OpenPrice = big.NewDecimal(data.Open)
		candle.MaxPrice = big.NewDecimal(data.High)
		candle.MinPrice = big.NewDecimal(data.Low)
		candle.ClosePrice = big.NewDecimal(data.Close)
		candle.Volume = big.NewDecimal(data.Volume)
		series.AddCandle(candle)
	}

	// Create close price indicator
	closeIndicator := techan.NewClosePriceIndicator(series)

	// Calculate Bollinger Bands (typically 20 periods, 2 standard deviations)
	lowerBollingerBandsIndicator := techan.NewBollingerLowerBandIndicator(closeIndicator, 20, 2.0)
	upperBollingerBandsIndicator := techan.NewBollingerUpperBandIndicator(closeIndicator, 20, 2.0)

	// Get the last index
	lastIndex := len(*marketDatas) - 1

	// Calculate Bollinger Bands values
	lowerBollingerBandsValue := lowerBollingerBandsIndicator.Calculate(lastIndex)
	upperBollingerBandsValue := upperBollingerBandsIndicator.Calculate(lastIndex)

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with Bollinger Bands values (convert from Decimal to float64)
	if !lowerBollingerBandsValue.IsZero() {
		// Middle band (SMA)
		middleFloat := (lowerBollingerBandsValue.Float() + upperBollingerBandsValue.Float()) / 2
		latestData.BollingerBands = &middleFloat

		// Upper band
		upperFloat := upperBollingerBandsValue.Float()
		latestData.BollingerBandsUpper = &upperFloat

		// Lower band
		lowerFloat := lowerBollingerBandsValue.Float()
		latestData.BollingerBandsLower = &lowerFloat

		// Band width (upper - lower)
		widthFloat := upperBollingerBandsValue.Sub(lowerBollingerBandsValue).Float()
		latestData.BollingerBandsWidth = &widthFloat
	}

	return &latestData, nil
}

// On Balance Volume

func (s *DefaultMarketDataService) CalculateOBV(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 14 {
		return nil, errors.ErrInsufficientDataForOBV
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Create a time series from market data
	series := techan.NewTimeSeries()
	for _, data := range *marketDatas {
		candle := techan.NewCandle(techan.NewTimePeriod(data.Timestamp, time.Minute))
		candle.OpenPrice = big.NewDecimal(data.Open)
		candle.MaxPrice = big.NewDecimal(data.High)
		candle.MinPrice = big.NewDecimal(data.Low)
		candle.ClosePrice = big.NewDecimal(data.Close)
		candle.Volume = big.NewDecimal(data.Volume)
		series.AddCandle(candle)
	}

	// Calculate OBV (On Balance Volume)
	obvIndicator := techan.NewVolumeIndicator(series)

	// Calculate indicator values
	obvValue := obvIndicator.Calculate(len(*marketDatas) - 1)

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with OBV value (convert from Decimal to float64)
	if !obvValue.IsZero() {
		obvFloat := obvValue.Float()
		latestData.OBV = &obvFloat
	}

	return &latestData, nil
}

// Average Directional Index

func (s *DefaultMarketDataService) CalculateADX(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	if len(*marketDatas) < 14 {
		return nil, errors.ErrInsufficientDataForADX
	}

	// Sort by created_at desc to ensure proper order
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(*marketDatas)-1; i < j; i, j = i+1, j-1 {
		(*marketDatas)[i], (*marketDatas)[j] = (*marketDatas)[j], (*marketDatas)[i]
	}

	// Calculate ADX components
	period := 14
	trueRanges := make([]float64, len(*marketDatas))
	directionalMoves := make([]float64, len(*marketDatas))
	directionalMovesMinus := make([]float64, len(*marketDatas))

	// Calculate True Range and Directional Movement
	for i := 1; i < len(*marketDatas); i++ {
		current := (*marketDatas)[i]
		previous := (*marketDatas)[i-1]

		// True Range = max(high-low, |high-prevClose|, |low-prevClose|)
		tr1 := current.High - current.Low
		tr2 := math.Abs(current.High - previous.Close)
		tr3 := math.Abs(current.Low - previous.Close)
		trueRanges[i] = math.Max(tr1, math.Max(tr2, tr3))

		// Directional Movement
		upMove := current.High - previous.High
		downMove := previous.Low - current.Low

		if upMove > downMove && upMove > 0 {
			directionalMoves[i] = upMove
			directionalMovesMinus[i] = 0
		} else if downMove > upMove && downMove > 0 {
			directionalMoves[i] = 0
			directionalMovesMinus[i] = downMove
		} else {
			directionalMoves[i] = 0
			directionalMovesMinus[i] = 0
		}
	}

	// Calculate smoothed averages
	atr := s.calculateEMA(trueRanges[1:], period)
	plusDI := s.calculateEMA(directionalMoves[1:], period)
	minusDI := s.calculateEMA(directionalMovesMinus[1:], period)

	// Calculate DI+ and DI-
	plusDIPercent := (plusDI / atr) * 100
	minusDIPercent := (minusDI / atr) * 100

	// Calculate DX (Directional Index)
	dx := math.Abs(plusDIPercent-minusDIPercent) / (plusDIPercent + minusDIPercent) * 100

	// Calculate ADX (Average of DX over the period)
	// For simplicity, we'll use the current DX value
	// In a full implementation, you'd calculate the EMA of DX values
	adx := dx

	// Get the last market data entry (most recent)
	latestData := (*marketDatas)[len(*marketDatas)-1]

	// Update with ADX values
	latestData.ADX = &adx
	latestData.ADXIndex = &adx
	latestData.ADXPositive = &plusDIPercent
	latestData.ADXNegative = &minusDIPercent

	return &latestData, nil
}

// Helper function to calculate EMA (Exponential Moving Average)
func (s *DefaultMarketDataService) calculateEMA(values []float64, period int) float64 {
	if len(values) == 0 {
		return 0
	}

	multiplier := 2.0 / float64(period+1)
	ema := values[0]

	for i := 1; i < len(values); i++ {
		ema = (values[i] * multiplier) + (ema * (1 - multiplier))
	}

	return ema
}

// General Technical Indicators Caculation

func (s *DefaultMarketDataService) CalculateGeneralTechnicalIndicators(
	ctx echo.Context,
	marketDatas *entities.MarketDatas,
) (*entities.MarketData, error) {
	macd, err := s.CalculateMACD(ctx, marketDatas)
	if err != nil {
		return nil, err
	}

	rsi, err := s.CalculateRSI(ctx, marketDatas)
	if err != nil {
		return nil, err
	}

	sma, err := s.CalculateSMA(ctx, marketDatas)
	if err != nil {
		return nil, err
	}

	atr, err := s.CalculateATR(ctx, marketDatas)
	if err != nil {
		return nil, err
	}

	bollingerBands, err := s.CalculateBollingerBands(marketDatas)
	if err != nil {
		return nil, err
	}

	obv, err := s.CalculateOBV(ctx, marketDatas)
	if err != nil {
		return nil, err
	}

	adx, err := s.CalculateADX(ctx, marketDatas)
	if err != nil {
		return nil, err
	}

	// Sort by created_at desc to ensure proper order
	// Take the last market data entry (most recent)
	sort.Slice(*marketDatas, func(i, j int) bool {
		return (*marketDatas)[i].CreatedAt.After((*marketDatas)[j].CreatedAt)
	})
	latestMarketData := (*marketDatas)[len(*marketDatas)-1]

	// Assign technical indicators to the latest market data entry
	latestMarketData.MACD = macd.MACD
	latestMarketData.MACDSignal = macd.MACDSignal
	latestMarketData.MACDHist = macd.MACDHist
	latestMarketData.RSI6 = rsi.RSI6
	latestMarketData.RSI12 = rsi.RSI12
	latestMarketData.RSI24 = rsi.RSI24
	latestMarketData.SMA20 = sma.SMA20
	latestMarketData.SMA50 = sma.SMA50
	latestMarketData.SMA200 = sma.SMA200
	latestMarketData.ATR = atr.ATR
	latestMarketData.BollingerBands = bollingerBands.BollingerBands
	latestMarketData.BollingerBandsUpper = bollingerBands.BollingerBandsUpper
	latestMarketData.BollingerBandsLower = bollingerBands.BollingerBandsLower
	latestMarketData.BollingerBandsWidth = bollingerBands.BollingerBandsWidth
	latestMarketData.OBV = obv.OBV
	latestMarketData.ADX = adx.ADX
	latestMarketData.ADXIndex = adx.ADXIndex
	latestMarketData.ADXPositive = adx.ADXPositive
	latestMarketData.ADXNegative = adx.ADXNegative

	return &latestMarketData, nil
}

// Opportunity Score calculations

func (s *DefaultMarketDataService) CalculateOpportunityScore(
	ctx echo.Context,
	marketData *entities.MarketData,
) (*entities.MarketData, error) {
	score := 0.0
	totalWeight := 0.0

	// 1. MACD Analysis (Weight: 20%)
	if marketData.MACD != nil && marketData.MACDSignal != nil && marketData.MACDHist != nil {
		macdScore := s.CalculateMACDScore(*marketData.MACD, *marketData.MACDSignal, *marketData.MACDHist)
		score += macdScore * 0.20
		totalWeight += 0.20
	}

	// 2. RSI Analysis (Weight: 15%)
	if marketData.RSI6 != nil && marketData.RSI12 != nil && marketData.RSI24 != nil {
		rsiScore := s.CalculateRSIScore(*marketData.RSI6, *marketData.RSI12, *marketData.RSI24)
		score += rsiScore * 0.15
		totalWeight += 0.15
	}

	// 3. Moving Average Analysis (Weight: 20%)
	if marketData.SMA20 != nil && marketData.SMA50 != nil && marketData.SMA200 != nil {
		smaScore := s.CalculateSMAScore(marketData.Close, *marketData.SMA20, *marketData.SMA50, *marketData.SMA200)
		score += smaScore * 0.20
		totalWeight += 0.20
	}

	// 4. Bollinger Bands Analysis (Weight: 15%)
	if marketData.BollingerBandsUpper != nil && marketData.BollingerBandsLower != nil && marketData.BollingerBandsWidth != nil {
		bbScore := s.CalculateBollingerBandsScore(marketData.Close, *marketData.BollingerBandsUpper, *marketData.BollingerBandsLower, *marketData.BollingerBandsWidth)
		score += bbScore * 0.15
		totalWeight += 0.15
	}

	// 5. Volume Analysis (Weight: 10%)
	if marketData.OBV != nil {
		volumeScore := s.CalculateVolumeScore(*marketData.OBV, marketData.Volume)
		score += volumeScore * 0.10
		totalWeight += 0.10
	}

	// 6. Trend Strength Analysis (Weight: 10%)
	if marketData.ADX != nil && marketData.ADXPositive != nil && marketData.ADXNegative != nil {
		trendScore := s.CalculateTrendScore(*marketData.ADX, *marketData.ADXPositive, *marketData.ADXNegative)
		score += trendScore * 0.10
		totalWeight += 0.10
	}

	// 7. Volatility Analysis (Weight: 10%)
	if marketData.ATR != nil {
		volatilityScore := s.CalculateVolatilityScore(*marketData.ATR, marketData.Close)
		score += volatilityScore * 0.10
		totalWeight += 0.10
	}

	// Normalize score to 0-100 range
	if totalWeight > 0 {
		score = (score / totalWeight) * 100
	}

	// Ensure score is within bounds
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Create a copy of the market data with the opportunity score
	result := *marketData
	// Note: If you want to store the opportunity score in the database,
	// you would need to add an OpportunityScore field to the MarketData entity

	return &result, nil
}

// Helper methods for calculating individual component scores

func (s *DefaultMarketDataService) CalculateMACDScore(macd, signal, histogram float64) float64 {
	score := 0.0

	// MACD crossover analysis
	if macd > signal {
		score += 30 // Bullish crossover
	} else {
		score += 10 // Bearish crossover
	}

	// Histogram analysis
	if histogram > 0 {
		// If histogram is weak, no score
		if histogram < 0.1 {
			score += 0
		} else {
			score += 20 // Positive histogram (bullish momentum)
		}
	} else {
		// No score
		score += 0
	}

	// MACD strength analysis
	macdStrength := math.Abs(macd)
	if macdStrength > 0.5 {
		score += 25 // Strong MACD signal
	} else if macdStrength > 0.2 {
		score += 15 // Moderate MACD signal
	} else {
		// No score
		score += 0
	}

	return score / 100.0 // Normalize to 0-1
}

func (s *DefaultMarketDataService) CalculateRSIScore(rsi6, rsi12, rsi24 float64) float64 {
	score := 0.0

	// RSI oversold/overbought analysis
	if rsi6 < 30 && rsi12 < 35 && rsi24 < 40 {
		score += 50 // Strong oversold condition (bullish opportunity)
	} else if rsi6 > 70 && rsi12 > 65 && rsi24 > 60 {
		score += 0 // Strong overbought condition (bearish)
	} else if rsi6 < 40 && rsi12 < 45 {
		score += 25 // Moderate oversold condition
	} else if rsi6 > 60 && rsi12 > 55 {
		score += 10 // Moderate overbought condition
	} else {
		score += 20 // Neutral condition
	}

	// RSI divergence analysis (simplified)
	if rsi6 > rsi12 && rsi12 > rsi24 {
		score += 30 // Bullish alignment
	} else if rsi6 < rsi12 && rsi12 < rsi24 {
		score += 10 // Bearish alignment
	} else {
		score += 20 // Mixed signals
	}

	return score / 100.0 // Normalize to 0-1
}

func (s *DefaultMarketDataService) CalculateSMAScore(close, sma20, sma50, sma200 float64) float64 {
	score := 0.0

	// Price relative to moving averages
	if close > sma20 && close > sma50 && close > sma200 {
		score += 40 // Price above all MAs (strong bullish)
	} else if close > sma20 && close > sma50 {
		score += 30 // Price above short-term MAs
	} else if close > sma20 {
		score += 20 // Price above 20-day MA only
	} else if close < sma20 && close < sma50 && close < sma200 {
		score += 5 // Price below all MAs (bearish)
	} else {
		score += 15 // Mixed signals
	}

	// Moving average alignment
	if sma20 > sma50 && sma50 > sma200 {
		score += 30 // Bullish alignment (golden cross)
	} else if sma20 < sma50 && sma50 < sma200 {
		score += 10 // Bearish alignment (death cross)
	} else {
		score += 20 // Mixed alignment
	}

	// Distance from moving averages (mean reversion opportunity)
	avgDistance := (math.Abs(close-sma20) + math.Abs(close-sma50) + math.Abs(close-sma200)) / 3
	normalizedDistance := avgDistance / close
	if normalizedDistance > 0.05 {
		score += 20 // Significant deviation (mean reversion opportunity)
	} else if normalizedDistance > 0.02 {
		score += 10 // Moderate deviation
	} else {
		score += 5 // Small deviation
	}

	return score / 100.0 // Normalize to 0-1
}

func (s *DefaultMarketDataService) CalculateBollingerBandsScore(close, upper, lower, width float64) float64 {
	score := 0.0

	// Position within Bollinger Bands
	bandRange := upper - lower
	if bandRange > 0 {
		position := (close - lower) / bandRange

		if position < 0.2 {
			score += 40 // Near lower band (oversold, bullish opportunity)
		} else if position > 0.8 {
			score += 10 // Near upper band (overbought, bearish)
		} else if position > 0.4 && position < 0.6 {
			score += 25 // Middle of bands (neutral)
		} else {
			score += 15 // Between middle and extremes
		}
	}

	// Bollinger Band width (volatility)
	normalizedWidth := width / close
	if normalizedWidth > 0.05 {
		score += 30 // High volatility (opportunity for large moves)
	} else if normalizedWidth > 0.03 {
		score += 20 // Moderate volatility
	} else {
		score += 10 // Low volatility
	}

	// Squeeze detection (low volatility before breakout)
	if normalizedWidth < 0.02 {
		score += 20 // Potential squeeze (breakout opportunity)
	}

	return score / 100.0 // Normalize to 0-1
}

func (s *DefaultMarketDataService) CalculateVolumeScore(obv, volume float64) float64 {
	score := 0.0

	// OBV trend analysis (simplified)
	if obv > 0 {
		score += 30 // Positive OBV (bullish volume)
	} else {
		score += 10 // Negative OBV (bearish volume)
	}

	// Volume relative analysis (would need historical data for better analysis)
	// For now, we'll use a simple heuristic
	if volume > 1000 {
		score += 25 // High volume (stronger signals)
	} else if volume > 500 {
		score += 15 // Moderate volume
	} else {
		score += 5 // Low volume
	}

	// Volume-price relationship (simplified)
	// In a real implementation, you'd compare current volume to average volume
	score += 20 // Neutral assumption

	return score / 100.0 // Normalize to 0-1
}

func (s *DefaultMarketDataService) CalculateTrendScore(adx, adxPositive, adxNegative float64) float64 {
	score := 0.0

	// ADX strength (trend strength)
	if adx > 25 {
		score += 40 // Strong trend
	} else if adx > 20 {
		score += 25 // Moderate trend
	} else {
		score += 10 // Weak trend
	}

	// Directional movement
	if adxPositive > adxNegative {
		score += 30 // Bullish trend
	} else {
		score += 10 // Bearish trend
	}

	// Trend strength vs direction balance
	trendStrength := math.Abs(adxPositive - adxNegative)
	if trendStrength > 10 {
		score += 20 // Clear directional bias
	} else if trendStrength > 5 {
		score += 15 // Moderate directional bias
	} else {
		score += 10 // Weak directional bias
	}

	return score / 100.0 // Normalize to 0-1
}

func (s *DefaultMarketDataService) CalculateVolatilityScore(atr, close float64) float64 {
	score := 0.0

	// ATR relative to price
	normalizedATR := atr / close

	if normalizedATR > 0.03 {
		score += 35 // High volatility (opportunity for large moves)
	} else if normalizedATR > 0.02 {
		score += 25 // Moderate volatility
	} else if normalizedATR > 0.01 {
		score += 15 // Low volatility
	} else {
		score += 5 // Very low volatility
	}

	// Volatility opportunity (mean reversion vs trend following)
	// High volatility can indicate both risk and opportunity
	if normalizedATR > 0.04 {
		score += 25 // Very high volatility (high risk/reward)
	} else if normalizedATR > 0.025 {
		score += 20 // High volatility
	} else {
		score += 15 // Lower volatility
	}

	// Volatility stability (would need historical ATR for better analysis)
	score += 20 // Neutral assumption

	return score / 100.0 // Normalize to 0-1
}
