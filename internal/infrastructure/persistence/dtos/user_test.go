package dtos

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestUser_ToEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()
	lastLogin := now.Add(-24 * time.Hour)

	dto := &User{
		ID:        id,
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "John",
		LastName:  "Doe",
		LastLogin: lastLogin,
		LoggedIn:  true,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	entity := dto.ToEntity()

	// Assert
	assert.NotNil(t, entity)
	assert.Equal(t, id, entity.ID)
	assert.Equal(t, "test@example.com", entity.Email)
	assert.Equal(t, "hashed_password", entity.Password)
	assert.Equal(t, "John", entity.FirstName)
	assert.Equal(t, "Doe", entity.LastName)
	assert.Equal(t, lastLogin, entity.LastLogin)
	assert.True(t, entity.LoggedIn)
	assert.True(t, entity.Enabled)
	assert.Equal(t, now, entity.CreatedAt)
	assert.Equal(t, now, entity.UpdatedAt)
}

func TestUser_FromEntity(t *testing.T) {
	// Arrange
	id := uuid.New()
	now := time.Now()
	lastLogin := now.Add(-24 * time.Hour)

	entity := &entities.User{
		ID:        id,
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "John",
		LastName:  "Doe",
		LastLogin: lastLogin,
		LoggedIn:  true,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := &User{}

	// Act
	dto.FromEntity(entity)

	// Assert
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, "test@example.com", dto.Email)
	assert.Equal(t, "hashed_password", dto.Password)
	assert.Equal(t, "John", dto.FirstName)
	assert.Equal(t, "Doe", dto.LastName)
	assert.Equal(t, lastLogin, dto.LastLogin)
	assert.True(t, dto.LoggedIn)
	assert.True(t, dto.Enabled)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

func TestUsers_ToEntities(t *testing.T) {
	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()
	now := time.Now()
	lastLogin1 := now.Add(-24 * time.Hour)
	lastLogin2 := now.Add(-12 * time.Hour)

	dtos := Users{
		{
			ID:        id1,
			Email:     "john@example.com",
			Password:  "hashed_password_1",
			FirstName: "John",
			LastName:  "Doe",
			LastLogin: lastLogin1,
			LoggedIn:  true,
			Enabled:   true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        id2,
			Email:     "jane@example.com",
			Password:  "hashed_password_2",
			FirstName: "Jane",
			LastName:  "Smith",
			LastLogin: lastLogin2,
			LoggedIn:  false,
			Enabled:   true,
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
	assert.Equal(t, "john@example.com", (*entities)[0].Email)
	assert.Equal(t, "hashed_password_1", (*entities)[0].Password)
	assert.Equal(t, "John", (*entities)[0].FirstName)
	assert.Equal(t, "Doe", (*entities)[0].LastName)
	assert.Equal(t, lastLogin1, (*entities)[0].LastLogin)
	assert.True(t, (*entities)[0].LoggedIn)
	assert.True(t, (*entities)[0].Enabled)
	assert.Equal(t, now, (*entities)[0].CreatedAt)
	assert.Equal(t, now, (*entities)[0].UpdatedAt)

	// Check second entity
	assert.Equal(t, id2, (*entities)[1].ID)
	assert.Equal(t, "jane@example.com", (*entities)[1].Email)
	assert.Equal(t, "hashed_password_2", (*entities)[1].Password)
	assert.Equal(t, "Jane", (*entities)[1].FirstName)
	assert.Equal(t, "Smith", (*entities)[1].LastName)
	assert.Equal(t, lastLogin2, (*entities)[1].LastLogin)
	assert.False(t, (*entities)[1].LoggedIn)
	assert.True(t, (*entities)[1].Enabled)
	assert.Equal(t, now, (*entities)[1].CreatedAt)
	assert.Equal(t, now, (*entities)[1].UpdatedAt)
}
