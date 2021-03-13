package infrastructures

import (
	"context"
	"errors"

	"github.com/edy4c7/darkpot-school-works/internal/entities"
	"gorm.io/gorm"
)

type WorksRepositoryImpl struct {
	db *gorm.DB
}

func NewWorksRepositoryImpl(db *gorm.DB) *WorksRepositoryImpl {
	return &WorksRepositoryImpl{
		db: db,
	}
}

func (r *WorksRepositoryImpl) GetAll(ctx context.Context) ([]*entities.Work, error) {
	works := make([]*entities.Work, 0)
	r.db.WithContext(ctx).Find(&works)
	err := r.db.Error
	return works, err
}

func (r *WorksRepositoryImpl) FindByID(ctx context.Context, id uint64) (*entities.Work, error) {
	var work entities.Work
	r.db.WithContext(ctx).First(&work, id)
	err := r.db.Error
	return &work, err
}

func (r *WorksRepositoryImpl) Save(ctx context.Context, work *entities.Work) error {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.Save(work).Error
	}
	return errors.New("Not in transaction")
}

func (r *WorksRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx.Delete(&entities.Work{}, id).Error
	}
	return errors.New("Not in transaction")
}
