package infrastructures

import (
	"context"
	"errors"

	"github.com/edy4c7/works-uploader/internal/entities"
	wuErr "github.com/edy4c7/works-uploader/internal/errors"
	"gorm.io/gorm"
)

const notInTransactionMessage = "not in transaction"

type WorksRepositoryImpl struct {
	db *gorm.DB
}

func NewWorksRepositoryImpl(db *gorm.DB) *WorksRepositoryImpl {
	return &WorksRepositoryImpl{
		db: db,
	}
}

func (r *WorksRepositoryImpl) GetAll(ctx context.Context, offset int, limit int) ([]*entities.Work, error) {
	works := make([]*entities.Work, 0)
	err := r.db.WithContext(ctx).Preload("Author").Offset(offset).Limit(limit).Find(&works).Error
	return works, err
}

func (r *WorksRepositoryImpl) CountAll(ctx context.Context) (int64, error) {
	works := make([]*entities.Work, 0)
	var count int64
	err := r.db.WithContext(ctx).Find(&works).Count(&count).Error
	return count, err
}

func (r *WorksRepositoryImpl) FindByID(ctx context.Context, id uint64) (*entities.Work, error) {
	var work entities.Work
	err := r.db.WithContext(ctx).Preload("Author").First(&work, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, wuErr.NewRecordNotFoundError(err.Error(), err)
	}
	return &work, err
}

func (r *WorksRepositoryImpl) Create(ctx context.Context, work *entities.Work) error {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.WithContext(ctx).Create(work).Error
	}
	return errors.New(notInTransactionMessage)
}

func (r *WorksRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.WithContext(ctx).Delete(&entities.Work{}, id).Error
	}
	return errors.New(notInTransactionMessage)
}
