package events

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/pubsub"
)

// Event structures

type MarketDataEvent struct {
	BaseEvent
	DatapointID   uuid.UUID `json:"datapoint_id"`
	Symbol        string    `json:"symbol"`
	DataTimestamp time.Time `json:"data_timestamp"`
	Open          float64   `json:"open"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Close         float64   `json:"close"`
	Volume        float64   `json:"volume"`
	CandleClose   bool      `json:"candle_close"`
}

// Factories

type MarketDataEventFactory struct{}

func NewMarketDataEventFactory() *MarketDataEventFactory {
	return &MarketDataEventFactory{}
}

// Factory receivers

func (f *MarketDataEventFactory) NewMarketDataEvent(
	datapointID uuid.UUID,
	symbol string,
	timestamp time.Time,
	open float64,
	high float64,
	low float64,
	close float64,
	volume float64,
	candleClose bool,
) *MarketDataEvent {
	return &MarketDataEvent{
		BaseEvent: BaseEvent{
			ID:        uuid.New(),
			Domain:    constants.MarketDataEventDomain,
			Type:      constants.IncomingMarketDataEvent,
			Timestamp: time.Now().UTC(),
		},
		DatapointID:   datapointID,
		Symbol:        symbol,
		DataTimestamp: timestamp,
		Open:          open,
		High:          high,
		Low:           low,
		Close:         close,
		Volume:        volume,
		CandleClose:   candleClose,
	}
}

func (f *MarketDataEventFactory) NewRawMarketDataEvent(
	datapointID uuid.UUID,
	symbol string,
	timestamp time.Time,
	open string,
	high string,
	low string,
	close string,
	volume string,
	candleClose bool,
) (*MarketDataEvent, error) {
	openFloat, err := strconv.ParseFloat(open, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketOpen
	}
	highFloat, err := strconv.ParseFloat(high, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketHigh
	}
	lowFloat, err := strconv.ParseFloat(low, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketLow
	}
	closeFloat, err := strconv.ParseFloat(close, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketClose
	}
	volumeFloat, err := strconv.ParseFloat(volume, 64)
	if err != nil {
		return nil, errors.ErrInvalidMarketVolume
	}
	return f.NewMarketDataEvent(
		datapointID,
		symbol,
		timestamp,
		openFloat,
		highFloat,
		lowFloat,
		closeFloat,
		volumeFloat,
		candleClose,
	), nil
}

func (f *MarketDataEventFactory) NewPartialMarketDataEvent(
	datapointID uuid.UUID,
) *MarketDataEvent {
	return &MarketDataEvent{
		BaseEvent: BaseEvent{
			ID:        uuid.New(),
			Domain:    constants.MarketDataEventDomain,
			Type:      constants.PartialMarketDataEvent,
			Timestamp: time.Now().UTC(),
		},
		DatapointID:   datapointID,
		Symbol:        "",
		DataTimestamp: time.Time{},
		Open:          0,
		High:          0,
		Low:           0,
		Close:         0,
		Volume:        0,
		CandleClose:   false,
	}
}

// Domain events receivers

func (e MarketDataEvent) Dispatch(ps *pubsub.EventsPubSub) error {
	conf := config.GetConfig()
	err := ps.Publish(conf.NATS.NatsMarketDataDomainEventsSubject, e)
	if err != nil {
		return err
	}
	return nil
}
