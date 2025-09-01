package keys

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/domain/errors"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/filtering"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/persistence/repositories/key"
)

// Structs

type DefaultKeyService struct {
	KeyRepository key.KeyRepository
	UacService    uacs.UacService
}

// Factories

func NewDefaultKeyService(
	keyRepository key.KeyRepository,
	uacService uacs.UacService,
) *DefaultKeyService {
	return &DefaultKeyService{KeyRepository: keyRepository, UacService: uacService}
}

// KeyService implementation

func (d *DefaultKeyService) GetByID(
	ctx echo.Context,
	id uuid.UUID,
) (*entities.ApiKey, error) {
	key, err := d.KeyRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = d.UacService.IsResourceOwner(ctx, key.UserID)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (d *DefaultKeyService) GetAll(
	ctx echo.Context,
	filters filtering.ComplexFilters,
) (*entities.ApiKeys, error) {
	filters.SetMetaParameters()
	filters.NarrowUserFilters("user_id")
	return d.KeyRepository.GetAll(ctx, filters)
}

func (d *DefaultKeyService) Create(
	ctx echo.Context,
	apiKey *entities.ApiKey,
) (*entities.ApiKey, error) {
	err := d.UacService.IsResourceOwner(ctx, apiKey.UserID)
	if err != nil {
		return nil, err
	}
	err = apiKey.Validate()
	if err != nil {
		return nil, err
	}
	return d.KeyRepository.Create(ctx, apiKey)
}

func (d *DefaultKeyService) Update(
	ctx echo.Context,
	apiKey *entities.ApiKey,
) (*entities.ApiKey, error) {
	_, err := d.GetByID(ctx, apiKey.ID)
	if err != nil {
		return nil, err
	}
	err = apiKey.Validate()
	if err != nil {
		return nil, err
	}
	return d.KeyRepository.Update(ctx, apiKey)
}

func (d *DefaultKeyService) Delete(
	ctx echo.Context,
	id uuid.UUID,
) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return d.KeyRepository.Delete(ctx, id)
}

func (d *DefaultKeyService) GetTelegramKeys(
	ctx echo.Context,
) (*entities.ApiKey, error) {
	filters := filtering.NewComplexFilter(
		ctx,
		map[string]interface{}{
			"type": "telegram",
		},
		"created_at",
		"desc",
		1,
		1,
	)
	keys, err := d.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(*keys) == 0 {
		return nil, errors.ErrApiKeyTelegramNotFound
	}
	return &(*keys)[0], nil
}
