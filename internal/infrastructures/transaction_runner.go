package infrastructures

import (
	"context"

	"github.com/edy4c7/darkpot-school-works/internal/repositories"
	"gorm.io/gorm"
)

type TransactionKeyType string

const transactionKey TransactionKeyType = "transaction"

type TransactionRunnerImpl struct {
	db *gorm.DB
}

func NewTransactionRunnerImpl(db *gorm.DB) *TransactionRunnerImpl {
	return &TransactionRunnerImpl{
		db: db,
	}
}

func (r *TransactionRunnerImpl) Run(ctx context.Context, tranFunc repositories.TransactionFunction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tranFunc(context.WithValue(ctx, transactionKey, tx.WithContext(ctx)))
	})
}
