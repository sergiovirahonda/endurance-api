package key

import (
	"os"
	"testing"

	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"gorm.io/gorm"
)

var (
	database      *gorm.DB
	keyRepository *DefaultKeyRepository
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running key repository tests...")
	logger.Info("Instantiating test database...")
	database = db.NewTestConnection()
	logger.Info("Test DB connection established.")
	models := []interface{}{
		&dtos.ApiKey{},
	}
	logger.Info("Attempting to run migrations on test database...")
	db.Migrate(database, models)
	logger.Info("Migrations completed. Running tests...")
	keyRepository = NewDefaultKeyRepository(database)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
