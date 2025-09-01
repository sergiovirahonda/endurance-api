package constants

const (
	TradingPreferenceRiskLevelLow    = "low"
	TradingPreferenceRiskLevelMedium = "medium"
	TradingPreferenceRiskLevelHigh   = "high"

	// Trading algorithms
	TradingAlgorithmSwingTrading = "swing_trading"
	TradingAlgorithmScalping     = "scalping"
	TradingAlgorithmDayTrading   = "day_trading"

	// Order statuses
	OrderStatusOpen    = "open"
	OrderStatusFilled  = "filled"
	OrderStatusPending = "pending"

	// Order types
	OrderTypeStopLoss   = "stop_loss"
	OrderTypeTakeProfit = "take_profit"

	// Holding statuses
	HoldingStatusOpen   = "open"
	HoldingStatusClosed = "closed"

	// Pull back trade signals
	PullBackTradeSignalHold = "hold"
	PullBackTradeSignalSell = "sell"
)

var (
	TradingPreferenceAlgorithms = []string{
		TradingAlgorithmSwingTrading,
		TradingAlgorithmScalping,
		TradingAlgorithmDayTrading,
	}
	TradingPreferenceRiskLevels = []string{
		TradingPreferenceRiskLevelLow,
		TradingPreferenceRiskLevelMedium,
		TradingPreferenceRiskLevelHigh,
	}
	OrderStatuses = []string{
		OrderStatusOpen,
		OrderStatusFilled,
		OrderStatusPending,
	}
	OrderTypes = []string{
		OrderTypeStopLoss,
		OrderTypeTakeProfit,
	}
	HoldingStatuses = []string{
		HoldingStatusOpen,
		HoldingStatusClosed,
	}
)
