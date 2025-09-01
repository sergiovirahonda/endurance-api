package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestConnection() *gorm.DB {
	zaplogger := config.GetLogger()
	zaplogger.Info("Instantiating test database...")
	configurations := config.GetConfig()
	if !configurations.Logger.LogTestQueries {
		zaplogger.Info("Test queries logging is disabled.")
		db, err := gorm.Open(
			sqlite.Open("file::memory:?cache=shared"),
			&gorm.Config{},
		)
		if err != nil {
			e := fmt.Sprintf("Failed to connect to test database: %v", err)
			panic(e)
		}
		return db
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Use standard Go logger or custom writer
		logger.Config{
			SlowThreshold: time.Microsecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		e := fmt.Sprintf("Failed to connect to test database: %v", err)
		panic(e)
	}

	return db
}

func Migrate(db *gorm.DB, models []interface{}) {
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic(err)
		}
	}
}

func DropDatabase(db *gorm.DB, models []interface{}) {
	for _, model := range models {
		err := db.Migrator().DropTable(model)
		if err != nil {
			panic(err)
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
}
