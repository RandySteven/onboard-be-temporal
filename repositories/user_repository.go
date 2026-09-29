package repositories

import (
	"context"

	db_client "github.com/RandySteven/go-cook/db"
	"github.com/RandySteven/onboard-be/entities/models"
	"github.com/RandySteven/onboard-be/queries"
)

type (
	UserRepository interface {
		Repository[models.User]
	}

	userRepository struct {
		dbx db_client.DBX
	}
)

func (u *userRepository) Save(ctx context.Context, entity *models.User) (*models.User, error) {
	additionalInfo := entity.AdditionalInfo
	if len(additionalInfo) == 0 {
		additionalInfo = []byte(`{}`)
	}
	id, err := db_client.Save[models.User](
		ctx,
		u.dbx(ctx),
		queries.InsertUser.ToString(),
		entity.OnboardingID,
		entity.UserName,
		entity.Status,
		string(additionalInfo),
	)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (u *userRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	result := &models.User{}
	if err := db_client.FindByID(ctx, u.dbx(ctx), queries.SelectUserByID.ToString(), id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userRepository) FindAll(ctx context.Context, skip uint64, take uint64) ([]*models.User, error) {
	return nil, nil
}

func (u *userRepository) Update(ctx context.Context, entity *models.User) (*models.User, error) {
	err := db_client.Update[models.User](
		ctx,
		u.dbx(ctx),
		queries.UpdateUser.ToString(),
		entity.OnboardingID,
		entity.UserName,
		entity.Status,
		entity.ID,
	)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *userRepository) DeleteByID(ctx context.Context, id uint64) error {
	return nil
}

func newUserRepository(dbx db_client.DBX) *userRepository {
	return &userRepository{dbx: dbx}
}

var _ UserRepository = &userRepository{}
