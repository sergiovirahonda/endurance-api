package users

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
)

type UserService interface {
	// User services
	GetUserByID(ctx echo.Context, id uuid.UUID) (*entities.User, error)
	GetUserByEmail(ctx echo.Context, email string) (*entities.User, error)
	GetAllUsers(ctx echo.Context, filters filtering.ComplexFilters) (*entities.Users, error)
	CreateUser(ctx echo.Context, user *entities.User) (*entities.User, error)
	UpdateUser(ctx echo.Context, user *entities.User) (*entities.User, error)
	DeleteUser(ctx echo.Context, id uuid.UUID) error
	RegisterUser(ctx echo.Context, email, password, firstName, lastName, role string) (*entities.User, error)
	PasswordReset(ctx echo.Context, password string) error
	CheckUserPassword(ctx echo.Context, user *entities.User, password string) error
	// Authentication services
	Authenticate(ctx echo.Context, credentials valueobjects.LoginCredentials) (*entities.User, error)
	Logout(ctx echo.Context) error
	GetClaimsFromToken(ctx echo.Context, tokenString string) (*entities.JWTClaim, error)
	GenerateAccessToken(ctx echo.Context, user *entities.User) (string, error)
	GenerateRefreshToken(ctx echo.Context, user *entities.User) (string, error)
	GetTokens(ctx echo.Context, user *entities.User) (*valueobjects.Tokens, error)
	Login(ctx echo.Context, credentials valueobjects.LoginCredentials) (*valueobjects.Tokens, error)
	RefreshToken(ctx echo.Context, refreshToken string) (*valueobjects.Tokens, error)
	ValidateAccessToken(ctx echo.Context, tokenString string) (*entities.JWTClaim, error)
}
