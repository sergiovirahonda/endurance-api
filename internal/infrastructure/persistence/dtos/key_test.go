package dtos

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestApiKey_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()

	dto := &ApiKey{
		ID:        id,
		UserID:    userId,
		Service:   "test-service",
		Key:       "test-key",
		Secret:    "test-secret",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, userId, entity.UserID)
	assert.Equal(t, "test-service", entity.Service)
	assert.Equal(t, "test-key", entity.Key)
	assert.Equal(t, "test-secret", entity.Secret)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestApiKey_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()

	entity := &entities.ApiKey{
		ID:        id,
		UserID:    userId,
		Service:   "test-service",
		Key:       "test-key",
		Secret:    "test-secret",
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := &ApiKey{}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, userId, dto.UserID)
	assert.Equal(t, "test-service", dto.Service)
	assert.Equal(t, "test-key", dto.Key)
	assert.Equal(t, "test-secret", dto.Secret)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestApiKeys_ToEntities(t *testing.T) {
	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()
	userId := uuid.New()
	now := time.Now()

	dtos := ApiKeys{
		{
			ID:        id1,
			UserID:    userId,
			Service:   "service-1",
			Key:       "key-1",
			Secret:    "secret-1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        id2,
			UserID:    userId,
			Service:   "service-2",
			Key:       "key-2",
			Secret:    "secret-2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Act
	entities := dtos.ToEntities()

	// Assert
	assert.NotNil(t, entities)
	assert.Len(t, *entities, 2)

	// Check first entity
	assert.Equal(t, id1, (*entities)[0].ID)
	assert.Equal(t, userId, (*entities)[0].UserID)
	assert.Equal(t, "service-1", (*entities)[0].Service)
	assert.Equal(t, "key-1", (*entities)[0].Key)
	assert.Equal(t, "secret-1", (*entities)[0].Secret)
	assert.Equal(t, now, (*entities)[0].CreatedAt)
	assert.Equal(t, now, (*entities)[0].UpdatedAt)

	// Check second entity
	assert.Equal(t, id2, (*entities)[1].ID)
	assert.Equal(t, userId, (*entities)[1].UserID)
	assert.Equal(t, "service-2", (*entities)[1].Service)
	assert.Equal(t, "key-2", (*entities)[1].Key)
	assert.Equal(t, "secret-2", (*entities)[1].Secret)
	assert.Equal(t, now, (*entities)[1].CreatedAt)
	assert.Equal(t, now, (*entities)[1].UpdatedAt)
}
