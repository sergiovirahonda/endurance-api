package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/lib"
)

type ApiKey struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApiKeys []ApiKey

// Validators

func (a *ApiKey) Validate() error {
	if a.UserID == uuid.Nil {
		return errors.ErrApiKeyUserIDRequired
	}
	if a.Service == "" {
		return errors.ErrApiKeyServiceRequired
	}
	if !lib.SliceContains(constants.ApiKeyServiceTypes, a.Service) {
		return errors.ErrApiKeyServiceInvalid
	}
	if a.Key == "" {
		return errors.ErrApiKeyKeyRequired
	}
	if a.Secret == "" {
		return errors.ErrApiKeySecretRequired
	}
	return nil
}

// Factories

type ApiKeyFactory struct{}

func (f *ApiKeyFactory) NewApiKey(
	userID uuid.UUID,
	service string,
	key string,
	secret string,
) *ApiKey {
	return &ApiKey{
		ID:        uuid.New(),
		UserID:    userID,
		Service:   service,
		Key:       key,
		Secret:    secret,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (f *ApiKeyFactory) Clone(
	apiKey *ApiKey,
	service string,
	key string,
	secret string,
) *ApiKey {
	return &ApiKey{
		ID:        apiKey.ID,
		UserID:    apiKey.UserID,
		Service:   service,
		Key:       key,
		Secret:    secret,
		CreatedAt: apiKey.CreatedAt,
		UpdatedAt: apiKey.UpdatedAt,
	}
}
