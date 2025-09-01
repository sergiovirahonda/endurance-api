package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUserWithInvalidIDReturnsError(t *testing.T) {
	// Arrange
	id := uuid.New()
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := repository.GetByID(ctx, id)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetUserByIDReturnsUser(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}
	user := userFactory.NewUser(
		"test1@example.com",
		"password",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)
	createdUser, err := repository.Create(ctx, user)
	assert.NoError(t, err)

	// Act
	foundUser, err := repository.GetByID(ctx, createdUser.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, createdUser.Email, foundUser.Email)
	assert.Equal(t, createdUser.FirstName, foundUser.FirstName)
	assert.Equal(t, createdUser.LastName, foundUser.LastName)
}

func TestGetUserByEmailReturnsUser(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}
	user := userFactory.NewUser(
		"test2@example.com",
		"password",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)
	createdUser, err := repository.Create(ctx, user)
	assert.NoError(t, err)

	// Act
	foundUser, err := repository.GetByEmail(ctx, createdUser.Email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, createdUser.Email, foundUser.Email)
	assert.Equal(t, createdUser.FirstName, foundUser.FirstName)
	assert.Equal(t, createdUser.LastName, foundUser.LastName)
}

func TestGetUserByInvalidEmailReturnsError(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)

	// Act
	_, err := repository.GetByEmail(ctx, "nonexistent@example.com")

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetAllUsersReturnsUsers(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}

	// Create multiple users
	user1 := userFactory.NewUser(
		"test3@example.com",
		"password1",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)
	user2 := userFactory.NewUser(
		"test4@example.com",
		"password2",
		"Jane",
		"Smith",
		"admin",
		time.Now(),
		true,
		true,
	)

	_, err := repository.Create(ctx, user1)
	assert.NoError(t, err)
	_, err = repository.Create(ctx, user2)
	assert.NoError(t, err)

	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{}, "created_at", "desc", 0, 10)

	// Act
	users, err := repository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.GreaterOrEqual(t, len(*users), 2)
}

func TestGetAllUsersReturnsUsers_SpecificFilters(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}

	user1 := userFactory.NewUser(
		"user1@example.com",
		"password",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)
	user2 := userFactory.NewUser(
		"user2@example.com",
		"password",
		"Jane",
		"Smith",
		"admin",
		time.Now(),
		true,
		true,
	)

	_, err := repository.Create(ctx, user1)
	assert.NoError(t, err)
	_, err = repository.Create(ctx, user2)
	assert.NoError(t, err)

	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"email": "user1@example.com",
	}, "created_at", "desc", 0, 10)

	// Act
	users, err := repository.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 1, len(*users))
	assert.Equal(t, user1.ID, (*users)[0].ID)
}

func TestCreateUserSucceeds(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}
	user := userFactory.NewUser(
		"test5@example.com",
		"password",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)

	// Act
	createdUser, err := repository.Create(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.FirstName, createdUser.FirstName)
	assert.Equal(t, user.LastName, createdUser.LastName)
	assert.NotEmpty(t, createdUser.ID)
}

func TestUpdateUserSucceeds(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}
	user := userFactory.NewUser(
		"test6@example.com",
		"password",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)
	createdUser, err := repository.Create(ctx, user)
	assert.NoError(t, err)

	// Update user fields
	createdUser.FirstName = "Updated"
	createdUser.LastName = "Name"

	// Act
	updatedUser, err := repository.Update(ctx, createdUser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "Updated", updatedUser.FirstName)
	assert.Equal(t, "Name", updatedUser.LastName)
	assert.Equal(t, createdUser.ID, updatedUser.ID)
}

func TestDeleteUserSucceeds(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userFactory := &entities.UserFactory{}
	user := userFactory.NewUser(
		"test7@example.com",
		"password",
		"John",
		"Doe",
		"admin",
		time.Now(),
		true,
		true,
	)
	createdUser, err := repository.Create(ctx, user)
	assert.NoError(t, err)

	// Act
	err = repository.Delete(ctx, createdUser.ID)

	// Assert
	assert.NoError(t, err)

	// Verify user is deleted
	_, err = repository.GetByID(ctx, createdUser.ID)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteNonExistentUserReturnsError(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	nonExistentID := uuid.New()

	// Act
	err := repository.Delete(ctx, nonExistentID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// Authentication tests
