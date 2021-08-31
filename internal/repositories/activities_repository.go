package repositories

import (
	"context"

	"github.com/edy4c7/works-uploader/internal/entities"
)

type ActivitiesRepository interface {
	GetAll(context.Context, int) ([]*entities.Activity, error)
	FindByUserID(context.Context, string, int) ([]*entities.Activity, error)
	Create(context.Context, *entities.Activity) error
}
