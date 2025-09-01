package aggregate

import "github.com/sergiovirahonda/endurance-api/internal/domain/entities"

type TradingPositionAggregate struct {
	Holding           *entities.Holding
	TradingPreference *entities.TradingPreference
}

type TradingPositionAggregates []TradingPositionAggregate
