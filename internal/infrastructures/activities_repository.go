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

func (r *ActivitiesRepositoryImpl) GetAll(ctx context.Context, limit int) ([]*entities.Activity, error) {
	acts := make([]*entities.Activity, 0)
	err := r.db.WithContext(ctx).Preload("Work").Limit(limit).Order("created_at desc").Find(&acts).Error
	return acts, err
}

func (r *ActivitiesRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit int) ([]*entities.Activity, error) {
	acts := make([]*entities.Activity, 0)
	err := r.db.WithContext(ctx).Preload("Work").Limit(limit).Where("activities.user = ?", userID).Order("created_at desc").Find(&acts).Error
	return acts, err
}

func (r *ActivitiesRepositoryImpl) Create(ctx context.Context, act *entities.Activity) error {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.WithContext(ctx).Create(act).Error
	}

	return errors.New(notInTransactionMessage)
}
