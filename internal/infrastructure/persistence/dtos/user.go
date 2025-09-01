package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Password  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"type:varchar(100);unique_index"`
	FirstName string    `gorm:"type:varchar(100)"`
	LastName  string    `gorm:"type:varchar(100)"`
	Role      string    `gorm:"type:varchar(100)"`
	LastLogin time.Time `gorm:"type:timestamp"`
	LoggedIn  bool      `gorm:"type:boolean;default:false"`
	Enabled   bool      `gorm:"type:boolean;default:false"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;"`
}

type Users []User

// Receivers

func (u *User) ToEntity() *entities.User {
	return &entities.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		LastLogin: u.LastLogin,
		LoggedIn:  u.LoggedIn,
		Enabled:   u.Enabled,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) FromEntity(user *entities.User) {
	u.ID = user.ID
	u.Email = user.Email
	u.Password = user.Password
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Role = user.Role
	u.LastLogin = user.LastLogin
	u.LoggedIn = user.LoggedIn
	u.Enabled = user.Enabled
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

func (u *Users) ToEntities() *entities.Users {
	entities := make(entities.Users, len(*u))
	for i, user := range *u {
		entities[i] = *user.ToEntity()
	}
	return &entities
}
