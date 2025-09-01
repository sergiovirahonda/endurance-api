package market

import (
	"os"
	"testing"

	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"gorm.io/gorm"
)

var (
	database             *gorm.DB
	marketRepository     MarketRepository
	marketDataRepository MarketDataRepository
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running market data repository tests...")
	logger.Info("Instantiating test database...")
	database = db.NewTestConnection()
	logger.Info("Test DB connection established.")
	models := []interface{}{
		&dtos.Market{},
		&dtos.MarketData{},
	}
	logger.Info("Attempting to run migrations on test database...")
	db.Migrate(database, models)
	logger.Info("Migrations completed. Running tests...")
	marketRepository = NewDefaultMarketRepository(database)
	marketDataRepository = NewDefaultMarketDataRepository(database)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
