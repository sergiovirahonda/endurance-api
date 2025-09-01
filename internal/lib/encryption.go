package lib

import (
	"golang.org/x/crypto/bcrypt"
)

// Structs

type Hasher struct{}

// Factories

func NewHasher() Hasher {
	return Hasher{}
}

// Receivers

func (h Hasher) HashString(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	return string(bytes), err
}

func (h Hasher) CheckStringHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
