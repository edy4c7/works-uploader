package repositories

import (
	"context"

	"github.com/edy4c7/darkpot-school-works/internal/entities"
)

type ActivitiesRepository interface {
	GetAll(context.Context) ([]*entities.Activity, error)
	Save(context.Context, *entities.Activity) error
}
