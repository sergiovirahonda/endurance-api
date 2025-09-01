package entities

import (
	"time"

	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
)

type ExchangeConversionQuote struct {
	ID           string  `json:"id"`
	FromAsset    string  `json:"from_asset"`
	ToAsset      string  `json:"to_asset"`
	FromAmount   float64 `json:"from_amount"`
	ToAmount     float64 `json:"to_amount"`
	Ratio        float64 `json:"ratio"`
	InverseRatio float64 `json:"inverse_ratio"`
	ValidTime    int64   `json:"valid_time"`
	Fee          float64 `json:"fee"`
	FeeAsset     string  `json:"fee_asset"`
}

type ExchangeConversionOrder struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

// Validations

func (e *ExchangeConversionQuote) Validate() error {
	if e.FromAmount < 0 {
		return errors.ErrInvalidFromAmount
	}
	if e.ToAmount < 0 {
		return errors.ErrInvalidToAmount
	}
	if e.Ratio < 0 {
		return errors.ErrInvalidRatio
	}
	if e.InverseRatio < 0 {
		return errors.ErrInvalidInverseRatio
	}
	if e.ValidTime < 0 {
		return errors.ErrInvalidValidTime
	}
	if e.Fee < 0 {
		return errors.ErrInvalidFee
	}
	return nil
}

func (e *ExchangeConversionQuote) ValidateConversionDrift(
	originAmountUSDT float64,
	conversionTicker float64,
) error {
	if originAmountUSDT == 0 {
		return errors.ErrInvalidPrice
	}
	if conversionTicker == 0 {
		return errors.ErrInvalidPrice
	}
	destinationAmountUSDT := conversionTicker * e.ToAmount
	drift := (destinationAmountUSDT - originAmountUSDT) / originAmountUSDT
	if drift > 0.01 {
		return errors.ErrInvalidConversionDrift
	}
	return nil
}
