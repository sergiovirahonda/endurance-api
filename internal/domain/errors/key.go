package errors

import "errors"

var (
	ErrApiKeyAlreadyExists    = errors.New("api key already exists")
	ErrApiKeyNotFound         = errors.New("api key not found")
	ErrApiKeyInvalid          = errors.New("api key is invalid")
	ErrApiKeyUserIDRequired   = errors.New("user_id is required")
	ErrApiKeyServiceRequired  = errors.New("service is required")
	ErrApiKeyKeyRequired      = errors.New("key is required")
	ErrApiKeySecretRequired   = errors.New("secret is required")
	ErrApiKeyServiceInvalid   = errors.New("service is invalid")
	ErrApiKeyTelegramNotFound = errors.New("telegram key not found")
)
