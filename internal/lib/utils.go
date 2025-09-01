package lib

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sergiovirahonda/endurance-api/internal/config"
)

// Defer helpers to avoid linting issues

func Check(f func() error) {
	if err := f(); err != nil {
		logger := config.GetLogger()
		logger.Error(fmt.Sprintf("%v", err))
	}
}

// Slice helpers

func UUIDSliceContains(slice []uuid.UUID, item uuid.UUID) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func SliceContains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
