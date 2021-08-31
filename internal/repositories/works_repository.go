package repositories

import (
	"context"

	"github.com/edy4c7/works-uploader/internal/entities"
)

type WorksRepository interface {
	GetAll(ctx context.Context, offset int, limit int) ([]*entities.Work, error)
	CountAll(context.Context) (int64, error)
	FindByID(context.Context, uint64) (*entities.Work, error)
	Save(context.Context, *entities.Work) error
	DeleteByID(context.Context, uint64) error
}
