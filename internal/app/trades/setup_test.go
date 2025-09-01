package trades

import (
	"os"
	"testing"

	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/trade"
	"gorm.io/gorm"
)

var (
	database                  *gorm.DB
	tradePreferenceRepository trade.TradingPreferenceRepository
	holdingRepository         trade.HoldingRepository
	orderRepository           trade.OrderRepository
	tradingPreferenceService  TradingPreferenceService
	holdingService            HoldingService
	orderService              OrderService
	uacService                uacs.UacService
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running trades service tests...")
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
	tradePreferenceRepository = trade.NewDefaultTradingPreferenceRepository(database)
	holdingRepository = trade.NewDefaultHoldingRepository(database)
	orderRepository = trade.NewDefaultOrderRepository(database)
	uacService = uacs.NewDefaultUacService()
	tradingPreferenceService = NewDefaultTradingPreferenceService(tradePreferenceRepository, uacService)
	holdingService = NewDefaultHoldingService(holdingRepository, uacService)
	orderService = NewDefaultOrderService(orderRepository, uacService)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
