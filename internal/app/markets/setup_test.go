package markets

import (
	"os"
	"testing"

	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/market"
	"gorm.io/gorm"
)

var (
	database             *gorm.DB
	MarketDataRepository market.MarketDataRepository
	marketDataService    MarketDataService
	uacService           uacs.UacService
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running markets service tests...")
	logger.Info("Instantiating test database...")
	database = db.NewTestConnection()
	logger.Info("Test DB connection established.")
	models := []interface{}{
		&dtos.MarketData{},
	}
	logger.Info("Attempting to run migrations on test database...")
	db.Migrate(database, models)
	logger.Info("Migrations completed. Running tests...")
	MarketDataRepository = market.NewDefaultMarketDataRepository(database)
	uacService = uacs.NewDefaultUacService()
	marketDataService = NewDefaultMarketDataService(MarketDataRepository, uacService)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
