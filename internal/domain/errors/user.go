package errors

import "errors"

var (
	// User errors
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotEnabled     = errors.New("user not enabled")
	ErrEmailAddressInUse  = errors.New("email address already in use")

	// Authentication errors
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrTokenExpired            = errors.New("token expired")
	ErrUserLoggedOut           = errors.New("user logged out")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)
