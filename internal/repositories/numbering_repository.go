package repositories

import (
	"context"

	"github.com/edy4c7/darkpot-school-works/internal/entities"
)

type NumberingRepository interface {
	FindByType(ctx context.Context, typ int) (*entities.Numbering, error)
	Save(ctx context.Context, e *entities.Numbering) error
}
