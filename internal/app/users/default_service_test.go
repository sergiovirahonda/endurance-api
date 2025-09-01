package users

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/constants"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/sergiovirahonda/endurance-api/internal/lib"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	regularUser = &entities.User{
		ID:        uuid.New(),
		Email:     "regular@user.com",
		Password:  "Test Password",
		FirstName: "Test",
		LastName:  "User",
		Role:      constants.RoleUser,
		LastLogin: time.Now().UTC(),
		LoggedIn:  true,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
)

func TestGetUserByIDReturnsErrorIfUserDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	ctx.Set("user", user)

	// Act
	user, err := userService.GetUserByID(ctx, user.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, user)
}

func TestGetUserReturnsUserIfUserExistsAndOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test1@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx.Set("user", user)

	// Act
	user, err := userService.GetUserByID(ctx, user.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, user.ID)
	assert.Equal(t, user.Email, user.Email)
	assert.Equal(t, user.FirstName, user.FirstName)
	assert.Equal(t, user.LastName, user.LastName)
	assert.Equal(t, user.LastLogin, user.LastLogin)
	assert.Equal(t, user.LoggedIn, user.LoggedIn)
	assert.Equal(t, user.Enabled, user.Enabled)
	assert.Equal(t, user.CreatedAt, user.CreatedAt)
	assert.Equal(t, user.UpdatedAt, user.UpdatedAt)
}

func TestGetUserReturnsErrorIfUserExistsButNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test2@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	nonExistingUser := factory.NewUser(
		"test@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)

	ctx.Set("user", nonExistingUser)

	// Act
	user, err := userService.GetUserByID(ctx, user.ID)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
	assert.Nil(t, user)
}

func TestGetUserByEmailReturnsErrorIfUserDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)

	ctx.Set("user", user)

	// Act
	user, err := userService.GetUserByEmail(ctx, user.Email)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, user)
}

func TestGetUserByEmailReturnsErrorIfUserExistsButNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test3@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	nonExistingUser := factory.NewUser(
		"test@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)

	ctx.Set("user", nonExistingUser)

	// Act
	user, err := userService.GetUserByEmail(ctx, user.Email)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
	assert.Nil(t, user)
}

func TestGetUserByEmailReturnsUserIfUserExistsAndOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test4@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx.Set("user", user)

	// Act
	user, err := userService.GetUserByEmail(ctx, user.Email)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, user.ID)
	assert.Equal(t, user.Email, user.Email)
	assert.Equal(t, user.FirstName, user.FirstName)
	assert.Equal(t, user.LastName, user.LastName)
	assert.Equal(t, user.LastLogin, user.LastLogin)
	assert.Equal(t, user.LoggedIn, user.LoggedIn)
	assert.Equal(t, user.Enabled, user.Enabled)
	assert.Equal(t, user.CreatedAt, user.CreatedAt)
	assert.Equal(t, user.UpdatedAt, user.UpdatedAt)
}

func TestGetAllUsersReturnsEmptySliceIfFiltersDontMatchCriteria(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test5@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx.Set("user", user)

	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"last_name": "stark",
		},
		"created_at",
		"desc",
		0,
		10,
	)

	// Act
	users, err := userService.GetAllUsers(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*users))
}

func TestGetAllUsersReturnsUsersIfFiltersMatchCriteria(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test6@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx.Set("user", user)

	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{},
		"created_at",
		"desc",
		0,
		10,
	)

	// Act
	users, err := userService.GetAllUsers(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*users))

	assert.Equal(t, user.ID, (*users)[0].ID)
	assert.Equal(t, user.Email, (*users)[0].Email)
	assert.Equal(t, user.FirstName, (*users)[0].FirstName)
	assert.Equal(t, user.LastName, (*users)[0].LastName)
	assert.Equal(t, user.LastLogin, (*users)[0].LastLogin)
	assert.Equal(t, user.LoggedIn, (*users)[0].LoggedIn)
	assert.Equal(t, user.Enabled, (*users)[0].Enabled)
	assert.Equal(t, user.CreatedAt, (*users)[0].CreatedAt)
	assert.Equal(t, user.UpdatedAt, (*users)[0].UpdatedAt)
}

func TestCreateUserCreatesUser(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test7@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		false,
		true,
	)

	// Act
	user, err := userService.CreateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, user.ID)
	assert.Equal(t, user.Email, user.Email)
	assert.Equal(t, user.FirstName, user.FirstName)
	assert.Equal(t, user.LastName, user.LastName)
	assert.NotEqual(t, user.Password, "test")
	assert.Equal(t, user.LastLogin, user.LastLogin)
	assert.Equal(t, user.LoggedIn, user.LoggedIn)
	assert.Equal(t, user.Enabled, user.Enabled)
	assert.Equal(t, user.CreatedAt, user.CreatedAt)
	assert.Equal(t, user.UpdatedAt, user.UpdatedAt)
}

func TestUpdateUserReturnsErrorIfUserDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)

	// Act
	user, err := userService.UpdateUser(ctx, user)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, user)
}

func TestUpdateUserReturnsErrorIfUserExistsButNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test8@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto1 := dtos.User{}
	dto1.FromEntity(user)
	database.Create(&dto1)

	otherUser := factory.NewUser(
		"test9@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto2 := dtos.User{}
	dto2.FromEntity(otherUser)
	database.Create(&dto2)

	ctx.Set("user", user)

	// Act
	user, err := userService.UpdateUser(ctx, otherUser)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
	assert.Nil(t, user)
}

func TestUpdateUserUpdatesUser(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test10@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx.Set("user", user)

	// Act
	user.FirstName = "John"
	user.LastName = "Doe"
	user.Email = "john.doe@example.com"
	user, err := userService.UpdateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, user.ID)
	assert.Equal(t, user.Email, "john.doe@example.com")
	assert.Equal(t, user.FirstName, "John")
	assert.Equal(t, user.LastName, "Doe")
}

func TestDeleteUserReturnsErrorIfUserDoesNotExist(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)

	ctx.Set("user", user)

	// Act
	err := userService.DeleteUser(ctx, user.ID)

	// Assert
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteUserReturnsErrorIfUserExistsButNotOwner(t *testing.T) {
	// Arrange
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test11@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto1 := dtos.User{}
	dto1.FromEntity(user)
	database.Create(&dto1)

	otherUser := factory.NewUser(
		"test12@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto2 := dtos.User{}
	dto2.FromEntity(otherUser)
	database.Create(&dto2)

	ctx.Set("user", user)

	// Act
	err := userService.DeleteUser(ctx, otherUser.ID)

	// Assert
	assert.Equal(t, errors.ErrForbidden, err)
}

func TestDeleteUserDeletesUser(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test13@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto1 := dtos.User{}
	dto1.FromEntity(user)
	database.Create(&dto1)

	ctx.Set("user", user)

	// Act
	err := userService.DeleteUser(ctx, user.ID)

	// Assert
	assert.NoError(t, err)
}

// Authentication tests

func TestGenerateAccessTokenGeneratesExpectedToken(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	token, err := userService.GenerateAccessToken(ctx, regularUser)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, token, "")
}

func TestGenerateRefreshTokenGeneratesExpectedToken(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	token, err := userService.GenerateRefreshToken(ctx, regularUser)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, token, "")
	splitted := strings.Split(token, ".")
	assert.Equal(t, len(splitted), 3)
}

func TestGetTokensReturnsBothTokensSuccessfully(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	tokens, err := userService.GetTokens(ctx, regularUser)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, tokens.AccessToken, "")
	assert.NotEqual(t, tokens.RefreshToken, "")
	splittedAccess := strings.Split(tokens.AccessToken, ".")
	assert.Equal(t, len(splittedAccess), 3)
	splittedRefresh := strings.Split(tokens.RefreshToken, ".")
	assert.Equal(t, len(splittedRefresh), 3)
}

func TestGetClaimsFromTokenReturnsClaimsSuccessfully(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	tokens, err := userService.GetTokens(ctx, regularUser)
	assert.Equal(t, err, nil)
	claims, err := userService.GetClaimsFromToken(ctx, tokens.AccessToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, claims.ID, regularUser.ID.String())
	assert.Equal(t, claims.Email, regularUser.Email)
	assert.Equal(t, claims.FirstName, regularUser.FirstName)
	assert.Equal(t, claims.LastName, regularUser.LastName)
	assert.NotEqual(t, claims.ExpiresAt, 0)
	assert.Equal(t, claims.Audience, "access")
}

func TestGetClaimsFromTokenFailsIfTokenIsInvalid(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	tokens, err := userService.GetTokens(ctx, regularUser)
	assert.Equal(t, err, nil)
	invalidToken := tokens.AccessToken + "invalid"
	_, err = userService.GetClaimsFromToken(ctx, invalidToken)
	assert.NotEqual(t, err, nil)
}

func TestRefreshTokenFailsIfNotRefreshTokenPassed(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	tokens, err := userService.GetTokens(ctx, regularUser)
	assert.Equal(t, err, nil)
	_, err = userService.RefreshToken(ctx, tokens.AccessToken)
	assert.Equal(t, err, errors.ErrInvalidToken)
}

func TestRefreshTokenFailsIfUserInTokenDoesNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	factory := entities.UserFactory{}
	fakeUser := factory.NewUser(
		"random@random.com",
		"Test Password",
		"Test FirstName",
		"Test LastName",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	tokens, err := userService.GetTokens(ctx, fakeUser)
	assert.Equal(t, err, nil)
	_, err = userService.RefreshToken(ctx, tokens.RefreshToken)
	assert.Equal(t, err, gorm.ErrRecordNotFound)
}

func TestRefreshTokenFailsIfUserNotEnabled(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test16@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		false,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx.Set("user", user)

	tokens, err := userService.GetTokens(ctx, user)
	assert.Equal(t, err, nil)
	_, err = userService.RefreshToken(ctx, tokens.RefreshToken)
	assert.Equal(t, err, errors.ErrUserNotEnabled)
	user.Enabled = true
	_, err = userService.UpdateUser(ctx, user)
	assert.Equal(t, err, nil)
}

func TestRefreshTokenFailsIfUserNotLoggedIn(t *testing.T) {
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"stark@tonyindustries.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		false,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", user)

	tokens, err := userService.GetTokens(ctx, user)
	assert.Equal(t, err, nil)
	_, err = userService.RefreshToken(ctx, tokens.RefreshToken)
	assert.Equal(t, err, errors.ErrUserLoggedOut)
	user.LoggedIn = true
	_, err = userService.UpdateUser(ctx, user)
	assert.Equal(t, err, nil)
}

func TestRefreshTokenSucceedsIfUserExistsAndIsEnabled(t *testing.T) {
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test15@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", user)
	tokens, err := userService.GetTokens(ctx, user)
	assert.Equal(t, err, nil)
	refreshedTokens, err := userService.RefreshToken(ctx, tokens.RefreshToken)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, refreshedTokens.AccessToken, "")
	assert.NotEqual(t, refreshedTokens.RefreshToken, "")
	splittedAccess := strings.Split(refreshedTokens.AccessToken, ".")
	assert.Equal(t, len(splittedAccess), 3)
	splittedRefresh := strings.Split(refreshedTokens.RefreshToken, ".")
	assert.Equal(t, len(splittedRefresh), 3)
}

func TestLogoutFailsIfUserDoesNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	err := userService.Logout(ctx)
	assert.Equal(t, err, gorm.ErrRecordNotFound)
}

func TestLogoutSucceedsIfUserExists(t *testing.T) {
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test14@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", user)
	err := userService.Logout(ctx)
	assert.Equal(t, err, nil)
	instance, err := userService.GetUserByID(ctx, user.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, instance.LoggedIn, false)
	user.LoggedIn = true
	_, err = userService.UpdateUser(ctx, user)
	assert.Equal(t, err, nil)
}

func TestAuthenticateFailsIfUserDoesNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	ctx.Set("user", regularUser)
	credentials := valueobjects.LoginCredentials{
		Email:    "something@anything.com",
		Password: "password",
	}
	_, err := userService.Authenticate(ctx, credentials)
	assert.Equal(t, err, errors.ErrInvalidCredentials)
}

func TestAuthenticateFailsIfPasswordIsInvalid(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	credentials := valueobjects.LoginCredentials{
		Email:    regularUser.Email,
		Password: "invalid",
	}
	_, err := userService.Authenticate(ctx, credentials)
	assert.Equal(t, err, errors.ErrInvalidCredentials)
}

func TestAuthenticateSucceedsIfUserValid(t *testing.T) {
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test-auth@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	ctx := echo.New().NewContext(nil, nil)
	_, err := userService.CreateUser(ctx, user)

	assert.Nil(t, err)

	credentials := valueobjects.LoginCredentials{
		Email:    user.Email,
		Password: "test",
	}
	authenticatedUser, err := userService.Authenticate(ctx, credentials)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, authenticatedUser, entities.User{})
	assert.Equal(t, authenticatedUser.Email, "test-auth@test.com")
}

func TestLoginSucceeds(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)

	factory := entities.UserFactory{}
	user := factory.NewUser(
		"tony@starkindustries.com.mx",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	hasher := lib.NewHasher()
	hash, err := hasher.HashString(user.Password)
	if err != nil {
		t.Fatal(err)
	}
	user.Password = hash
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	credentials := valueobjects.LoginCredentials{
		Email:    user.Email,
		Password: "test",
	}
	tokens, err := userService.Login(ctx, credentials)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, tokens.AccessToken, "")
	assert.NotEqual(t, tokens.RefreshToken, "")
	splittedAccess := strings.Split(tokens.AccessToken, ".")
	assert.Equal(t, len(splittedAccess), 3)
	splittedRefresh := strings.Split(tokens.RefreshToken, ".")
	assert.Equal(t, len(splittedRefresh), 3)
}

func TestRegisterServiceFailsIfEmailAlreadyExists(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	factory := entities.UserFactory{}
	user := factory.NewUser(
		"test-register@test.com",
		"test",
		"test",
		"test",
		constants.RoleUser,
		time.Now().UTC(),
		true,
		true,
	)
	dto := dtos.User{}
	dto.FromEntity(user)
	database.Create(&dto)

	registration := valueobjects.RegisterCredentials{
		Email:     user.Email,
		Password:  "Test Password",
		FirstName: "Test",
		LastName:  "Test",
	}
	_, err := userService.RegisterUser(
		ctx,
		registration.Email,
		registration.Password,
		registration.FirstName,
		registration.LastName,
		constants.RoleUser,
	)
	assert.Equal(t, err, errors.ErrEmailAddressInUse)
}

func TestRegistrationServiceSucceedsIfEmailDoesNotExist(t *testing.T) {
	ctx := echo.New().NewContext(nil, nil)
	registration := valueobjects.RegisterCredentials{
		Email:     "test-register-2@test.com",
		Password:  "test",
		FirstName: "test",
		LastName:  "test",
	}
	_, err := userService.RegisterUser(
		ctx,
		registration.Email,
		registration.Password,
		registration.FirstName,
		registration.LastName,
		constants.RoleUser,
	)
	assert.Equal(t, err, nil)
	instance, err := repository.GetByEmail(ctx, registration.Email)
	assert.Equal(t, err, nil)
	assert.Equal(t, instance.Email, instance.Email)
}
