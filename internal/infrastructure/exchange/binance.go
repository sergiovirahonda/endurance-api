package exchange

import (
	binance "github.com/adshao/go-binance/v2"
	binanceSapiConnector "github.com/binance/binance-connector-go"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
)

func NewSapiClient(
	cfg *config.Config,
	credentials *valueobjects.ExchangeCredentials,
) *binanceSapiConnector.Client {
	client := binanceSapiConnector.NewClient(
		credentials.APIKey,
		credentials.APISecret,
		cfg.Binance.BaseURL,
	)
	return client
}

func NewGeneralClient(
	cfg *config.Config,
) *binance.Client {
	client := binance.NewClient(
		cfg.Binance.APIKey,
		cfg.Binance.APISecret,
	)
	return client
}
