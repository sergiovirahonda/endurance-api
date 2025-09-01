package trade

import (
	"os"
	"testing"

	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"gorm.io/gorm"
)

var (
	database                  *gorm.DB
	tradePreferenceRepository TradingPreferenceRepository
	holdingRepository         HoldingRepository
	orderRepository           OrderRepository
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running tradePreferences, holdings and orders repository tests...")
	logger.Info("Instantiating test database...")
	database = db.NewTestConnection()
	logger.Info("Test DB connection established.")
	models := []interface{}{
		&dtos.TradingPreference{},
		&dtos.Holding{},
		&dtos.Order{},
	}
	logger.Info("Attempting to run migrations on test database...")
	db.Migrate(database, models)
	logger.Info("Migrations completed. Running tests...")
	tradePreferenceRepository = NewDefaultTradingPreferenceRepository(database)
	holdingRepository = NewDefaultHoldingRepository(database)
	orderRepository = NewDefaultOrderRepository(database)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
