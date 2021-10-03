package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edy4c7/works-uploader/internal/entities"
	"github.com/edy4c7/works-uploader/internal/mocks"
	"github.com/edy4c7/works-uploader/internal/common/constants"
	"github.com/gin-gonic/gin"
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

func TestNewActivitiesController(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMockActivitiesService(ctrl)
		actCtrl := NewActivitiesController(service)

		assert.Same(t, service, actCtrl.service)
	})

	t.Run("service is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		assert.Panics(t, func() {
			NewActivitiesController(nil)
		})
	})
}

func TestGetAllActivities(t *testing.T) {
	const path = "/"

	t.Run("is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockActivitiesService(ctrl)
		service.EXPECT().GetAll(ctx).Return(activitiesTestData, nil)
		actCtrl := NewActivitiesController(service)
		r.GET("/", actCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, path, nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(activitiesTestData)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		errExpect := errors.New("error")
		service := mocks.NewMockActivitiesService(ctrl)
		service.EXPECT().GetAll(ctx).Return(nil, errExpect)
		actCtrl := NewActivitiesController(service)
		r.GET("/", actCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, path, nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestFindActivitiesByUserID(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockActivitiesService(ctrl)
		userid := "user"
		expect := []*entities.Activity{
			activitiesTestData[0],
		}
		service.EXPECT().FindByUserID(ctx, userid).Return(expect, nil)
		actCtrl := NewActivitiesController(service)
		r.GET("/", actCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, path+"?user="+userid, nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(expect)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("service error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockActivitiesService(ctrl)
		errExpect := errors.New("error")
		service.EXPECT().FindByUserID(ctx, gomock.Any()).Return(nil, errExpect)
		actCtrl := NewActivitiesController(service)
		r.GET("/", actCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, path+"?user=hogeuser", nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}
