package services

import (
	"context"
	"errors"
	"testing"

	"github.com/edy4c7/works-uploader/internal/common/constants"
	"github.com/edy4c7/works-uploader/internal/entities"
	myErr "github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var activitiesTestData = []*entities.Activity{
	{
		ID:     1,
		Type:   constants.ActivityAdded,
		UserID:   "hogeuser",
		WorkID: 2,
		Work: &entities.Work{
			ID:           2,
			Title:        "hogetitle",
			Description:  "hogedesc",
			ThumbnailURL: "https://example.com/hoge/thumb",
			ContentURL:   "https://example.com/hoge",
		},
	},
	{
		ID:     3,
		Type:   constants.ActivityUpdated,
		UserID:   "fugauser",
		WorkID: 4,
		Work: &entities.Work{
			ID:           4,
			Title:        "fugatitle",
			Description:  "fugadesc",
			ThumbnailURL: "https://example.com/fuga/thumb",
			ContentURL:   "https://example.com/fuga",
		},
	},
}

func TestNewAcitivitiesServiceImpl(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		actRepo := mocks.NewMockActivitiesRepository(ctrl)

		service := NewActivitiesServiceImpl(actRepo)

		assert.Same(t, actRepo, service.repository)
	})

	t.Run("repository is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewActivitiesServiceImpl(nil)
		})
	})
}

func TestGetAllActivities(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		repo := mocks.NewMockActivitiesRepository(ctrl)
		expect := activitiesTestData
		limit := 20
		repo.EXPECT().GetAll(ctx, limit).Return(expect, nil)

		service := &ActivitiesServiceImpl{
			repository: repo,
		}

		actual, err := service.GetAll(ctx)

		assert.Equal(t, expect, actual)
		assert.Nil(t, err)
	})

	t.Run("is error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		repo := mocks.NewMockActivitiesRepository(ctrl)
		errExpect := errors.New("error")
		repo.EXPECT().GetAll(ctx, gomock.Any()).Return(nil, errExpect)

		service := &ActivitiesServiceImpl{
			repository: repo,
		}

		result, errActual := service.GetAll(ctx)

		assert.Nil(t, result)
		assert.True(t, errors.Is(errActual, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(errActual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", errActual)
		}
	})
}

func TestFindActivitiesByUserID(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		repo := mocks.NewMockActivitiesRepository(ctrl)
		expect := []*entities.Activity{
			activitiesTestData[0],
		}
		userID := "user"
		limit := 10
		repo.EXPECT().FindByUserID(ctx, userID, limit).Return(expect, nil)

		service := &ActivitiesServiceImpl{
			repository: repo,
		}

		actual, err := service.FindByUserID(ctx, userID)

		assert.Equal(t, expect, actual)
		assert.Nil(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		repo := mocks.NewMockActivitiesRepository(ctrl)
		errExpect := myErr.NewRecordNotFoundError("", nil)
		userID := "user"
		repo.EXPECT().FindByUserID(ctx, userID, gomock.Any()).Return(nil, errExpect)

		service := &ActivitiesServiceImpl{
			repository: repo,
		}

		result, errActual := service.FindByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.True(t, errors.Is(errActual, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(errActual, &appErr) {
			assert.Equal(t, myErr.WUE01, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", errActual)
		}
	})

	t.Run("is error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		repo := mocks.NewMockActivitiesRepository(ctrl)
		errExpect := errors.New("error")
		userID := "user"
		repo.EXPECT().FindByUserID(ctx, userID, gomock.Any()).Return(nil, errExpect)

		service := &ActivitiesServiceImpl{
			repository: repo,
		}

		result, errActual := service.FindByUserID(ctx, userID)

		assert.Nil(t, result)
		assert.True(t, errors.Is(errActual, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(errActual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", errActual)
		}
	})
}
