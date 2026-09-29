package repositories

import (
	"context"

	db_client "github.com/RandySteven/go-cook/db"
	"github.com/RandySteven/onboard-be/entities/models"
	"github.com/RandySteven/onboard-be/queries"
)

type (
	OnboardingRepository interface {
		Repository[models.Onboarding]
	}

	onboardingRepository struct {
		dbx db_client.DBX
	}
)

func (o *onboardingRepository) DeleteByID(ctx context.Context, id uint64) error {
	return nil
}

func (o *onboardingRepository) FindAll(ctx context.Context, skip uint64, take uint64) ([]*models.Onboarding, error) {
	return nil, nil
}

func (o *onboardingRepository) FindByID(ctx context.Context, id uint64) (*models.Onboarding, error) {
	return nil, nil
}

func (o *onboardingRepository) Save(ctx context.Context, entity *models.Onboarding) (*models.Onboarding, error) {
	additionalInfo := entity.AdditionalInfo
	if len(additionalInfo) == 0 {
		additionalInfo = []byte(`{}`)
	}
	id, err := db_client.Save[models.Onboarding](
		ctx,
		o.dbx(ctx),
		queries.InsertOnboarding.ToString(),
		entity.Name,
		entity.Email,
		entity.Password,
		entity.PhoneNumber,
		entity.Address,
		entity.RegisterAs,
		entity.Status,
		string(additionalInfo),
	)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (o *onboardingRepository) Update(ctx context.Context, entity *models.Onboarding) (*models.Onboarding, error) {
	err := db_client.Update[models.Onboarding](
		ctx,
		o.dbx(ctx),
		queries.UpdateOnboarding.ToString(),
		entity.Email,
		entity.Password,
		entity.PhoneNumber,
		entity.Name,
		entity.Address,
		entity.RegisterAs,
		entity.Status,
		entity.ID,
	)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

var _ OnboardingRepository = &onboardingRepository{}

func newOnboardingRepository(dbx db_client.DBX) *onboardingRepository {
	return &onboardingRepository{dbx: dbx}
}
