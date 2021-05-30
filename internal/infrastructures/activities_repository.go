package infrastructures

import (
	"context"
	"errors"

	"github.com/edy4c7/works-uploader/internal/entities"
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
	acts := make([]*entities.Activity, 0)
	r.db.WithContext(ctx).Preload("Work").Find(&acts)
	err := r.db.Error
	return acts, err
}

func (r *ActivitiesRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*entities.Activity, error) {
	acts := make([]*entities.Activity, 0)
	r.db.WithContext(ctx).Preload("Work").Where("user = ?", userID).Find(&acts)
	err := r.db.Error
	return acts, err
}

func (r *ActivitiesRepositoryImpl) Create(ctx context.Context, act *entities.Activity) error {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.Create(act).Error
	}

	return errors.New(notInTransactionMessage)
}
