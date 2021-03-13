package infrastructures

import (
	"context"

	"github.com/edy4c7/darkpot-school-works/internal/entities"
	"gorm.io/gorm"
)

type ActivitiesRepositoryImpl struct {
	db *gorm.DB
}

func NewActivitiesRepositoryImpl(db *gorm.DB) *ActivitiesRepositoryImpl {
	return &ActivitiesRepositoryImpl{
		db: db,
	}
}

func (r *ActivitiesRepositoryImpl) GetAll(ctx context.Context) ([]*entities.Activity, error) {
	return nil, nil
}

func (r *ActivitiesRepositoryImpl) Save(ctx context.Context, work *entities.Activity) error {
	return nil
}
