package keys

import (
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// --- KeyService Tests ---

func TestGetApiKeyByIDReturnsErrorIfNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	// Act
	_, err := keyService.GetByID(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetApiKeyByIDReturnsApiKey(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	// Act
	foundApiKey, err := keyService.GetByID(ctx, apiKey.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundApiKey)
	assert.Equal(t, apiKey.ID, foundApiKey.ID)
	assert.Equal(t, apiKey.UserID, foundApiKey.UserID)
	assert.Equal(t, apiKey.Service, foundApiKey.Service)
	assert.Equal(t, apiKey.Key, foundApiKey.Key)
	assert.Equal(t, apiKey.Secret, foundApiKey.Secret)
}

func TestGetApiKeyByIDReturnsErrorIfNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	otherUserID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		otherUserID, // Different user
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	// Act
	_, err := keyService.GetByID(ctx, apiKey.ID)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
}

func TestGetAllApiKeysReturnsApiKeys(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey1 := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key-1",
		"test-api-secret-1",
	)
	apiKey2 := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeTelegram,
		"test-api-key-2",
		"test-api-secret-2",
	)

	dto1 := dtos.ApiKey{}
	dto1.FromEntity(apiKey1)
	database.Create(&dto1)

	dto2 := dtos.ApiKey{}
	dto2.FromEntity(apiKey2)
	database.Create(&dto2)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"user_id": userID,
	}, "created_at", "desc", 0, 10)
	foundApiKeys, err := keyService.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundApiKeys)
	assert.Equal(t, 2, len(*foundApiKeys))
	assert.Equal(t, apiKey1.ID, (*foundApiKeys)[0].ID)
	assert.Equal(t, apiKey2.ID, (*foundApiKeys)[1].ID)
}

func TestGetAllApiKeysWithUnmatchingUserIDReturnsEmpty(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"user_id": uuid.New(),
	}, "created_at", "desc", 0, 10)
	foundApiKeys, err := keyService.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundApiKeys)
	assert.Equal(t, 0, len(*foundApiKeys))
}

func TestGetAllApiKeysWithServiceFilter(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey1 := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key-1",
		"test-api-secret-1",
	)
	apiKey2 := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeTelegram,
		"test-api-key-2",
		"test-api-secret-2",
	)

	dto1 := dtos.ApiKey{}
	dto1.FromEntity(apiKey1)
	database.Create(&dto1)

	dto2 := dtos.ApiKey{}
	dto2.FromEntity(apiKey2)
	database.Create(&dto2)

	// Act
	filters := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"service": constants.ApiKeyServiceTypeBinance,
	}, "created_at", "desc", 0, 10)
	foundApiKeys, err := keyService.GetAll(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundApiKeys)
	assert.Equal(t, 1, len(*foundApiKeys))
	assert.Equal(t, apiKey1.ID, (*foundApiKeys)[0].ID)
	assert.Equal(t, constants.ApiKeyServiceTypeBinance, (*foundApiKeys)[0].Service)
}

func TestCreateApiKeyReturnsApiKey(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	// Act
	createdApiKey, err := keyService.Create(ctx, apiKey)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdApiKey)
	assert.Equal(t, apiKey.ID, createdApiKey.ID)
	assert.Equal(t, apiKey.UserID, createdApiKey.UserID)
	assert.Equal(t, apiKey.Service, createdApiKey.Service)
	assert.Equal(t, apiKey.Key, createdApiKey.Key)
	assert.Equal(t, apiKey.Secret, createdApiKey.Secret)

	// Assert that the api key was created in the database
	ak, err := keyService.GetByID(ctx, apiKey.ID)
	assert.NoError(t, err)
	assert.NotNil(t, ak)
	assert.Equal(t, apiKey.ID, ak.ID)
}

func TestCreateApiKeyReturnsErrorIfNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	otherUserID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		otherUserID, // Different user
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	// Act
	_, err := keyService.Create(ctx, apiKey)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
}

func TestCreateApiKeyReturnsErrorIfInvalidService(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		"invalid-service", // Invalid service
		"test-api-key",
		"test-api-secret",
	)

	// Act
	_, err := keyService.Create(ctx, apiKey)

	// Assert
	assert.Equal(t, errors.ErrApiKeyServiceInvalid, err)
}

func TestCreateApiKeyReturnsErrorIfEmptyKey(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"", // Empty key
		"test-api-secret",
	)

	// Act
	_, err := keyService.Create(ctx, apiKey)

	// Assert
	assert.Equal(t, errors.ErrApiKeyKeyRequired, err)
}

func TestCreateApiKeyReturnsErrorIfEmptySecret(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"", // Empty secret
	)

	// Act
	_, err := keyService.Create(ctx, apiKey)

	// Assert
	assert.Equal(t, errors.ErrApiKeySecretRequired, err)
}

func TestUpdateApiKeyReturnsApiKey(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	apiKey.Service = constants.ApiKeyServiceTypeTelegram
	apiKey.Key = "updated-api-key"
	apiKey.Secret = "updated-api-secret"

	// Act
	updatedApiKey, err := keyService.Update(ctx, apiKey)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedApiKey)
	assert.Equal(t, apiKey.ID, updatedApiKey.ID)
	assert.Equal(t, apiKey.Service, updatedApiKey.Service)
	assert.Equal(t, apiKey.Key, updatedApiKey.Key)
	assert.Equal(t, apiKey.Secret, updatedApiKey.Secret)
}

func TestUpdateApiKeyReturnsErrorIfNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	// Act
	_, err := keyService.Update(ctx, apiKey)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUpdateApiKeyReturnsErrorIfNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	otherUserID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		otherUserID, // Different user
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	// Act
	_, err := keyService.Update(ctx, apiKey)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
}

func TestUpdateApiKeyReturnsErrorIfInvalidService(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	apiKey.Service = "invalid-service" // Invalid service

	// Act
	_, err := keyService.Update(ctx, apiKey)

	// Assert
	assert.Equal(t, errors.ErrApiKeyServiceInvalid, err)
}

func TestDeleteApiKeyReturnsErrorIfIDDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	// Act
	err := keyService.Delete(ctx, uuid.New())

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteApiKeyReturnsErrorIfNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	otherUserID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		otherUserID, // Different user
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	// Act
	err := keyService.Delete(ctx, apiKey.ID)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
}

func TestDeleteApiKeyDeletesApiKey(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID := uuid.New()
	ctx.Set("user", &entities.User{ID: userID})

	apiKeyFactory := &entities.ApiKeyFactory{}
	apiKey := apiKeyFactory.NewApiKey(
		userID,
		constants.ApiKeyServiceTypeBinance,
		"test-api-key",
		"test-api-secret",
	)

	dto := dtos.ApiKey{}
	dto.FromEntity(apiKey)
	database.Create(&dto)

	// Act
	err := keyService.Delete(ctx, apiKey.ID)

	// Assert
	assert.NoError(t, err)

	foundApiKey, err := keyService.GetByID(ctx, apiKey.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, foundApiKey)
}

func TestCreateApiKeyWithDifferentUsers(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	userID1 := uuid.New()
	userID2 := uuid.New()

	apiKeyFactory := &entities.ApiKeyFactory{}

	apiKey1 := apiKeyFactory.NewApiKey(
		userID1,
		constants.ApiKeyServiceTypeBinance,
		"user1-api-key",
		"user1-api-secret",
	)
	apiKey2 := apiKeyFactory.NewApiKey(
		userID2,
		constants.ApiKeyServiceTypeBinance,
		"user2-api-key",
		"user2-api-secret",
	)

	// Act - Create for user 1
	ctx.Set("user", &entities.User{ID: userID1})
	createdApiKey1, err1 := keyService.Create(ctx, apiKey1)

	// Act - Create for user 2
	ctx.Set("user", &entities.User{ID: userID2})
	createdApiKey2, err2 := keyService.Create(ctx, apiKey2)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotNil(t, createdApiKey1)
	assert.NotNil(t, createdApiKey2)
	assert.Equal(t, userID1, createdApiKey1.UserID)
	assert.Equal(t, userID2, createdApiKey2.UserID)

	// Verify they can be retrieved separately
	ctx.Set("user", &entities.User{ID: userID1})
	filters1 := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"user_id": userID1,
	}, "created_at", "desc", 0, 10)
	user1Keys, err := keyService.GetAll(ctx, filters1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*user1Keys))

	ctx.Set("user", &entities.User{ID: userID2})
	filters2 := filtering.NewComplexFilter(ctx, map[string]interface{}{
		"user_id": userID2,
	}, "created_at", "desc", 0, 10)
	user2Keys, err := keyService.GetAll(ctx, filters2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*user2Keys))
}
