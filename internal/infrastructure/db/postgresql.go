package db

import (
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) *gorm.DB {
	logger := config.GetLogger()
	logger.Info("Creating new database connection...")
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	logger.Info("Database connection created successfully.")
	return db
}
