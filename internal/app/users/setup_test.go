package users

import (
	"os"
	"testing"

	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/db"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/dtos"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/user"
	"gorm.io/gorm"
)

var (
	database    *gorm.DB
	repository  user.UserRepository
	userService UserService
	uacService  uacs.UacService
)

func TestMain(m *testing.M) {
	logger := config.GetLogger()
	logger.Info("Running user service tests...")
	logger.Info("Instantiating test database...")
	database = db.NewTestConnection()
	logger.Info("Test DB connection established.")
	models := []interface{}{
		&dtos.User{},
	}
	logger.Info("Attempting to run migrations on test database...")
	db.Migrate(database, models)
	logger.Info("Migrations completed. Running tests...")
	repository = user.NewDefaultUserRepository(database)
	uacService = uacs.NewDefaultUacService()
	userService = NewDefaultUserService(repository, uacService)
	os.Exit(m.Run())
	logger.Info("Tests completed.")
}
