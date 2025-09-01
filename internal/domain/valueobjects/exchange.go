package valueobjects

import "time"

type ExchangeCredentials struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

type ExchangeKline struct {
	OpenTime time.Time `json:"open_time"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   float64   `json:"volume"`
	Trades   int64     `json:"trades"`
}

type WebSocketKline struct {
	Symbol    string    `json:"symbol"`
	Interval  string    `json:"interval"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Trades    int64     `json:"trades"`
	IsFinal   bool      `json:"is_final"`
}

type ExchangeBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}

type ExchangeTicker struct {
	Symbol                string  `json:"symbol"`
	Price                 float64 `json:"price"`
	Volume                float64 `json:"volume"`
	PricePercentageChange float64 `json:"price_percentage_change"`
}

type ExchangeAvailableSymbol struct {
	Symbol         string `json:"symbol"`
	BaseAsset      string `json:"base_asset"`
	QuoteAsset     string `json:"quote_asset"`
	Status         string `json:"status"`
	QuotePrecision int64  `json:"quote_precision"`
}

type WebSocketSubscription struct {
	Symbols   []string `json:"symbols"`
	Intervals []string `json:"intervals"`
}
