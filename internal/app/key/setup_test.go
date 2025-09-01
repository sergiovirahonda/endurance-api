package keys

import (
	"os"
	"testing"

	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/key"
	"gorm.io/gorm"
)

var (
	database      *gorm.DB
	keyRepository key.KeyRepository
	uacService    uacs.UacService
	keyService    KeyService
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running key service tests...")
	logger.Info("Instantiating test database...")
	database = db.NewTestConnection()
	logger.Info("Test DB connection established.")
	models := []interface{}{
		&dtos.ApiKey{},
	}
	logger.Info("Attempting to run migrations on test database...")
	db.Migrate(database, models)
	logger.Info("Migrations completed. Running tests...")
	keyRepository = key.NewDefaultKeyRepository(database)
	uacService = uacs.NewDefaultUacService()
	keyService = NewDefaultKeyService(keyRepository, uacService)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
