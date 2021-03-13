package repositories

import "context"

type TransactionFunction func(context.Context) error

type TransactionRunner interface {
	Run(context.Context, TransactionFunction) error
}
