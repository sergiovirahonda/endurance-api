package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
)

type ApiKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;"`
	Service   string    `gorm:"type:varchar(100);not null;"`
	Key       string    `gorm:"type:varchar(100);not null;"`
	Secret    string    `gorm:"type:varchar(100);not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;"`
}

type ApiKeys []ApiKey

// Receivers

func (a *ApiKey) ToEntity() *entities.ApiKey {
	return &entities.ApiKey{
		ID:        a.ID,
		UserID:    a.UserID,
		Service:   a.Service,
		Key:       a.Key,
		Secret:    a.Secret,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (a *ApiKey) FromEntity(apiKey *entities.ApiKey) {
	a.ID = apiKey.ID
	a.UserID = apiKey.UserID
	a.Service = apiKey.Service
	a.Key = apiKey.Key
	a.Secret = apiKey.Secret
	a.CreatedAt = apiKey.CreatedAt
	a.UpdatedAt = apiKey.UpdatedAt
}

func (a *ApiKeys) ToEntities() *entities.ApiKeys {
	entities := make(entities.ApiKeys, len(*a))
	for i, apiKey := range *a {
		entities[i] = *apiKey.ToEntity()
	}
	return &entities
}
