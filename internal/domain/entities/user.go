package entities

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	LastLogin time.Time `json:"last_login"`
	LoggedIn  bool      `json:"logged_in"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Users []User

// Authentication

type JWTClaim struct {
	// User Information
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
	jwt.StandardClaims
}

// Factories

type UserFactory struct{}

func (f *UserFactory) NewUser(
	email string,
	password string,
	firstName string,
	lastName string,
	role string,
	lastLogin time.Time,
	loggedIn bool,
	enabled bool,
) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		LastLogin: lastLogin,
		LoggedIn:  loggedIn,
		Enabled:   enabled,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (f *UserFactory) Clone(
	user *User,
	email string,
	password string,
	firstName string,
	lastName string,
	role string,
	lastLogin time.Time,
	loggedIn bool,
	enabled bool,
) *User {
	return &User{
		ID:        user.ID,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		LastLogin: lastLogin,
		LoggedIn:  loggedIn,
		Enabled:   enabled,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
