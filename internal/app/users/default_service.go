package users

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/domain/valueobjects"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/user"
	"github.com/sergiovirahonda/endurance-api/internal/lib"
)

type DefaultUserService struct {
	UserRepository user.UserRepository
	UacService     uacs.UacService
}

func NewDefaultUserService(
	userRepository user.UserRepository,
	uacService uacs.UacService,
) *DefaultUserService {
	return &DefaultUserService{
		UserRepository: userRepository,
		UacService:     uacService,
	}
}

func (s *DefaultUserService) GetUserByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.User, error) {
	user, err := s.UserRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsRegularUser(ctx); err == nil {
		if err := s.UacService.IsResourceOwner(ctx, user.ID); err != nil {
			return nil, err
		}
		return user, nil
	}
	return user, nil
}

func (s *DefaultUserService) GetUserByEmail(
	ctx echo.Context,
	email string,
) (*entities.User, error) {
	user, err := s.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsRegularUser(ctx); err == nil {
		if err := s.UacService.IsResourceOwner(ctx, user.ID); err != nil {
			return nil, err
		}
		return user, nil
	}
	return user, nil
}

func (s *DefaultUserService) GetAllUsers(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.Users, error) {
	filters.SetMetaParameters()
	if err := s.UacService.IsRegularUser(ctx); err == nil {
		filters.NarrowUserFilters("id")
	}
	return s.UserRepository.GetAll(ctx, filters)
}

func (s *DefaultUserService) CreateUser(
	ctx echo.Context,
	user *entities.User,
) (*entities.User, error) {
	_, err := s.UserRepository.GetByEmail(ctx, user.Email)
	if err == nil {
		return nil, errors.ErrEmailAddressInUse
	}
	hasher := lib.NewHasher()
	hash, err := hasher.HashString(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hash
	_, err = s.UserRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *DefaultUserService) UpdateUser(
	ctx echo.Context,
	user *entities.User,
) (*entities.User, error) {
	u, err := s.UserRepository.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.UacService.IsRegularUser(ctx); err == nil {
		if err := s.UacService.IsResourceOwner(ctx, u.ID); err != nil {
			return nil, err
		}
	}
	if user.Password != u.Password {
		hasher := lib.NewHasher()
		hash, err := hasher.HashString(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hash
	}
	return s.UserRepository.Update(ctx, user)
}

func (s *DefaultUserService) DeleteUser(
	ctx echo.Context,
	id uuid.UUID,
) error {
	u, err := s.UserRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.UacService.IsRegularUser(ctx); err == nil {
		if err := s.UacService.IsResourceOwner(ctx, u.ID); err != nil {
			return err
		}
	}
	return s.UserRepository.Delete(ctx, u.ID)
}

func (s *DefaultUserService) RegisterUser(
	ctx echo.Context,
	email string,
	password string,
	firstName string,
	lastName string,
	role string,
) (*entities.User, error) {
	factory := entities.UserFactory{}
	user := factory.NewUser(
		email,
		password,
		firstName,
		lastName,
		role,
		time.Time{},
		false,
		true,
	)
	return s.CreateUser(ctx, user)
}

func (s *DefaultUserService) PasswordReset(
	ctx echo.Context,
	password string,
) error {
	user := s.UacService.GetUser(ctx)
	u, err := s.UserRepository.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	hasher := lib.NewHasher()
	hash, err := hasher.HashString(password)
	if err != nil {
		return err
	}
	u.Password = hash
	_, err = s.UpdateUser(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) CheckUserPassword(
	ctx echo.Context,
	user *entities.User,
	password string,
) error {
	hasher := lib.NewHasher()
	valid := hasher.CheckStringHash(password, user.Password)
	if !valid {
		return errors.ErrInvalidPassword
	}
	return nil
}

// User authentication services

func (s DefaultUserService) Authenticate(
	ctx echo.Context,
	credentials valueobjects.LoginCredentials,
) (*entities.User, error) {
	hasher := lib.NewHasher()
	user, err := s.UserRepository.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}
	valid := hasher.CheckStringHash(credentials.Password, user.Password)
	if !valid {
		return nil, errors.ErrInvalidCredentials
	}
	if !user.Enabled {
		return nil, errors.ErrUserNotEnabled
	}
	user.LastLogin = time.Now().UTC()
	user.LoggedIn = true
	_, err = s.UserRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s DefaultUserService) Logout(ctx echo.Context) error {
	user := s.UacService.GetUser(ctx)
	user, err := s.UserRepository.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	user.LoggedIn = false
	_, err = s.UserRepository.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultUserService) GetClaimsFromToken(
	ctx echo.Context,
	tokenString string,
) (*entities.JWTClaim, error) {
	logger := config.GetLogger()
	signingKey := []byte(config.GetConfig().JWT.SigningKey)
	token, err := jwt.ParseWithClaims(tokenString, &entities.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("unexpected signing method", "method", token.Header["alg"])
			return nil, errors.ErrUnexpectedSigningMethod
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.ErrInvalidToken
	}
	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}
	err = token.Claims.Valid()
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*entities.JWTClaim)
	if !ok {
		return nil, errors.ErrInvalidToken
	}
	customClaims := entities.JWTClaim{
		ID:        claims.ID,
		Email:     claims.Email,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Role:      claims.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: claims.ExpiresAt,
			Audience:  claims.Audience,
		},
	}
	return &customClaims, nil
}

func (s DefaultUserService) GenerateAccessToken(
	ctx echo.Context,
	user *entities.User,
) (string, error) {
	accessTokenClaim := &entities.JWTClaim{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Audience:  "access",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	jwtToken, err := token.SignedString([]byte(config.GetConfig().JWT.SigningKey))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (s DefaultUserService) GenerateRefreshToken(
	ctx echo.Context,
	user *entities.User,
) (string, error) {

	accessTokenClaim := &entities.JWTClaim{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Audience:  "refresh",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	jwtToken, err := token.SignedString([]byte(config.GetConfig().JWT.SigningKey))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (s DefaultUserService) GetTokens(
	ctx echo.Context,
	user *entities.User,
) (*valueobjects.Tokens, error) {
	tokens := valueobjects.Tokens{}
	accessToken, err := s.GenerateAccessToken(ctx, user)
	if err != nil {
		return &valueobjects.Tokens{}, err
	}
	refreshToken, err := s.GenerateRefreshToken(ctx, user)
	if err != nil {
		return &valueobjects.Tokens{}, err
	}
	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return &tokens, nil
}

func (s DefaultUserService) Login(
	ctx echo.Context,
	credentials valueobjects.LoginCredentials,
) (*valueobjects.Tokens, error) {
	user, err := s.Authenticate(ctx, credentials)
	if err != nil {
		return &valueobjects.Tokens{}, err
	}
	tokens, err := s.GetTokens(ctx, user)
	if err != nil {
		return tokens, err
	}
	return tokens, nil
}

func (s DefaultUserService) RefreshToken(
	ctx echo.Context,
	refreshToken string,
) (*valueobjects.Tokens, error) {
	claims, err := s.GetClaimsFromToken(ctx, refreshToken)
	if err != nil {
		return &valueobjects.Tokens{}, err
	}
	if claims.StandardClaims.Audience != "refresh" {
		return &valueobjects.Tokens{}, errors.ErrInvalidToken
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return &valueobjects.Tokens{}, errors.ErrTokenExpired
	}
	userID := claims.ID
	id, err := uuid.Parse(userID)
	if err != nil {
		return &valueobjects.Tokens{}, errors.ErrInvalidToken
	}
	user := &entities.User{
		ID:        id,
		Email:     claims.Email,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Role:      claims.Role,
	}
	ctx.Set("user", user)
	instance, err := s.GetUserByID(ctx, id)
	if err != nil {
		return &valueobjects.Tokens{}, err
	}
	if !instance.Enabled {
		return &valueobjects.Tokens{}, errors.ErrUserNotEnabled
	}
	if !instance.LoggedIn {
		return &valueobjects.Tokens{}, errors.ErrUserLoggedOut
	}
	accessToken, err := s.GenerateAccessToken(ctx, instance)
	if err != nil {
		return &valueobjects.Tokens{}, err
	}
	tokens := valueobjects.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &tokens, nil
}

func (s DefaultUserService) ValidateAccessToken(
	ctx echo.Context,
	tokenString string,
) (*entities.JWTClaim, error) {
	claims, err := s.GetClaimsFromToken(ctx, tokenString)
	if err != nil {
		return &entities.JWTClaim{}, err
	}
	if claims.Audience != "access" {
		return &entities.JWTClaim{}, errors.ErrInvalidToken
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return &entities.JWTClaim{}, errors.ErrTokenExpired
	}
	userID := claims.ID
	id, err := uuid.Parse(userID)
	if err != nil {
		return &entities.JWTClaim{}, errors.ErrInvalidToken
	}
	user := &entities.User{
		ID:        id,
		Email:     claims.Email,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Role:      claims.Role,
	}
	ctx.Set("user", user)
	_, err = s.GetUserByID(ctx, id)
	if err != nil {
		return &entities.JWTClaim{}, err
	}
	return claims, nil
}
