package services

import (
	"context"
	"errors"

	"github.com/edy4c7/works-uploader/internal/entities"
	myErr "github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/repositories"
)

type ActivitiesService interface {
	GetAll(context.Context) ([]*entities.Activity, error)
	FindByUserID(context.Context, string) ([]*entities.Activity, error)
}

type ActivitiesServiceImpl struct {
	repository repositories.ActivitiesRepository
}

func NewActivitiesServiceImpl(repo repositories.ActivitiesRepository) *ActivitiesServiceImpl {
	if repo == nil {
		panic("repository can't be null")
	}
	return &ActivitiesServiceImpl{
		repository: repo,
	}
}

func (r *ActivitiesServiceImpl) GetAll(ctx context.Context) ([]*entities.Activity, error) {
	result, err := r.repository.GetAll(ctx)
	if err != nil {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}
	return result, nil
}

func (r *ActivitiesServiceImpl) FindByUserID(ctx context.Context, userID string) ([]*entities.Activity, error) {
	result, err := r.repository.FindByUserID(ctx, userID)
	if err != nil {
		var rnfErr *myErr.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE01), myErr.Cause(err))
		}

		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	return result, nil
}
